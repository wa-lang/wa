// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appgo2wz

import (
	"errors"
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

var CmdGo2wz = &cli.Command{
	Hidden: true,
	Name:   "go2wz",
	Usage:  "convert wago to wa chinese",
	Action: CmdGo2wzAction,
}

func CmdGo2wzAction(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input file")
		os.Exit(1)
	}

	code, err := go2wz(c.Args().First())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(string(code))
	return nil
}

func go2wz(filename string) ([]byte, error) {
	if !appbase.HasExt(filename, ".go") {
		return nil, fmt.Errorf("%q is not Go file", filename)
	}

	if !appbase.PathExists(filename) {
		return nil, fmt.Errorf("%q not found", filename)
	}
	if !appbase.IsNativeFile(filename) {
		return nil, fmt.Errorf("%q must be file", filename)
	}

	src, _ := os.ReadFile(filename)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	if true {
		return nil, errors.New("TODO")
	}

	f.Name.Name = ""
	return format.DevFormat(fset, f, src)
}
