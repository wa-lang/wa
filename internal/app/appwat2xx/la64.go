// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2xx

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/wat2la"
)

var CmdWat2XX_la64 = &cli.Command{
	Name:      "la64",
	Usage:     "convert a WebAssembly text file to LoongArch",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set code output file",
			Value:   "a.out.was",
		},
		&cli.BoolFlag{
			Name:  "zh",
			Usage: "set zh mode",
		},
	},
	Action: CmdWat2XXAction_la64,
}

func CmdWat2XXAction_la64(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input file")
		os.Exit(1)
	}

	infile := c.Args().First()
	outfile := c.String("output")

	isZhLang := false
	if c.Bool("zh") {
		isZhLang = true
	}

	if outfile == "" {
		outfile = infile
		if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
			if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
				outfile = outfile[:n1-n2]
			}
		}
		outfile += ".was"
	}
	if !strings.HasSuffix(outfile, ".was") {
		outfile += ".was"
	}

	source, err := os.ReadFile(infile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, code, err := wat2la.Wat2LA64(infile, source, "", isZhLang)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 输出汇编格式
	err = os.WriteFile(outfile, code, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
