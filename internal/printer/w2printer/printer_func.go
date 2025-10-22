// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) funcDecl_pkg(d *ast.FuncDecl) {
	p.setComment(d.Doc)
	p.print(d.Pos(), token.Zh_函数, token.K_点)

	if d.Recv != nil {
		var thisTypeIdent *ast.Ident
		if s := d.Recv.List[0].Names[0].Name; s == token.K_this || s == token.K_我的 {
			if typ, ok := d.Recv.List[0].Type.(*ast.StarExpr); ok {
				if ident, ok := typ.X.(*ast.Ident); ok {
					thisTypeIdent = ident
				}
			}
		}
		if thisTypeIdent != nil {
			p.print(thisTypeIdent.Name)
			p.print(token.K_点)
		} else {
			p.parameters(d.Recv, funcParam) // method: print receiver
			p.print(blank)
		}
	}

	p.print(d.Name.Name)

	p.signature(d.Type)
	p.funcBody(p.distanceFrom(d.Pos()), vtab, d.Body)
}

type paramMode int

const (
	funcParam paramMode = iota
	funcTParam
	typeTParam
)

func (p *printer) parameters(fields *ast.FieldList, mode paramMode) {
	openTok, closeTok := token.LPAREN, token.RPAREN
	if mode != funcParam {
		openTok, closeTok = token.LBRACK, token.RBRACK
	}

	p.print(fields.Opening, openTok)
	if len(fields.List) > 0 {
		prevLine := p.lineFor(fields.Opening)
		ws := indent
		for i, par := range fields.List {
			// determine par begin and end line (may be different
			// if there are multiple parameter names for this par
			// or the type is on a separate line)
			var parLineBeg int
			if len(par.Names) > 0 {
				parLineBeg = p.lineFor(par.Names[0].Pos())
			} else {
				parLineBeg = p.lineFor(par.Type.Pos())
			}
			var parLineEnd = p.lineFor(par.Type.End())
			// separating "," if needed
			needsLinebreak := 0 < prevLine && prevLine < parLineBeg
			if i > 0 {
				// use position of parameter following the comma as
				// comma position for correct comma placement, but
				// only if the next parameter is on the same line
				if !needsLinebreak {
					p.print(par.Pos())
				}
				p.print(token.COMMA)
			}
			// separator if needed (linebreak or blank)
			if needsLinebreak && p.linebreak(parLineBeg, 0, ws, true) > 0 {
				// break line if the opening "(" or previous parameter ended on a different line
				ws = ignore
			} else if i > 0 {
				p.print(blank)
			}
			// parameter names
			if len(par.Names) > 0 {
				// Very subtle: If we indented before (ws == ignore), identList
				// won't indent again. If we didn't (ws == indent), identList will
				// indent if the identList spans multiple lines, and it will outdent
				// again at the end (and still ws == indent). Thus, a subsequent indent
				// by a linebreak call after a type, or in the next multi-line identList
				// will do the right thing.
				p.identList(par.Names, ws == indent)
				p.print(token.COLON)
				p.print(blank)
			}
			// parameter type
			p.expr(stripParensAlways(par.Type))
			prevLine = parLineEnd
		}
		// if the closing ")" is on a separate line from the last parameter,
		// print an additional "," and line break
		if closing := p.lineFor(fields.Closing); 0 < prevLine && prevLine < closing {
			p.print(token.COMMA)
			p.linebreak(closing, 0, ignore, true)
		} else if mode == typeTParam && fields.NumFields() == 1 && combinesWithName(fields.List[0].Type) {
			// A type parameter list [P T] where the name P and the type expression T syntactically
			// combine to another valid (value) expression requires a trailing comma, as in [P *T,]
			// (or an enclosing interface as in [P interface(*T)]), so that the type parameter list
			// is not parsed as an array length [P*T].
			p.print(token.COMMA)
		}
		// unindent if we indented
		if ws == ignore {
			p.print(unindent)
		}
	}
	p.print(fields.Closing, closeTok)
}

