// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"math"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

type exprListMode uint

const (
	commaTerm exprListMode = 1 << iota // list is optionally terminated by a comma
	noIndent                           // no extra indentation in multi-line lists
)

// If indent is set, a multi-line identifier list is indented after the
// first linebreak encountered.
func (p *printer) identList(list []*ast.Ident, indent bool) {
	// convert into an expression list so we can re-use exprList formatting
	xlist := make([]ast.Expr, len(list))
	for i, x := range list {
		xlist[i] = x
	}
	var mode exprListMode
	if !indent {
		mode = noIndent
	}
	p.exprList(token.NoPos, xlist, 1, mode, token.NoPos, false)
}

const filteredMsg = "contains filtered or unexported fields"

// Print a list of expressions. If the list spans multiple
// source lines, the original line breaks are respected between
// expressions.
//
// TODO(gri) Consider rewriting this to be independent of []ast.Expr
//
//	so that we can use the algorithm for any kind of list
//	(e.g., pass list via a channel over which to range).
func (p *printer) exprList(prev0 token.Pos, list []ast.Expr, depth int, mode exprListMode, next0 token.Pos, isIncomplete bool) {
	if len(list) == 0 {
		if isIncomplete {
			prev := p.posFor(prev0)
			next := p.posFor(next0)
			if prev.IsValid() && prev.Line == next.Line {
				p.print("/* " + filteredMsg + " */")
			} else {
				p.print(newline)
				p.print(indent, "// "+filteredMsg, unindent, newline)
			}
		}
		return
	}

	prev := p.posFor(prev0)
	next := p.posFor(next0)
	line := p.lineFor(list[0].Pos())
	endLine := p.lineFor(list[len(list)-1].End())

	if prev.IsValid() && prev.Line == line && line == endLine {
		// all list entries on a single line
		for i, x := range list {
			if i > 0 {
				// use position of expression following the comma as
				// comma position for correct comment placement
				p.print(x.Pos(), token.COMMA, blank)
			}
			p.expr0(x, depth)
		}
		if isIncomplete {
			p.print(token.COMMA, blank, "/* "+filteredMsg+" */")
		}
		return
	}

	// list entries span multiple lines;
	// use source code positions to guide line breaks

	// Don't add extra indentation if noIndent is set;
	// i.e., pretend that the first line is already indented.
	ws := ignore
	if mode&noIndent == 0 {
		ws = indent
	}

	// The first linebreak is always a formfeed since this section must not
	// depend on any previous formatting.
	prevBreak := -1 // index of last expression that was followed by a linebreak
	if prev.IsValid() && prev.Line < line && p.linebreak(line, 0, ws, true) > 0 {
		ws = ignore
		prevBreak = 0
	}

	// initialize expression/key size: a zero value indicates expr/key doesn't fit on a single line
	size := 0

	// We use the ratio between the geometric mean of the previous key sizes and
	// the current size to determine if there should be a break in the alignment.
	// To compute the geometric mean we accumulate the ln(size) values (lnsum)
	// and the number of sizes included (count).
	lnsum := 0.0
	count := 0

	// print all list elements
	prevLine := prev.Line
	for i, x := range list {
		line = p.lineFor(x.Pos())

		// Determine if the next linebreak, if any, needs to use formfeed:
		// in general, use the entire node size to make the decision; for
		// key:value expressions, use the key size.
		// TODO(gri) for a better result, should probably incorporate both
		//           the key and the node size into the decision process
		useFF := true

		// Determine element size: All bets are off if we don't have
		// position information for the previous and next token (likely
		// generated code - simply ignore the size in this case by setting
		// it to 0).
		prevSize := size
		const infinity = 1e6 // larger than any source line
		size = p.nodeSize(x, infinity)
		pair, isPair := x.(*ast.KeyValueExpr)
		if size <= infinity && prev.IsValid() && next.IsValid() {
			// x fits on a single line
			if isPair {
				size = p.nodeSize(pair.Key, infinity) // size <= infinity
			}
		} else {
			// size too large or we don't have good layout information
			size = 0
		}

		// If the previous line and the current line had single-
		// line-expressions and the key sizes are small or the
		// ratio between the current key and the geometric mean
		// if the previous key sizes does not exceed a threshold,
		// align columns and do not use formfeed.
		if prevSize > 0 && size > 0 {
			const smallSize = 40
			if count == 0 || prevSize <= smallSize && size <= smallSize {
				useFF = false
			} else {
				const r = 2.5                               // threshold
				geomean := math.Exp(lnsum / float64(count)) // count > 0
				ratio := float64(size) / geomean
				useFF = r*ratio <= 1 || r <= ratio
			}
		}

		needsLinebreak := 0 < prevLine && prevLine < line
		if i > 0 {
			// Use position of expression following the comma as
			// comma position for correct comment placement, but
			// only if the expression is on the same line.
			if !needsLinebreak {
				p.print(x.Pos())
			}
			p.print(token.COMMA)
			needsBlank := true
			if needsLinebreak {
				// Lines are broken using newlines so comments remain aligned
				// unless useFF is set or there are multiple expressions on
				// the same line in which case formfeed is used.
				nbreaks := p.linebreak(line, 0, ws, useFF || prevBreak+1 < i)
				if nbreaks > 0 {
					ws = ignore
					prevBreak = i
					needsBlank = false // we got a line break instead
				}
				// If there was a new section or more than one new line
				// (which means that the tabwriter will implicitly break
				// the section), reset the geomean variables since we are
				// starting a new group of elements with the next element.
				if nbreaks > 1 {
					lnsum = 0
					count = 0
				}
			}
			if needsBlank {
				p.print(blank)
			}
		}

		if len(list) > 1 && isPair && size > 0 && needsLinebreak {
			// We have a key:value expression that fits onto one line
			// and it's not on the same line as the prior expression:
			// Use a column for the key such that consecutive entries
			// can align if possible.
			// (needsLinebreak is set if we started a new line before)
			p.expr(pair.Key)
			p.print(pair.Colon, token.COLON, vtab)
			p.expr(pair.Value)
		} else {
			p.expr0(x, depth)
		}

		if size > 0 {
			lnsum += math.Log(float64(size))
			count++
		}

		prevLine = line
	}

	if mode&commaTerm != 0 && next.IsValid() && p.pos.Line < next.Line {
		// Print a terminating comma if the next token is on a new line.
		p.print(token.COMMA)
		if isIncomplete {
			p.print(newline)
			p.print("// " + filteredMsg)
		}
		if ws == ignore && mode&noIndent == 0 {
			// unindent if we indented
			p.print(unindent)
		}
		p.print(formfeed) // terminating comma needs a line break to look good
		return
	}

	if isIncomplete {
		p.print(token.COMMA, newline)
		p.print("// "+filteredMsg, newline)
	}

	if ws == ignore && mode&noIndent == 0 {
		// unindent if we indented
		p.print(unindent)
	}
}

