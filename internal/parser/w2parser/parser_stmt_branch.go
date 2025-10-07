// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseBranchStmt(tok token.Token) *ast.BranchStmt {
	if p.trace {
		defer un(trace(p, "BranchStmt"))
	}

	pos := p.expect(tok)
	var label *ast.Ident
	if p.tok == token.IDENT {
		label = p.parseIdent()
		// add to list of unresolved targets
		n := len(p.targetStack) - 1
		p.targetStack[n] = append(p.targetStack[n], label)
	}
	p.expectSemi()

	return &ast.BranchStmt{TokPos: pos, Tok: tok, Label: label}
}
