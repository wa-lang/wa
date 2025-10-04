// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseTypeSpec_zh(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
	if p.trace {
		defer un(trace(p, "TypeSpec"))
	}

	ident := p.parseIdent()

	// Go spec: The scope of a type identifier declared inside a function begins
	// at the identifier in the TypeSpec and ends at the end of the innermost
	// containing block.
	// (Global identifiers are resolved in a separate phase after parsing.)
	spec := &ast.TypeSpec{Doc: doc, Name: ident}
	p.declare(spec, nil, p.topScope, ast.Typ, ident)
	if p.tok == token.COLON {
		spec.ColonPos = p.pos
		p.next()
	}
	spec.Type = p.parseType_zh()
	p.expectSemi_zh() // call before accessing p.linecomment
	spec.Comment = p.lineComment

	return spec
}

func (p *parser) parseType_zh() ast.Expr {
	if p.trace {
		defer un(trace(p, "类型"))
	}

	typ := p.tryType()

	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "类型")
		p.advance(exprEnd)
		return &ast.BadExpr{From: pos, To: p.pos}
	}

	return typ
}
