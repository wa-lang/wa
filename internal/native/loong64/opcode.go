// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

//
// 31                                                           10 9            5 4            0
// +-------------------------------------------------------------+ +------------+ +------------+
// | opcode                                                      | | rj         | | rd         | 2R-type
// +-------------------------------------------------------------+ +------------+ +------------+
//
// 31                                            15 14          10 9            5 4            0
// +----------------------------------------------+ +------------+ +------------+ +------------+
// | opcode                                       | | rk         | | rj         | | rd         | 3R-type
// +----------------------------------------------+ +------------+ +------------+ +------------+
//
// 31                             20 19          15 14          10 9            5 4            0
// +-------------------------------+ +------------+ +------------+ +------------+ +------------+
// | opcode                        | | ra         | | rk         | | rj         | | rd         | 4R-type
// +-------------------------------+ +------------+ +------------+ +------------+ +------------+
//
// 31                                   18 17                   10 9            5 4            0
// +-------------------------------------+ +---------------------+ +------------+ +------------+
// | opcode                              | | I8                  | | rj         | | rd         | 2RI8-type
// +-------------------------------------+ +---------------------+ +------------+ +------------+
//
// 31                      22 21                                10 9            5 4            0
// +------------------------+ +----------------------------------+ +------------+ +------------+
// | opcode                 | | I12                              | | rj         | | rd         | 2RI12-type
// +------------------------+ +----------------------------------+ +------------+ +------------+
//
// 31                 24 23                                     10 9            5 4            0
// +-------------------+ +---------------------------------------+ +------------+ +------------+
// | opcode            | | I14                                   | | rj         | | rd         | 2RI14-type
// +-------------------+ +---------------------------------------+ +------------+ +------------+
//
// 31            26 25                                          10 9            5 4            0
// +--------------+ +--------------------------------------------+ +------------+ +------------+
// | opcode       | | I16                                        | | rj         | | rd         | 2RI16-type
// +--------------+ +--------------------------------------------+ +------------+ +------------+
//
// 31            26 25                                          10 9            5 4            0
// +--------------+ +--------------------------------------------+ +------------+ +------------+
// | opcode       | | I21[15:0]                                  | |   rj       | | I21[20:16] | 1RI21-type
// +--------------+ +--------------------------------------------+ +------------+ +------------+
//
// 31            26 25                                          10 9                           0
// +--------------+ +--------------------------------------------+ +---------------------------+
// | opcode       | | I21[15:0]                                  | | I21[25:16]                | I26-type
// +--------------+ +--------------------------------------------+ +---------------------------+
//

// 指令格式类型9种
type OpFormatType int

const (
	_ OpFormatType = iota + 1
	OpFormatType_2R
	OpFormatType_3R
	OpFormatType_4R
	OpFormatType_2RI8
	OpFormatType_2RI12
	OpFormatType_2RI14
	OpFormatType_2RI16
	OpFormatType_1RI21
	OpFormatType_I26
)

type instArg uint16

// 指令每个参数
const (
	_ instArg = iota // 0 是无效参数, 表示列表结束
	// 1-5 bit
	arg_fd
	arg_fj
	arg_fk
	arg_fa
	arg_rd
	// 6-10 bit
	arg_rj
	arg_rk
	arg_op_4_0
	arg_fcsr_4_0
	arg_fcsr_9_5
	// 11-15 bit
	arg_csr_23_10
	arg_cd
	arg_cj
	arg_ca
	arg_sa2_16_15
	// 16-20 bit
	arg_sa3_17_15
	arg_code_4_0
	arg_code_14_0
	arg_ui5_14_10
	arg_ui6_15_10
	// 21-25 bit
	arg_ui12_21_10
	arg_lsbw
	arg_msbw
	arg_lsbd
	arg_msbd
	// 26-30 bit
	arg_hint_4_0
	arg_hint_14_0
	arg_level_14_0
	arg_level_17_10
	arg_seq_17_10
	// 31-35 bit
	arg_si12_21_10
	arg_si14_23_10
	arg_si16_25_10
	arg_si20_24_5
	arg_offset_20_0
	// 36~
	arg_offset_25_0
	arg_offset_15_0
)

