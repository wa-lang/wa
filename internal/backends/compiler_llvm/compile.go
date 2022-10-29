// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"
	"strings"

	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/ssa"
)

type FmtStr struct {
	fmt  string
	size int
}

type Compiler struct {
	ssaPkg *ssa.Package
	target string
	output strings.Builder
	fmts   []FmtStr
}

func New(target string) *Compiler {
	p := new(Compiler)
	p.target = target
	return p
}

func (p *Compiler) Compile(prog *loader.Program) (output string, err error) {
	if prog == nil || prog.SSAMainPkg == nil {
		return "", errors.New("invalid or empty input program")
	}
	p.ssaPkg = prog.SSAMainPkg
	p.output.WriteString("declare i32 @printf(i8* readonly, ...)\n\n")
	if err := p.compilePackage(); err != nil {
		return "", err
	}
	return p.output.String(), nil
}

func (p *Compiler) compilePackage() error {
	var ts []*ssa.Type       // global types
	var cs []*ssa.NamedConst // global constants
	var gs []*ssa.Global     // global variables
	var fns []*ssa.Function  // global functions & methods

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
			panic("unknown global object")
		}
	}

	// Collect methods which are not treated as functions.
	for _, v := range p.ssaPkg.GetValues() {
		if f, ok := v.(*ssa.Function); ok {
			if !findFunc(f, fns) {
				fns = append(fns, f)
			}
		}
	}

	// Generate LLVM-IR for each global function.
	for _, v := range fns {
		// TODO: Support the builtin function 'init'.
		if v.Name() == "init" {
			continue
		}
		if err := p.compileFunction(v); err != nil {
			return err
		}
	}

	return nil
}

func findFunc(f *ssa.Function, fns []*ssa.Function) bool {
	for _, m := range fns {
		if m.Object() == f.Object() {
			return true
		}
	}
	return false
}
