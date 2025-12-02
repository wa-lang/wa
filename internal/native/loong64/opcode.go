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
	OpFormatType_NULL OpFormatType = iota // 其他格式数量较少
	OpFormatType_2R
	OpFormatType_3R
	OpFormatType_4R
	OpFormatType_2RI8
	OpFormatType_2RI12
	OpFormatType_2RI14
	OpFormatType_2RI16
	OpFormatType_1RI20 // 新加, 参考手册没有
	OpFormatType_1RI21
	OpFormatType_I26  // OpFormatType_offset
	OpFormatType_0_2R // 后面都是新加
	OpFormatType_3R_s2
	OpFormatType_3R_s3
	OpFormatType_code
	OpFormatType_code_1R_si12
	OpFormatType_msbw_lsbw
	OpFormatType_msbd_lsbd
	OpFormatType_fcsr_1R
	OpFormatType_1R_fcsr
	OpFormatType_cd_1R
	OpFormatType_cd_2R
	OpFormatType_1R_cj
	OpFormatType_1R_csr
	OpFormatType_2R_csr
	OpFormatType_2R_level
	OpFormatType_level
	OpFormatType_0_1R_seq
	OpFormatType_3R_ca
	OpFormatType_hint_1R_si
	OpFormatType_hint_2R
	OpFormatType_hint
	OpFormatType_cj_offset
	OpFormatType_rj_offset
	OpFormatType_rj_rd_offset
	OpFormatType_rd_rj_offset
	OpFormatType_offset
)

type InstArg uint16

// 指令每个参数
const (
	_ InstArg = iota // 0 是无效参数, 表示列表结束
	// 1-5 bit
	Arg_fd
	Arg_fj
	Arg_fk
	Arg_fa
	Arg_rd
	// 6-10 bit
	Arg_rj
	Arg_rk
	Arg_op_4_0
	Arg_fcsr_4_0
	Arg_fcsr_9_5
	// 11-15 bit
	Arg_csr_23_10
	Arg_cd
	Arg_cj
	Arg_ca
	Arg_sa2_16_15
	// 16-20 bit
	Arg_sa3_17_15
	Arg_code_4_0
	Arg_code_14_0
	Arg_ui5_14_10
	Arg_ui6_15_10
	// 21-25 bit
	Arg_ui12_21_10
	Arg_lsbw
	Arg_msbw
	Arg_lsbd
	Arg_msbd
	// 26-30 bit
	Arg_hint_4_0
	Arg_hint_14_0
	Arg_level_14_0
	Arg_level_17_10
	Arg_seq_17_10
	// 31-35 bit
	Arg_si12_21_10
	Arg_si14_23_10
	Arg_si16_25_10
	Arg_si20_24_5
	Arg_offset_20_0
	// 36~
	Arg_offset_25_0
	Arg_offset_15_0
)

// 指令参数列表打包为数组结构, 可以简化初始化
type instArgs [5]InstArg

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

// 指令的编码格式
func (opcode _OpContextType) FormatType() OpFormatType {
	regCount := 0
	for i, argTyp := range opcode.args {
		if argTyp == 0 {
			break
		}
		switch argTyp {
		case Arg_fd, Arg_fj, Arg_fk, Arg_fa:
			regCount++
		case Arg_rd, Arg_rj, Arg_rk:
			regCount++
		case Arg_ui5_14_10:
			return OpFormatType_2RI8
		case Arg_ui6_15_10:
			return OpFormatType_2RI8
		case Arg_si12_21_10:
			return OpFormatType_2RI12
		case Arg_si14_23_10:
			return OpFormatType_2RI14
		case Arg_si16_25_10:
			return OpFormatType_1RI21
		case Arg_si20_24_5:
			return OpFormatType_1RI20
		case Arg_offset_20_0:
			switch opcode.args[0] {
			case Arg_rj:
				return OpFormatType_rj_offset
			case Arg_cj:
				return OpFormatType_cj_offset
			default:
				panic("unreachable")
			}
		case Arg_offset_25_0:
			return OpFormatType_offset
		case Arg_offset_15_0:
			assert(i == 2)
			switch {
			case opcode.args[0] == Arg_rd:
				assert(opcode.args[1] == Arg_rj)
				return OpFormatType_rd_rj_offset
			case opcode.args[0] == Arg_rj:
				assert(opcode.args[1] == Arg_rd)
				return OpFormatType_rj_rd_offset
			default:
				panic("unreachable")
			}
		case Arg_sa2_16_15:
			return OpFormatType_3R_s2
		case Arg_sa3_17_15:
			return OpFormatType_3R_s3
		case Arg_code_4_0:
			return OpFormatType_code_1R_si12
		case Arg_code_14_0:
			return OpFormatType_code

		case Arg_lsbw, Arg_msbw:
			return OpFormatType_msbw_lsbw
		case Arg_lsbd, Arg_msbd:
			return OpFormatType_msbd_lsbd

		case Arg_fcsr_4_0:
			if opcode.args[0] == Arg_fcsr_4_0 {
				return OpFormatType_fcsr_1R
			} else {
				return OpFormatType_1R_fcsr
			}

		case Arg_cd:
			if arg2 := opcode.args[2]; arg2 == Arg_fk {
				return OpFormatType_cd_2R
			} else {
				return OpFormatType_cd_1R
			}

		case Arg_cj:
			assert(opcode.args[0] == Arg_fd || opcode.args[0] == Arg_rd)
			return OpFormatType_1R_cj

		case Arg_csr_23_10:
			if i == 1 {
				return OpFormatType_1R_csr
			} else {
				assert(i == 2)
				return OpFormatType_2R_csr
			}

		case Arg_level_14_0:
			return OpFormatType_level
		case Arg_level_17_10:
			return OpFormatType_2R_level

		case Arg_seq_17_10:
			return OpFormatType_0_1R_seq

		case Arg_ca:
			assert(opcode.args[1] == Arg_fd)
			assert(opcode.args[1] == Arg_fj)
			assert(opcode.args[2] == Arg_fk)
			assert(i == 3)
			return OpFormatType_3R_ca

		case Arg_hint_4_0:
			assert(i == 0)
			if opcode.args[2] == Arg_rk {
				return OpFormatType_hint_2R
			} else {
				assert(opcode.args[2] == Arg_si12_21_10)
				return OpFormatType_hint_1R_si
			}
		case Arg_hint_14_0:
			return OpFormatType_hint

		default:
			panic("unreachable")
		}
	}

	switch regCount {
	case 0:
		return OpFormatType_NULL
	case 2:
		if opcode.args[0] != Arg_rd {
			assert(opcode.op == AASRTLE_D || opcode.op == AASRTGT_D)
			return OpFormatType_0_2R
		}
		return OpFormatType_2R
	case 3:
		return OpFormatType_3R
	case 4:
		return OpFormatType_4R
	}

	return 0
}

