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
	if file.Name != nil && file.Name.Name != "" {
		p.print(file.Pos(), token.PACKAGE, blank)
		p.expr(file.Name)
	}

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
			p.print(d.Pos(), d.Tok, blank)

			assert(len(d.Specs) == 1)
			switch s := d.Specs[0].(type) {
			case *ast.ImportSpec:
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
				p.spec_ValueSpec(s, 1, true)

			case *ast.TypeSpec:
				// 结构/接口
				p.setComment(s.Doc)
				p.expr(s.Name)
				p.exprTypeSpec(s.Type)
				p.setComment(s.Comment)

			default:
				panic("unreachable")
			}

		case *ast.FuncDecl:
			p.funcDecl(d)

		default:
			panic("unreachable")
		}
	}

	p.print(newline)
	return nil
}
