// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import "wa-lang.org/wa/internal/token"

// A bailout panic is raised to indicate early termination.
type bailout struct{}

func (p *parser) error(pos token.Pos, msg string) {
	epos := p.file.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	if p.mode&AllErrors == 0 {
		n := len(p.errors)
		if n > 0 && p.errors[n-1].Pos.Line == epos.Line {
			return // discard - likely a spurious error
		}
		if n > 10 {
			panic(bailout{})
		}
	}

	p.errors.Add(epos, msg)
}

func (p *parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			msg += ", found newline"
		case p.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

// expectClosing is like expect but provides a better error message
// for the common case of a missing comma before a newline.
func (p *parser) expectClosing(tok token.Token, context string) token.Pos {
	if p.tok != tok && p.tok == token.SEMICOLON && p.lit == "\n" {
		p.error(p.pos, "missing ',' before newline in "+context)
		p.next()
	}
	return p.expect(tok)
}

func (p *parser) expectSemi() {
	// semicolon is optional before a closing ')' or '完毕'
	if p.tok != token.RPAREN && p.tok != token.Zh_完毕 {
		switch p.tok {
		case token.COMMA:
			// permit a ',' instead of a ';' but complain
			p.errorExpected(p.pos, "';'")
			fallthrough
		case token.SEMICOLON:
			p.next()
		default:
			p.errorExpected(p.pos, "';'")
			p.advance(stmtStart)
		}
	}
}
