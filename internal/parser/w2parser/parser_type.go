// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseTypeSpec(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
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
	spec.Type = p.parseType()
	p.expectSemi() // call before accessing p.linecomment
	spec.Comment = p.lineComment

	return spec
}

func (p *parser) parseTypeList() (list []ast.Expr) {
	if p.trace {
		defer un(trace(p, "TypeList"))
	}

	list = append(list, p.parseType())
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.parseType())
	}

	return
}

func (p *parser) parseInterfaceType(keyword token.Token) *ast.InterfaceType {
	if p.trace {
		defer un(trace(p, "InterfaceType"))
	}

	pos := p.expect(keyword)
	lbrace := p.expect(token.LBRACE)
	scope := ast.NewScope(nil) // interface scope
	var list []*ast.Field
	for p.tok == token.FUNC || p.tok == token.IDENT {
		list = append(list, p.parseMethodSpec(scope))
	}
	rbrace := p.expect(token.RBRACE)

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

func (p *parser) parseMapType() *ast.MapType {
	if p.trace {
		defer un(trace(p, "MapType"))
	}

	pos := p.expect(token.MAP)
	p.expect(token.LBRACK)
	key := p.parseType()
	p.expect(token.RBRACK)
	value := p.parseType()

	return &ast.MapType{Map: pos, Key: key, Value: value}
}

func (p *parser) parseTypeInstance(typ ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "TypeInstance"))
	}

	opening := p.expect(token.LBRACK)
	typeParam := p.parseType()
	closing := p.expectClosing(token.RBRACK, "type argument")

	return &ast.IndexExpr{
		X:      typ,
		Lbrack: opening,
		Index:  typeParam,
		Rbrack: closing,
	}
}

// If the result is an identifier, it is not resolved.
func (p *parser) tryIdentOrType() ast.Expr {
	switch p.tok {
	case token.IDENT:
		typ := p.parseTypeName()
		if p.tok == token.LBRACK {
			typ = p.parseTypeInstance(typ)
		}
		return typ
	case token.LBRACK:
		return p.parseArrayType()
	case token.STRUCT:
		return p.parseStructType(p.tok)
	case token.MUL:
		return p.parsePointerType()
	case token.FUNC:
		typ, _ := p.parseFuncType(p.tok)
		return typ
	case token.INTERFACE:
		return p.parseInterfaceType(p.tok)
	case token.MAP:
		return p.parseMapType()
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
	}

	// no type found
	return nil
}

func (p *parser) tryType() ast.Expr {
	typ := p.tryIdentOrType()
	if typ != nil {
		p.resolve(typ)
	}
	return typ
}

func (p *parser) parseFieldDecl(scope *ast.Scope) *ast.Field {
	if p.trace {
		defer un(trace(p, "FieldDecl"))
	}

	doc := p.leadComment

	// 1st FieldDecl
	// A type name used as an anonymous field looks like a field identifier.
	var list []ast.Expr
	for {
		list = append(list, p.parseVarType(false))
		if p.tok != token.COMMA {
			break
		}
		p.next()
	}

	var colonPos token.Pos
	if p.tok == token.COLON {
		colonPos = p.pos
		p.next()
	}

	typ := p.tryVarType(false)

	// analyze case
	var idents []*ast.Ident
	if typ != nil {
		// IdentifierList Type
		idents = p.makeIdentList(list)
	} else {
		// ["*"] TypeName (AnonymousField)
		typ = list[0] // we always have at least one element
		if n := len(list); n > 1 {
			p.errorExpected(p.pos, "type")
			typ = &ast.BadExpr{From: p.pos, To: p.pos}
		} else if !isTypeName(deref(typ)) {
			p.errorExpected(typ.Pos(), "anonymous field")
			typ = &ast.BadExpr{From: typ.Pos(), To: p.safePos(typ.End())}
		}
	}

	// Tag
	var tag *ast.BasicLit
	if p.tok == token.STRING {
		tag = &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
	}

	p.expectSemi() // call before accessing p.linecomment

	field := &ast.Field{
		Doc:      doc,
		Names:    idents,
		ColonPos: colonPos,
		Type:     typ,
		Tag:      tag,
		Comment:  p.lineComment,
	}
	p.declare(field, nil, scope, ast.Var, idents...)
	p.resolve(typ)

	return field
}

