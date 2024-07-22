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

	scanner scanner.Scanner
	module  *ast.Module

	pos token.Pos
	tok token.Token
	lit string

	err error
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
	}

	return p
}

// 解析 wat 模块
func (p *parser) ParseModule() (m *ast.Module, err error) {
	p.module = &ast.Module{
		File: token.NewFile(p.filename, len(p.src)),
	}

	defer func() {
		if r := recover(); r != nil {
			if errx, ok := r.(*parserError); ok {
				m = p.module
				err = errx
			} else {
				panic(r)
			}
		}
	}()

	p.scanner.Init(p.module.File, p.src,
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
		pos: p.module.File.Position(pos),
		msg: fmt.Sprintf(format, a...),
	}
	panic(p.err)
}

// 下一个 token
func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
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

	// 解析开头的注释
	p.consumeComments()

	// 解析 module 节点
	p.parseModule()

	// 解析结尾的注释
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

// 解析整数常量面值
func (p *parser) parseIntLit() int {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	n, err := strconv.Atoi(lit)
	if err != nil {
		p.errorf(pos, "expect int, got %q", lit)
	}
	return n
}

// 解析浮点数常量面值
func (p *parser) parseFloatLit() float64 {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.FLOAT, token.INT)

	n, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		p.errorf(pos, "expect int, got %q", lit)
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
	s := p.lit // todo: 解码
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
