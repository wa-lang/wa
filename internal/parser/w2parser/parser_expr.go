// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"fmt"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) makeExpr(s ast.Stmt, want string) ast.Expr {
	if s == nil {
		return nil
	}
	if es, isExpr := s.(*ast.ExprStmt); isExpr {
		return p.checkExpr(es.X)
	}
	found := "simple statement"
	if _, isAss := s.(*ast.AssignStmt); isAss {
		found = "assignment"
	}
	p.error(s.Pos(), fmt.Sprintf("expected %s, found %s (missing parentheses around composite literal?)", want, found))
	return &ast.BadExpr{From: s.Pos(), To: p.safePos(s.End())}
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *parser) parseBinaryExpr(lhs bool, prec1 int) ast.Expr {
	if p.trace {
		defer un(trace(p, "BinaryExpr"))
	}

	x := p.parseUnaryExpr(lhs)
	for {
		op, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		pos := p.expect(op)
		if lhs {
			p.resolve(x)
			lhs = false
		}
		y := p.parseBinaryExpr(false, oprec+1)
		x = &ast.BinaryExpr{X: p.checkExpr(x), OpPos: pos, Op: op, Y: p.checkExpr(y)}
	}
}

// If lhs is set and the result is an identifier, it is not resolved.
// The result may be a type or even a raw type ([...]int). Callers must
// check the result (using checkExpr or checkExprOrType), depending on
// context.
func (p *parser) parseExpr(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Expression"))
	}

	return p.parseBinaryExpr(lhs, token.LowestPrec+1)
}

func (p *parser) parseRhs() ast.Expr {
	old := p.inRhs
	p.inRhs = true
	x := p.checkExpr(p.parseExpr(false))
	p.inRhs = old
	return x
}

func (p *parser) parseRhsOrType() ast.Expr {
	old := p.inRhs
	p.inRhs = true
	x := p.checkExprOrType(p.parseExpr(false))
	p.inRhs = old
	return x
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *parser) parseUnaryExpr(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "UnaryExpr"))
	}

	switch p.tok {
	case token.ADD, token.SUB, token.NOT, token.XOR, token.AND:
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr(false)
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: p.checkExpr(x)}

	case token.MUL:
		// pointer type or unary "*" expression
		pos := p.pos
		p.next()
		x := p.parseUnaryExpr(false)
		return &ast.StarExpr{Star: pos, X: p.checkExprOrType(x)}
	}

	return p.parsePrimaryExpr(lhs)
}

func (p *parser) tokPrec() (token.Token, int) {
	tok := p.tok
	if p.inRhs && tok == token.ASSIGN {
		tok = token.EQL
	}
	return tok, tok.Precedence()
}

// If x is of the form *T, deref returns T, otherwise it returns x.
func deref(x ast.Expr) ast.Expr {
	if p, isPtr := x.(*ast.StarExpr); isPtr {
		x = p.X
	}
	return x
}

// If x is of the form (T), unparen returns unparen(T), otherwise it returns x.
func unparen(x ast.Expr) ast.Expr {
	if p, isParen := x.(*ast.ParenExpr); isParen {
		x = unparen(p.X)
	}
	return x
}

// checkExprOrType checks that x is an expression or a type
// (and not a raw type such as [...]T).
func (p *parser) checkExprOrType(x ast.Expr) ast.Expr {
	switch t := unparen(x).(type) {
	case *ast.ParenExpr:
		panic("unreachable")
	case *ast.UnaryExpr:
	case *ast.ArrayType:
		if len, isEllipsis := t.Len.(*ast.Ellipsis); isEllipsis {
			p.error(len.Pos(), "expected array length, found '...'")
			x = &ast.BadExpr{From: x.Pos(), To: p.safePos(x.End())}
		}
	}

	// all other nodes are expressions or types
	return x
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *parser) parsePrimaryExpr(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "PrimaryExpr"))
	}

	x := p.parseOperand(lhs)
L:
	for {
		switch p.tok {
		case token.PERIOD:
			p.next()
			if lhs {
				p.resolve(x)
			}
			switch p.tok {
			case token.IDENT:
				x = p.parseSelector(p.checkExprOrType(x))
			case token.LPAREN:
				x = p.parseTypeAssertion(p.checkExpr(x))
			default:
				pos := p.pos
				p.errorExpected(pos, "selector or type assertion")
				p.next() // make progress
				sel := &ast.Ident{NamePos: pos, Name: "_"}
				x = &ast.SelectorExpr{X: x, Sel: sel}
			}
		case token.LBRACK:
			if lhs {
				p.resolve(x)
			}
			x = p.parseIndexOrSlice(p.checkExpr(x))
		case token.LPAREN:
			if lhs {
				p.resolve(x)
			}
			x = p.parseCallOrConversion(p.checkExprOrType(x))
		case token.LBRACE:
			if isLiteralType(x) && (p.exprLev >= 0 || !isTypeName(x)) {
				if lhs {
					p.resolve(x)
				}
				x = p.parseLiteralValue(x)
			} else {
				break L
			}
		default:
			break L
		}
		lhs = false // no need to try to resolve again
	}

	return x
}

