// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

// 寄存器类型
type RegType int16

const (
	// 通用寄存器
	REG_X0 RegType = iota
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15 // RV32M 嵌入式版本只有 16 个寄存器
	REG_X16
	REG_X17
	REG_X18
	REG_X19
	REG_X20
	REG_X21
	REG_X22
	REG_X23
	REG_X24
	REG_X25
	REG_X26
	REG_X27
	REG_X28
	REG_X29
	REG_X30
	REG_X31

	// 浮点数寄存器(F/D扩展)
	REG_F0
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_F8
	REG_F9
	REG_F10
	REG_F11
	REG_F12
	REG_F13
	REG_F14
	REG_F15
	REG_F16
	REG_F17
	REG_F18
	REG_F19
	REG_F20
	REG_F21
	REG_F22
	REG_F23
	REG_F24
	REG_F25
	REG_F26
	REG_F27
	REG_F28
	REG_F29
	REG_F30
	REG_F31

	// 寄存器编号结束
	REG_END

	// 寄存器的 ABI 使用约定
	REG_ZERO = REG_X0  // 零寄存器
	REG_RA   = REG_X1  // 返回地址
	REG_SP   = REG_X2  // 栈指针
	REG_GP   = REG_X3  // 全局基指针
	REG_TP   = REG_X4  // 线程指针
	REG_T0   = REG_X5  // 临时变量
	REG_T1   = REG_X6  // 临时变量
	REG_T2   = REG_X7  // 临时变量
	REG_S0   = REG_X8  // Saved register, 帧指针
	REG_S1   = REG_X9  // Saved register
	REG_A0   = REG_X10 // 函数参数/返回值
	REG_A1   = REG_X11 // 函数参数/返回值
	REG_A2   = REG_X12 // 函数参数
	REG_A3   = REG_X13 // 函数参数
	REG_A4   = REG_X14 // 函数参数
	REG_A5   = REG_X15 // 函数参数
	REG_A6   = REG_X16 // 函数参数
	REG_A7   = REG_X17 // 函数参数
	REG_S2   = REG_X18 // Saved register
	REG_S3   = REG_X19 // Saved register
	REG_S4   = REG_X20 // Saved register
	REG_S5   = REG_X21 // Saved register
	REG_S6   = REG_X22 // Saved register
	REG_S7   = REG_X23 // Saved register
	REG_S8   = REG_X24 // Saved register
	REG_S9   = REG_X25 // Saved register
	REG_S10  = REG_X26 // Saved register
	REG_S11  = REG_X27 // Saved register
	REG_T3   = REG_X28 // 临时变量
	REG_T4   = REG_X29 // 临时变量
	REG_T5   = REG_X30 // 临时变量
	REG_T6   = REG_X31 // 临时变量

	// 浮点数寄存器使用约定
	REG_FT0  = REG_F0  // 临时变量
	REG_FT1  = REG_F1  // 临时变量
	REG_FT2  = REG_F2  // 临时变量
	REG_FT3  = REG_F3  // 临时变量
	REG_FT4  = REG_F4  // 临时变量
	REG_FT5  = REG_F5  // 临时变量
	REG_FT6  = REG_F6  // 临时变量
	REG_FT7  = REG_F7  // 临时变量
	REG_FS0  = REG_F8  // Saved register
	REG_FS1  = REG_F9  // Saved register
	REG_FA0  = REG_F10 // 函数参数/返回值
	REG_FA1  = REG_F11 // 函数参数/返回值
	REG_FA2  = REG_F12 // 函数参数
	REG_FA3  = REG_F13 // 函数参数
	REG_FA4  = REG_F14 // 函数参数
	REG_FA5  = REG_F15 // 函数参数
	REG_FA6  = REG_F16 // 函数参数
	REG_FA7  = REG_F17 // 函数参数
	REG_FS2  = REG_F18 // Saved register
	REG_FS3  = REG_F19 // Saved register
	REG_FS4  = REG_F20 // Saved register
	REG_FS5  = REG_F21 // Saved register
	REG_FS6  = REG_F22 // Saved register
	REG_FS7  = REG_F23 // Saved register
	REG_FS8  = REG_F24 // Saved register
	REG_FS9  = REG_F25 // Saved register
	REG_FS10 = REG_F26 // Saved register
	REG_FS11 = REG_F27 // Saved register
	REG_FT8  = REG_F28 // 临时变量
	REG_FT9  = REG_F29 // 临时变量
	REG_FT10 = REG_F30 // 临时变量
	REG_FT11 = REG_F31 // 临时变量
)

