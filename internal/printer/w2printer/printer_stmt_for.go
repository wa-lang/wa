// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) stmtFor(s *ast.ForStmt) {
	p.print(token.Zh_循环)
	p.stmtFor_controlClause(s.Init, s.Cond, s.Post)
	p.block(s.Body, 1, token.Zh_完毕)
}

func (p *printer) stmtForRange(s *ast.RangeStmt) {
	p.print(token.Zh_循环, blank)
	if s.Key != nil {
		p.expr(s.Key)
		if s.Value != nil {
			// use position of value following the comma as
			// comma position for correct comment placement
			p.print(s.Value.Pos(), token.COMMA, blank)
			p.expr(s.Value)
		}
		p.print(blank, s.TokPos, s.Tok, blank)
	}
	p.print(token.Zh_迭代, blank)
	p.expr(stripParens(s.X))
	p.block(s.Body, 1, token.Zh_完毕)
}

func (p *printer) stmtFor_controlClause(init ast.Stmt, expr ast.Expr, post ast.Stmt) {
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

		p.print(token.SEMICOLON, blank)

		if post != nil {
			p.stmt(post, false)
		}
	}
}
