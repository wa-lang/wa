// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"fmt"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir"
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/backends/target_spec"
	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/ssa"
)

type Compiler struct {
	ssaPkg *ssa.Package

	module *wir.Module
}

func New() *Compiler {
	p := new(Compiler)
	p.module = wir.NewModule()
	wir.SetCurrentModule(p.module)
	return p
}

func (p *Compiler) Compile(prog *loader.Program, target target_spec.Machine) (output string, err error) {
	switch target {
	case target_spec.Machine_Wasm32_wa, "":
		p.module.BaseWat = modBaseWat_wa
	case target_spec.Machine_Wasm32_wasi:
		p.module.BaseWat = modBaseWat_wasi
	default:
		return "", fmt.Errorf("compiler_wat.Compiler: unsupport target: %v", target)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgConst(pkg.SSAPkg)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgGlocal(pkg.SSAPkg)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgFunc(pkg.SSAPkg)
	}

	{
		var f wir.Function
		f.Name = "_start"
		f.Insts = append(f.Insts, wat.NewInstCall("$waGlobalAlloc"))
		n := wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path()) + "init"
		f.Insts = append(f.Insts, wat.NewInstCall(n))
		n = wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path()) + "main"
		f.Insts = append(f.Insts, wat.NewInstCall(n))
		p.module.AddFunc(&f)
	}

	return p.module.ToWatModule().String(), nil
}

func (p *Compiler) CompilePkgConst(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if con, ok := m.(*ssa.NamedConst); ok {
			p.compileConst(con)
		}
	}
}

func (p *Compiler) CompilePkgGlocal(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if global, ok := m.(*ssa.Global); ok {
			p.compileGlobal(global)
		}
	}
}

func (p *Compiler) CompilePkgFunc(ssaPkg *ssa.Package) {
	var fns []*ssa.Function

	for _, m := range p.ssaPkg.Members {
		if fn, ok := m.(*ssa.Function); ok {
			fns = append(fns, fn)
		}
	}

	for _, v := range ssaPkg.GetValues() {
		if f, ok := v.(*ssa.Function); ok {
			found := false
			for _, m := range fns {
				if m.Object() == f.Object() {
					found = true
				}
			}
			if found {
				continue
			}
			fns = append(fns, f)
		}
	}

	for _, v := range fns {
		p.module.AddFunc(newFunctionGenerator(p).genFunction(v))
	}
}
