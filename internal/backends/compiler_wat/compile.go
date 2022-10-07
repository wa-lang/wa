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

	p.CompilePackage(prog.SSAMainPkg)

	return p.module.ToWatModule().String(), nil
}

func (p *Compiler) CompilePackage(ssaPkg *ssa.Package) {

	p.ssaPkg = ssaPkg

	var ts []*ssa.Type
	var cs []*ssa.NamedConst
	var gs []*ssa.Global
	var fns []*ssa.Function

	{
		var sig wat.FuncSig
		sig.Params = []wat.ValueType{wat.I32{}}
		//p.module.Imports = append(p.module.Imports, wat.NewImpFunc("js", "print_i32", "$$print_i32", sig))
		//p.module.Imports = append(p.module.Imports, wat.NewImpFunc("js", "print_rune", "$$print_rune", sig))
	}

	for _, m := range p.ssaPkg.Members {
		switch member := m.(type) {
		case *ssa.Type:
			ts = append(ts, member)
		case *ssa.NamedConst:
			cs = append(cs, member)
		case *ssa.Global:
			gs = append(gs, member)
		case *ssa.Function:
			//fns = append(fns, member)
		default:
			panic("Unreachable")
		}
	}

	for _, v := range ts {
		p.compileType(v)
	}

	for _, v := range cs {
		p.compileConst(v)
	}

	for _, v := range gs {
		p.compileGlobal(v)
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
		p.module.Funcs = append(p.module.Funcs, newFunctionGenerator(p).genFunction(v))
	}

	//println(p.module.String())
}
