// 版权 @2019 凹语言 作者。保留所有权利。

// 凹语言™ 命令行程序。
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"

	"github.com/wa-lang/wa/internal/3rdparty/cli"
	"github.com/wa-lang/wa/internal/app"
	"github.com/wa-lang/wa/internal/config"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Wa"
	cliApp.Usage = "Wa is a tool for managing Wa source code."
	cliApp.Version = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" {
				return info.Main.Version
			}
		}
		return "(devel)"
	}()

	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{Name: "os", Usage: "set target OS", Value: runtime.GOOS},
		&cli.StringFlag{Name: "arch", Usage: "set target Arch", Value: runtime.GOARCH},
		&cli.StringFlag{Name: "backend", Usage: "set backend code generator"},
		&cli.StringFlag{Name: "clang", Usage: "set clang"},
		&cli.StringFlag{Name: "wasm-llc", Usage: "set wasm-llc"},
		&cli.StringFlag{Name: "wasm-ld", Usage: "set wasm-ld"},
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "set debug mode"},
		&cli.StringFlag{Name: "trace", Aliases: []string{"t"}, Usage: "set trace mode (*|app|compiler|loader)"},
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

	// 没有参数时对应 run 命令
	cliApp.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		waApp := app.NewApp(build_Options(c))
		data, err := waApp.Run(c.Args().First(), nil)
		if len(data) != 0 {
			fmt.Print(string(data))
		}
		if errx, ok := err.(*exec.ExitError); ok {
			os.Exit(errx.ExitCode())
		}
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	cliApp.Commands = []*cli.Command{
		{
			Name:      "init",
			Usage:     "init a sketch app",
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
				waApp := app.NewApp(build_Options(c))
				err := waApp.InitApp(c.String("name"), c.String("pkgpath"), c.Bool("update"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "fmt",
			Usage: "format Wa package sources",
			Action: func(c *cli.Context) error {
				waApp := app.NewApp(build_Options(c))
				err := waApp.Fmt(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "compile and run Wa program",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				waApp := app.NewApp(build_Options(c))
				data, err := waApp.Run(c.Args().First(), nil)
				if len(data) != 0 {
					fmt.Print(string(data))
				}
				if errx, ok := err.(*exec.ExitError); ok {
					os.Exit(errx.ExitCode())
				}
				if err != nil {
					fmt.Println(err)
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
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				opt := build_Options(c)

				outfile := c.String("output")
				if outfile == "" {
					outfile = "a.out"
					if opt.TargetOS == "windows" {
						outfile += ".exe"
					}
				}

				waApp := app.NewApp(opt)
				data, err := waApp.Build(c.Args().First(), nil, outfile)
				if len(data) != 0 {
					fmt.Print(string(data))
				}
				if errx, ok := err.(*exec.ExitError); ok {
					os.Exit(errx.ExitCode())
				}
				if err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
		{
			Name:  "lex",
			Usage: "lex Wa source code and print token list",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				waApp := app.NewApp(build_Options(c))
				err := waApp.Lex(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "ast",
			Usage: "parse Wa source code and print ast",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				waApp := app.NewApp(build_Options(c))
				err := waApp.AST(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "ssa",
			Usage: "print Wa ssa code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := app.NewApp(build_Options(c))
				err := ctx.SSA(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "cir",
			Usage: "print cir code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := app.NewApp(build_Options(c))
				err := ctx.CIR(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "asm",
			Usage: "parse Wa and print ouput assembly code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := app.NewApp(build_Options(c))
				err := ctx.ASM(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "wasm",
			Usage: "parse Wa and print ouput WebAssembly text",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := app.NewApp(build_Options(c))
				err := ctx.WASM(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "test",
			Usage: "test packages",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO")
				return nil
			},
		},
		{
			Name:  "doc",
			Usage: "show documentation for package or symbol",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO")
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}

func build_Options(c *cli.Context) *app.Option {
	return &app.Option{
		Debug:      c.Bool("debug"),
		TargetOS:   c.String("os"),
		TargetArch: c.String("arch"),
		Clang:      c.String("clang"),
		WasmLLC:    c.String("wasm-llc"),
		WasmLD:     c.String("wasm-ld"),
	}
}