func (p *printer) signature(sig *ast.FuncType) {
	if sig.Params != nil {
		if sig.Params.Opening != token.NoPos {
			p.parameters(sig.Params, funcParam)
		}
	} else {
		p.print(token.LPAREN, token.RPAREN)
	}
	n := sig.Results.NumFields()
	if n > 0 {
		// result != nil
		p.print(blank)
		p.print(token.ARROW)
		p.print(blank)

		if n == 1 && sig.Results.List[0].Names == nil {
			// single anonymous result; no ()'s
			p.expr(stripParensAlways(sig.Results.List[0].Type))
			return
		}
		p.parameters(sig.Results, funcParam)
	}
}

// funcBody prints a function body following a function header of given headerSize.
// If the header's and block's size are "small enough" and the block is "simple enough",
// the block is printed on the current line, without line breaks, spaced from the header
// by sep. Otherwise the block's opening "{" is printed on the current line, followed by
// lines for the block's statements and its closing "}".
func (p *printer) funcBody(headerSize int, sep whiteSpace, b *ast.BlockStmt) {
	if b == nil {
		return
	}

	// save/restore composite literal nesting level
	defer func(level int) {
		p.level = level
	}(p.level)
	p.level = 0

	const maxSize = 100
	if headerSize+p.funcBodySize(b, maxSize) <= maxSize {
		p.print(sep, b.Lbrace, token.COLON)
		if len(b.List) > 0 {
			p.print(blank)
			for i, s := range b.List {
				if i > 0 {
					p.print(token.SEMICOLON, blank)
				}
				p.stmt(s, i == len(b.List)-1)
			}
			p.print(blank)
		}
		p.print(noExtraLinebreak, b.Rbrace, token.Zh_完毕, noExtraLinebreak)
		return
	}

	if sep != ignore {
		// 忽略
	}
	p.block(b, 1, token.Zh_完毕)
}

// funcBodySize is like nodeSize but it is specialized for *ast.BlockStmt's.
func (p *printer) funcBodySize(b *ast.BlockStmt, maxSize int) int {
	pos1 := b.Pos()
	pos2 := b.Rbrace
	if pos1.IsValid() && pos2.IsValid() && p.lineFor(pos1) != p.lineFor(pos2) {
		// opening and closing brace are on different lines - don't make it a one-liner
		return maxSize + 1
	}
	if len(b.List) > 5 {
		// too many statements - don't make it a one-liner
		return maxSize + 1
	}
	// otherwise, estimate body size
	bodySize := p.commentSizeBefore(p.posFor(pos2))
	for i, s := range b.List {
		if bodySize > maxSize {
			break // no need to continue
		}
		if i > 0 {
			bodySize += 2 // space for a semicolon and blank
		}
		bodySize += p.nodeSize(s, maxSize)
	}
	return bodySize
}

func stripParensAlways(x ast.Expr) ast.Expr {
	if x, ok := x.(*ast.ParenExpr); ok {
		return stripParensAlways(x.X)
	}
	return x
}

// combinesWithName reports whether a name followed by the expression x
// syntactically combines to another valid (value) expression. For instance
// using *T for x, "name *T" syntactically appears as the expression x*T.
// On the other hand, using  P|Q or *P|~Q for x, "name P|Q" or name *P|~Q"
// cannot be combined into a valid (value) expression.
func combinesWithName(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.StarExpr:
		// name *x.X combines to name*x.X if x.X is not a type element
		return !isTypeElem(x.X)
	case *ast.BinaryExpr:
		return combinesWithName(x.X) && !isTypeElem(x.Y)
	case *ast.ParenExpr:
		// name(x) combines but we are making sure at
		// the call site that x is never parenthesized.
		panic("unexpected parenthesized expression")
	}
	return false
}

// isTypeElem reports whether x is a (possibly parenthesized) type element expression.
// The result is false if x could be a type element OR an ordinary (value) expression.
func isTypeElem(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType:
		return true
	case *ast.BinaryExpr:
		return isTypeElem(x.X) || isTypeElem(x.Y)
	case *ast.ParenExpr:
		return isTypeElem(x.X)
	}
	return false
}
