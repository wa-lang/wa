// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appast

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/native/abi"
	native_parser "wa-lang.org/wa/internal/native/parser"
	native_token "wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
	wat_parser "wa-lang.org/wa/internal/wat/parser"
)

var CmdAst = &cli.Command{
	Hidden: true,
	Name:   "ast",
	Usage:  "parse Wa/wat source code and print ast",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "arch",
			Usage: "set target architecture (loong64|riscv32|riscv64)",
			Value: "loong64",
		},
	},
	Action: CmdAstAction,
}

func CmdAstAction(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input file")
		os.Exit(1)
	}

	filename := c.Args().First()

	if appbase.IsNativeFile(filename, ".wat") {
		err := PrintWatAST(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if appbase.HasExt(filename, ".wa.s", ".wz.s") {
		native_parser.DebugMode = c.Bool("debug")

		cpuType := abi.ParseCPUType(c.String("arch"))

		err := PrintNasmAST(filename, cpuType)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	err := PrintAST(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func PrintAST(filename string) error {
	if !appbase.HasExt(filename, ".wa", ".wz") {
		return fmt.Errorf("%q is not Wa file", filename)
	}

	if !appbase.PathExists(filename) {
		return fmt.Errorf("%q not found", filename)
	}
	if !appbase.IsNativeFile(filename) {
		return fmt.Errorf("%q must be file", filename)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}

func PrintWatAST(filename string) error {
	m, err := wat_parser.ParseModule(filename, nil)
	if err != nil {
		return err
	}

	fmt.Println(m)
	return nil
}

func PrintNasmAST(filename string, cpuType abi.CPUType) error {
	fset := native_token.NewFileSet()
	m, err := native_parser.ParseFile(cpuType, fset, filename, nil)
	if err != nil {
		return err
	}

	fmt.Println(m)
	return nil
}
