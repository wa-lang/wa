// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseReturnStmt(keyword token.Token) *ast.ReturnStmt {
	if p.trace {
		defer un(trace(p, "ReturnStmt"))
	}

	pos := p.pos
	p.expect(keyword)
	var x []ast.Expr
	if p.tok != token.SEMICOLON && p.tok != token.Zh_完毕 {
		x = p.parseRhsList()
	}
	p.expectSemi()

	return &ast.ReturnStmt{TokPos: pos, Tok: keyword, Results: x}
}
