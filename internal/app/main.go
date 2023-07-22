// 版权 @2019 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

// 凹语言，The Wa Programming Language.
package app

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appast"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/app/appcir"
	"wa-lang.org/wa/internal/app/appfmt"
	"wa-lang.org/wa/internal/app/appinit"
	"wa-lang.org/wa/internal/app/applex"
	"wa-lang.org/wa/internal/app/appllvm"
	"wa-lang.org/wa/internal/app/applogo"
	"wa-lang.org/wa/internal/app/appplay"
	"wa-lang.org/wa/internal/app/apprun"
	"wa-lang.org/wa/internal/app/appssa"
	"wa-lang.org/wa/internal/app/apptest"
	"wa-lang.org/wa/internal/app/appyacc"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/lsp"
	"wa-lang.org/wa/internal/wabt"
	"wa-lang.org/wa/waroot"
)

func Main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Wa"
	cliApp.Usage = "Wa is a tool for managing Wa source code."
	cliApp.Copyright = "Copyright 2018 The Wa Authors. All rights reserved."
	cliApp.Version = waroot.GetVersion()
	cliApp.EnableBashCompletion = true

	cliApp.CustomAppHelpTemplate = cli.AppHelpTemplate +
		"\nSee \"https://wa-lang.org\" for more information.\n"

	cliApp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "set debug mode",
		},
		&cli.StringFlag{
			Name:    "trace",
			Aliases: []string{"t"},
			Usage:   "set trace mode (*|app|compiler|loader)",
		},
	}

	cliApp.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			config.SetDebugMode()
		}
		if parten := c.String("trace"); parten != "" {
			config.SetEnableTrace(parten)
		}
		return nil
	}

	// 没有参数时显示 help 信息
	cliApp.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	cliApp.Commands = []*cli.Command{
		{
			// go run main.go debug
			Hidden: true,
			Name:   "debug",
			Usage:  "only for dev/debug",
			Action: func(c *cli.Context) error {
				wat, err := api.BuildFile(
					config.DefaultConfig(),
					"hello.wa", "func main() { println(123) }",
				)
				if err != nil {
					if len(wat) != 0 {
						fmt.Println(string(wat))
					}
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(string(wat))
				return nil
			},
		},

		{
			Name:  "play",
			Usage: "start Wa playground",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "http",
					Value: ":2023",
					Usage: "set http address",
				},
			},
			Action: func(c *cli.Context) error {
				err := appplay.RunPlayground(c.String("http"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},

		{
			Name:  "init",
			Usage: "init a sketch Wa module",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "set app name",
					Value:   "hello",
				},
				&cli.StringFlag{
					Name:    "pkgpath",
					Aliases: []string{"p"},
					Usage:   "set pkgpath file",
					Value:   "myapp",
				},
				&cli.BoolFlag{
					Name:    "update",
					Aliases: []string{"u"},
					Usage:   "update example",
				},
			},

			Action: func(c *cli.Context) error {
				err := appinit.InitApp(c.String("name"), c.String("pkgpath"), c.Bool("update"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
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
				watOutput, err := appbuild.Build(opt, input)
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
		},
		{
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
			Action: func(c *cli.Context) error {
				apprun.Run(c)
				return nil
			},
		},
		{
			Name:  "fmt",
			Usage: "format Wa source code file",
			Action: func(c *cli.Context) error {
				err := appfmt.Fmt(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},

		{
			Name:  "test",
			Usage: "test Wa packages",
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
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowAppHelpAndExit(c, 0)
				}
				appArgs := c.Args().Slice()[1:]
				opt := appbase.BuildOptions(c)
				apptest.RunTest(opt.Config(), c.Args().First(), appArgs...)
				return nil
			},
		},

		{
			Hidden: true,
			Name:   "native",
			Usage:  "compile Wa source code to native executable",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "set output file",
					Value:   "",
				},
				&cli.StringFlag{
					Name:  "target",
					Usage: "set native target",
					Value: "",
				},
				&cli.StringFlag{
					Name:  "tags",
					Usage: "set build tags",
				},
				&cli.BoolFlag{
					Name:  "debug",
					Usage: "dump orginal intermediate representation",
				},
				&cli.StringFlag{
					Name:  "clang",
					Usage: "set llvm/clang path",
				},
				&cli.StringFlag{
					Name:  "llc",
					Usage: "set llvm/llc path",
				},
			},
			Action: func(c *cli.Context) error {
				outfile := c.String("output")
				target := c.String("target")
				debug := c.Bool("debug")
				infile := ""

				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				infile = c.Args().First()

				opt := appbase.BuildOptions(c, config.WaBackend_llvm)
				if err := appllvm.LLVMRun(opt, infile, outfile, target, debug); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				return nil
			},
		},

		{
			Hidden: true,
			Name:   "lex",
			Usage:  "lex Wa source code and print token list",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				err := applex.Lex(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Hidden: true,
			Name:   "ast",
			Usage:  "parse Wa source code and print ast",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				err := appast.PrintAST(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Hidden: true,
			Name:   "ssa",
			Usage:  "print Wa ssa code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				opt := appbase.BuildOptions(c)
				err := appssa.SSARun(opt, c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Hidden: true,
			Name:   "cir",
			Usage:  "print cir code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				opt := appbase.BuildOptions(c)
				err := appcir.PrintCIR(opt, c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},

		{
			Hidden: true,
			Name:   "doc",
			Usage:  "show documentation for package or symbol",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO")
				return nil
			},
		},

		{
			Name:      "yacc",
			Usage:     "generates parsers for LALR(1) grammars",
			ArgsUsage: "<input>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "l",
					Usage: "disable line directives",
				},
				&cli.StringFlag{
					Name:  "o",
					Usage: "set parser output file",
					Value: "y.wa",
				},
				&cli.StringFlag{
					Name:  "p",
					Usage: "name prefix to use in generated code",
					Value: "yy",
				},
				&cli.StringFlag{
					Name:  "v",
					Usage: "create parsing tables",
					Value: "y.output",
				},
				&cli.StringFlag{
					Name:  "c",
					Usage: "set copyright file",
					Value: "",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowSubcommandHelpAndExit(c, 1)
				}
				appyacc.InitFlags(appyacc.Flags{
					Oflag:     c.String("o"),
					Vflag:     c.String("v"),
					Lflag:     c.Bool("l"),
					Prefix:    c.String("p"),
					Copyright: loadCopyright(c.String("c")),
				})
				appyacc.Main(c.Args().First())
				return nil
			},
		},

		{
			Hidden: true,
			Name:   "lsp",
			Usage:  "run Wa langugage server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "log-file",
					Usage: "set log file",
					Value: "",
				},
			},
			Action: func(c *cli.Context) error {
				lsp.NewLSPServer(&lsp.Option{}).Run()
				return nil
			},
		},

		{
			Name:  "logo",
			Usage: "print Wa text format logo",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "more",
					Aliases: []string{"m"},
					Usage:   "print more logos",
				},
			},
			Action: func(c *cli.Context) error {
				applogo.PrintLogo(c.Bool("more"))
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}

func loadCopyright(filename string) string {
	data, _ := os.ReadFile(filename)
	return string(data)
}