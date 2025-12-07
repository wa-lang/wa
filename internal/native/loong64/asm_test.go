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

func ExampleEncodeLA64_xALui12iwAndOri() {
	// $message = 0x80000028
	//
	// # a0 = 字符串地址
	// lu12i.w a0, %hi($message)     # 高20位
	// ori     a0, a0, %lo($message) # 低12位(不能使用addi.w)

	// 0b 10000000000000000000000000101000
	// 0b 10000000000000000000
	// 0b 00000101000
	// as 0001010
	// => 0001010 10000000000000000000
	// 0b 0001010 10000000000000000000 00100

	message := uint32(0x80000028)
	{
		// si20 直接用整数表示, 超出了范围, 必须要转换为负数吗？
		// 那样的话对于编写者来说, 这是一个很大的负担
		as, arg := loong64.ALU12I_W, &abi.AsArgument{Rd: loong64.REG_A0, Imm: int32(message >> 12)}

		x, err := loong64.EncodeLA64(as, arg)
		if err != nil {
			panic(err)
		}

		fmt.Printf("0x%08X # %v\n", x, loong64.AsmSyntax(as, "", arg))
		fmt.Printf("0b%032b # %v\n", x, loong64.AsmSyntax(as, "", arg))
		fmt.Printf("0b%032b # %s\n", message, "message")
	}
	{
		as, arg := loong64.AORI, &abi.AsArgument{Rd: loong64.REG_A0, Rs1: loong64.REG_A0, Imm: int32(message & ((1 << 12) - 1))}

		x, err := loong64.EncodeLA64(as, arg)
		if err != nil {
			panic(err)
		}

		fmt.Printf("0x%08X # %v\n", x, loong64.AsmSyntax(as, "", arg))
		fmt.Printf("0b%032b # %v\n", x, loong64.AsmSyntax(as, "", arg))
		fmt.Printf("0b%032b # %s\n", message, "message")
	}

	// Output:
	// 0x15000004 # LU12I.W A0, 524288
	// 0b00010101000000000000000000000100 # LU12I.W A0, 524288
	// 0b10000000000000000000000000101000 # message
	// 0x0380A084 # ORI A0, A0, 40
	// 0b00000011100000001010000010000100 # ORI A0, A0, 40
	// 0b10000000000000000000000000101000 # message
}

func ExampleDecode_xALui12iwAndOri() {
	x := uint32(0x15000004)
	as, arg, err := loong64.Decode(x)
	if err != nil {
		panic(err)
	}
	fmt.Println(loong64.AsmSyntax(as, "", arg))

	// TODO: Output:
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

	// Output:
	// 0x14000004 # LU12I.W A0, 0
	// 0x0280F084 # ADDI.W A0, A0, 60
	// 0x2A000085 # LD.BU A1, A0, 0
	// 0x58006005 # BEQ ZERO, A1, 24
	// 0x1420000C # LU12I.W T0, 65536
	// 0x0280018C # ADDI.W T0, T0, 0
	// 0x29000185 # ST.B A1, T0, 0
	// 0x02800484 # ADDI.W A0, A0, 1
	// 0x53FFA3FF # B -24
	// 0x53FFF3FF # B -4
}

func ExampleDecode() {
	var instData = []uint32{
		0x14000004, // # LU12I.W A0, 0
		0x0280F084, // # ADDI.W A0, A0, 60
		0x2A000085, // # LD.BU A1, A0, 0
		0x58006005, // # BEQ ZERO, A1, 24
		0x1420000C, // # LU12I.W T0, 65536
		0x0280018C, // # ADDI.W T0, T0, 0
		0x29000185, // # ST.B A1, T0, 0
		0x02800484, // # ADDI.W A0, A0, 1
		0x53FFA3FF, // # B -24
		0x53FFF3FF, // # B -4
	}
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

	// Output:
	// 0x14000004 # LU12I.W A0, 0
	// 0x0280F084 # ADDI.W A0, A0, 60
	// 0x2A000085 # LD.BU A1, A0, 0
	// 0x58006005 # BEQ ZERO, A1, 24
	// 0x1420000C # LU12I.W T0, 65536
	// 0x0280018C # ADDI.W T0, T0, 0
	// 0x29000185 # ST.B A1, T0, 0
	// 0x02800484 # ADDI.W A0, A0, 1
	// 0x53FFA3FF # B -24
	// 0x53FFF3FF # B -4
}