// 指令参数列表打包为数组结构, 可以简化初始化
type instArgs [5]instArg

// 指令编码格式
type _OpContextType struct {
	mask  uint32   // opcode 掩码
	value uint32   // opcode 值
	op    abi.As   // 操作码定义
	args  instArgs // args[i] 表示结束
}

func (ctx *_OpContextType) encodeRaw(as abi.As, arg *abi.AsArgument) (uint32, error) {
	assert(ctx != nil)
	assert(ctx.op == as)
	assert(arg != nil)

	var x = ctx.mask & ctx.value
	for _, argTyp := range ctx.args {
		x &= ctx.encodeArg(argTyp, arg)
	}

	return x, nil
}

func (ctx *_OpContextType) encodeArg(argTyp instArg, arg *abi.AsArgument) uint32 {
	// 根据参数类型分别编码
	switch argTyp {
	case 0:
		return 0

	// 1-5 bit
	case arg_fd:
		return ctx.regF(arg.Rd)
	case arg_fj:
		return ctx.regF(arg.Rd) << 5
	case arg_fk:
		return ctx.regF(arg.Rd) << 10
	case arg_fa:
		return ctx.regF(arg.Rd) << 15
	case arg_rd:
		return ctx.regI(arg.Rd)

	// 6-10 bit
	case arg_rj:
		return ctx.regI(arg.Rd) << 5
	case arg_rk:
		return ctx.regI(arg.Rd) << 10
	case arg_op_4_0:
		panic("TODO")
	case arg_fcsr_4_0:
		panic("TODO")
	case arg_fcsr_9_5:
		panic("TODO")

	// 11-15 bit
	case arg_csr_23_10:
		panic("TODO")
	case arg_cd:
		panic("TODO")
	case arg_cj:
		panic("TODO")
	case arg_ca:
		panic("TODO")
	case arg_sa2_16_15:
		panic("TODO")

	// 16-20 bit
	case arg_sa3_17_15:
		panic("TODO")
	case arg_code_4_0:
		panic("TODO")
	case arg_code_14_0:
		panic("TODO")
	case arg_ui5_14_10:
		panic("TODO")
	case arg_ui6_15_10:
		panic("TODO")

	// 21-25 bit
	case arg_ui12_21_10:
		panic("TODO")
	case arg_lsbw:
		panic("TODO")
	case arg_msbw:
		panic("TODO")
	case arg_lsbd:
		panic("TODO")
	case arg_msbd:
		panic("TODO")

	// 26-30 bit
	case arg_hint_4_0:
		panic("TODO")
	case arg_hint_14_0:
		panic("TODO")
	case arg_level_14_0:
		panic("TODO")
	case arg_level_17_10:
		panic("TODO")
	case arg_seq_17_10:
		panic("TODO")

	// 31-35 bit
	case arg_si12_21_10:
		panic("TODO")
	case arg_si14_23_10:
		panic("TODO")
	case arg_si16_25_10:
		panic("TODO")
	case arg_si20_24_5:
		panic("TODO")
	case arg_offset_20_0:
		panic("TODO")

	// 36~
	case arg_offset_25_0:
		panic("TODO")
	case arg_offset_15_0:
		panic("TODO")
	}

	panic("unreachable")
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regI(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_R0, REG_R31)
}

// 返回浮点数寄存器机器码编号
func (ctx *_OpContextType) regF(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_F0, REG_F31)
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regVal(r, min, max abi.RegType) uint32 {
	if r < min || r > max {
		panic(fmt.Sprintf("register out of range, want %d <= %d <= %d", min, r, max))
	}
	return uint32(r - min)
}
