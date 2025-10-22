// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) printFile(file *ast.File) error {
	p.comments = file.Comments
	p.useNodeComments = p.comments == nil

	p.nextComment()

	// #!...
	if file.Shebang != "" {
		p.print(file.Shebang)
		p.print(newline)
	}

	p.setComment(file.Doc)

	tok := token.ILLEGAL
	for _, d := range file.Decls {
		prev := tok

		switch d := d.(type) {
		case *ast.GenDecl:
			tok = d.Tok
		case *ast.FuncDecl:
			tok = d.Type.Tok
		default:
			tok = token.ILLEGAL
		}

		// If the declaration token changed (e.g., from CONST to TYPE)
		// or the next declaration has documentation associated with it,
		// print an empty line between top-level declarations.
		// (because p.linebreak is called with the position of d, which
		// is past any documentation, the minimum requirement is satisfied
		// even w/o the extra getDoc(d) nil-check - leave it in case the
		// linebreak logic improves - there's already a TODO).
		if len(p.output) > 0 {
			// only print line break if we are not at the beginning of the output
			// (i.e., we are not printing only a partial program)
			min := 1
			if prev != tok || getDoc(d) != nil {
				min = 2
			}
			// start a new section if the next declaration is a function
			// that spans multiple lines (see also issue #19544)
			p.linebreak(p.lineFor(d.Pos()), min, ignore, tok == token.Zh_函数 && p.numLines(d) > 1)
		}

		switch d := d.(type) {
		case *ast.BadDecl:
			p.print(d.Pos(), "BadDecl")

		case *ast.GenDecl:
			p.setComment(d.Doc)

			switch s := d.Specs[0].(type) {
			case *ast.ImportSpec:
				assert(len(d.Specs) == 1)

				p.print(d.Pos(), token.Zh_引入, blank)

				p.setComment(s.Doc)
				p.expr(sanitizeImportPath(s.Path))
				if s.Name != nil {
					p.print(blank)
					p.print(token.ARROW)
					p.print(blank)

					p.expr(s.Name)
				}
				p.setComment(s.Comment)
				p.print(s.EndPos)

			case *ast.ValueSpec:
				if len(d.Specs) > 1 {
					var tok token.Token
					switch d.Tok {
					case token.VAR, token.GLOBAL, token.Zh_全局:
						tok = token.Zh_全局
					case token.CONST, token.Zh_常量:
						tok = token.Zh_常量
					default:
						panic("unreachale")
					}
					p.print(d.Lparen, tok, token.COLON)
					if n := len(d.Specs); n > 0 {
						p.print(indent, formfeed)
						if n > 1 {
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
								p.spec_ValueSpec(s.(*ast.ValueSpec), n, false)
							}
						}
						p.print(unindent, formfeed)
					}
					p.print(d.Rparen, token.Zh_完毕)
				} else {
					switch d.Tok {
					case token.CONST, token.Zh_常量:
						p.print(d.Pos(), token.Zh_常量, token.K_点)
						p.spec_ValueSpec(s, 1, true)
					case token.VAR, token.GLOBAL, token.Zh_全局:
						p.print(d.Pos(), token.Zh_全局, token.K_点)
						p.spec_ValueSpec(s, 1, true)
					default:
						panic("unreachable")
					}
				}

			case *ast.TypeSpec:
				assert(len(d.Specs) == 1)

				switch s.Type.(type) {
				case *ast.StructType:
					p.declStructType(s)
				case *ast.InterfaceType:
					p.declInterfaceType(s)
				default:
					panic("unreachable")
				}

			default:
				panic("unreachable")
			}

		case *ast.FuncDecl:
			p.funcDecl_pkg(d)

		default:
			panic("unreachable")
		}
	}

	p.print(newline)
	return nil
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

func (p *printer) valueSpec(s *ast.ValueSpec, keepType bool) {
	p.setComment(s.Doc)
	p.identList(s.Names, false) // always present
	extraTabs := 3
	if s.Type != nil || keepType {
		p.print(vtab)
		extraTabs--
	}
	if s.Type != nil {
		p.print(token.COLON)
		p.expr(s.Type)
	}
	if s.Values != nil {
		p.print(vtab, token.ASSIGN, blank)
		p.exprList(token.NoPos, s.Values, 1, 0, token.NoPos, false)
		extraTabs--
	}
	if s.Comment != nil {
		for ; extraTabs > 0; extraTabs-- {
			p.print(vtab)
		}
		p.setComment(s.Comment)
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
