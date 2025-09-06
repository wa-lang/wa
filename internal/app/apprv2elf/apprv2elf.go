// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package apprv2elf

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
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
		_ = source

		// 2. 处理汇编中的符号和跳转的地址
		// 3. 生成对应cpu的机器码
		// 4. 以elf格式保持

		var elfBytes = []byte("todo")

		err = os.WriteFile(outfile, elfBytes, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if true {
			fmt.Fprintln(os.Stderr, "todo")
			os.Exit(1)
		}

		return nil
	},
}
