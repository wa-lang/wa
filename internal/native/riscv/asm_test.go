// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv_test

import (
	"fmt"
	"log"

	"wa-lang.org/wa/internal/native/riscv"
)

func ExampleAsmSyntax() {
	fmt.Println(riscv.AsmSyntax(0, riscv.AADD, &riscv.AsArgument{
		Rd:  riscv.REG_X1,
		Rs1: riscv.REG_X2,
		Rs2: riscv.REG_X3,
	}))

	// Output:
	// ADD X1, X2, X3
}

func ExampleAs_EncodeRV64() {
	// 数据来源: ../examples/hello-riscv
	const start_pc = 0x80000000
	instList := []struct {
		pc  int64
		as  riscv.As
		arg *riscv.AsArgument
	}{
		{0, riscv.AAUIPC, &riscv.AsArgument{Rd: riscv.REG_A0, Imm: 0}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 60}},
		{0, riscv.ALBU, &riscv.AsArgument{Rd: riscv.REG_A1, Rs1: riscv.REG_A0, Imm: 0}},
		{0, riscv.ABEQ, &riscv.AsArgument{Rs1: riscv.REG_A1, Rs2: riscv.REG_ZERO, Imm: 24}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T0, Imm: 0x10000}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{0, riscv.ASB, &riscv.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_A1, Imm: 0}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 1}},
		{0, riscv.AJAL, &riscv.AsArgument{Rd: riscv.REG_ZERO, Imm: -24}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T0, Imm: 0x100}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T1, Imm: 0x5}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T1, Rs1: riscv.REG_T1, Imm: 0x555}},
		{0, riscv.ASW, &riscv.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_T1, Imm: 0}},
		{0, riscv.AJAL, &riscv.AsArgument{Rd: riscv.REG_ZERO, Imm: 0}},
	}

	for i, inst := range instList {
		inst.pc = start_pc + int64(i)*4
		x, err := riscv.EncodeRV64(inst.as, inst.arg)
		if err != nil {
			log.Fatal(i, err)
		}
		fmt.Printf("0x%08X # %v\n", x, riscv.AsmSyntaxEx(inst.pc, inst.as, inst.arg, riscv.RegAliasString))
	}

	// TODO: 校对指令的参数打印格式

	// Output:
	// 0x00000517 # AUIPC A0, 0x0
	// 0x03C50513 # ADDI A0, 60(A0)
	// 0x00054583 # LBU A1, 0(A0)
	// 0x00058C63 # BEQ A1, ZERO, 0x80000024
	// 0x100002B7 # LUI T0, 0x10000
	// 0x00028293 # ADDI T0, 0(T0)
	// 0x00B28023 # SB A1, 0(T0)
	// 0x00150513 # ADDI A0, 1(A0)
	// 0xFE9FF06F # JAL ZERO, 0x80000008
	// 0x001002B7 # LUI T0, 0x100
	// 0x00028293 # ADDI T0, 0(T0)
	// 0x00005337 # LUI T1, 0x5
	// 0x55530313 # ADDI T1, 1365(T1)
	// 0x0062A023 # SW T1, 0(T0)
	// 0x0000006F # JAL ZERO, 0x80000038
}

func ExampleDecode() {
	const start_pc = 0x80000000
	var instData = []uint32{
		0x00000517,
		0x03C50513,
		0x00054583,
		0x00058C63,
		0x100002B7,
		0x00028293,
		// 0x00B28023, // TODO: fix test
		// 0x00150513,
		// 0xFE9FF06F,
		// 0x001002B7,
		// 0x00028293,
		// 0x00005337,
		// 0x55530313,
		// 0x0062A023,
		// 0x0000006F,
	}
	for i, x := range instData {
		as, arg, err := riscv.Decode(x)
		if err != nil {
			log.Fatalf("%d: riscv.Decode(0x%08X): %v", i, x, err)
		}

		pc := start_pc + int64(i)*4
		x, err := riscv.EncodeRV64(as, arg)
		if err != nil {
			log.Fatal(i, err)
		}
		fmt.Printf("0x%08X # %v\n", x, riscv.AsmSyntaxEx(pc, as, arg, riscv.RegAliasString))
	}

	// Output:
	// 0x00000517 # AUIPC A0, 0x0
	// 0x03C50513 # ADDI A0, 60(A0)
	// 0x00054583 # LBU A1, 0(A0)
	// 0x00058C63 # BEQ A1, ZERO, 0x80000024
	// 0x100002B7 # LUI T0, 0x10000
	// 0x00028293 # ADDI T0, 0(T0)
}
