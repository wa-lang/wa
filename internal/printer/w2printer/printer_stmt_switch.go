// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) stmtSwitch(s *ast.SwitchStmt) {
	p.print(token.Zh_找辙)
	p.stmtSwitch_controlClause(s.Init, s.Tag, nil)
	p.block(s.Body, 0, token.Zh_完毕)
}

func (p *printer) stmtSwitchType(s *ast.TypeSwitchStmt) {
	p.print(token.Zh_找辙)
	if s.Init != nil {
		p.print(blank)
		p.stmt(s.Init, false)
		p.print(token.SEMICOLON)
	}
	p.print(blank)
	p.stmt(s.Assign, false)
	p.print(blank)
	p.block(s.Body, 0, token.Zh_完毕)
}

func (p *printer) stmtSwitch_controlClause(init ast.Stmt, expr ast.Expr, post ast.Stmt) {
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
