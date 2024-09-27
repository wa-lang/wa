// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/scanner"
	"wa-lang.org/wa/internal/wat/token"
)

type parser struct {
	filename string
	src      []byte

	fset    *token.File
	scanner scanner.Scanner
	module  *ast.Module

	pos token.Pos
	tok token.Token
	lit string

	err error

	trace bool
}

type parserError struct {
	pos token.Position
	msg string
}

func (e *parserError) Error() string {
	return fmt.Sprintf("%v: %s", e.pos, e.msg)
}

func newParser(path string, src []byte) *parser {
	p := &parser{
		filename: path,
		src:      src,
		fset:     token.NewFile(path, len(src)),
	}

	return p
}

// 解析 wat 模块
func (p *parser) ParseModule() (m *ast.Module, err error) {
	p.module = &ast.Module{}

	if p.trace {
		fmt.Println("filename:", p.filename)
		fmt.Printf("code: %s\n", string(p.src))
	}

	defer func() {
		if p.trace {
			// 调试模型不捕获异常
			return
		}
		if r := recover(); r != nil {
			if errx, ok := r.(*parserError); ok {
				m = p.module
				err = errx
			} else {
				panic(r)
			}
		}
	}()

	p.scanner.Init(p.fset, p.src,
		func(pos token.Position, msg string) {
			if p.err != nil {
				p.err = &parserError{pos, msg}
			}
		},
		scanner.ScanComments,
	)

	p.parseFile()

	m = p.module
	err = p.err
	return
}

// 报错, 只报第一个错误
func (p *parser) errorf(pos token.Pos, format string, a ...interface{}) {
	p.err = &parserError{
		pos: p.fset.Position(pos),
		msg: fmt.Sprintf(format, a...),
	}
	panic(p.err)
}

// 下一个 token
func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
	if p.trace {
		fmt.Println(p.fset.Position(p.pos), p.tok, p.lit)
	}
}

// 吃掉一个预期的 token
func (p *parser) acceptToken(expectToken token.Token, moreExpectTokens ...token.Token) {
	if p.tok == expectToken {
		p.next()
		return
	}
	for _, tok := range moreExpectTokens {
		if p.tok == tok {
			p.next()
			return
		}
	}

	if len(moreExpectTokens) > 0 {
		all := append([]token.Token{expectToken}, moreExpectTokens...)
		p.errorf(p.pos, "expect %v, got %v", all, p.tok)
	} else {
		p.errorf(p.pos, "expect %v, got %v", expectToken, p.tok)
	}
}

// 解析 wat 文件
func (p *parser) parseFile() {
	// 读取第一个Token
	p.next()

	// 解析 module
	p.parseModule()

	// 忽略尾部注释
	p.consumeComments()

	// 结束解析
	p.acceptToken(token.EOF)
}

// 跳过注释
func (p *parser) consumeComments() {
	for p.tok == token.COMMENT {
		p.next()
	}
}

// 解析标别符
func (p *parser) parseIdent() string {
	s := p.lit
	p.acceptToken(token.IDENT)
	return s
}

// 解析整数常量面值解析整数常量面值
func (p *parser) parseIntLit() int {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseInt(lit[2:], 16, 64)
		if err != nil {
			p.errorf(pos, "expect int, got %q", lit)
		}
		return int(n)
	}

	n, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		p.errorf(pos, "expect int, got %q", lit)
	}
	return int(n)
}
func (p *parser) parseInt32Lit() int32 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseInt(lit[2:], 16, 32)
		if err != nil {
			p.errorf(pos, "expect int32, got %q", lit)
		}
		return int32(n)
	}

	// 需要支持 u32 和 -1 两种格式
	n, err := strconv.ParseInt(lit, 10, 32)
	if err != nil {
		if n, errx := strconv.ParseUint(lit, 10, 32); errx == nil {
			return int32(uint32(n))
		} else {
			p.errorf(pos, "expect int32, got %q, err = %v", lit, err)
		}

	}
	return int32(n)
}
func (p *parser) parseUint32Lit() uint32 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseInt(lit[2:], 16, 32)
		if err != nil {
			p.errorf(pos, "expect uint32, got %q", lit)
		}
		return uint32(n)
	}

	n, err := strconv.ParseUint(lit, 10, 32)
	if err != nil {
		p.errorf(pos, "expect uint32, got %q", lit)
	}
	return uint32(n)
}
func (p *parser) parseInt64Lit() int64 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseInt(lit[2:], 16, 64)
		if err != nil {
			p.errorf(pos, "expect int64, got %q", lit)
		}
		return int64(n)
	}

	n, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		p.errorf(pos, "expect int64, got %q", lit)
	}
	return int64(uint64(n))
}
func (p *parser) parseUint64Lit() uint64 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseInt(lit[2:], 16, 64)
		if err != nil {
			p.errorf(pos, "expect uint64, got %q", lit)
		}
		return uint64(n)
	}

	n, err := strconv.ParseUint(lit, 10, 64)
	if err != nil {
		p.errorf(pos, "expect uint64, got %q", lit)
	}
	return n
}

// 解析浮点数常量面值
func (p *parser) parseFloat32Lit() float32 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.FLOAT, token.INT)

	n, err := strconv.ParseFloat(lit, 32)
	if err != nil {
		p.errorf(pos, "expect float32, got %q", lit)
	}
	return float32(n)
}
func (p *parser) parseFloat64Lit() float64 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.FLOAT, token.INT)

	n, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		p.errorf(pos, "expect float64, got %q", lit)
	}
	return n
}

// 解析字符常量面值
func (p *parser) parseCharLit() rune {
	s := []rune(p.lit)
	p.acceptToken(token.CHAR)
	return s[0]
}

// 解析字符串常量面值(含二进制数据)
func (p *parser) parseStringLit() string {
	s := p.lit
	p.acceptToken(token.STRING)
	return s
}

// 解析索引, 标识符或整数
func (p *parser) parseIdentOrIndex() string {
	s := p.lit
	p.acceptToken(token.IDENT, token.INT)
	return s
}

func (p *parser) parseIdentOrIndexList() (ss []string) {
	for {
		if p.tok == token.IDENT || p.tok == token.INT {
			ss = append(ss, p.lit)
			p.acceptToken(token.IDENT, token.INT)
			continue
		}
		break
	}
	if len(ss) == 0 {
		p.errorf(p.pos, "expect token.IDENT or token.INT, got %q", p.tok)
	}
	return
}

func (p *parser) parseNumberType() token.Token {
	switch p.tok {
	case token.I32, token.I64, token.F32, token.F64:
		tok := p.tok
		p.next()
		return tok
	default:
		p.errorf(p.pos, "export %v, got %v", "i32|i64|f32|f64", p.tok)
	}
	panic("unreachable")
}

func (p *parser) parseNumberTypeList() []token.Token {
	var tokens []token.Token

Loop:
	for {
		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			tokens = append(tokens, p.tok)
			p.next()
		default:
			break Loop
		}
	}
	if len(tokens) == 0 {
		p.errorf(p.pos, "export %v, got %v", "i32|i64|f32|f64", p.tok)
	}

	return tokens
}
