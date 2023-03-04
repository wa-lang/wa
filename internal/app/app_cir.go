// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"wa-lang.org/wa/internal/backends/compiler_c"
	"wa-lang.org/wa/internal/loader"
)

func (p *App) CIR(filename string) error {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	var c compiler_c.CompilerC
	c.CompilePackage(prog.SSAMainPkg)
	print("\n\n")
	print(c.String())

	return nil
}
