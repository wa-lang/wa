// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/types"
)

func (p *Compiler) compileGlobal(g *ssa.Global) {
	if len(g.LinkName()) > 0 {
		p.module.AddGlobal(g.LinkName(), "", p.tLib.compile(g.Type().(*types.Pointer).Elem()), false, g)
	} else {
		pkg_name, _ := wir.GetPkgMangleName(g.Pkg.Pkg.Path())
		if g.Name() == "init$guard" || g.ForceRegister() {
			p.module.AddGlobal(pkg_name+"."+g.Name(), "", p.tLib.compile(g.Type().(*types.Pointer).Elem()), false, g)
		} else {
			name_exp := ""
			if g.Pkg.Pkg.Path() == p.prog.Manifest.MainPkg && g.Object().Exported() {
				name_exp = pkg_name + "." + g.Name()
			}
			p.module.AddGlobal(pkg_name+"."+wir.GenSymbolName(g.Name()), name_exp, p.tLib.compile(g.Type()), true, g)
		}
	}
	//logger.Fatal("Todo")
}
