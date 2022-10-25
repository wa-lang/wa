// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"
	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/ssa"
)

type Compiler struct {
	ssaPkg *ssa.Package
	target string
	output string
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
	if err := p.CompilePackage(); err != nil {
		return "", err
	}
	return p.output, nil
}

func (p *Compiler) CompilePackage() error {
	return nil
}
