// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// X64 操作数类型
type X64OperandType int16

const (
	X64Operand_Reg X64OperandType = iota + 1
	X64Operand_Mem
	X64Operand_Imm
)

// X64 指针类型
type X64PtrType int64

// 指针类型
const (
	X64BytePtr X64PtrType = iota + 1
	X64WordPtr
	X64DWordPtr
	X64QWordPtr
)

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
