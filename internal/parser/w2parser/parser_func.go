// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseFuncDecl(keyword token.Token) *ast.FuncDecl {
	if p.trace {
		defer un(trace(p, "FunctionDecl"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)
	scope := ast.NewScope(p.topScope) // function scope

	var recv *ast.FieldList
	if p.tok == token.LPAREN {
		recv = p.parseParameters(scope, false)
	}

	ident := p.parseIdent()

	// func Type.method()
	if p.tok == token.PERIOD {
		thisIdent := &ast.Ident{Name: token.K_我的}
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

	params, results, arrowPos := p.parseSignature(scope)

	var body *ast.BlockStmt
	if p.tok == token.COLON {
		body = p.parseFuncBody(scope)
	}
	p.expectSemi()

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
		// Go spec: The scope of an identifier denoting a constant, type,
		// variable, or function (but not method) declared at top level
		// (outside any function) is the package block.
		//
		// init() functions cannot be referred to and there may
		// be more than one - don't put them in the pkgScope
		if ident.Name != token.K_准备 {
			p.declare(decl, nil, p.pkgScope, ast.Fun, ident)
		}
	}

	return decl
}

func (p *parser) parseFuncBody(scope *ast.Scope) *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, "Body"))
	}

	lbrace := p.expect(token.COLON)
	p.topScope = scope // open function scope
	p.openLabelScope()
	list := p.parseStmtList()
	p.closeLabelScope()
	p.closeScope()
	rbrace := p.expect(token.Zh_完毕)

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

func (p *parser) parseParameterList(scope *ast.Scope, ellipsisOk bool) (params []*ast.Field) {
	if p.trace {
		defer un(trace(p, "ParameterList"))
	}

	// 1st ParameterDecl
	// A list of identifiers looks like a list of type names.
	var list []ast.Expr
	for {
		list = append(list, p.parseVarType(ellipsisOk))
		// 开头的 Ident 列表, 遇到非 ',' 时结束
		if p.tok != token.COMMA {
			break
		}
		p.next()
		if p.tok == token.RPAREN {
			break
		}
	}

	// Wa 必须加 ':'
	var colonPos token.Pos
	if p.tok == token.COLON {
		colonPos = p.pos
		p.next()
	} else {
		if p.tok != token.RPAREN {
			p.expect(token.COLON)
		}
	}

	// analyze case
	if typ := p.tryVarType(ellipsisOk); typ != nil {
		// IdentifierList Type
		idents := p.makeIdentList(list)
		field := &ast.Field{Names: idents, ColonPos: colonPos, Type: typ}
		params = append(params, field)
		// Go spec: The scope of an identifier denoting a function
		// parameter or result variable is the function body.
		p.declare(field, nil, scope, ast.Var, idents...)
		p.resolve(typ)
		if !p.atComma("parameter list", token.RPAREN) {
			return
		}
		p.next()
		for p.tok != token.RPAREN && p.tok != token.EOF {
			idents := p.parseIdentList()
			if p.tok == token.COLON {
				colonPos = p.pos
				p.next()
			} else {
				p.expect(token.COLON)
			}

			typ := p.parseVarType(ellipsisOk)
			field := &ast.Field{Names: idents, ColonPos: colonPos, Type: typ}
			params = append(params, field)
			// Go spec: The scope of an identifier denoting a function
			// parameter or result variable is the function body.
			p.declare(field, nil, scope, ast.Var, idents...)
			p.resolve(typ)
			if !p.atComma("parameter list", token.RPAREN) {
				break
			}
			p.next()
		}
		return
	}

	// 缺少类型信息
	if colonPos != token.NoPos {
		p.errorExpected(p.pos, "type")
		return
	}

	// Type { "," Type } (anonymous parameters)
	params = make([]*ast.Field, len(list))
	for i, typ := range list {
		p.resolve(typ)
		params[i] = &ast.Field{Type: typ}
	}
	return
}

func (p *parser) parseParameters(scope *ast.Scope, ellipsisOk bool) *ast.FieldList {
	if p.trace {
		defer un(trace(p, "Parameters"))
	}

	var params []*ast.Field
	lparen := p.expect(token.LPAREN)
	if p.tok != token.RPAREN {
		params = p.parseParameterList(scope, ellipsisOk)
	}
	rparen := p.expect(token.RPAREN)

	return &ast.FieldList{Opening: lparen, List: params, Closing: rparen}
}

func (p *parser) parseResult(scope *ast.Scope) *ast.FieldList {
	if p.trace {
		defer un(trace(p, "Result"))
	}

	if p.tok == token.LPAREN {
		return p.parseParameters(scope, false)
	}

	typ := p.tryType()
	if typ != nil {
		list := make([]*ast.Field, 1)
		list[0] = &ast.Field{Type: typ}
		return &ast.FieldList{List: list}
	}

	return nil
}

func (p *parser) parseSignature(scope *ast.Scope) (params, results *ast.FieldList, arrowPos token.Pos) {
	if p.trace {
		defer un(trace(p, "Signature"))
	}

	// 无参数和返回值时可省略 `()`
	// func main
	// func main { ... }
	// 定义 主控
	// 定义 主控: 完毕
	if p.tok == token.COLON || p.tok == token.SEMICOLON {
		params = new(ast.FieldList)
		results = new(ast.FieldList)
		return
	}

	// func answer => i32 {}
	if p.tok == token.ARROW {
		params = new(ast.FieldList)
		p.next()
	} else {
		params = p.parseParameters(scope, true)

		if p.tok == token.ARROW {
			arrowPos = p.pos
			p.next()
		}
	}

	results = p.parseResult(scope)

	return
}

func (p *parser) parseFuncType(keyword token.Token) (*ast.FuncType, *ast.Scope) {
	if p.trace {
		defer un(trace(p, "FuncType"))
	}

	pos := p.expect(keyword)
	scope := ast.NewScope(p.topScope) // function scope
	params, results, arrowPos := p.parseSignature(scope)

	return &ast.FuncType{
		TokPos:   pos,
		Tok:      keyword,
		Params:   params,
		ArrowPos: arrowPos,
		Results:  results,
	}, scope
}

func (p *parser) parseMethodSpec(scope *ast.Scope) *ast.Field {
	if p.trace {
		defer un(trace(p, "MethodSpec"))
	}

	var keyword token.Token
	if p.tok == token.Zh_函数 {
		keyword = p.tok
		p.next()
	}

	doc := p.leadComment
	var idents []*ast.Ident
	var typ ast.Expr
	x := p.parseTypeName()
	if ident, isIdent := x.(*ast.Ident); isIdent && p.tok == token.LPAREN {
		// method
		idents = []*ast.Ident{ident}
		scope := ast.NewScope(nil) // method scope
		params, results, arrowPos := p.parseSignature(scope)
		typ = &ast.FuncType{
			TokPos:   token.NoPos,
			Tok:      keyword,
			Params:   params,
			ArrowPos: arrowPos,
			Results:  results,
		}
	} else {
		// embedded interface
		typ = x
		p.resolve(typ)
	}
	p.expectSemi() // call before accessing p.linecomment

	spec := &ast.Field{Doc: doc, Names: idents, Type: typ, Comment: p.lineComment}
	p.declare(spec, nil, scope, ast.Fun, idents...)

	return spec
}
