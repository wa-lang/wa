// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseStructDecl(keyword token.Token) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	ident := p.parseIdent()
	spec := &ast.TypeSpec{Doc: doc, Name: ident}
	p.declare(spec, nil, p.topScope, ast.Typ, ident)

	lbrace := p.expect(token.COLON)

	scope := ast.NewScope(nil) // struct scope
	var list []*ast.Field
	for p.tok == token.IDENT || p.tok == token.MUL || p.tok == token.LPAREN {
		// a field declaration cannot start with a '(' but we accept
		// it here for more robust parsing and better error messages
		// (parseFieldDecl will check and complain if necessary)
		list = append(list, p.parseFieldDecl(scope))
	}
	rbrace := p.expect(token.Zh_完毕)
	p.expectSemi() // call before accessing p.linecomment
	spec.Comment = p.lineComment

	spec.Type = &ast.StructType{
		TokPos: pos,
		Tok:    keyword,
		Fields: &ast.FieldList{
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
