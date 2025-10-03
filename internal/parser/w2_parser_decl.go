// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseDecl_zh(sync map[token.Token]bool) ast.Decl {
	if p.trace {
		defer un(trace(p, "Declaration"))
	}

	var f parseSpecFunction
	switch p.tok {
	case token.Zh_常量, token.Zh_定义, token.Zh_全局:
		f = p.parseValueSpec_zh

	case token.Zh_类型:
		f = p.parseTypeSpec_zh

	case token.Zh_算始, token.Zh_函始:
		return p.parseFuncDecl_zh(p.tok)

	default:
		pos := p.pos
		p.errorExpected_zh(pos, "宣告:"+p.lit+p.tok.String())
		p.advance(sync)
		return &ast.BadDecl{From: pos, To: p.pos}
	}

	return p.parseGenDecl(p.tok, f)
}

func (p *parser) parseGenDecl_zh(keyword token.Token, f parseSpecFunction) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	var lparen, rparen token.Pos
	var list []ast.Spec
	if p.tok == token.LPAREN {
		lparen = p.pos
		p.next()
		for iota := 0; p.tok != token.RPAREN && p.tok != token.EOF; iota++ {
			list = append(list, f(p.leadComment, keyword, iota))
		}
		rparen = p.expect(token.RPAREN)
		p.expectSemi()
	} else {
		list = append(list, f(nil, keyword, 0))
	}

	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    keyword,
		Lparen: lparen,
		Specs:  list,
		Rparen: rparen,
	}
}

func (p *parser) parseFuncDecl_zh(keyword token.Token) *ast.FuncDecl {
	if p.trace {
		defer un(trace(p, "FunctionDecl"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)
	scope := ast.NewScope(p.topScope) // function scope

	ident := p.parseIdent()

	// func Type.method()
	var recv *ast.FieldList
	{
		if p.tok == token.PERIOD {
			thisIdent := &ast.Ident{Name: "自身"}
			thisField := &ast.Field{
				Names: []*ast.Ident{thisIdent},
				Type:  &ast.StarExpr{X: ident},
			}
			recv = &ast.FieldList{
				List: []*ast.Field{thisField},
			}

			p.declare(thisField, nil, scope, ast.Var, thisIdent)

			p.next()
			ident = p.parseIdent()
		}
	}

	params, results, arrowPos := p.parseSignature(scope)

	body := p.parseBody_zh(keyword, scope)
	p.expectSemi_zh()

	decl := &ast.FuncDecl{
		Doc:  doc,
		Recv: recv,
		Name: ident,
		Type: &ast.FuncType{
			TokPos:   pos,
			Tok:      keyword,
			Params:   params,
			ArrowPos: arrowPos,
			Results:  results,
		},
		Body: body,
	}
	if recv == nil {
		// Wa spec: The scope of an identifier denoting a constant, type,
		// variable, or function (but not method) declared at top level
		// (outside any function) is the package block.
		//
		// 准备() functions cannot be referred to and there may
		// be more than one - don't put them in the pkgScope
		if ident.Name != "准备" {
			p.declare(decl, nil, p.pkgScope, ast.Fun, ident)
		}
	}

	return decl
}
