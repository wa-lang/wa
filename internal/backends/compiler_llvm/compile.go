// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/ssa"
)

type FmtStr struct {
	fmt  string
	size int
}

type Compiler struct {
	target string
	output strings.Builder
	debug  bool
	fmts   []FmtStr
	anofn  []*ssa.Function
}

func New(target string, debug bool) *Compiler {
	p := new(Compiler)
	p.target = target
	p.debug = debug
	return p
}

func (p *Compiler) Compile(prog *loader.Program) (output string, err error) {
	if prog == nil || prog.SSAMainPkg == nil {
		return "", errors.New("invalid or empty input program")
	}

	// Compile each package.
	for _, pkg := range prog.Pkgs {
		// TODO: need a better way to supress the useless 'runtime' package.
		if pkg.Pkg.Name() == "runtime" {
			continue
		}
		if err := p.compilePackage(pkg.SSAPkg); err != nil {
			return "", err
		}
	}

	// Emit each format string as a global variable.
	for i, f := range p.fmts {
		fvar := fmt.Sprintf("@printfmt%d = private constant [%d x i8] c\"%s\"\n", i, f.size, f.fmt)
		p.output.WriteString(fvar)
		if i == len(p.fmts)-1 {
			p.output.WriteString("\n")
		}
	}

	// Emit some auxiliary functions.
	p.output.WriteString("define void @wa_main() {\n")
	p.output.WriteString("  call void @")
	p.output.WriteString(getNormalName(prog.SSAMainPkg.Pkg.Path() + ".init()\n"))
	p.output.WriteString("  call void @")
	p.output.WriteString(getNormalName(prog.SSAMainPkg.Pkg.Path() + ".main()\n"))
	p.output.WriteString("  ret void\n")
	p.output.WriteString("}\n\n")

	// Emit some target specific functions.
	switch getArch(p.target) {
	case "", "x86_64", "aarch64":
		p.output.WriteString("define void @main() {\n")
		p.output.WriteString("  call void @wa_main()\n")
		p.output.WriteString("  ret void\n")
		p.output.WriteString("}\n\n")
		p.output.WriteString("declare i32 @printf(i8* readonly, ...)\n")
	case "avr":
		p.output.WriteString("define void @__avr_write_port__(i16 %0, i8 zeroext %1) {\n")
		p.output.WriteString("  %3 = inttoptr i16 %0 to i8*\n")
		p.output.WriteString("  store volatile i8 %1, i8* %3, align 1\n")
		p.output.WriteString("  ret void\n")
		p.output.WriteString("}\n\n")
		p.output.WriteString("define zeroext i8 @__avr_read_port__(i16 %0) {\n")
		p.output.WriteString("  %2 = inttoptr i16 %0 to i8*\n")
		p.output.WriteString("  %3 = load volatile i8, i8* %2, align 1\n")
		p.output.WriteString("  ret i8 %3\n")
		p.output.WriteString("}\n\n")
	default:
	}

	return p.output.String(), nil
}

func (p *Compiler) compilePackage(pkg *ssa.Package) error {
	var gs []*ssa.Global    // global variables
	var fns []*ssa.Function // global functions & methods

	for _, m := range pkg.Members {
		switch member := m.(type) {
		case *ssa.Global:
			gs = append(gs, member)
		case *ssa.Function:
			fns = append(fns, member)
		case *ssa.Type, *ssa.NamedConst:
			// Omit name consts and type definitions
		default:
			panic("unknown global object")
		}
	}

	// Collect methods which are not treated as functions.
	for _, v := range pkg.GetValues() {
		if f, ok := v.(*ssa.Function); ok {
			if !findFunc(f, fns) {
				fns = append(fns, f)
			}
		}
	}

	// Emit all global variables.
	for _, gv := range gs {
		// Dump each global variable in IR form.
		if p.debug {
			p.output.WriteString("; ")
			p.output.WriteString(gv.String())
			p.output.WriteString("\n")
		}
		// Compile each global variable to LLVM-IR form.
		p.output.WriteString("@")
		p.output.WriteString(getNormalName(gv.Pkg.Pkg.Path() + "." + gv.Name()))
		p.output.WriteString(" = global ")
		tys := getTypeStr(gv.Type(), p.target)
		p.output.WriteString(tys[0:(len(tys) - 1)])
		p.output.WriteString(" zeroinitializer\n")
	}
	if len(gs) > 0 {
		p.output.WriteString("\n")
	}

	// Generate LLVM-IR for each global function.
	for _, v := range fns {
		if err := p.compileFunction(v); err != nil {
			return err
		}
	}

	// Generate LLVM-IR for each internal function.
	for _, v := range p.anofn {
		if err := p.compileFunction(v); err != nil {
			return err
		}
	}
	p.anofn = []*ssa.Function{}

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
