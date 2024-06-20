// 版权 @2023 凹语言 作者。保留所有权利。

package apprun

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/wazero"
)

var CmdRunWasm = &cli.Command{
	Name:      "run-wasm",
	Usage:     "run wasm program",
	ArgsUsage: "file.wasm",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "main-func",
			Usage: "set main func",
			Value: "_main",
		},
	},
	Action: CmdRunAction,
}

func CmdRunAction(c *cli.Context) error {
	input := c.Args().First()

	if !appbase.HasExt(input, ".wasm") {
		fmt.Printf("%q is not valid wasm file\n", input)
		os.Exit(1)
	}

	wasmBytes, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	var appArgs []string
	if c.NArg() > 2 {
		appArgs = c.Args().Slice()[2:]
	}

	var mainFunc = c.String("main-func")

	stdout, stderr, err := wazero.RunWasm(input, wasmBytes, mainFunc, appArgs...)
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
		return nil
	}
	if len(stdout) > 0 {
		fmt.Fprint(os.Stdout, string(stdout))
	}
	if len(stderr) > 0 {
		fmt.Fprint(os.Stderr, string(stderr))
	}
	return nil
}
