// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2rv

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/printer"
	"wa-lang.org/wa/internal/native/wat2rv"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
)

var CmdWat2rv = &cli.Command{
	Hidden:    true,
	Name:      "wat2rv",
	Usage:     "convert a WebAssembly text file to RISC-V",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set code output file",
			Value:   "a.out.txt",
		},
		&cli.StringFlag{
			Name:  "arch",
			Usage: "set target architecture (riscv32|riscv64)",
			Value: "riscv64",
		},
		&cli.Int64Flag{
			Name:  "Ttext",
			Usage: "set address of .text segment",
			Value: dram.DRAM_BASE,
		},
		&cli.Int64Flag{
			Name:  "Tdata",
			Usage: "set address of .data segment",
			Value: 0,
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")

		opt := wat2rv.Options{}
		opt.Ttext = uint64(c.Int64("Ttext"))
		opt.Tdata = uint64(c.Int64("Tdata"))

		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".c"
		}
		if !strings.HasSuffix(outfile, ".c") {
			outfile += ".c"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var f *ast.File
		switch arch := c.String("arch"); true {
		case strings.EqualFold(arch, "riscv64"):
			_, f, err = wat2rv.Wat2rv64(infile, source, opt)
		case strings.EqualFold(arch, "riscv32"):
			_, f, err = wat2rv.Wat2rv32(infile, source, opt)
		default:
			fmt.Fprintln(os.Stderr, "unknown arch: "+arch)
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 输出汇编格式
		printer.Fprint(os.Stdout, int64(opt.Ttext), f)
		return nil
	},
}
