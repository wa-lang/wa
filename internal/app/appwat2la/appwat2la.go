// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2la

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/wat2la"
)

var CmdWat2la = &cli.Command{
	Hidden:    true,
	Name:      "wat2la",
	Usage:     "convert a WebAssembly text file to LoongArch",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set code output file",
			Value:   "a.out.ws",
		},
		&cli.StringFlag{
			Name:    "prefix",
			Aliases: []string{"p"},
			Usage:   "name prefix to use in generated code",
			Value:   "app",
		},
		&cli.StringSliceFlag{
			Name:  "exports",
			Usage: "set export func list (K1=V1,K2=V2,...)",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")
		prefix := c.String("prefix")
		exports := c.StringSliceAsMap("exports")

		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".ws"
		}
		if !strings.HasSuffix(outfile, ".ws") {
			outfile += ".ws"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, code, err := wat2la.Wat2LA64(infile, source, wat2la.Options{
			Prefix:  prefix,
			Exports: exports,
		})
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
	},
}