func (p *printer) expr1(expr ast.Expr, prec1, depth int) {
	p.print(expr.Pos())

	switch x := expr.(type) {
	case *ast.BadExpr:
		p.print("BadExpr")

	case *ast.Ident:
		p.print(x)

	case *ast.BinaryExpr:
		if depth < 1 {
			p.internalError("depth < 1:", depth)
			depth = 1
		}
		p.binaryExpr(x, prec1, cutoff(x, depth), depth)

	case *ast.KeyValueExpr:
		p.expr(x.Key)
		p.print(x.Colon, token.COLON, blank)
		p.expr(x.Value)

	case *ast.StarExpr:
		const prec = token.UnaryPrec
		if prec < prec1 {
			// parenthesis needed
			p.print(token.LPAREN)
			p.print(token.MUL)
			p.expr(x.X)
			p.print(token.RPAREN)
		} else {
			// no parenthesis needed
			p.print(token.MUL)
			p.expr(x.X)
		}

	case *ast.UnaryExpr:
		const prec = token.UnaryPrec
		if prec < prec1 {
			// parenthesis needed
			p.print(token.LPAREN)
			p.expr(x)
			p.print(token.RPAREN)
		} else {
			// no parenthesis needed
			p.print(x.Op)
			if x.Op == token.RANGE || x.Op == token.Zh_迭代 {
				// TODO(gri) Remove this code if it cannot be reached.
				p.print(blank)
			}
			p.expr1(x.X, prec, depth)
		}

	case *ast.BasicLit:
		p.print(x)

	case *ast.FuncLit:
		p.expr(x.Type)
		p.funcBody(p.distanceFrom(x.Type.Pos()), blank, x.Body)

	case *ast.ParenExpr:
		if _, hasParens := x.X.(*ast.ParenExpr); hasParens {
			// don't print parentheses around an already parenthesized expression
			// TODO(gri) consider making this more general and incorporate precedence levels
			p.expr0(x.X, depth)
		} else {
			p.print(token.LPAREN)
			p.expr0(x.X, reduceDepth(depth)) // parentheses undo one level of depth
			p.print(x.Rparen, token.RPAREN)
		}

	case *ast.SelectorExpr:
		p.selectorExpr(x, depth, false)

	case *ast.TypeAssertExpr:
		p.expr1(x.X, token.HighestPrec, depth)
		p.print(token.PERIOD, x.Lparen, token.LPAREN)
		if x.Type != nil {
			p.expr(x.Type)
		} else {
			p.print(token.Zh_类型)
		}
		p.print(x.Rparen, token.RPAREN)

	case *ast.IndexExpr:
		// TODO(gri): should treat[] like parentheses and undo one level of depth
		p.expr1(x.X, token.HighestPrec, 1)
		p.print(x.Lbrack, token.LBRACK)
		p.expr0(x.Index, depth+1)
		p.print(x.Rbrack, token.RBRACK)

	case *ast.IndexListExpr:
		p.expr1(x.X, token.HighestPrec, 1)
		p.print(x.Lbrack, token.LBRACK)
		p.exprList(x.Lbrack, x.Indices, depth+1, commaTerm, x.Rbrack, false)
		p.print(x.Rbrack, token.RBRACK)

	case *ast.SliceExpr:
		// TODO(gri): should treat[] like parentheses and undo one level of depth
		p.expr1(x.X, token.HighestPrec, 1)
		p.print(x.Lbrack, token.LBRACK)
		indices := []ast.Expr{x.Low, x.High}
		if x.Max != nil {
			indices = append(indices, x.Max)
		}
		// determine if we need extra blanks around ':'
		var needsBlanks bool
		if depth <= 1 {
			var indexCount int
			var hasBinaries bool
			for _, x := range indices {
				if x != nil {
					indexCount++
					if isBinary(x) {
						hasBinaries = true
					}
				}
			}
			if indexCount > 1 && hasBinaries {
				needsBlanks = true
			}
		}
		for i, x := range indices {
			if i > 0 {
				if indices[i-1] != nil && needsBlanks {
					p.print(blank)
				}
				p.print(token.COLON)
				if x != nil && needsBlanks {
					p.print(blank)
				}
			}
			if x != nil {
				p.expr0(x, depth+1)
			}
		}
		p.print(x.Rbrack, token.RBRACK)

	case *ast.CallExpr:
		if len(x.Args) > 1 {
			depth++
		}
		var wasIndented bool
		if _, ok := x.Fun.(*ast.FuncType); ok {
			// conversions to literal function types require parentheses around the type
			p.print(token.LPAREN)
			wasIndented = p.possibleSelectorExpr(x.Fun, token.HighestPrec, depth)
			p.print(token.RPAREN)
		} else {
			wasIndented = p.possibleSelectorExpr(x.Fun, token.HighestPrec, depth)
		}
		p.print(x.Lparen, token.LPAREN)
		if x.Ellipsis.IsValid() {
			p.exprList(x.Lparen, x.Args, depth, 0, x.Ellipsis, false)
			p.print(x.Ellipsis, token.ELLIPSIS)
			if x.Rparen.IsValid() && p.lineFor(x.Ellipsis) < p.lineFor(x.Rparen) {
				p.print(token.COMMA, formfeed)
			}
		} else {
			p.exprList(x.Lparen, x.Args, depth, commaTerm, x.Rparen, false)
		}
		p.print(x.Rparen, token.RPAREN)
		if wasIndented {
			p.print(unindent)
		}

	case *ast.CompositeLit:
		// composite literal elements that are composite literals themselves may have the type omitted
		if x.Type != nil {
			p.expr1(x.Type, token.HighestPrec, depth)
		}
		p.level++
		p.print(x.Lbrace, token.LBRACE)
		p.exprList(x.Lbrace, x.Elts, 1, commaTerm, x.Rbrace, x.Incomplete)
		// do not insert extra line break following a /*-style comment
		// before the closing '}' as it might break the code if there
		// is no trailing ','
		mode := noExtraLinebreak
		// do not insert extra blank following a /*-style comment
		// before the closing '}' unless the literal is empty
		if len(x.Elts) > 0 {
			mode |= noExtraBlank
		}
		// need the initial indent to print lone comments with
		// the proper level of indentation
		p.print(indent, unindent, mode, x.Rbrace, token.RBRACE, mode)
		p.level--

	case *ast.Ellipsis:
		p.print(token.ELLIPSIS)
		if x.Elt != nil {
			p.expr(x.Elt)
		}

	case *ast.ArrayType:
		p.print(token.LBRACK)
		if x.Len != nil {
			p.expr(x.Len)
		}
		p.print(token.RBRACK)
		p.expr(x.Elt)

	case *ast.StructType:
		// skip: 已经打印过了
		p.fieldList(x.Fields, true, x.Incomplete)

	case *ast.FuncType:
		p.print(token.Zh_函数)
		p.signature(x)

	case *ast.InterfaceType:
		// skip: 已经打印过了
		p.fieldList(x.Methods, false, x.Incomplete)

	case *ast.MapType:
		p.print(token.Zh_字典, token.LBRACK)
		p.expr(x.Key)
		p.print(token.RBRACK)
		p.expr(x.Value)

	default:
		panic("unreachable")
	}
}

