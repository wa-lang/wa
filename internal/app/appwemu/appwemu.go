// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwemu

import (
	"debug/elf"
	"fmt"
	"io"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/abi"
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
func readELF(filename string) (prog *abi.LinkedProgram, err error) {
	prog = new(abi.LinkedProgram)

	f, err := elf.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD || p.Flags&elf.PF_R == 0 {
			continue // 跳过不可读部分
		}

		if p.Flags&elf.PF_X != 0 {
			prog.TextAddr = int64(p.Vaddr)
			prog.TextData = make([]byte, p.Filesz)
			_, err := io.ReadFull(p.Open(), prog.TextData)
			if err != nil {
				return nil, err
			}
		} else {
			prog.DataAddr = int64(p.Vaddr)
			prog.DataData = make([]byte, p.Filesz)
			_, err := io.ReadFull(p.Open(), prog.DataData)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}
