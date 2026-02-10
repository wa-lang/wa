// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import (
	"wa-lang.org/wa/internal/native/abi"
)

const (
	// 通用寄存器
	REG_X0 abi.RegType = iota + 1 // 0 是无效的编号
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
	REG_X15
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

	// 浮点数寄存器
	REG_V0
	REG_V1
	REG_V2
	REG_V3
	REG_V4
	REG_V5
	REG_V6
	REG_V7
	REG_V8
	REG_V9
	REG_V10
	REG_V11
	REG_V12
	REG_V13
	REG_V14
	REG_V15
	REG_V16
	REG_V17
	REG_V18
	REG_V19
	REG_V20
	REG_V21
	REG_V22
	REG_V23
	REG_V24
	REG_V25
	REG_V26
	REG_V27
	REG_V28
	REG_V29
	REG_V30
	REG_V31

	// 寄存器编号结束
	REG_END

	// TODO: 修改正确
	// 寄存器的 ABI 使用约定
	REG_ZERO = REG_X0  // 零寄存器
	REG_RA   = REG_X1  // 返回地址
	REG_TP   = REG_X2  // 线程指针
	REG_SP   = REG_X3  // 栈指针
	REG_GP   = REG_X8  // S8 复用于 WASM 的 Memory 地址, 临时方案
	REG_A0   = REG_X4  // 函数参数/返回值
	REG_A1   = REG_X5  // 函数参数/返回值
	REG_A2   = REG_X6  // 函数参数
	REG_A3   = REG_X7  // 函数参数
	REG_A4   = REG_X8  // 函数参数
	REG_A5   = REG_X9  // 函数参数
	REG_A6   = REG_X10 // 函数参数
	REG_A7   = REG_X11 // 函数参数
	REG_T0   = REG_X12 // 临时变量
	REG_T1   = REG_X13 // 临时变量
	REG_T2   = REG_X14 // 临时变量
	REG_T3   = REG_X15 // 临时变量
	REG_T4   = REG_X16 // 临时变量
	REG_T5   = REG_X17 // 临时变量
	REG_T6   = REG_X18 // 临时变量
	REG_T7   = REG_X19 // 临时变量
	REG_T8   = REG_X20 // 临时变量
	REG_X    = REG_X21 // 保留寄存器
	REG_FP   = REG_X22 // 帧指针
	REG_S0   = REG_X23 // Saved register
	REG_S1   = REG_X24 // Saved register
	REG_S2   = REG_X25 // Saved register
	REG_S3   = REG_X26 // Saved register
	REG_S4   = REG_X27 // Saved register
	REG_S5   = REG_X28 // Saved register
	REG_S6   = REG_X29 // Saved register
	REG_S7   = REG_X30 // Saved register
	REG_S8   = REG_X31 // Saved register

	// TODO: 修改正确
	REG_FA0  = REG_V0  // 函数参数/返回值
	REG_FA1  = REG_V1  // 函数参数/返回值
	REG_FA2  = REG_V2  // 函数参数
	REG_FA3  = REG_V3  // 函数参数
	REG_FA4  = REG_V4  // 函数参数
	REG_FA5  = REG_V5  // 函数参数
	REG_FA6  = REG_V6  // 函数参数
	REG_FA7  = REG_V7  // 函数参数
	REG_FT0  = REG_V8  // 临时变量
	REG_FT1  = REG_V9  // 临时变量
	REG_FT2  = REG_V10 // 临时变量
	REG_FT3  = REG_V11 // 临时变量
	REG_FT4  = REG_V12 // 临时变量
	REG_FT5  = REG_V13 // 临时变量
	REG_FT6  = REG_V14 // 临时变量
	REG_FT7  = REG_V15 // 临时变量
	REG_FT8  = REG_V16 // 临时变量
	REG_FT9  = REG_V17 // 临时变量
	REG_FT10 = REG_V18 // 临时变量
	REG_FT11 = REG_V19 // 临时变量
	REG_FT12 = REG_V20 // 临时变量
	REG_FT13 = REG_V21 // 临时变量
	REG_FT14 = REG_V22 // 临时变量
	REG_FT15 = REG_V23 // 临时变量
	REG_FS0  = REG_V24 // Saved register
	REG_FS1  = REG_V25 // Saved register
	REG_FS2  = REG_V26 // Saved register
	REG_FS3  = REG_V27 // Saved register
	REG_FS4  = REG_V28 // Saved register
	REG_FS5  = REG_V29 // Saved register
	REG_FS6  = REG_V30 // Saved register
	REG_FS7  = REG_V31 // Saved register
)