// checkExpr checks that x is an expression (and not a type).
func (p *parser) checkExpr(x ast.Expr) ast.Expr {
	switch unparen(x).(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.BasicLit:
	case *ast.FuncLit:
	case *ast.CompositeLit:
	case *ast.ParenExpr:
		panic("unreachable")
	case *ast.SelectorExpr:
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
		// If t.Type == nil we have a type assertion of the form
		// y.(type), which is only allowed in type switch expressions.
		// It's hard to exclude those but for the case where we are in
		// a type switch. Instead be lenient and test this in the type
		// checker.
	case *ast.CallExpr:
	case *ast.StarExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
	default:
		// all other nodes are not proper expressions
		p.errorExpected(x.Pos(), "expression")
		x = &ast.BadExpr{From: x.Pos(), To: p.safePos(x.End())}
	}
	return x
}

// isTypeName reports whether x is a (qualified) TypeName.
func isTypeName(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.SelectorExpr:
		_, isIdent := t.X.(*ast.Ident)
		return isIdent
	default:
		return false // all other nodes are not type names
	}
	return true
}

// isLiteralType reports whether x is a legal composite literal type.
func isLiteralType(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.SelectorExpr:
		_, isIdent := t.X.(*ast.Ident)
		return isIdent
	case *ast.ArrayType:
	case *ast.StructType:
	case *ast.MapType:
	default:
		return false // all other nodes are not legal composite literal types
	}
	return true
}

func (p *parser) parseElement() ast.Expr {
	if p.trace {
		defer un(trace(p, "Element"))
	}

	x := p.parseValue(true)
	if p.tok == token.COLON {
		colon := p.pos
		p.next()
		x = &ast.KeyValueExpr{Key: x, Colon: colon, Value: p.parseValue(false)}
	}

	return x
}

func (p *parser) parseElementList() (list []ast.Expr) {
	if p.trace {
		defer un(trace(p, "ElementList"))
	}

	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseElement())
		if !p.atComma("composite literal", token.RBRACE) {
			break
		}
		p.next()
	}

	return
}

func (p *parser) parseLiteralValue(typ ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "LiteralValue"))
	}

	lbrace := p.expect(token.LBRACE)
	var elts []ast.Expr
	p.exprLev++
	if p.tok != token.RBRACE {
		elts = p.parseElementList()
	}
	p.exprLev--
	rbrace := p.expectClosing(token.RBRACE, "composite literal")
	return &ast.CompositeLit{Type: typ, Lbrace: lbrace, Elts: elts, Rbrace: rbrace}
}

func (p *parser) parseTypeAssertion(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "TypeAssertion"))
	}

	lparen := p.expect(token.LPAREN)
	var typ ast.Expr
	if p.tok == token.TYPE {
		// type switch: typ == nil
		p.next()
	} else {
		typ = p.parseType()
	}
	rparen := p.expect(token.RPAREN)

	return &ast.TypeAssertExpr{X: x, Type: typ, Lparen: lparen, Rparen: rparen}
}

func (p *parser) parseIndexOrSlice(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "IndexOrSlice"))
	}

	const N = 3 // change the 3 to 2 to disable 3-index slices
	lbrack := p.expect(token.LBRACK)
	p.exprLev++
	var index [N]ast.Expr
	var colons [N - 1]token.Pos
	if p.tok != token.COLON {
		index[0] = p.parseRhs()
	}
	ncolons := 0
	for p.tok == token.COLON && ncolons < len(colons) {
		colons[ncolons] = p.pos
		ncolons++
		p.next()
		if p.tok != token.COLON && p.tok != token.RBRACK && p.tok != token.EOF {
			index[ncolons] = p.parseRhs()
		}
	}
	p.exprLev--
	rbrack := p.expect(token.RBRACK)

	if ncolons > 0 {
		// slice expression
		slice3 := false
		if ncolons == 2 {
			slice3 = true
			// Check presence of 2nd and 3rd index here rather than during type-checking
			// to prevent erroneous programs from passing through gofmt (was issue 7305).
			if index[1] == nil {
				p.error(colons[0], "2nd index required in 3-index slice")
				index[1] = &ast.BadExpr{From: colons[0] + 1, To: colons[1]}
			}
			if index[2] == nil {
				p.error(colons[1], "3rd index required in 3-index slice")
				index[2] = &ast.BadExpr{From: colons[1] + 1, To: rbrack}
			}
		}
		return &ast.SliceExpr{X: x, Lbrack: lbrack, Low: index[0], High: index[1], Max: index[2], Slice3: slice3, Rbrack: rbrack}
	}

	return &ast.IndexExpr{X: x, Lbrack: lbrack, Index: index[0], Rbrack: rbrack}
}

