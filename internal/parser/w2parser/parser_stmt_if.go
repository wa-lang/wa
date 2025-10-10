// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseIfStmt(keyword token.Token) *ast.IfStmt {
	if p.trace {
		defer un(trace(p, "IfStmt"))
	}

	pos := p.expect(keyword)
	p.openScope()
	defer p.closeScope()

	init, cond := p.parseIfHeader()
	body := p.parseBlockStmt(token.Zh_完毕, token.Zh_或者, token.Zh_否则)

	var else_ ast.Stmt
	switch p.tok {
	case token.Zh_或者: // else if
		else_ = p.parseIfStmt(p.tok)
	case token.Zh_否则: // else
		p.expect(p.tok)
		else_ = p.parseBlockStmt(token.Zh_完毕)
		p.expectSemi()
	default:
		p.errorExpected(p.pos, "if statement or block:"+p.tok.String())
		else_ = &ast.BadStmt{From: p.pos, To: p.pos}
	}

	return &ast.IfStmt{TokePos: pos, Tok: keyword, Init: init, Cond: cond, Body: body, Else: else_}
}

// parseIfHeader is an adjusted version of parser.header
// in cmd/compile/internal/syntax/parser.go, which has
// been tuned for better error handling.
func (p *parser) parseIfHeader() (init ast.Stmt, cond ast.Expr) {
	if p.tok == token.COLON {
		p.error(p.pos, "missing condition in if statement")
		cond = &ast.BadExpr{From: p.pos, To: p.pos}
		return
	}

	outer := p.exprLev
	p.exprLev = -1

	if p.tok != token.SEMICOLON {
		init, _ = p.parseSimpleStmt(token.Zh_如果, basic)
	}

	var condStmt ast.Stmt
	var semi struct {
		pos token.Pos
		lit string // ";" or "\n"; valid if pos.IsValid()
	}
	if p.tok != token.COLON {
		if p.tok == token.SEMICOLON {
			semi.pos = p.pos
			semi.lit = p.lit
			p.next()
		} else {
			p.expect(token.SEMICOLON)
		}
		if p.tok != token.COLON {
			condStmt, _ = p.parseSimpleStmt(token.Zh_如果, basic)
		}
	} else {
		condStmt = init
		init = nil
	}

	if condStmt != nil {
		cond = p.makeExpr(condStmt, "boolean expression")
	} else if semi.pos.IsValid() {
		if semi.lit == "\n" {
			p.error(semi.pos, "unexpected newline, expecting { after if clause")
		} else {
			p.error(semi.pos, "missing condition in if statement")
		}
	}

	// make sure we have a valid AST
	if cond == nil {
		cond = &ast.BadExpr{From: p.pos, To: p.pos}
	}

	p.exprLev = outer
	return
}
