// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import "wa-lang.org/wa/internal/native/abi"

const (
	// 通用寄存器
	REG_RAX abi.RegType = iota + 1 // 0 是无效的编号
	REG_RCX
	REG_RDX
	REG_RBX
	REG_RSP
	REG_RBP
	REG_RSI
	REG_RDI
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15

	// 浮点数寄存器
	REG_XMM0
	REG_XMM1
	REG_XMM2
	REG_XMM3
	REG_XMM4
	REG_XMM5
	REG_XMM6
	REG_XMM7

	// 寄存器编号结束
	REG_END
)
