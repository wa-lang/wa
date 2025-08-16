// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv64

// Riscv64虚拟机
type Riscv64VM struct {
	Mem  Memory      // 内存
	Sys  Syscall     // 系统调用
	RegX [32]uint64  // 整数寄存器
	RegF [32]float64 // 浮点数寄存器
	PC   uint64      // PC指针
}

// 构造新的虚拟机
func New(mem Memory, sys Syscall) *Riscv64VM {
	return &Riscv64VM{Mem: mem, Sys: sys}
}

// 单步执行
func (p *Riscv64VM) StepRun() error {
	panic("TODO")
}
