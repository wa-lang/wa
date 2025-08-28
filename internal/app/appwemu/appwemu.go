// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwemu

import (
	"debug/elf"
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/wemu"
)

var CmdWEmu = &cli.Command{
	Hidden:    true,
	Name:      "wemu",
	Usage:     "debug run elf program with WEnum (riscv64 emulator)",
	ArgsUsage: "<file.elf>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "wemu debug run mode",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		debugRun := c.Bool("debug")

		// 1. 解析 elf 文件
		prog, err := readELF(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. 创建 VM
		vm := wemu.NewWEmu(prog)

		// 3. 执行
		if debugRun {
			if err := vm.DebugRun(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			if err := vm.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		return nil
	},
}

// 读取 elf 文件
func readELF(filename string) (prog *link.Program, err error) {
	f, err := elf.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	prog = new(link.Program)
	for _, sec := range f.Sections {
		switch sec.Name {
		case ".text":
			prog.TextAddr = int64(sec.Addr)
			if prog.TextData, err = sec.Data(); err != nil {
				return nil, err
			}
		case ".data":
			prog.DataAddr = int64(sec.Addr)
			if prog.DataData, err = sec.Data(); err != nil {
				return nil, err
			}
		case ".rodata":
			prog.RoDataAddr = int64(sec.Addr)
			if prog.RoDataData, err = sec.Data(); err != nil {
				return nil, err
			}
		case ".sdata":
			prog.SDataAddr = int64(sec.Addr)
			if prog.SDataData, err = sec.Data(); err != nil {
				return nil, err
			}
		}
	}

	return
}
