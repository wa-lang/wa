// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// 寄存器类型
type RegType int16

// 指令类型
type As int16

// 指令参数
type AsArgument struct {
	Rd      RegType // 目标寄存器
	Rs1     RegType // 原寄存器1
	Rs2     RegType // 原寄存器2
	Rs3     RegType // 原寄存器3
	Imm     int32   // 立即数
	ImmName string  // 立即数名字, 可能是 Label/符号, 用于打印
}
