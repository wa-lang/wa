// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"sort"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/waroot"
)

type Compiler struct {
	prog   *loader.Program
	ssaPkg *ssa.Package

	module *wir.Module
	tLib   *typeLib
}

func New() *Compiler {
	p := new(Compiler)
	p.module = wir.NewModule()
	p.module.AddGlobal("$wa.runtime.closure_data", p.module.GenValueType_Ref(p.module.VOID), false, nil)
	wir.SetCurrentModule(p.module)
	return p
}

func (p *Compiler) Compile(prog *loader.Program, mainFunc string) (output string, err error) {
	p.prog = prog
	p.CompileWsFiles(prog)

	p.tLib = newTypeLib(p.module, prog)

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgConst(pkg.SSAPkg)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgType(pkg.SSAPkg)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgGlobal(pkg.SSAPkg)
	}

	for _, pkg := range prog.Pkgs {
		p.ssaPkg = pkg.SSAPkg
		p.CompilePkgFunc(pkg.SSAPkg)
	}

	p.tLib.finish()

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

func (p *Compiler) CompileWsFiles(prog *loader.Program) {
	var sb strings.Builder

	sb.WriteString(waroot.GetBaseWsCode(config.WaBackend_wat))
	sb.WriteString("\n")

	var pkgpathList = make([]string, 0, len(prog.Pkgs))
	for pkgpath := range prog.Pkgs {
		pkgpathList = append(pkgpathList, pkgpath)
	}
	sort.Strings(pkgpathList)

	var lineCommentSep = ";; -" + strings.Repeat("-", 60-4) + "\n"

	for _, pkgpath := range pkgpathList {
		pkg := prog.Pkgs[pkgpath]
		if len(pkg.WsFiles) == 0 {
			continue
		}

		func() {
			sb.WriteString(lineCommentSep)
			sb.WriteString(";; package: " + pkgpath + "\n")
			sb.WriteString(lineCommentSep)
			sb.WriteString("\n")

			for _, sf := range pkg.WsFiles {
				sb.WriteString(";; file: " + sf.Name + "\n")
				sb.WriteString("\n")

				sb.WriteString(strings.TrimSpace(sf.Code))
				sb.WriteString("\n")
			}
		}()
	}

	p.module.BaseWat = sb.String()
}

func (p *Compiler) CompilePkgConst(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if con, ok := m.(*ssa.NamedConst); ok {
			p.compileConst(con)
		}
	}
}

func (p *Compiler) CompilePkgType(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if t, ok := m.(*ssa.Type); ok {
			p.compileType(t)
		}
	}
}

func (p *Compiler) CompilePkgGlobal(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if global, ok := m.(*ssa.Global); ok {
			p.compileGlobal(global)
		}
	}
}

func (p *Compiler) CompilePkgFunc(ssaPkg *ssa.Package) {
	for _, m := range p.ssaPkg.Members {
		if fn, ok := m.(*ssa.Function); ok {
			CompileFunc(fn, p.prog, p.tLib, p.module)
		}
	}
}

func CompileFunc(f *ssa.Function, prog *loader.Program, tLib *typeLib, module *wir.Module) {
	if len(f.Blocks) < 1 {
		if f.RuntimeGetter() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genGetter(f))
		} else if f.RuntimeSetter() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genSetter(f))
		} else if f.RuntimeSizer() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genSizer(f))
		} else if iname0, iname1 := f.ImportName(); len(iname0) > 0 && len(iname1) > 0 {
			var fn_name string
			if len(f.LinkName()) > 0 {
				fn_name = f.LinkName()
			} else {
				fn_name, _ = wir.GetFnMangleName(f)
			}

			sig := tLib.GenFnSig(f.Signature)
			module.AddImportFunc(iname0, iname1, fn_name, sig)
		}
		return
	}
	module.AddFunc(newFunctionGenerator(prog, module, tLib).genFunction(f))
}