func (p *parser) parseStructType(keyword token.Token) *ast.StructType {
	if p.trace {
		defer un(trace(p, "StructType"))
	}

	pos := p.expect(keyword)
	lbrace := token.NoPos
	if keyword == token.STRUCT {
		lbrace = p.expect(token.LBRACE)
	}
	scope := ast.NewScope(nil) // struct scope
	var list []*ast.Field
	for p.tok == token.IDENT || p.tok == token.MUL || p.tok == token.LPAREN {
		// a field declaration cannot start with a '(' but we accept
		// it here for more robust parsing and better error messages
		// (parseFieldDecl will check and complain if necessary)
		list = append(list, p.parseFieldDecl(scope))
	}
	rbrace := p.expect(token.RBRACE)

	return &ast.StructType{
		TokPos: pos,
		Tok:    keyword,
		Fields: &ast.FieldList{
			Opening: lbrace,
			List:    list,
			Closing: rbrace,
		},
	}
}

func (p *parser) parsePointerType() *ast.StarExpr {
	if p.trace {
		defer un(trace(p, "PointerType"))
	}

	star := p.expect(token.MUL)
	base := p.parseType()

	return &ast.StarExpr{Star: star, X: base}
}

// If the result is an identifier, it is not resolved.
func (p *parser) tryVarType(isParam bool) ast.Expr {
	if isParam && p.tok == token.ELLIPSIS {
		pos := p.pos
		p.next()
		typ := p.tryIdentOrType() // don't use parseType so we can provide better error message
		if typ != nil {
			p.resolve(typ)
		} else {
			p.error(pos, "'...' parameter is missing type")
			typ = &ast.BadExpr{From: pos, To: p.pos}
		}
		return &ast.Ellipsis{Ellipsis: pos, Elt: typ}
	}
	return p.tryIdentOrType()
}

// If the result is an identifier, it is not resolved.
func (p *parser) parseVarType(isParam bool) ast.Expr {
	typ := p.tryVarType(isParam)
	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.next() // make progress
		typ = &ast.BadExpr{From: pos, To: p.pos}
	}
	return typ
}

func (p *parser) parseType() ast.Expr {
	if p.trace {
		defer un(trace(p, "Type"))
	}

	typ := p.tryType()

	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.advance(exprEnd)
		return &ast.BadExpr{From: pos, To: p.pos}
	}

	return typ
}

// If the result is an identifier, it is not resolved.
func (p *parser) parseTypeName() ast.Expr {
	if p.trace {
		defer un(trace(p, "TypeName"))
	}

	ident := p.parseIdent()
	// don't resolve ident yet - it may be a parameter or field name

	if p.tok == token.PERIOD {
		// ident is a package name
		p.next()
		p.resolve(ident)
		sel := p.parseIdent()
		return &ast.SelectorExpr{X: ident, Sel: sel}
	}

	return ident
}

func (p *parser) parseArrayType() ast.Expr {
	if p.trace {
		defer un(trace(p, "ArrayType"))
	}

	lbrack := p.expect(token.LBRACK)
	p.exprLev++
	var len ast.Expr
	// always permit ellipsis for more fault-tolerant parsing
	if p.tok == token.ELLIPSIS {
		len = &ast.Ellipsis{Ellipsis: p.pos}
		p.next()
	} else if p.tok != token.RBRACK {
		len = p.parseRhs()
	}
	p.exprLev--
	p.expect(token.RBRACK)
	elt := p.parseType()

	return &ast.ArrayType{Lbrack: lbrack, Len: len, Elt: elt}
}

func (p *parser) makeIdentList(list []ast.Expr) []*ast.Ident {
	idents := make([]*ast.Ident, len(list))
	for i, x := range list {
		ident, isIdent := x.(*ast.Ident)
		if !isIdent {
			if _, isBad := x.(*ast.BadExpr); !isBad {
				// only report error if it's a new one
				p.errorExpected(x.Pos(), "identifier")
			}
			ident = &ast.Ident{NamePos: x.Pos(), Name: "_"}
		}
		idents[i] = ident
	}
	return idents
}
