// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) decl(decl ast.Decl, isFileScope bool) {
	switch d := decl.(type) {
	case *ast.BadDecl:
		p.print(d.Pos(), "BadDecl")
	case *ast.GenDecl:
		p.genDecl(d, isFileScope)
	case *ast.FuncDecl:
		p.funcDecl(d)
	default:
		panic("unreachable")
	}
}

func (p *printer) genDecl(d *ast.GenDecl, isFileScope bool) {
	p.setComment(d.Doc)

	if isFileScope || d.Tok != token.VAR || d.Lparen != token.NoPos {
		tok := d.Tok
		if d.Tok == token.VAR {
			if isFileScope {
				tok = token.Zh_全局
			} else {
				tok = token.Zh_设定
			}
		}
		p.print(d.Pos(), tok, blank)
	}

	if d.Lparen.IsValid() || len(d.Specs) > 1 {
		// group of parenthesized declarations
		p.print(d.Lparen, token.LPAREN)
		if n := len(d.Specs); n > 0 {
			p.print(indent, formfeed)
			if n > 1 && (d.Tok == token.CONST || d.Tok == token.VAR || d.Tok == token.GLOBAL || d.Tok == token.Zh_全局 || d.Tok == token.Zh_设定) {
				// two or more grouped const/var declarations:
				// determine if the type column must be kept
				keepType := keepTypeColumn(d.Specs)
				var line int
				for i, s := range d.Specs {
					if i > 0 {
						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
					}
					p.recordLine(&line)
					p.valueSpec(s.(*ast.ValueSpec), keepType[i])
				}
			} else {
				var line int
				for i, s := range d.Specs {
					if i > 0 {
						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
					}
					p.recordLine(&line)
					p.spec(s, n, false)
				}
			}
			p.print(unindent, formfeed)
		}
		p.print(d.Rparen, token.RPAREN)

	} else if len(d.Specs) > 0 {
		// single declaration
		p.spec(d.Specs[0], 1, true)
	}
}

// The keepTypeColumn function determines if the type column of a series of
// consecutive const or var declarations must be kept, or if initialization
// values (V) can be placed in the type column (T) instead. The i'th entry
// in the result slice is true if the type column in spec[i] must be kept.
//
// For example, the declaration:
//
//		const (
//			foobar int = 42 // comment
//			x          = 7  // comment
//			foo
//	             bar = 991
//		)
//
// leads to the type/values matrix below. A run of value columns (V) can
// be moved into the type column if there is no type for any of the values
// in that column (we only move entire columns so that they align properly).
//
//		matrix        formatted     result
//	                   matrix
//		T  V    ->    T  V     ->   true      there is a T and so the type
//		-  V          -  V          true      column must be kept
//		-  -          -  -          false
//		-  V          V  -          false     V is moved into T column
func keepTypeColumn(specs []ast.Spec) []bool {
	m := make([]bool, len(specs))

	populate := func(i, j int, keepType bool) {
		if keepType {
			for ; i < j; i++ {
				m[i] = true
			}
		}
	}

	i0 := -1 // if i0 >= 0 we are in a run and i0 is the start of the run
	var keepType bool
	for i, s := range specs {
		t := s.(*ast.ValueSpec)
		if t.Values != nil {
			if i0 < 0 {
				// start of a run of ValueSpecs with non-nil Values
				i0 = i
				keepType = false
			}
		} else {
			if i0 >= 0 {
				// end of a run
				populate(i0, i, keepType)
				i0 = -1
			}
		}
		if t.Type != nil {
			keepType = true
		}
	}
	if i0 >= 0 {
		// end of a run
		populate(i0, len(specs), keepType)
	}

	return m
}

// numLines returns the number of lines spanned by node n in the original source.
func (p *printer) numLines(n ast.Node) int {
	if from := n.Pos(); from.IsValid() {
		if to := n.End(); to.IsValid() {
			return p.lineFor(to) - p.lineFor(from) + 1
		}
	}
	return infinity
}
