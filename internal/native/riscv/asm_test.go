// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv_test

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/riscv"
)

func ExampleAsmSyntax() {
	fmt.Println(riscv.AsmSyntaxEx(
		riscv.AADD, "", &abi.AsArgument{
			Rd:  riscv.REG_X1,
			Rs1: riscv.REG_X2,
			Rs2: riscv.REG_X3,
		},
		func(r abi.RegType) string {
			return riscv.RegString(r)
		},
		func(x abi.As, xName string) string {
			return riscv.AsString(x, xName)
		},
	))

	// Output:
	// ADD X1, X2, X3
}

func ExampleEncodeRV64() {
	// 数据来源: ../examples/hello-riscv
	//const start_pc = 0x80000000
	instList := []struct {
		as  abi.As
		arg *abi.AsArgument
	}{
		{riscv.AAUIPC, &abi.AsArgument{Rd: riscv.REG_A0, Imm: 0}},
		{riscv.AADDI, &abi.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 60}},
		{riscv.ALBU, &abi.AsArgument{Rd: riscv.REG_A1, Rs1: riscv.REG_A0, Imm: 0}},
		{riscv.ABEQ, &abi.AsArgument{Rs1: riscv.REG_A1, Rs2: riscv.REG_ZERO, Imm: 24}},
		{riscv.ALUI, &abi.AsArgument{Rd: riscv.REG_T0, Imm: 0x10000}},
		{riscv.AADDI, &abi.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{riscv.ASB, &abi.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_A1, Imm: 0}},
		{riscv.AADDI, &abi.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 1}},
		{riscv.AJAL, &abi.AsArgument{Rd: riscv.REG_ZERO, Imm: -24}},
		{riscv.ALUI, &abi.AsArgument{Rd: riscv.REG_T0, Imm: 0x100}},
		{riscv.AADDI, &abi.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{riscv.ALUI, &abi.AsArgument{Rd: riscv.REG_T1, Imm: 0x5}},
		{riscv.AADDI, &abi.AsArgument{Rd: riscv.REG_T1, Rs1: riscv.REG_T1, Imm: 0x555}},
		{riscv.ASW, &abi.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_T1, Imm: 0}},
		{riscv.AJAL, &abi.AsArgument{Rd: riscv.REG_ZERO, Imm: 0}},
	}

	for i, inst := range instList {
		x, err := riscv.EncodeRV64(inst.as, inst.arg)
		if err != nil {
			panic(fmt.Sprint(i, err))
		}
		fmt.Printf("0x%08X # %v\n", x, riscv.AsmSyntax(inst.as, "", inst.arg))
	}

	// Output:
	// 0x00000517 # AUIPC A0, 0x0
	// 0x03C50513 # ADDI A0, A0, 60
	// 0x00054583 # LBU A1, 0(A0)
	// 0x00058C63 # BEQ A1, ZERO, 24
	// 0x100002B7 # LUI T0, 0x10000
	// 0x00028293 # ADDI T0, T0, 0
	// 0x00B28023 # SB A1, 0(T0)
	// 0x00150513 # ADDI A0, A0, 1
	// 0xFE9FF06F # JAL ZERO, -24
	// 0x001002B7 # LUI T0, 0x100
	// 0x00028293 # ADDI T0, T0, 0
	// 0x00005337 # LUI T1, 0x5
	// 0x55530313 # ADDI T1, T1, 1365
	// 0x0062A023 # SW T1, 0(T0)
	// 0x0000006F # JAL ZERO, 0
}

func ExampleDecode() {
	var instData = []uint32{
		0x00000517,
		0x03C50513,
		0x00054583,
		0x00058C63,
		0x100002B7,
		0x00028293,
		0x00B28023,
		0x00150513,
		0xFE9FF06F,
		0x001002B7,
		0x00028293,
		0x00005337,
		0x55530313,
		0x0062A023,
		0x0000006F,
	}
	for i, x := range instData {
		as, arg, err := riscv.Decode(x)
		if err != nil {
			panic(fmt.Errorf("%d: riscv.Decode(0x%08X): %v", i, x, err))
		}

		x, err := riscv.EncodeRV64(as, arg)
		if err != nil {
			panic(fmt.Sprint(i, err))
		}
		fmt.Printf("0x%08X # %v\n", x, riscv.AsmSyntax(as, "", arg))
	}

	// Output:
	// 0x00000517 # AUIPC A0, 0x0
	// 0x03C50513 # ADDI A0, A0, 60
	// 0x00054583 # LBU A1, 0(A0)
	// 0x00058C63 # BEQ A1, ZERO, 24
	// 0x100002B7 # LUI T0, 0x10000
	// 0x00028293 # ADDI T0, T0, 0
	// 0x00B28023 # SB A1, 0(T0)
	// 0x00150513 # ADDI A0, A0, 1
	// 0xFE9FF06F # JAL ZERO, -24
	// 0x001002B7 # LUI T0, 0x100
	// 0x00028293 # ADDI T0, T0, 0
	// 0x00005337 # LUI T1, 0x5
	// 0x55530313 # ADDI T1, T1, 1365
	// 0x0062A023 # SW T1, 0(T0)
	// 0x0000006F # JAL ZERO, 0
}
