// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"fmt"

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

func (p *parser) ParseModule() (*ast.Module, error) {
	p.module = &ast.Module{
		File: token.NewFile(p.filename, len(p.src)),
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*parserError); !ok {
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

	return p.module, p.err
}

func (p *parser) errorf(pos token.Pos, format string, a ...interface{}) {
	p.err = &parserError{
		pos: p.module.File.Position(pos),
		msg: fmt.Sprintf(format, a...),
	}
	panic(p.err)
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) acceptToken(expectToken token.Token) {
	if p.tok != expectToken {
		p.errorf(p.pos, "expect %v, got %v", expectToken, p.tok)
	}
	p.next()
}

func (p *parser) parseFile() {
	p.next()

	for {
		if p.tok == token.EOF {
			break
		}

		switch p.tok {
		case token.EOF:
			return

		case token.COMMENT:
			p.parseComment()

		case token.LPAREN: // ()
			p.parseModule()

		default:
			p.errorf(p.pos, "bad token: %v, lit: %s", p.tok, p.lit)
		}
	}
}

func (p *parser) parseComment() {
	p.acceptToken(token.COMMENT)
}

func (p *parser) parseModule() {
	p.acceptToken(token.LPAREN)
	p.acceptToken(token.MODULE)

	if p.tok == token.IDENT {
		p.module.Name = p.lit
		p.next()
	}

	p.acceptToken(token.RPAREN)
}