func (ctx *_OpContextType) encodeArg(argTyp InstArg, arg *abi.AsArgument) uint32 {
	// 根据参数类型分别编码
	switch argTyp {
	case 0: // 空参数
		return 0

	// 1-5 bit
	case Arg_fd:
		return ctx.regF(arg.Rd)
	case Arg_fj:
		return ctx.regF(arg.Rs1) << 5
	case Arg_fk:
		return ctx.regF(arg.Rs2) << 10
	case Arg_fa:
		return ctx.regF(arg.Rs3) << 15
	case Arg_rd:
		return ctx.regI(arg.Rd)

	// 6-10 bit
	case Arg_rj:
		return ctx.regI(arg.Rs1) << 5
	case Arg_rk:
		return ctx.regI(arg.Rs2) << 10
	case Arg_op_4_0:
		panic("TODO")
	case Arg_fcsr_4_0:
		return (uint32(arg.Imm) & 0b11)
	case Arg_fcsr_9_5:
		return (uint32(arg.Imm) & 0b_1_1111_1111) << 5

	// 11-15 bit
	case Arg_csr_23_10:
		return (uint32(arg.Imm) & 0b11) << 10
	case Arg_cd:
		return (ctx.regI(arg.Rd) | 0b_11000)
	case Arg_cj:
		return (ctx.regI(arg.Rs1) | 0b_11000) << 5
	case Arg_ca:
		return (ctx.regI(arg.Rs2) | 0b_11000) << 10
	case Arg_sa2_16_15:
		return (uint32(arg.Imm) & 0b11) << 15

	// 16-20 bit
	case Arg_sa3_17_15:
		return (uint32(arg.Imm) & 0b111) << 15
	case Arg_code_4_0:
		return (uint32(arg.Imm) & 0b11)
	case Arg_code_14_0:
		return (uint32(arg.Imm) & 0b_0111_1111_1111_1111)
	case Arg_ui5_14_10:
		return (uint32(arg.Imm) & 0b_0111_1111_1111_1111) << 10
	case Arg_ui6_15_10:
		return (uint32(arg.Imm) & 0b_1111_1111_1111_1111) << 10

	// 21-25 bit
	case Arg_ui12_21_10:
		return (uint32(arg.Imm) & 0b_1111_1111_1111) << 10
	case Arg_lsbw:
		panic("TODO") // 需要2个立即数
	case Arg_msbw:
		panic("TODO") // 需要2个立即数
	case Arg_lsbd:
		panic("TODO") // 需要2个立即数
	case Arg_msbd:
		panic("TODO") // 需要2个立即数

	// 26-30 bit
	case Arg_hint_4_0:
		return uint32(arg.Imm) & 0b_1_1111
	case Arg_hint_14_0:
		return uint32(arg.Imm) & 0b_0111_1111_1111_1111
	case Arg_level_14_0:
		return uint32(arg.Imm) & 0b_0111_1111_1111_1111
	case Arg_level_17_10:
		return uint32(arg.Imm) << 10

	case Arg_seq_17_10:
		return uint32(arg.Imm) << 10

	// 31-35 bit
	case Arg_si12_21_10:
		return (uint32(arg.Imm) & 0b_1111_1111_1111) << 10
	case Arg_si14_23_10:
		return (uint32(arg.Imm) & 0b_0011_1111_1111_1111) << 10
	case Arg_si16_25_10:
		return (uint32(arg.Imm) & 0xFFFF) << 10
	case Arg_si20_24_5:
		return (uint32(arg.Imm) & 0xFFFFF) << 5
	case Arg_offset_20_0:
		offs := uint32(arg.Imm)
		return (offs&0xFFFF)<<10 | ((offs >> 16) & 0b_11111)

	// 36~
	case Arg_offset_25_0:
		// op, offs[15:0], offs[25:16]
		offs := uint32(arg.Imm)
		bit0_9 := (offs >> 16) & 0b_11111_11111
		bit10_25 := (offs & 0xFFFF) << 10
		return bit0_9 | bit10_25
	case Arg_offset_15_0:
		return (uint32(arg.Imm) & 0xFFFF) << 10
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
