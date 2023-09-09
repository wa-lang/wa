// 版权 @2023 凹语言 作者。保留所有权利。

package appdoc

import (
	"fmt"
	"log"

	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/types"
)

type PkgDoc struct {
	Prog    *loader.Program
	Consts  []*types.Const
	Types   []*types.TypeName
	Globals []*types.Var
	Funcs   []*types.Func
}

func BuildPkgDoc(prog *loader.Program) (p *PkgDoc) {
	p = &PkgDoc{Prog: prog}
	p.buildDoc()
	return p
}

func (p *PkgDoc) Show(names ...string) {
	if len(names) != 0 {
		fmt.Println("TODO:", names)
		return
	}

	for _, c := range p.Consts {
		// typ := c.Type().(*types.Basic)
		fmt.Printf("const %s %v\n", c.Name(), c.Type())
	}
	for _, c := range p.Globals {
		fmt.Printf("global %s %v\n", c.Name(), c.Type())
	}
	for _, fn := range p.Funcs {
		sig := fn.Type().(*types.Signature)
		if sig.Results().Len() > 0 {
			fmt.Printf("func %s%v => %v\n", fn.Name(), sig.Params(), sig.Results())
		} else {
			fmt.Printf("func %s%v\n", fn.Name(), sig.Params())
		}
	}
}

func (p *PkgDoc) buildDoc() {
	mainPkg := p.Prog.Pkgs[p.Prog.Manifest.MainPkg]
	scope := mainPkg.Pkg.Scope()

	for _, s := range scope.Names() {
		if obj := scope.Lookup(s); obj.Exported() {
			switch x := obj.(type) {
			case *types.Const:
				p.Consts = append(p.Consts, x)
			case *types.TypeName:
				p.Types = append(p.Types, x)
			case *types.Var:
				p.Globals = append(p.Globals, x)
			case *types.Func:
				p.Funcs = append(p.Funcs, x)

			default:
				log.Fatalf("unknown %[1]T %[1]v\n", obj)
			}
		}
	}
}
