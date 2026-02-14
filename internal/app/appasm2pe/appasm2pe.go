// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appasm2pe

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/asm"
	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
)

var CmdAsm2pe = &cli.Command{
	Hidden:    true,
	Name:      "asm2pe",
	Usage:     "convert wa native assembly code to pe binary format",
	ArgsUsage: "<file.s>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
		},
		&cli.Int64Flag{
			Name:  "DRAM-base",
			Usage: "set DRAM address",
			Value: 0,
		},
		&cli.Int64Flag{
			Name:  "DRAM-size",
			Usage: "set DRAM size",
			Value: dram.DRAM_SIZE,
		},
		&cli.StringFlag{
			Name:  "entry",
			Usage: "set entry func name",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		if c.Bool("debug") {
			parser.DebugMode = true
		}

		infile := c.Args().First()
		outfile := c.String("output")

		opt := &abi.LinkOptions{}
		opt.DRAMBase = c.Int64("DRAM-base")
		opt.DRAMSize = c.Int64("DRAM-size")
		opt.CPU = abi.X64Windows
		if opt.DRAMBase == 0 {
			opt.DRAMBase = dram.DRAM_BASE_X64_WINDOWS
		}

		if s := c.String("entry"); s != "" {
			opt.EntryFunc = s
		}

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

		// 解析汇编程序, 并生成对应cpu的机器码
		prog, err := asm.AssembleFile(infile, source, opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 保持到PE格式文件
		exeBytes, err := link.LinkEXE(prog)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := os.WriteFile(outfile, exeBytes, 0777); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
