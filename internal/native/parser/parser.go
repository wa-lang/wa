// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"fmt"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/scanner"
	"wa-lang.org/wa/internal/native/token"
)

type parser struct {
	filename string
	src      []byte

	fset    *token.FileSet
	file    *token.File
	scanner scanner.Scanner
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

func newParser(fset *token.FileSet, filename string, src []byte) *parser {
	p := &parser{
		fset:     fset,
		file:     fset.AddFile(filename, -1, len(src)),
		filename: filename,
		src:      src,
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

	p.scanner.Init(p.file, p.src,
		func(pos token.Position, msg string) {
			if p.err != nil {
				p.err = &parserError{pos, msg}
			}
		},
		scanner.ScanComments,
	)

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

// 下一个 token
func (p *parser) next() {
	p.next0()
	if p.tok == token.COMMENT {
		p.consumeComments()
	}
}
func (p *parser) next0() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
	if p.trace {
		fmt.Println(p.fset.Position(p.pos), p.tok, p.lit)
	}
}

// 跳过注释
func (p *parser) consumeComments() {
	for p.tok == token.COMMENT {
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

func (p *parser) parseFile() {
	// 读取第一个Token
	p.next()

	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}

		switch p.tok {
		case token.CONST:
			p.parseConst()
		case token.GLOBAL:
			p.parseGlobal()
		case token.FUNC:
			p.parseFunc()
		default:
			p.err = fmt.Errorf("unkonw token: %v", p.tok)
			return
		}
	}
}

func (p *parser) parseConst() {
	// TODO
}

func (p *parser) parseGlobal() {
	// TODO
}
func (p *parser) parseFunc() {
	// TODO
}
