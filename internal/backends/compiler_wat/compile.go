// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/ssa"
)

type Compiler struct {
	ssaPkg *ssa.Package

	module *wir.Module
}

func New() *Compiler {
	p := new(Compiler)
	p.module = wir.NewModule()
	p.module.AddGlobal("$wa.RT.closure_data", wir.NewRef(wir.VOID{}), false, nil)
	wir.SetCurrentModule(p.module)
	return p
}

func (p *Compiler) Compile(prog *loader.Program, mainFunc string) (output string, err error) {
	p.module.BaseWat = modBaseWat_wasi

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
		f.InternalName, f.ExternalName = "_start", "_start"
		f.Insts = append(f.Insts, wat.NewInstCall("$waGlobalAlloc"))
		n, _ := wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path())
		n += ".init"
		f.Insts = append(f.Insts, wat.NewInstCall(n))
		n, _ = wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path())
		n += "."
		n += mainFunc
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
		if len(v.Blocks) < 1 {
			if v.RuntimeGetter() {
				p.module.AddFunc(newFunctionGenerator(p.module).genGetter(v))
			} else if v.RuntimeSetter() {
				p.module.AddFunc(newFunctionGenerator(p.module).genSetter(v))
			} else if iname0, iname1 := v.ImportName(); len(iname0) > 0 && len(iname1) > 0 {
				var fn_name string
				if len(v.LinkName()) > 0 {
					fn_name = v.LinkName()
				} else {
					fn_name, _ = GetFnMangleName(v)
				}

				sig := wir.NewFnSigFromSignature(v.Signature)
				p.module.AddImportFunc(iname0, iname1, fn_name, sig)
			}
			continue
		}
		p.module.AddFunc(newFunctionGenerator(p.module).genFunction(v))
	}
}
