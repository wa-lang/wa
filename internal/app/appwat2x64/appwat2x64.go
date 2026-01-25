// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2x64

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/wat2x64"
)

var CmdWat2x64 = &cli.Command{
	Hidden:    true,
	Name:      "wat2x64",
	Usage:     "convert a WebAssembly text file to X64 assembly code",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set code output file",
			Value:   "a.out.wa.s",
		},
		&cli.StringFlag{
			Name:  "os",
			Usage: "set os (linux|windows)",
			Value: "linux",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")
		osName := c.String("os")

		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".wa.s"
		}
		if !strings.HasSuffix(outfile, ".wa.s") {
			outfile += ".wa.s"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cpuType := abi.X64Windows
		if strings.EqualFold(osName, config.WaOS_linux) {
			cpuType = abi.X64Unix
		}
		_, code, err := wat2x64.Wat2X64(infile, source, cpuType)
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
