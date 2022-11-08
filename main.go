// 版权 @2019 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

// 凹语言，The Wa Programming Language.
package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/wa-lang/wa/api"
	"github.com/wa-lang/wa/internal/3rdparty/cli"
	"github.com/wa-lang/wa/internal/app"
	"github.com/wa-lang/wa/internal/app/apputil"
	"github.com/wa-lang/wa/internal/backends/target_spec"
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
		return "devel:" + time.Now().Format("2006-01-02+15:04:05")
	}()

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

	// 没有参数时对应 run 命令
	cliApp.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 0)
		}

		ctx := app.NewApp(build_Options(c))
		output, err := ctx.WASM(c.Args().First(), target_spec.Machine_Wasm32_wa)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		outfile := "a.out.wat"
		if !c.Bool("debug") {
			defer os.Remove(outfile)
		}

		if err = os.WriteFile(outfile, []byte(output), 0666); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		stdoutStderr, err := apputil.RunWasm(outfile)
		if err != nil {
			if len(stdoutStderr) > 0 {
				fmt.Println(string(stdoutStderr))
			}
			fmt.Println(err)
			os.Exit(1)
		}
		if len(stdoutStderr) > 0 {
			fmt.Println(string(stdoutStderr))
		}
		return nil
	}

	cliApp.Commands = []*cli.Command{
		{
			// go run main.go debug
			Hidden: true,
			Name:   "debug",
			Usage:  "only for dev/debug",
			Action: func(c *cli.Context) error {
				wat, err := api.BuildFile("hello.wa", "fn main() { println(123) }", "")
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
			Name:      "init",
			Usage:     "init a sketch wa module",
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
			Name:  "run",
			Usage: "compile and run Wa program",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "html",
					Usage: "output html",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := app.NewApp(build_Options(c))
				output, err := ctx.WASM(c.Args().First(), target_spec.Machine_Wasm32_wa)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				outfile := "a.out.wat"
				if !c.Bool("debug") {
					defer os.Remove(outfile)
				}

				if err = os.WriteFile(outfile, []byte(output), 0666); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if c.Bool("html") {
					// todo
				} else {
					stdoutStderr, err := apputil.RunWasm(outfile)
					if err != nil {
						if len(stdoutStderr) > 0 {
							fmt.Println(string(stdoutStderr))
						}
						fmt.Println(err)
						os.Exit(1)
					}
					if len(stdoutStderr) > 0 {
						fmt.Println(string(stdoutStderr))
					}
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
					Value:   "a.out",
				},
				&cli.BoolFlag{
					Name:  "html",
					Usage: "output html",
				},
				&cli.StringFlag{
					Name:  "target",
					Usage: "set target",
				},
			},
			Action: func(c *cli.Context) error {
				outfile := c.String("output")

				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				var target target_spec.Machine
				if s := c.String("target"); s != "" {
					if t, ok := api.ParseMachine(s); ok {
						target = t
					} else {
						fmt.Printf("invalid target: %q", s)
						os.Exit(1)
					}
				}
				if target == "" {
					target = target_spec.Machine_Wasm32_wa
				}

				ctx := app.NewApp(build_Options(c))
				output, err := ctx.WASM(c.Args().First(), target)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if outfile != "" && outfile != "-" {
					if !strings.HasSuffix(outfile, ".wat") {
						outfile += ".wat"
					}
					err := os.WriteFile(outfile, []byte(output), 0666)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				} else {
					fmt.Println(string(output))
				}

				return nil
			},
		},
		{
			Name:  "native",
			Usage: "compile wa source code to native executable",
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
				&cli.BoolFlag{
					Name:  "debug",
					Usage: "dump orginal intermediate representation",
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

				ctx := app.NewApp(build_Options(c))
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
			Hidden: true,
			Name:   "ast",
			Usage:  "parse Wa source code and print ast",
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
			Hidden: true,
			Name:   "ssa",
			Usage:  "print Wa ssa code",
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
			Hidden: true,
			Name:   "cir",
			Usage:  "print cir code",
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
			Name:  "test",
			Usage: "test packages",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO")
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
			Name:  "doc",
			Usage: "show documentation for package or symbol",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO")
				return nil
			},
		},
		{
			Hidden: true,
			Name:   "install-wat2wasm",
			Usage:  "install-wat2wasm tool",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "set output dir",
					Value: "",
				},
			},
			Action: func(c *cli.Context) error {
				outdir := c.String("dir")
				if err := apputil.InstallWat2wasm(outdir); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			},
		},

		{
			Name:  "logo",
			Usage: "print logo",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "more",
					Usage: "print more logos",
				},
			},
			Action: func(c *cli.Context) error {
				app.PrintLogo(c.Bool("more"))
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}

func build_Options(c *cli.Context) *app.Option {
	return &app.Option{
		Debug: c.Bool("debug"),
	}
}
