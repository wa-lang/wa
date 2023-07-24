// 版权 @2023 凹语言 作者。保留所有权利。

package apprun

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/wazero"
)

var CmdRun = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "target",
			Usage: fmt.Sprintf("set target os (%s)", strings.Join(config.WaOS_List, "|")),
			Value: config.WaOS_Default,
		},
		&cli.StringFlag{
			Name:  "tags",
			Usage: "set build tags",
		},
	},
	Action: CmdRunAction,
}

func CmdRunAction(c *cli.Context) error {
	input := c.Args().First()
	outfile := ""

	if input == "" {
		input, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	wasmBytes, err := appbuild.BuildApp(opt, input, outfile)
	if err != nil {
		return err
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()[1:]
	}

	stdout, stderr, err := wazero.RunWasm(opt.Config(), input, wasmBytes, appArgs...)
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
	return nil
}
