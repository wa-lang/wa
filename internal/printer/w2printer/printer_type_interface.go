// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) declInterfaceType(s *ast.TypeSpec) {
	p.print(s.Pos(), token.Zh_接口, token.K_点)
	p.setComment(s.Doc)
	p.expr(s.Name)
	p.interface_exprTypeSpec(s.Type.(*ast.InterfaceType))
	p.setComment(s.Comment)
}

func (p *printer) interface_exprTypeSpec(x *ast.InterfaceType) {
	const depth = 1
	p.expr1(x, token.LowestPrec, depth)
}
