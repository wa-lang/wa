// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) stmtIf(s *ast.IfStmt, nextIsRBrace bool) {
	p.print(s.Tok)
	p.stmtIf_controlClause(s.Init, s.Cond, nil)
	p.stmtIf_bodyBlock(s.Body, s.Else, 1, nextIsRBrace)
}

func (p *printer) stmtIf_controlClause(init ast.Stmt, expr ast.Expr, post ast.Stmt) {
	p.print(blank)
	if init == nil && post == nil {
		// no semicolons required
		if expr != nil {
			p.expr(stripParens(expr))
		}
	} else {
		// all semicolons required
		// (they are not separators, print them explicitly)
		if init != nil {
			p.stmt(init, false)
		}
		p.print(token.SEMICOLON, blank)
		if expr != nil {
			p.expr(stripParens(expr))
		}
	}
}

// block prints an *ast.BlockStmt; it always spans at least two lines.
func (p *printer) stmtIf_bodyBlock(body *ast.BlockStmt, else_ ast.Stmt, nindent int, nextIsRBrace bool) {
	p.print(body.Lbrace, token.COLON)
	p.stmtList(body.List, nindent, true)
	p.linebreak(p.lineFor(body.Rbrace), 1, ignore, true)

	if else_ != nil {
		switch s := else_.(type) {
		case *ast.IfStmt:
			p.stmtIf(s, nextIsRBrace)
		case *ast.BlockStmt:
			p.print(token.Zh_否则)
			p.stmtIf_elseBlock(s, nindent)
		default:
			panic("unreachable")
		}
	} else {
		p.print(body.Rbrace, token.Zh_完毕)
	}
}

func (p *printer) stmtIf_elseBlock(b *ast.BlockStmt, nindent int) {
	p.print(b.Lbrace, token.COLON)
	p.stmtList(b.List, nindent, true)
	p.linebreak(p.lineFor(b.Rbrace), 1, ignore, true)
	p.print(b.Rbrace, token.Zh_完毕)
}
