// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// Print the statement list indented, but without a newline after the last statement.
// Extra line breaks between statements in the source are respected but at most one
// empty line is printed between statements.
func (p *printer) stmtList(list []ast.Stmt, nindent int, nextIsRBrace bool) {
	if nindent > 0 {
		p.print(indent)
	}
	var line int
	i := 0
	for _, s := range list {
		// ignore empty statements (was issue 3466)
		if _, isEmpty := s.(*ast.EmptyStmt); !isEmpty {
			// nindent == 0 only for lists of switch/select case clauses;
			// in those cases each clause is a new section
			if len(p.output) > 0 {
				// only print line break if we are not at the beginning of the output
				// (i.e., we are not printing only a partial program)
				p.linebreak(p.lineFor(s.Pos()), 1, ignore, i == 0 || nindent == 0 || p.linesFrom(line) > 0)
			}
			p.recordLine(&line)
			p.stmt(s, nextIsRBrace && i == len(list)-1)
			// labeled statements put labels on a separate line, but here
			// we only care about the start line of the actual statement
			// without label - correct line for each label
			for t := s; ; {
				lt, _ := t.(*ast.LabeledStmt)
				if lt == nil {
					break
				}
				line++
				t = lt.Stmt
			}
			i++
		}
	}
	if nindent > 0 {
		p.print(unindent)
	}
}

