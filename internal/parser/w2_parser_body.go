// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseBody_zh(keyword token.Token, scope *ast.Scope) *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, keyword.String()))
	}

	lbrace := token.NoPos
	p.topScope = scope // open function scope
	p.openLabelScope()
	list := p.parseStmtList_zh()
	p.closeLabelScope()
	p.closeScope()

	rbrace := token.NoPos
	switch keyword {
	case token.Zh_算始:
		rbrace = p.expect(token.Zh_算终)
	case token.Zh_函始:
		rbrace = p.expect(token.Zh_函终)
	default:
		panic("unreachable")
	}

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}