// 指令类型
type As int32

// 各平台通用的伪指令
const (
	AXXX           As = iota // 无效或未初始化的指令
	AGLOBL                   // 全局变量定义(类似于 .globl)
	ADATA                    // 静态数据段的数据定义
	ATEXT                    // 函数定义入口标记, 指定函数名和属性
	A_ARCHSPECIFIC           // 架构专属操作码的起点
)

// 优先支持最小指令集
//
// https://github.com/riscv/riscv-isa-manual
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
// 35. RV32/64G Instruction Set Listings
const (
	//
	// Unprivileged ISA (version 20240411)
	//

	// 2.4: Integer Computational Instructions (RV32I)
	AADDI As = A_ARCHSPECIFIC + iota
	ASLTI
	ASLTIU
	AANDI
	AORI
	AXORI
	ASLLI
	ASRLI
	ASRAI
	ALUI
	AAUIPC
	AADD
	ASLT
	ASLTU
	AAND
	AOR
	AXOR
	ASLL
	ASRL
	ASUB
	ASRA

	// 2.5: Control Transfer Instructions (RV32I)
	AJAL
	AJALR
	ABEQ
	ABNE
	ABLT
	ABLTU
	ABGE
	ABGEU

	// 2.6: Load and Store Instructions (RV32I)
	ALW
	ALH
	ALHU
	ALB
	ALBU
	ASW
	ASH
	ASB

	// 2.7: Memory Ordering Instructions (RV32I)
	AFENCE

	// 3.3.1: Environment Call and Breakpoint
	AECALL
	AEBREAK

	// 4.2: Integer Computational Instructions (RV64I)
	AADDIW
	ASLLIW
	ASRLIW
	ASRAIW
	AADDW
	ASLLW
	ASRLW
	ASUBW
	ASRAW

	// 4.3: Load and Store Instructions (RV64I)
	ALWU
	ALD
	ASD

	// 7.1: CSR Instructions (Zicsr)
	ACSRRW
	ACSRRS
	ACSRRC
	ACSRRWI
	ACSRRSI
	ACSRRCI

	// 13.1: Multiplication Operations (RV32M/RV64M)
	AMUL
	AMULH
	AMULHU
	AMULHSU
	AMULW // RV64M

	// 13.2: Division Operations (RV32M/RV64M)
	ADIV
	ADIVU
	AREM
	AREMU
	ADIVW  // RV64M
	ADIVUW // RV64M
	AREMW  // RV64M
	AREMUW // RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	AFLW
	AFSW

	// 20.6: Single-Precision Floating-Point Computational Instructions
	AFADDS
	AFSUBS
	AFMULS
	AFDIVS
	AFMINS
	AFMAXS
	AFSQRTS
	AFMADDS
	AFMSUBS
	AFNMADDS
	AFNMSUBS

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVTWS
	AFCVTLS
	AFCVTSW
	AFCVTSL
	AFCVTWUS
	AFCVTLUS
	AFCVTSWU
	AFCVTSLU
	AFSGNJS
	AFSGNJNS
	AFSGNJXS
	AFMVXW
	AFMVWX

	// 20.8: Single-Precision Floating-Point Compare Instructions
	AFEQS
	AFLTS
	AFLES

	// 20.9: Single-Precision Floating-Point Classify Instruction
	AFCLASSS

	// 21.3: Double-Precision Load and Store Instructions (D)
	AFLD
	AFSD

	// 21.4: Double-Precision Floating-Point Computational Instructions
	AFADDD
	AFSUBD
	AFMULD
	AFDIVD
	AFMIND
	AFMAXD
	AFSQRTD
	AFMADDD
	AFMSUBD
	AFNMADDD
	AFNMSUBD

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVTWD
	AFCVTLD
	AFCVTDW
	AFCVTDL
	AFCVTWUD
	AFCVTLUD
	AFCVTDWU
	AFCVTDLU
	AFCVTSD
	AFCVTDS
	AFSGNJD
	AFSGNJND
	AFSGNJXD
	AFMVXD
	AFMVDX

	// 21.6: Double-Precision Floating-Point Compare Instructions
	AFEQD
	AFLTD
	AFLED

	// 21.7: Double-Precision Floating-Point Classify Instruction
	AFCLASSD

	// End marker
	ALAST
)
