// 版权 @2023 凹语言 作者。保留所有权利。

package appbuild

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wabt"
)

var CmdBuild = &cli.Command{
	Name:  "build",
	Usage: "compile Wa source code",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
			Value:   "",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: fmt.Sprintf("set target os (%s)", strings.Join(config.WaOS_List, "|")),
			Value: config.WaOS_Default,
		},
		&cli.StringFlag{
			Name:  "tags",
			Usage: "set build tags",
		},
		&cli.IntFlag{
			Name:  "ld-stack-size",
			Usage: "set stack size",
		},
		&cli.IntFlag{
			Name:  "ld-max-memory",
			Usage: "set max memory size",
		},
	},
	Action: func(c *cli.Context) error {
		var input string
		if c.NArg() > 0 {
			input = c.Args().First()
		} else {
			input, _ = os.Getwd()
		}

		outfile := c.String("output")
		if outfile == "" {
			if fi, _ := os.Lstat(input); fi != nil && fi.IsDir() {
				// /xxx/yyy/ => /xxx/yyy/output/yyy.wat
			} else {
				// /xxx/yyy/zzz.wa => /xxx/yyy/zzz.wat
			}
			outfile = "a.out.wasm"
		}

		opt := appbase.BuildOptions(c)
		watOutput, err := Build(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if outfile != "" && outfile != "-" {
			watFilename := outfile
			if !strings.HasSuffix(watFilename, ".wat") {
				watFilename += ".wat"
			}
			err := os.WriteFile(watFilename, []byte(watOutput), 0666)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if strings.HasSuffix(outfile, ".wasm") {
				wasmBytes, err := wabt.Wat2Wasm(watOutput)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err := os.WriteFile(outfile, wasmBytes, 0666); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Remove(watFilename)
			}
		} else {
			fmt.Println(string(watOutput))
		}

		return nil
	},
}

func Build(opt *appbase.Option, filename string) ([]byte, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, fmt.Errorf("%q not found", filename)
	}
	cfg := opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return nil, err
	}

	output, err := compiler_wat.New().Compile(prog, "main")

	if err != nil {
		return nil, err
	}

	return []byte(output), nil
}
