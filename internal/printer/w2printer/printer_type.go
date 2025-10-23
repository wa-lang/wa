// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) declTypeAssign(s *ast.TypeSpec) {
	p.print(s.Pos(), token.Zh_类型, token.K_点)
	p.setComment(s.Doc)
	p.expr(s.Name)
	p.print(blank, token.ASSIGN, blank)
	switch x := s.Type.(type) {
	case *ast.Ident:
		p.print(x)
	case *ast.SelectorExpr:
		p.print(x.X.(*ast.Ident), token.K_点, x.Sel)
	default:
		panic("unreachable")
	}
	p.setComment(s.Comment)
}
