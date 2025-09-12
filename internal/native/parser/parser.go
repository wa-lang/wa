// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/scanner"
	"wa-lang.org/wa/internal/native/token"
)

type parser struct {
	cpu      abi.CPUType
	filename string
	src      []byte

	fset    *token.FileSet
	file    *token.File
	scanner *scanner.Scanner
	prog    *ast.File

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

func newParser(cpu abi.CPUType, fset *token.FileSet, filename string, src []byte) *parser {
	p := &parser{
		cpu:      cpu,
		fset:     fset,
		file:     fset.AddFile(filename, -1, len(src)),
		filename: filename,
		src:      src,
	}

	switch cpu {
	case abi.RISCV32, abi.RISCV64:
		p.scanner = scanner.NewScanner(
			func(ident string) token.Token {
				// 将原始的寄存器映射到 token.Token 编码
				if reg, ok := riscv.LookupRegister(ident); ok {
					return token.REG_RISCV_BEGIN + token.Token(reg)
				}
				// 将原始的指令映射到 token.Token 编码
				if reg, ok := riscv.LookupAs(ident); ok {
					return token.A_RISCV_BEGIN + token.Token(reg)
				}
				return token.NONE
			},
		)
	default:
		panic(fmt.Errorf("unknown cpu: %v", cpu))
	}

	p.scanner.Init(p.file, p.src,
		func(pos token.Position, msg string) {
			if p.err != nil {
				p.err = &parserError{pos, msg}
			}
		},
		scanner.ScanComments,
	)
	return p
}

func (p *parser) ParseFile() (prog *ast.File, err error) {
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
				prog = p.prog
				err = errx
			} else {
				panic(r)
			}
		}
	}()

	p.prog = &ast.File{CPU: p.cpu}

	p.scanner.Init(p.file, p.src,
		func(pos token.Position, msg string) {
			if p.err != nil {
				p.err = &parserError{pos, msg}
			}
		},
		scanner.ScanComments,
	)

	// 读取第一个Token
	p.next()

	p.parseFile()

	prog = p.prog
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

// pos 转行号
func (p *parser) posLine(pos token.Pos) int {
	return p.fset.Position(pos).Line
}

// 下一个 token
func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
	if p.trace {
		fmt.Println(p.fset.Position(p.pos), p.tok, p.lit)
	}
}

// 跳过分号列表
func (p *parser) consumeSemicolonList() {
	for p.tok == token.SEMICOLON {
		p.next()
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

// 吃掉其中一个
func (p *parser) acceptTokenAorB(expectTokenA, expectTokenB token.Token) {
	assert(expectTokenA != expectTokenB)
	switch p.tok {
	case expectTokenA, expectTokenB:
		p.next()
	default:
		p.errorf(p.pos, "expect %v or %v, got %v", expectTokenA, expectTokenB, p.tok)
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

	if p.tok == token.CHAR {
		p.acceptToken(token.CHAR)
		return int32(lit[1]) // '?'
	}

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

	if p.tok == token.CHAR {
		p.acceptToken(token.CHAR)
		return uint32(lit[1]) // '?'
	}

	p.acceptToken(token.INT)

	if len(lit) > 2 && lit[0] == '0' && (lit[1] == 'x' || lit[1] == 'X') {
		n, err := strconv.ParseUint(lit[2:], 16, 32)
		if err != nil {
			p.errorf(pos, "expect int32, got %q", lit)
		}
		return uint32(n)
	}

	// 需要支持 u32 和 -1 两种格式
	n, err := strconv.ParseUint(lit, 10, 32)
	if err != nil {
		if n, errx := strconv.ParseUint(lit, 10, 32); errx == nil {
			return uint32(n)
		} else {
			p.errorf(pos, "expect int32, got %q, err = %v", lit, err)
		}

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
		n, err := strconv.ParseUint(lit[2:], 16, 64)
		if err != nil {
			p.errorf(pos, "expect int64, got %q", lit)
		}
		return uint64(n)
	}

	n, err := strconv.ParseUint(lit, 10, 64)
	if err != nil {
		p.errorf(pos, "expect int64, got %q", lit)
	}
	return uint64(n)
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

// 解析字符常量面值(含二进制数据)
func (p *parser) parseCharLit() string {
	s := p.lit
	p.acceptToken(token.CHAR)
	return s
}

// 解析字符串常量面值(含二进制数据)
func (p *parser) parseStringLit() string {
	s := p.lit
	p.acceptToken(token.STRING)
	return s
}

// 解析寄存器
func (p *parser) parseRegister() abi.RegType {
	if !p.tok.IsRegister() {
		p.errorf(p.pos, "expect register, got %v", p.tok)
	}
	reg := p.tok.RawReg()
	p.acceptToken(p.tok)
	return reg
}

// 解析指令
func (p *parser) parseAs() abi.As {
	if !p.tok.IsAs() {
		p.errorf(p.pos, "expect as, got %v", p.tok)
	}
	as := p.tok.RawAs()
	p.acceptToken(p.tok)
	return as
}
