// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package compiler_c

import (
	"wa-lang.org/wa/internal/backends/compiler_c/cir"
	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/types"
)

func (p *CompilerC) compileGlobal(g *ssa.Global) {
	c_type := cir.ToCType(g.Type().(*types.Pointer).Elem())

	switch c_type.(type) {
	case *ctypes.Array:
		p.curScope.AddLocalVarDecl(g.Name(), c_type).Var.AssociatedSSAObj = g

	default:
		c_name := "$Global_" + g.Name()
		c_var := &p.curScope.AddLocalVarDecl(c_name, c_type).Var

		w_var_decl := p.curScope.AddLocalVarDecl(g.Name(), ctypes.NewPointerType(c_type))
		w_var_decl.InitVal = cir.NewGetaddrExpr(c_var)
		w_var_decl.Var.AssociatedSSAObj = g
	}
}
