package appcir

import (
	"fmt"

	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_c"
	"wa-lang.org/wa/internal/loader"
)

func PrintCIR(opt *appbase.Option, filename string) error {
	cfg := opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	var c compiler_c.CompilerC
	c.CompilePackage(prog.SSAMainPkg)
	fmt.Println(c.String())

	return nil
}
