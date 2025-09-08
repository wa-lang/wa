// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package apprv2elf

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/asm"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
)

var CmdRv2elf = &cli.Command{
	Hidden:    true,
	Name:      "rv2elf",
	Usage:     "convert riscv assembly code to elf binary format",
	ArgsUsage: "<file.s>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
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

		opt := &abi.LinkOptions{}
		opt.Ttext = uint64(c.Int64("Ttext"))
		opt.Tdata = uint64(c.Int64("Tdata"))

		parser.DebugMode = c.Bool("debug")

		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".s"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".s") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".elf.exe"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 1. 解析汇编程序
		var fset = token.NewFileSet()
		var f *ast.File
		switch arch := c.String("arch"); arch {
		case "riscv32":
			f, err = parser.ParseFile(abi.RISCV32, fset, infile, source)
		case "riscv64":
			f, err = parser.ParseFile(abi.RISCV64, fset, infile, source)
		default:
			fmt.Printf("unknown arch: %s\n", arch)
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 生成对应cpu的机器码
		prog, err := asm.AssembleFile(fset, f, opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 生成ELF文件
		elfBytes, err := link.LinkELF(prog)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := os.WriteFile(outfile, elfBytes, 0777); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
