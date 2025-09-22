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

// 链接参数
type LinkOptions struct {
	CPU      CPUType // CPU类型
	DRAMBase int64   // 代码段开始地址
	DRAMSize int64   // 数据段开始地址
}

// 链接后的符号
type LinkedSymbol struct {
	Name string // 名字
	Addr int64  // 内存地址
	Size int64  // 内存大小
	Data []byte // 内存数据
}

// 程链接后的程序
type LinkedProgram struct {
	CPU CPUType // CPU类型

	TextAddr int64  // 程序段地址
	TextData []byte // 程序段数据
	DataAddr int64  // 数据段地址
	DataData []byte // 数据段数据
}

// 寄存器类型
type RegType int16

// 指令类型
type As int16

// 内置的重定位函数
type BuiltinFn int16

const (
	BuiltinFn_HI       = iota + 1 // %hi(symbol) # 绝对地址 HI20, 指令 lui
	BuiltinFn_LO                  // %lo(symbol) # 绝对地址 LO12, 指令 load/store/add
	BuiltinFn_PCREL_HI            // %pcrel_hi(symbol) # PC相对地址 HI20, auipc
	BuiltinFn_PCREL_LO            // %pcrel_lo(label)  # label 对应的指令中, 计算出的PC相对地址的 LO12 部分, 参数必须是 label

	BuiltinFn_HI_zh
	BuiltinFn_LO_zh
	BuiltinFn_PCREL_HI_zh
	BuiltinFn_PCREL_LO_zh
)

// 指令参数
type AsArgument struct {
	Rd  RegType // 目标寄存器
	Rs1 RegType // 原寄存器1
	Rs2 RegType // 原寄存器2
	Rs3 RegType // 原寄存器3
	Imm int32   // 立即数

	Symbol      string    // 可能是 Label/全局符号, 用于重定位和输出文本
	SymbolDecor BuiltinFn // 符号的修饰函数, 可能要重新计算
}

// 指令原生参数
type AsRawArgument struct {
	Rd  uint32 // 目标寄存器
	Rs1 uint32 // 原寄存器1
	Rs2 uint32 // 原寄存器2
	Rs3 uint32 // 原寄存器3
	Imm int32  // 立即数
}
