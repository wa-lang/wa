// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

import (
	"fmt"
)

// X64 操作数类型
type X64OperandType int16

const (
	X64Operand_Reg X64OperandType = iota + 1
	X64Operand_Mem
	X64Operand_Imm
)

var _X64OperandType_strings = []string{
	X64Operand_Reg: "reg",
	X64Operand_Mem: "mem",
	X64Operand_Imm: "imm",
}

func (v X64OperandType) String() string {
	if v >= 0 && int(v) < len(_X64OperandType_strings) {
		if s := _X64OperandType_strings[v]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("abi.X64OperandType(%d)", int(v))
}

// X64 指针类型
type X64PtrType int64

// 指针类型
const (
	X64BytePtr X64PtrType = iota + 1
	X64WordPtr
	X64DWordPtr
	X64QWordPtr
)

var _X64PtrType_strings = []string{
	X64BytePtr:  "byte ptr",
	X64WordPtr:  "word ptr",
	X64DWordPtr: "dword ptr",
	X64QWordPtr: "qword ptr",
}

func (v X64PtrType) String() string {
	if v >= 0 && int(v) < len(_X64PtrType_strings) {
		if s := _X64PtrType_strings[v]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("abi.X64PtrType(%d)", int(v))
}

// X64 操作数
type X64Operand struct {
	Kind    X64OperandType
	Reg     RegType
	RegName string
	PtrTyp  X64PtrType
	Offset  int64
	Symbol  string
	Imm     int64
}

// X64 指令参数
type X64Argument struct {
	// 无参: Dst == nil && Src == nil (例如 RET)
	// 一元: Dst != nil && Src == nil (例如 PUSH, CALL, INC)
	// 二元: Dst != nil && Src != nil (例如 MOV, ADD, XOR)

	Dst *X64Operand
	Src *X64Operand
}
