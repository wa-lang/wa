// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

const (
	// 通用寄存器
	REG_R0 abi.RegType = iota + 1 // 0 是无效的编号
	REG_R1
	REG_R2
	REG_R3
	REG_R4
	REG_R5
	REG_R6
	REG_R7
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15
	REG_R16
	REG_R17
	REG_R18
	REG_R19
	REG_R20
	REG_R21
	REG_R22
	REG_R23
	REG_R24
	REG_R25
	REG_R26
	REG_R27
	REG_R28
	REG_R29
	REG_R30
	REG_R31

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
	REG_ZERO = REG_R0  // 零寄存器
	REG_RA   = REG_R1  // 返回地址
	REG_TP   = REG_R2  // 线程指针
	REG_SP   = REG_R3  // 栈指针
	REG_A0   = REG_R4  // 函数参数/返回值
	REG_A1   = REG_R5  // 函数参数/返回值
	REG_A2   = REG_R6  // 函数参数
	REG_A3   = REG_R7  // 函数参数
	REG_A4   = REG_R8  // 函数参数
	REG_A5   = REG_R9  // 函数参数
	REG_A6   = REG_R10 // 函数参数
	REG_A7   = REG_R11 // 函数参数
	REG_T0   = REG_R12 // 临时变量
	REG_T1   = REG_R13 // 临时变量
	REG_T2   = REG_R14 // 临时变量
	REG_T3   = REG_R15 // 临时变量
	REG_T4   = REG_R16 // 临时变量
	REG_T5   = REG_R17 // 临时变量
	REG_T6   = REG_R18 // 临时变量
	REG_T7   = REG_R19 // 临时变量
	REG_T8   = REG_R20 // 临时变量
	REG_X    = REG_R21 // 保留寄存器
	REG_FP   = REG_R22 // 帧指针
	REG_S0   = REG_R23 // Saved register
	REG_S1   = REG_R24 // Saved register
	REG_S2   = REG_R25 // Saved register
	REG_S3   = REG_R26 // Saved register
	REG_S4   = REG_R27 // Saved register
	REG_S5   = REG_R28 // Saved register
	REG_S6   = REG_R29 // Saved register
	REG_S7   = REG_R30 // Saved register
	REG_S8   = REG_R31 // Saved register

	REG_FA0  = REG_F0  // 函数参数/返回值
	REG_FA1  = REG_F1  // 函数参数/返回值
	REG_FA2  = REG_F2  // 函数参数
	REG_FA3  = REG_F3  // 函数参数
	REG_FA4  = REG_F4  // 函数参数
	REG_FA5  = REG_F5  // 函数参数
	REG_FA6  = REG_F6  // 函数参数
	REG_FA7  = REG_F7  // 函数参数
	REG_FT0  = REG_F8  // 临时变量
	REG_FT1  = REG_F9  // 临时变量
	REG_FT2  = REG_F10 // 临时变量
	REG_FT3  = REG_F11 // 临时变量
	REG_FT4  = REG_F12 // 临时变量
	REG_FT5  = REG_F13 // 临时变量
	REG_FT6  = REG_F14 // 临时变量
	REG_FT7  = REG_F15 // 临时变量
	REG_FT8  = REG_F16 // 临时变量
	REG_FT9  = REG_F17 // 临时变量
	REG_FT10 = REG_F18 // 临时变量
	REG_FT11 = REG_F19 // 临时变量
	REG_FT12 = REG_F20 // 临时变量
	REG_FT13 = REG_F21 // 临时变量
	REG_FT14 = REG_F22 // 临时变量
	REG_FT15 = REG_F23 // 临时变量
	REG_FS0  = REG_F24 // Saved register
	REG_FS1  = REG_F25 // Saved register
	REG_FS2  = REG_F26 // Saved register
	REG_FS3  = REG_F27 // Saved register
	REG_FS4  = REG_F28 // Saved register
	REG_FS5  = REG_F29 // Saved register
	REG_FS6  = REG_F30 // Saved register
	REG_FS7  = REG_F31 // Saved register
)

// 浮点数控制状态寄存器
type Fcsr uint8

const (
	FCSR0 Fcsr = iota
	FCSR1
	FCSR2
	FCSR3
)

func (f Fcsr) String() string {
	return fmt.Sprintf("$fcsr%d", uint8(f))
}

// 浮点数条件标志寄存器
type Fcc uint8

const (
	FCC0 Fcc = iota
	FCC1
	FCC2
	FCC3
	FCC4
	FCC5
	FCC6
	FCC7
)

func (f Fcc) String() string {
	return fmt.Sprintf("$fcc%d", uint8(f))
}
