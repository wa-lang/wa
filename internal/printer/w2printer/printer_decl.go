// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
)

func (p *printer) decl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.BadDecl:
		p.print(d.Pos(), "BadDecl")

	case *ast.GenDecl:
		p.setComment(d.Doc)
		p.print(d.Pos(), d.Tok, blank)
		assert(len(d.Specs) == 1)

		if s, ok := d.Specs[0].(*ast.ValueSpec); ok {
			p.spec_ValueSpec(s, 1, true)
		} else {
			panic("unreachable")
		}

	case *ast.FuncDecl:
		p.funcDecl(d)

	default:
		panic("unreachable")
	}
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
