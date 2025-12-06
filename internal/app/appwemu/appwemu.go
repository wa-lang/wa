// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwemu

import (
	"fmt"
	"io"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/link/elf"
	"wa-lang.org/wa/internal/native/wemu"
	"wa-lang.org/wa/internal/native/wemu/device/uart"
)

var CmdWEmu = &cli.Command{
	Hidden:    true,
	Name:      "wemu",
	Usage:     "debug run elf program with WEnum (loong64|riscv emulator)",
	ArgsUsage: "<file.elf>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "wemu debug run mode",
		},
		&cli.StringFlag{
			Name:  "uart",
			Usage: "set uart type(qemu|esp32c3)",
			Value: "qemu",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		debugRun := c.Bool("debug")
		opt := new(wemu.Option)

		// 串口映射的类型
		switch s := c.String("uart"); s {
		case "esp32c3":
			opt.UARTBase = uart.UART_BASE_ESP32_C3
		case "qemu":
			opt.UARTBase = uart.UART_BASE_QEMU
		default:
			fmt.Fprintf(os.Stderr, "invalid uart type: %q\n", s)
			os.Exit(1)
		}

		// 1. 解析 elf 文件
		prog, err := readELF(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. 创建 VM
		vm := wemu.NewWEmu(prog, opt)

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

	f, err := elfOpen(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 判断CPU类型
	switch elf.Machine(f.Machine) {
	case elf.EM_RISCV:
		// 判断机器指针长度
		switch elf.Class(f.Class) {
		case elf.ELFCLASS32:
			prog.CPU = abi.RISCV32
		case elf.ELFCLASS64:
			prog.CPU = abi.RISCV64
		default:
			return nil, fmt.Errorf("unsupported ELF class: %v", f.Class)
		}
	case elf.EM_LOONGARCH:
		// 判断机器指针长度
		switch elf.Class(f.Class) {
		case elf.ELFCLASS64:
			prog.CPU = abi.LOONG64
		default:
			return nil, fmt.Errorf("unsupported ELF class: %v", f.Class)
		}
	default:
		return nil, fmt.Errorf("wemu: donot support machine %v", f.Machine)
	}

	// 读取数据段和指令段
	for _, p := range f.Progs {
		if elf.ProgType(p.Type) != elf.PT_LOAD || elf.ProgFlag(p.Flags)&elf.PF_R == 0 {
			continue // 跳过不可读部分
		}

		if elf.ProgFlag(p.Flags)&elf.PF_X != 0 {
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
