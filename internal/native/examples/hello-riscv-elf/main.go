// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"flag"
	"os"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/asm"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/riscv"
)

var (
	flagOutput = flag.String("o", "a.out.exe", "set output file")
)

func main() {
	flag.Parse()

	data, err := link.LinkELF(prog)
	if err != nil {
		panic(err)
	}

	os.WriteFile(*flagOutput, data, 0777)
}

var prog = &abi.LinkedProgram{
	TextAddr: 0x80000000,
	TextData: asm.AssemblerRV64(fnBody),
	DataAddr: 0x8000003c,
	DataData: []byte("Hello wa-lang:native/RISC-V!\n\x00"),
}

var fnBody = &ast.FuncBody{
	Insts: []*ast.Instruction{
		{As: riscv.AAUIPC, Arg: &abi.AsArgument{Rd: riscv.REG_A0, Imm: 0}},
		{As: riscv.AADDI, Arg: &abi.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 60}},
		{As: riscv.ALBU, Arg: &abi.AsArgument{Rd: riscv.REG_A1, Rs1: riscv.REG_A0, Imm: 0}},
		{As: riscv.ABEQ, Arg: &abi.AsArgument{Rs1: riscv.REG_A1, Rs2: riscv.REG_ZERO, Imm: 24}},
		{As: riscv.ALUI, Arg: &abi.AsArgument{Rd: riscv.REG_T0, Imm: 0x10000}},
		{As: riscv.AADDI, Arg: &abi.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{As: riscv.ASB, Arg: &abi.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_A1, Imm: 0}},
		{As: riscv.AADDI, Arg: &abi.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 1}},
		{As: riscv.AJAL, Arg: &abi.AsArgument{Rd: riscv.REG_ZERO, Imm: -24}},
		{As: riscv.ALUI, Arg: &abi.AsArgument{Rd: riscv.REG_T0, Imm: 0x100}},
		{As: riscv.AADDI, Arg: &abi.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{As: riscv.ALUI, Arg: &abi.AsArgument{Rd: riscv.REG_T1, Imm: 0x5}},
		{As: riscv.AADDI, Arg: &abi.AsArgument{Rd: riscv.REG_T1, Rs1: riscv.REG_T1, Imm: 0x555}},
		{As: riscv.ASW, Arg: &abi.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_T1, Imm: 0}},
		{As: riscv.AJAL, Arg: &abi.AsArgument{Rd: riscv.REG_ZERO, Imm: 0}},
	},
}
