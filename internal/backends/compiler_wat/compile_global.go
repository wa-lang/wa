// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) compileGlobal(g *ssa.Global) {
	if len(g.LinkName()) > 0 {
		p.module.AddGlobal(g.LinkName(), wir.ToWType(g.Type().(*types.Pointer).Elem()), false, g)
	} else {
		if g.Name() == "init$guard" {
			p.module.AddGlobal(g.Name(), wir.ToWType(g.Type().(*types.Pointer).Elem()), false, g)
		} else {
			p.module.AddGlobal(g.Name(), wir.ToWType(g.Type().(*types.Pointer).Elem()), true, g)
		}
	}
	//logger.Fatal("Todo")
}
