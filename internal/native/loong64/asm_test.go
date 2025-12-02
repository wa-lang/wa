// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64_test

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
)

func ExampleAsmSyntax() {
	fmt.Println(loong64.AsmSyntaxEx(
		loong64.AADD_W, "", &abi.AsArgument{
			Rd:  loong64.REG_R1,
			Rs1: loong64.REG_R2,
			Rs2: loong64.REG_R3,
		},
		func(r abi.RegType) string {
			return loong64.RegString(r)
		},
		func(x abi.As, xName string) string {
			return loong64.AsString(x, xName)
		},
	))

	// Output:
	// ADD.W R1, R2, R3
}

func ExampleEncodeLA64() {
	// 数据来源: ../examples/hello-loong64
	//const start_pc = 0x80000000
	instList := []struct {
		as  abi.As
		arg *abi.AsArgument
	}{
		{loong64.ALU12I_W, &abi.AsArgument{Rd: loong64.REG_A0, Imm: 0}},
		{loong64.AADDI_W, &abi.AsArgument{Rd: loong64.REG_A0, Rs1: loong64.REG_A0, Imm: 60}},
		{loong64.ALD_BU, &abi.AsArgument{Rd: loong64.REG_A1, Rs1: loong64.REG_A0, Imm: 0}},
		{loong64.ABEQ, &abi.AsArgument{Rd: loong64.REG_A1, Rs1: loong64.REG_ZERO, Imm: 24}},
		{loong64.ALU12I_W, &abi.AsArgument{Rd: loong64.REG_T0, Imm: 0x10000}},
		{loong64.AADDI_W, &abi.AsArgument{Rd: loong64.REG_T0, Rs1: loong64.REG_T0, Imm: 0}},
		{loong64.AST_B, &abi.AsArgument{Rd: loong64.REG_A1, Rs1: loong64.REG_T0, Imm: 0}},
		{loong64.AADDI_W, &abi.AsArgument{Rd: loong64.REG_A0, Rs1: loong64.REG_A0, Imm: 1}},
		{loong64.AB, &abi.AsArgument{Imm: -24}},
		{loong64.AB, &abi.AsArgument{Imm: -4}},
	}

	for i, inst := range instList {
		x, err := loong64.EncodeLA64(inst.as, inst.arg)
		if err != nil {
			panic(fmt.Sprint(i, err))
		}
		fmt.Printf("0x%08X # %v\n", x, loong64.AsmSyntax(inst.as, "", inst.arg))
	}

	// TODO: Output:
}

func ExampleDecode() {
	var instData = []uint32{}
	for i, x := range instData {
		as, arg, err := loong64.Decode(x)
		if err != nil {
			panic(fmt.Errorf("%d: loong64.Decode(0x%08X): %v", i, x, err))
		}

		x, err := loong64.EncodeLA64(as, arg)
		if err != nil {
			panic(fmt.Sprint(i, err))
		}
		fmt.Printf("0x%08X # %v\n", x, loong64.AsmSyntax(as, "", arg))
	}

	// TODO: Output:
}