func (p *printer) stmt(stmt ast.Stmt, nextIsRBrace bool) {
	p.print(stmt.Pos())

	switch s := stmt.(type) {
	case *ast.BadStmt:
		p.print("BadStmt")

	case *ast.DeclStmt:
		p.decl(s.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		// a "correcting" unindent immediately following a line break
		// is applied before the line break if there is no comment
		// between (see writeWhitespace)
		p.print(unindent)
		p.expr(s.Label)
		p.print(s.Colon, token.COLON, indent)
		if e, isEmpty := s.Stmt.(*ast.EmptyStmt); isEmpty {
			if !nextIsRBrace {
				p.print(newline, e.Pos(), token.SEMICOLON)
				break
			}
		} else {
			p.linebreak(p.lineFor(s.Stmt.Pos()), 1, ignore, true)
		}
		p.stmt(s.Stmt, nextIsRBrace)

	case *ast.ExprStmt:
		const depth = 1
		p.expr0(s.X, depth)

	case *ast.IncDecStmt:
		const depth = 1
		p.expr0(s.X, depth+1)
		p.print(s.TokPos, s.Tok)

	case *ast.AssignStmt:
		var depth = 1
		if len(s.Lhs) > 1 && len(s.Rhs) > 1 {
			depth++
		}
		p.exprList(s.Pos(), s.Lhs, depth, 0, s.TokPos, false)
		p.print(blank, s.TokPos, s.Tok, blank)
		p.exprList(s.TokPos, s.Rhs, depth, 0, token.NoPos, false)

	case *ast.DeferStmt:
		p.print(token.Zh_押后, blank)
		p.expr(s.Call)

	case *ast.ReturnStmt:
		p.print(token.Zh_返回)
		if s.Results != nil {
			p.print(blank)
			// Use indentList heuristic to make corner cases look
			// better (issue 1207). A more systematic approach would
			// always indent, but this would cause significant
			// reformatting of the code base and not necessarily
			// lead to more nicely formatted code in general.
			if p.indentList(s.Results) {
				p.print(indent)
				p.exprList(s.Pos(), s.Results, 1, noIndent, token.NoPos, false)
				p.print(unindent)
			} else {
				p.exprList(s.Pos(), s.Results, 1, 0, token.NoPos, false)
			}
		}

	case *ast.BranchStmt:
		p.print(s.Tok)
		if s.Label != nil {
			p.print(blank)
			p.expr(s.Label)
		}

	case *ast.BlockStmt:
		p.print(token.Zh_区块)
		p.block(s, 1, token.Zh_完毕)

	case *ast.IfStmt:
		p.print(token.Zh_如果)
		p.controlClause(false, s.Init, s.Cond, nil)

		closeTok := token.RBRACE
		if s.Else != nil {
			if _, ok := s.Else.(*ast.IfStmt); ok {
				closeTok = token.Zh_或者
			} else {
				closeTok = token.Zh_否则
			}
		}

		p.block(s.Body, 1, closeTok)

		if s.Else != nil {
			switch s.Else.(type) {
			case *ast.IfStmt:
				p.print(blank, token.Zh_或者, blank)
				p.controlClause(false, s.Init, s.Cond, nil)
				p.stmt(s.Else, nextIsRBrace)
			default:
				p.print(blank, token.Zh_否则, blank)
				p.stmt(s.Else, nextIsRBrace)
			}
		}

	case *ast.CaseClause:
		if s.List != nil {
			p.print(token.Zh_有辙, blank)
			p.exprList(s.Pos(), s.List, 1, 0, s.Colon, false)
		} else {
			p.print(token.Zh_没辙)
		}
		p.print(s.Colon, token.COLON)
		p.stmtList(s.Body, 1, nextIsRBrace)

	case *ast.SwitchStmt:
		p.print(token.Zh_找辙)
		p.controlClause(false, s.Init, s.Tag, nil)
		p.block(s.Body, 0, token.Zh_完毕)

	case *ast.TypeSwitchStmt:
		p.print(token.Zh_找辙)
		if s.Init != nil {
			p.print(blank)
			p.stmt(s.Init, false)
			p.print(token.SEMICOLON)
		}
		p.print(blank)
		p.stmt(s.Assign, false)
		p.print(blank)
		p.block(s.Body, 0, token.Zh_完毕)

	case *ast.ForStmt:
		p.print(token.Zh_循环)
		p.controlClause(true, s.Init, s.Cond, s.Post)
		p.block(s.Body, 1, token.Zh_完毕)

	case *ast.RangeStmt:
		p.print(token.Zh_循环, blank)
		if s.Key != nil {
			p.expr(s.Key)
			if s.Value != nil {
				// use position of value following the comma as
				// comma position for correct comment placement
				p.print(s.Value.Pos(), token.COMMA, blank)
				p.expr(s.Value)
			}
			p.print(blank, s.TokPos, s.Tok, blank)
		}
		p.print(token.Zh_迭代, blank)
		p.expr(stripParens(s.X))
		p.print(blank)
		p.block(s.Body, 1, token.Zh_完毕)

	default:
		panic("unreachable")
	}
}

// block prints an *ast.BlockStmt; it always spans at least two lines.
func (p *printer) block(b *ast.BlockStmt, nindent int, closeTok token.Token) {
	p.print(b.Lbrace, token.COLON)
	p.stmtList(b.List, nindent, true)
	p.linebreak(p.lineFor(b.Rbrace), 1, ignore, true)
	p.print(b.Rbrace, closeTok)
}

// indentList reports whether an expression list would look better if it
// were indented wholesale (starting with the very first element, rather
// than starting at the first line break).
func (p *printer) indentList(list []ast.Expr) bool {
	// Heuristic: indentList reports whether there are more than one multi-
	// line element in the list, or if there is any element that is not
	// starting on the same line as the previous one ends.
	if len(list) >= 2 {
		var b = p.lineFor(list[0].Pos())
		var e = p.lineFor(list[len(list)-1].End())
		if 0 < b && b < e {
			// list spans multiple lines
			n := 0 // multi-line element count
			line := b
			for _, x := range list {
				xb := p.lineFor(x.Pos())
				xe := p.lineFor(x.End())
				if line < xb {
					// x is not starting on the same
					// line as the previous one ended
					return true
				}
				if xb < xe {
					// x is a multi-line element
					n++
				}
				line = xe
			}
			return n > 1
		}
	}
	return false
}

func (p *printer) controlClause(isForStmt bool, init ast.Stmt, expr ast.Expr, post ast.Stmt) {
	p.print(blank)
	needsBlank := false
	if init == nil && post == nil {
		// no semicolons required
		if expr != nil {
			p.expr(stripParens(expr))
			needsBlank = true
		}
	} else {
		// all semicolons required
		// (they are not separators, print them explicitly)
		if init != nil {
			p.stmt(init, false)
		}
		p.print(token.SEMICOLON, blank)
		if expr != nil {
			p.expr(stripParens(expr))
			needsBlank = true
		}
		if isForStmt {
			p.print(token.SEMICOLON, blank)
			needsBlank = false
			if post != nil {
				p.stmt(post, false)
				needsBlank = true
			}
		}
	}
	if needsBlank {
		p.print(blank)
	}
}

func stripParens(x ast.Expr) ast.Expr {
	if px, strip := x.(*ast.ParenExpr); strip {
		// parentheses must not be stripped if there are any
		// unparenthesized composite literals starting with
		// a type name
		ast.Inspect(px.X, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ParenExpr:
				// parentheses protect enclosed composite literals
				return false
			case *ast.CompositeLit:
				if isTypeName(x.Type) {
					strip = false // do not strip parentheses
				}
				return false
			}
			// in all other cases, keep inspecting
			return true
		})
		if strip {
			return stripParens(px.X)
		}
	}
	return x
}

func isTypeName(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isTypeName(t.X)
	}
	return false
}
