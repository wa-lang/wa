// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2wasm

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/wat/watutil"
)

var CmdWat2wasm = &cli.Command{
	Hidden:    true,
	Name:      "wat2wasm",
	Usage:     "convert wat format to wasm binary format",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")
		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".wasm"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wasmBytes, err := watutil.Wat2Wasm(infile, source)
		if err != nil {
			os.WriteFile(outfile, wasmBytes, 0666)
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.WriteFile(outfile, wasmBytes, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
