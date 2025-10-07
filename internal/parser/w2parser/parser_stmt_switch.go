// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseSwitchStmt(keyword token.Token) ast.Stmt {
	if p.trace {
		defer un(trace(p, "SwitchStmt"))
	}

	pos := p.expect(keyword)
	p.openScope()
	defer p.closeScope()

	var s1, s2 ast.Stmt
	if p.tok != token.COLON {
		prevLev := p.exprLev
		p.exprLev = -1
		if p.tok != token.SEMICOLON {
			s2, _ = p.parseSimpleStmt(basic)
		}
		if p.tok == token.SEMICOLON {
			p.next()
			s1 = s2
			s2 = nil
			if p.tok != token.COLON {
				// A TypeSwitchGuard may declare a variable in addition
				// to the variable declared in the initial SimpleStmt.
				// Introduce extra scope to avoid redeclaration errors:
				//
				//	switch t := 0; t := x.(T) { ... }
				//
				// (this code is not valid Go because the first t
				// cannot be accessed and thus is never used, the extra
				// scope is needed for the correct error message).
				//
				// If we don't have a type switch, s2 must be an expression.
				// Having the extra nested but empty scope won't affect it.
				p.openScope()
				defer p.closeScope()
				s2, _ = p.parseSimpleStmt(basic)
			}
		}
		p.exprLev = prevLev
	}

	typeSwitch := p.isTypeSwitchGuard(s2)
	lbrace := p.expect(token.COLON)
	var list []ast.Stmt
	for p.tok == token.Zh_有辙 || p.tok == token.Zh_没辙 {
		list = append(list, p.parseCaseClause(typeSwitch))
	}
	rbrace := p.expect(token.Zh_完毕)
	p.expectSemi()
	body := &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}

	if typeSwitch {
		return &ast.TypeSwitchStmt{TokPos: pos, Tok: keyword, Init: s1, Assign: s2, Body: body}
	}

	return &ast.SwitchStmt{TokPos: pos, Tok: keyword, Init: s1, Tag: p.makeExpr(s2, "switch expression"), Body: body}
}

func isTypeSwitchAssert(x ast.Expr) bool {
	a, ok := x.(*ast.TypeAssertExpr)
	return ok && a.Type == nil
}

func (p *parser) isTypeSwitchGuard(s ast.Stmt) bool {
	switch t := s.(type) {
	case *ast.ExprStmt:
		// x.(type)
		return isTypeSwitchAssert(t.X)
	case *ast.AssignStmt:
		// v := x.(type)
		if len(t.Lhs) == 1 && len(t.Rhs) == 1 && isTypeSwitchAssert(t.Rhs[0]) {
			switch t.Tok {
			case token.ASSIGN:
				// permit v = x.(type) but complain
				p.error(t.TokPos, "expected ':=', found '='")
				fallthrough
			case token.DEFINE:
				return true
			}
		}
	}
	return false
}

func (p *parser) parseCaseClause(typeSwitch bool) *ast.CaseClause {
	if p.trace {
		defer un(trace(p, "CaseClause"))
	}

	pos := p.pos
	var list []ast.Expr
	var keyword token.Token
	if p.tok == token.Zh_有辙 {
		keyword = p.tok
		p.next()
		if typeSwitch {
			list = p.parseTypeList()
		} else {
			list = p.parseRhsList()
		}
	} else {
		keyword = p.tok
		p.expect(token.Zh_没辙)
	}

	colon := p.expect(token.COLON)
	p.openScope()
	body := p.parseStmtList()
	p.closeScope()

	return &ast.CaseClause{TokPos: pos, Tok: keyword, List: list, Colon: colon, Body: body}
}
