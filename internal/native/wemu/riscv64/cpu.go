// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv64

import (
	"wa-lang.org/wa/internal/native/riscv"
)

// 处理器
type CPU struct {
	Bus  *Bus         // 外设总线
	Mode riscv.Mode   // 工作模式
	RegX [32]uint64   // 整数寄存器
	RegF [32]float64  // 浮点数寄存器
	CSR  [4096]uint64 // 状态寄存器
	PC   uint64       // PC指针
}

// 构建模拟器并初始化PC和SP
func NewCPU(bus *Bus, pc, sp uint64) *CPU {
	p := &CPU{Bus: bus}
	p.RegX[riscv.REG_SP] = sp
	p.PC = pc
	return p
}
