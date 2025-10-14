// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) file(src *ast.File) {
	// #!...
	if src.Shebang != "" {
		p.print(src.Shebang)
		p.print(newline)
	}

	p.setComment(src.Doc)
	if src.Name != nil && src.Name.Name != "" {
		p.print(src.Pos(), token.PACKAGE, blank)
		p.expr(src.Name)
	}
	p.declList(src.Decls, true)

	p.print(newline)
}
