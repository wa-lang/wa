// 版权 @2023 凹语言 作者。保留所有权利。

package appcir

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_c"
	"wa-lang.org/wa/internal/loader"
)

var CmdCir = &cli.Command{
	Hidden: true,
	Name:   "cir",
	Usage:  "print cir code",
	Flags: []cli.Flag{
		appbase.MakeFlag_output(),
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		opt := appbase.BuildOptions(c)
		err := PrintCIR(opt, c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

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
