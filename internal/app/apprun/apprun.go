// 版权 @2023 凹语言 作者。保留所有权利。

package apprun

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/wabt"
	"wa-lang.org/wa/internal/wazero"
)

func Run(c *cli.Context) {
	var infile string
	if c.NArg() > 0 {
		infile = c.Args().First()
	} else {
		infile, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	var watBytes []byte
	var wasmBytes []byte
	var err error

	switch {
	case strings.HasSuffix(infile, ".wat"):
		watBytes, err = os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wasmBytes, err = wabt.Wat2Wasm(watBytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case strings.HasSuffix(infile, ".wasm"):
		wasmBytes, err = os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		watBytes, err = appbuild.Build(opt, infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = os.WriteFile("a.out.wat", watBytes, 0666); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wasmBytes, err = wabt.Wat2Wasm(watBytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()[1:]
	}

	stdout, stderr, err := wazero.RunWasm(opt.Config(), infile, wasmBytes, appArgs...)
	if err != nil {
		if len(stdout) > 0 {
			fmt.Fprint(os.Stdout, string(stdout))
		}
		if len(stderr) > 0 {
			fmt.Fprint(os.Stderr, string(stderr))
		}
		if exitCode, ok := wazero.AsExitError(err); ok {
			os.Exit(exitCode)
		}
		fmt.Println(err)
	}
	if len(stdout) > 0 {
		fmt.Fprint(os.Stdout, string(stdout))
	}
	if len(stderr) > 0 {
		fmt.Fprint(os.Stderr, string(stderr))
	}
}
