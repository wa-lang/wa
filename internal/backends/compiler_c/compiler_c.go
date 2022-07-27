// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_c

import (
	"github.com/wa-lang/wa/internal/backends/compiler_c/cir"
	"github.com/wa-lang/wa/internal/ssa"
)

type CompilerC struct {
	ssaPkg *ssa.Package

	curScope *cir.Scope
}

func (p *CompilerC) CompilePackage(ssaPkg *ssa.Package) {

	p.ssaPkg = ssaPkg
	p.curScope = cir.NewScope(nil)

	var ts []*ssa.Type
	var cs []*ssa.NamedConst
	var gs []*ssa.Global
	var fns []*ssa.Function

	for _, m := range p.ssaPkg.Members {
		switch member := m.(type) {
		case *ssa.Type:
			ts = append(ts, member)
		case *ssa.NamedConst:
			cs = append(cs, member)
		case *ssa.Global:
			gs = append(gs, member)
		case *ssa.Function:
			fns = append(fns, member)
		default:
			panic("Unreachable")
		}
	}

	for _, t := range ts {
		p.compileType(t)
	}

	for _, c := range cs {
		p.compileConst(c)
	}

	for _, g := range gs {
		p.compileGlobal(g)
	}

	for _, fn := range fns {
		newFunctionGenerator(p).genFunction(fn)
	}
}

func (p *CompilerC) String() string {
	return p.curScope.CIRString()
}

/*
func (p *CompilerC) getMangledDefName(ident *ast.Ident) (mangled string) {
	obj, ok := p.pkg.Info.Defs[ident]
	if !ok {
		return fmt.Sprintf("ugo_%s_%s_pos%d", p.pkg.Pkg.Name(), ident.Name, ident.NamePos)
	}

	scope := obj.Parent()
	if scope == types.Universe {
		return fmt.Sprintf("ugo_%s_%s", "builtin", ident.Name)
	}
	if scope.Parent() == types.Universe {
		return fmt.Sprintf("ugo_%s_%s", p.pkg.Pkg.Name(), ident.Name)
	}

	return fmt.Sprintf("ugo_%s_%s_pos%d", p.pkg.Pkg.Name(), ident.Name, obj.Pos())
}

func (p *CompilerC) getMangledUseName(ident *ast.Ident) (mangled string) {
	obj, ok := p.pkg.Info.Uses[ident]
	if !ok {
		obj, ok = p.pkg.Info.Defs[ident]
		if !ok {
			return fmt.Sprintf("ugo_%s_%s_pos%d", p.pkg.Pkg.Name(), ident.Name, ident.NamePos)
		}
	}

	scope := obj.Parent()
	if scope == types.Universe {
		return fmt.Sprintf("ugo_%s_%s", "builtin", ident.Name)
	}
	if scope.Parent() == types.Universe {
		return fmt.Sprintf("ugo_%s_%s", p.pkg.Pkg.Name(), ident.Name)
	}

	return fmt.Sprintf("%s_pos%d", ident.Name, obj.Pos())
}  //*/

func (p *CompilerC) EnterScope(name string) {
	s := p.curScope.AddScope()
	p.curScope = s
}

func (p *CompilerC) LeaveScope() {
	if p.curScope.Parent == nil {
		return
	}

	p.curScope = p.curScope.Parent
}
