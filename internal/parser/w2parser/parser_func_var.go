// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseGenDecl_var(keyword token.Token) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	var lparen, rparen token.Pos
	var list []ast.Spec
	if p.tok == token.COLON {
		// XXX: ... 完毕
		lparen = p.pos
		p.next()
		for iota := 0; p.tok != token.Zh_完毕 && p.tok != token.EOF; iota++ {
			list = append(list, p.parseValueSpec_var(p.leadComment, keyword, iota))
		}
		rparen = p.expect(token.Zh_完毕)
		p.expectSemi()
	} else {
		list = append(list, p.parseValueSpec_var(nil, keyword, 0))
	}

	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    token.VAR, // 临时方案
		Lparen: lparen,
		Specs:  list,
		Rparen: rparen,
	}
}

func (p *parser) parseValueSpec_var(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec {
	if p.trace {
		defer un(trace(p, keyword.String()+"Spec"))
	}

	assert(keyword == token.Zh_设定, "var")

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

	if typ == nil && values == nil && iota == 0 {
		p.error(pos, "missing variable type or initialization")
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

	p.declare(spec, iota, p.topScope, ast.Var, idents...)

	return spec
}
