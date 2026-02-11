// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// X64 指令类型
type X64Mode int

const (
	X64Mode_NoArgs X64Mode = iota + 1 // 无参数: ret, nop, cdq
	X64Mode_Unary                     // 一元: push, pop, jmp, call, setcc
	X64Mode_Binary                    // 二元: mov, add, sub, lea, movzx
)

// X64 操作数类型
type X64OperandType int16

const (
	X64X64Operand_Reg X64OperandType = iota + 1
	X64X64Operand_Mem
	X64X64Operand_Imm
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
