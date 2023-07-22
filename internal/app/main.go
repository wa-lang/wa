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
	"wa-lang.org/wa/internal/app/yacc"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/lsp"
	"wa-lang.org/wa/internal/wabt"
	"wa-lang.org/wa/internal/wazero"
	"wa-lang.org/wa/waroot"
)

func Main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Wa"
	cliApp.Usage = "Wa is a tool for managing Wa source code."
	cliApp.Version = waroot.GetVersion()

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
			Usage: "the Wa playground",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "http",
					Value: ":2023",
					Usage: "set http address",
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
			},
			Action: func(c *cli.Context) error {
				waApp := NewApp(build_Options(c))
				err := waApp.Playground(c.String("http"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},

		{
			Name:      "init",
			Usage:     "init a sketch Wa module",
			ArgsUsage: "app-name",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "set app name",
					Value:   "_examples/hello",
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
				waApp := NewApp(build_Options(c))
				err := waApp.InitApp(c.String("name"), c.String("pkgpath"), c.Bool("update"))
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
					Value:   "a.out.wasm",
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
				outfile := c.String("output")

				var input string
				if c.NArg() > 0 {
					input = c.Args().First()
				} else {
					input, _ = os.Getwd()
				}

				ctx := NewApp(build_Options(c))
				watOutput, err := ctx.WASM(input)
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
				cliRun(c)
				return nil
			},
		},
		{
			Name:  "fmt",
			Usage: "format Wa source code file",
			Action: func(c *cli.Context) error {
				waApp := NewApp(new(Option))
				err := waApp.Fmt(c.Args().First())
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
				waApp := NewApp(build_Options(c))
				waApp.RunTest(c.Args().First(), appArgs...)
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

				ctx := NewApp(build_Options(c, config.WaBackend_llvm))
				if err := ctx.LLVM(infile, outfile, target, debug); err != nil {
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

				waApp := NewApp(build_Options(c))
				err := waApp.Lex(c.Args().First())
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

				waApp := NewApp(build_Options(c))
				err := waApp.AST(c.Args().First())
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

				ctx := NewApp(build_Options(c))
				err := ctx.SSA(c.Args().First())
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

				ctx := NewApp(build_Options(c))
				err := ctx.CIR(c.Args().First())
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
				yacc.InitFlags(yacc.Flags{
					Oflag:     c.String("o"),
					Vflag:     c.String("v"),
					Lflag:     c.Bool("l"),
					Prefix:    c.String("p"),
					Copyright: loadCopyright(c.String("c")),
				})
				yacc.Main(c.Args().First())
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
					Name:  "more",
					Usage: "print more logos",
				},
			},
			Action: func(c *cli.Context) error {
				PrintLogo(c.Bool("more"))
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}

func build_Options(c *cli.Context, waBackend ...string) *Option {
	opt := &Option{
		Debug:        c.Bool("debug"),
		WaBackend:    config.WaBackend_Default,
		BuilgTags:    strings.Fields(c.String("tags")),
		Clang:        c.String("clang"),
		Llc:          c.String("llc"),
		LD_StackSize: c.Int("ld-stack-size"),
		LD_MaxMemory: c.Int("ld-max-memory"),
	}

	opt.TargetArch = "wasm"
	if len(waBackend) > 0 {
		opt.WaBackend = waBackend[0]
	}

	if target := c.String("target"); !config.CheckWaOS(target) {
		fmt.Printf("unknown target: %s\n", c.String("target"))
		os.Exit(1)
	}

	switch c.String("target") {
	case "", "wa", "walang":
		opt.TargetOS = config.WaOS_Default
	case config.WaOS_wasi:
		opt.TargetOS = config.WaOS_wasi
	case config.WaOS_arduino:
		opt.TargetOS = config.WaOS_arduino
	case config.WaOS_chrome:
		opt.TargetOS = config.WaOS_chrome
	case config.WaOS_mvp:
		opt.TargetOS = config.WaOS_mvp
	default:
		fmt.Printf("unreachable: target: %s\n", c.String("target"))
		os.Exit(1)
	}
	return opt
}

func cliRun(c *cli.Context) {
	var infile string
	if c.NArg() > 0 {
		infile = c.Args().First()
	} else {
		infile, _ = os.Getwd()
	}

	var app = NewApp(build_Options(c))
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
		watBytes, err = app.WASM(infile)
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

	stdout, stderr, err := wazero.RunWasm(app.GetConfig(), infile, wasmBytes, appArgs...)
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

func loadCopyright(filename string) string {
	data, _ := os.ReadFile(filename)
	return string(data)
}
