// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// CPU类型
// 对应不同的汇编指令
type CPUType int

const (
	RISCV64 CPUType = iota + 1
	RISCV32
)

// 寄存器类型
type RegType int16

// 指令类型
type As int16

// 指令参数
type AsArgument struct {
	Rd  RegType // 目标寄存器
	Rs1 RegType // 原寄存器1
	Rs2 RegType // 原寄存器2
	Rs3 RegType // 原寄存器3
	Imm int32   // 立即数

	Symbol    string // 可能是 Label/全局符号, 用于重定位和输出文本
	SymbolVal int64  // 符号对应的地址, 用于和 PC 计算得到相对的地址, 用于 Imm 值
	PC        int64  // 涉及label跳转或者函数调用, 需要当前pc才能编码
}

// 指令原生参数
type AsRawArgument struct {
	Rd  uint32 // 目标寄存器
	Rs1 uint32 // 原寄存器1
	Rs2 uint32 // 原寄存器2
	Rs3 uint32 // 原寄存器3
	Imm int32  // 立即数
}