func (p *parser) parseCallOrConversion(fun ast.Expr) *ast.CallExpr {
	if p.trace {
		defer un(trace(p, "CallOrConversion"))
	}

	lparen := p.expect(token.LPAREN)
	p.exprLev++
	var list []ast.Expr
	var ellipsis token.Pos
	for p.tok != token.RPAREN && p.tok != token.EOF && !ellipsis.IsValid() {
		list = append(list, p.parseRhsOrType()) // builtins may expect a type: make(some type, ...)
		if p.tok == token.ELLIPSIS {
			ellipsis = p.pos
			p.next()
		}
		if !p.atComma("argument list", token.RPAREN) {
			break
		}
		p.next()
	}
	p.exprLev--
	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.CallExpr{Fun: fun, Lparen: lparen, Args: list, Ellipsis: ellipsis, Rparen: rparen}
}

func (p *parser) parseValue(keyOk bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Element"))
	}

	if p.tok == token.LBRACE {
		return p.parseLiteralValue(nil)
	}

	// Because the parser doesn't know the composite literal type, it cannot
	// know if a key that's an identifier is a struct field name or a name
	// denoting a value. The former is not resolved by the parser or the
	// resolver.
	//
	// Instead, _try_ to resolve such a key if possible. If it resolves,
	// it a) has correctly resolved, or b) incorrectly resolved because
	// the key is a struct field with a name matching another identifier.
	// In the former case we are done, and in the latter case we don't
	// care because the type checker will do a separate field lookup.
	//
	// If the key does not resolve, it a) must be defined at the top
	// level in another file of the same package, the universe scope, or be
	// undeclared; or b) it is a struct field. In the former case, the type
	// checker can do a top-level lookup, and in the latter case it will do
	// a separate field lookup.
	x := p.checkExpr(p.parseExpr(keyOk))
	if keyOk {
		if p.tok == token.COLON {
			// Try to resolve the key but don't collect it
			// as unresolved identifier if it fails so that
			// we don't get (possibly false) errors about
			// undeclared names.
			p.tryResolve(x, false)
		} else {
			// not a key
			p.resolve(x)
		}
	}

	return x
}

// parseOperand may return an expression or a raw type (incl. array
// types of the form [...]T. Callers must verify the result.
// If lhs is set and the result is an identifier, it is not resolved.
func (p *parser) parseOperand(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if !lhs {
			p.resolve(x)
		}
		return x

	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
		x := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x

	case token.LPAREN:
		lparen := p.pos
		p.next()
		p.exprLev++
		x := p.parseRhsOrType() // types may be parenthesized: (some type)
		p.exprLev--
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}

	case token.FUNC:
		return p.parseFuncTypeOrLit(p.tok)
	}

	if typ := p.tryIdentOrType(); typ != nil {
		// could be type for composite literal or conversion
		_, isIdent := typ.(*ast.Ident)
		assert(!isIdent, "type cannot be identifier")
		return typ
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	p.advance(stmtStart)
	return &ast.BadExpr{From: pos, To: p.pos}
}

func (p *parser) parseSelector(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Selector"))
	}

	sel := p.parseIdent()

	return &ast.SelectorExpr{X: x, Sel: sel}
}

func (p *parser) parseFuncTypeOrLit(keywork token.Token) ast.Expr {
	if p.trace {
		defer un(trace(p, "FuncTypeOrLit"))
	}

	typ, scope := p.parseFuncType(keywork)
	if p.tok != token.LBRACE {
		// function type only
		return typ
	}

	p.exprLev++
	body := p.parseFuncBody(scope)
	p.exprLev--

	return &ast.FuncLit{Type: typ, Body: body}
}

// If lhs is set, result list elements which are identifiers are not resolved.
func (p *parser) parseExprList(lhs bool) (list []ast.Expr) {
	if p.trace {
		defer un(trace(p, "ExpressionList"))
	}

	list = append(list, p.checkExpr(p.parseExpr(lhs)))
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.checkExpr(p.parseExpr(lhs)))
	}

	return
}

func (p *parser) parseLhsList() []ast.Expr {
	old := p.inRhs
	p.inRhs = false
	list := p.parseExprList(true)
	switch p.tok {
	case token.DEFINE:
		// lhs of a short variable declaration
		// but doesn't enter scope until later:
		// caller must call p.shortVarDecl(p.makeIdentList(list))
		// at appropriate time.
	case token.COLON:
		// lhs of a label declaration or a communication clause of a select
		// statement (parseLhsList is not called when parsing the case clause
		// of a switch statement):
		// - labels are declared by the caller of parseLhsList
		// - for communication clauses, if there is a stand-alone identifier
		//   followed by a colon, we have a syntax error; there is no need
		//   to resolve the identifier in that case
	default:
		// identifiers must be declared elsewhere
		for _, x := range list {
			p.resolve(x)
		}
	}
	p.inRhs = old
	return list
}

func (p *parser) parseRhsList() []ast.Expr {
	old := p.inRhs
	p.inRhs = true
	list := p.parseExprList(false)
	p.inRhs = old
	return list
}
