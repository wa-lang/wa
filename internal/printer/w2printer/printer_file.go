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

		p.decl(d, true)
	}

	p.print(newline)
	return nil
}
