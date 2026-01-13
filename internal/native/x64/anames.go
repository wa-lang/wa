// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

// 寄存器名字列表
var _Register = []string{
	// 通用寄存器
	REG_RAX: "rax",
	REG_RCX: "rcx",
	REG_RDX: "rdx",
	REG_RBX: "rbx",
	REG_RSP: "rsp",
	REG_RBP: "rbp",
	REG_RSI: "rsi",
	REG_RDI: "rdi",
	REG_R8:  "r8",
	REG_R9:  "r9",
	REG_R10: "r10",
	REG_R11: "r11",
	REG_R12: "r12",
	REG_R13: "r13",
	REG_R14: "r14",
	REG_R15: "r15",

	// 浮点数寄存器
	REG_XMM0: "xmm0",
	REG_XMM1: "xmm1",
	REG_XMM2: "xmm2",
	REG_XMM3: "xmm3",
	REG_XMM4: "xmm4",
	REG_XMM5: "xmm5",
	REG_XMM6: "xmm6",
	REG_XMM7: "xmm7",
}
