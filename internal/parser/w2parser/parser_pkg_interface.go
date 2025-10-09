// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseInterfaceDecl(keyword token.Token) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	ident := p.parseIdent()
	spec := &ast.TypeSpec{Doc: doc, Name: ident}
	p.declare(spec, nil, p.topScope, ast.Typ, ident)

	lbrace := p.expect(token.COLON)

	scope := ast.NewScope(nil) // interface scope
	var list []*ast.Field
	for p.tok == token.IDENT {
		list = append(list, p.parseMethodSpec(scope))
	}
	rbrace := p.expect(token.Zh_完毕)
	p.expectSemi() // call before accessing p.linecomment
	spec.Comment = p.lineComment

	spec.Type = &ast.InterfaceType{
		TokPos: pos,
		Tok:    keyword,
		Methods: &ast.FieldList{
			Opening: lbrace,
			List:    list,
			Closing: rbrace,
		},
	}

	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    keyword,
		Specs:  []ast.Spec{spec},
	}
}

func (p *parser) parseInterfaceType(keyword token.Token) *ast.InterfaceType {
	if p.trace {
		defer un(trace(p, "InterfaceType"))
	}

	pos := p.expect(keyword)
	lbrace := p.expect(token.COLON)
	scope := ast.NewScope(nil) // interface scope
	var list []*ast.Field
	for p.tok == token.Zh_函数 || p.tok == token.IDENT {
		list = append(list, p.parseMethodSpec(scope))
	}
	rbrace := p.expect(token.Zh_完毕)

	return &ast.InterfaceType{
		TokPos: pos,
		Tok:    keyword,
		Methods: &ast.FieldList{
			Opening: lbrace,
			List:    list,
			Closing: rbrace,
		},
	}
}
