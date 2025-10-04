// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseValueSpec_zh(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec {
	if p.trace {
		defer un(trace(p, keyword.String()+"Spec"))
	}

	pos := p.pos
	idents := p.parseIdentList()
	var colonPos token.Pos
	if p.tok == token.COLON {
		colonPos = p.pos
		p.next()
	} else {
		if p.tok != token.ASSIGN && p.tok != token.SEMICOLON {
			p.expect(token.COLON)
		}
	}

	typ := p.tryType()
	var values []ast.Expr
	// always permit optional initialization for more tolerant parsing
	if p.tok == token.ASSIGN {
		p.next()
		values = p.parseRhsList()
	}
	p.expectSemi() // call before accessing p.linecomment

	switch keyword {
	case token.VAR, token.GLOBAL, token.Zh_全局:
		if typ == nil && values == nil {
			p.error(pos, "缺少变量的类型或初始值")
		}
	case token.CONST, token.Zh_定义, token.Zh_常量:
		if typ == nil && values == nil && iota == 0 {
			p.error(pos, "缺少常量的类型或初始值")
		}

		// if values == nil && (iota == 0 || typ != nil) {
		//     p.error(pos, "missing constant value")
		// }
	}

	// Wa spec: The scope of a constant or variable identifier declared inside
	// a function begins at the end of the ConstSpec or VarSpec and ends at
	// the end of the innermost containing block.
	// (Global identifiers are resolved in a separate phase after parsing.)
	spec := &ast.ValueSpec{
		Doc:      doc,
		Names:    idents,
		ColonPos: colonPos,
		Type:     typ,
		Values:   values,
		Comment:  p.lineComment,
	}
	kind := ast.Con
	if keyword == token.VAR || keyword == token.GLOBAL || keyword == token.Zh_全局 {
		kind = ast.Var
	}
	p.declare(spec, iota, p.topScope, kind, idents...)

	return spec
}