func (p *printer) possibleSelectorExpr(expr ast.Expr, prec1, depth int) bool {
	if x, ok := expr.(*ast.SelectorExpr); ok {
		return p.selectorExpr(x, depth, true)
	}
	p.expr1(expr, prec1, depth)
	return false
}

// selectorExpr handles an *ast.SelectorExpr node and reports whether x spans
// multiple lines.
func (p *printer) selectorExpr(x *ast.SelectorExpr, depth int, isMethod bool) bool {
	p.expr1(x.X, token.HighestPrec, depth)
	p.print(token.PERIOD)
	if line := p.lineFor(x.Sel.Pos()); p.pos.IsValid() && p.pos.Line < line {
		p.print(indent, newline, x.Sel.Pos(), x.Sel)
		if !isMethod {
			p.print(unindent)
		}
		return true
	}
	p.print(x.Sel.Pos(), x.Sel)
	return false
}

func (p *printer) expr0(x ast.Expr, depth int) {
	p.expr1(x, token.LowestPrec, depth)
}

func (p *printer) expr(x ast.Expr) {
	const depth = 1
	p.expr1(x, token.LowestPrec, depth)
}

func (p *printer) exprTypeSpec(x ast.Expr) {
	const depth = 1
	p.expr1(x, token.LowestPrec, depth)
}

func cutoff(e *ast.BinaryExpr, depth int) int {
	has4, has5, maxProblem := walkBinary(e)
	if maxProblem > 0 {
		return maxProblem + 1
	}
	if has4 && has5 {
		if depth == 1 {
			return 5
		}
		return 4
	}
	if depth == 1 {
		return 6
	}
	return 4
}

func diffPrec(expr ast.Expr, prec int) int {
	x, ok := expr.(*ast.BinaryExpr)
	if !ok || prec != x.Op.Precedence() {
		return 1
	}
	return 0
}

func reduceDepth(depth int) int {
	depth--
	if depth < 1 {
		depth = 1
	}
	return depth
}

func isBinary(expr ast.Expr) bool {
	_, ok := expr.(*ast.BinaryExpr)
	return ok
}
