// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

//
// 以下的指令编码是龙芯官方手册简化后的版本(9种)
// 凹语言龙芯汇编器已经针对每个细粒度的不同编码作了定义区分处理(大约40+种)
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

// 龙芯指令格式类型
type OpFormatType int

const (
	OpFormatType_NULL OpFormatType = iota
	OpFormatType_2R
	OpFormatType_2F
	OpFormatType_1F_1R
	OpFormatType_1R_1F
	OpFormatType_3R
	OpFormatType_3F
	OpFormatType_1F_2R
	OpFormatType_4F
	OpFormatType_2R_ui5
	OpFormatType_2R_ui6
	OpFormatType_2R_si12
	OpFormatType_1F_1R_si12
	OpFormatType_2R_ui12
	OpFormatType_2R_si14
	OpFormatType_1R_si20
	OpFormatType_0_2R
	OpFormatType_3R_sa2
	OpFormatType_3R_sa3
	OpFormatType_code
	OpFormatType_code_1R_si12
	OpFormatType_2R_msbw_lsbw
	OpFormatType_2R_msbd_lsbd
	OpFormatType_fcsr_1R
	OpFormatType_1R_fcsr
	OpFormatType_cd_1R
	OpFormatType_cd_1F
	OpFormatType_cd_2F
	OpFormatType_1R_cj
	OpFormatType_1F_cj
	OpFormatType_1R_csr
	OpFormatType_2R_csr
	OpFormatType_2R_level
	OpFormatType_level
	OpFormatType_0_1R_seq
	OpFormatType_op_2R
	OpFormatType_3F_ca
	OpFormatType_hint_1R_si12
	OpFormatType_hint_2R
	OpFormatType_hint
	OpFormatType_cj_offset
	OpFormatType_rj_offset
	OpFormatType_rj_rd_offset
	OpFormatType_rd_rj_offset
	OpFormatType_offset
)

// 指令编码格式
type _OpContextType struct {
	mask  uint32       // opcode 掩码
	value uint32       // opcode 值
	op    abi.As       // 操作码定义
	fmt   OpFormatType // 指令格式
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regI(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_R0, REG_R31)
}

// 返回浮点数寄存器机器码编号
func (ctx *_OpContextType) regF(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_F0, REG_F31)
}

// 浮点数状态寄存器
func (ctx *_OpContextType) regFCSR(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_FCSR0, REG_FCSR3)
}

// 浮点数条件标志寄存器
func (ctx *_OpContextType) regFCC(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_FCC0, REG_FCC7)
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regVal(r, min, max abi.RegType) uint32 {
	if r < min || r > max {
		panic(fmt.Sprintf("register out of range, want %d <= %d <= %d", min, r, max))
	}
	return uint32(r - min)
}

func (x OpFormatType) String() string {
	switch x {
	case OpFormatType_NULL:
		return "OpFormatType_NULL"
	case OpFormatType_2R:
		return "OpFormatType_2R"
	case OpFormatType_2F:
		return "OpFormatType_2F"
	case OpFormatType_1F_1R:
		return "OpFormatType_1F_1R"
	case OpFormatType_1R_1F:
		return "OpFormatType_1R_1F"
	case OpFormatType_3R:
		return "OpFormatType_3R"
	case OpFormatType_3F:
		return "OpFormatType_3F"
	case OpFormatType_1F_2R:
		return "OpFormatType_1F_2R"
	case OpFormatType_4F:
		return "OpFormatType_4F"
	case OpFormatType_2R_ui5:
		return "OpFormatType_2R_ui5"
	case OpFormatType_2R_ui6:
		return "OpFormatType_2R_ui6"
	case OpFormatType_2R_si12:
		return "OpFormatType_2R_si12"
	case OpFormatType_1F_1R_si12:
		return "OpFormatType_1F_1R_si12"
	case OpFormatType_2R_ui12:
		return "OpFormatType_2R_ui12"
	case OpFormatType_2R_si14:
		return "OpFormatType_2R_si14"
	case OpFormatType_1R_si20:
		return "OpFormatType_1R_si20"
	case OpFormatType_0_2R:
		return "OpFormatType_0_2R"
	case OpFormatType_3R_sa2:
		return "OpFormatType_3R_sa2"
	case OpFormatType_3R_sa3:
		return "OpFormatType_3R_sa3"
	case OpFormatType_code:
		return "OpFormatType_code"
	case OpFormatType_code_1R_si12:
		return "OpFormatType_code_1R_si12"
	case OpFormatType_2R_msbw_lsbw:
		return "OpFormatType_2R_msbw_lsbw"
	case OpFormatType_2R_msbd_lsbd:
		return "OpFormatType_2R_msbd_lsbd"
	case OpFormatType_fcsr_1R:
		return "OpFormatType_fcsr_1R"
	case OpFormatType_1R_fcsr:
		return "OpFormatType_1R_fcsr"
	case OpFormatType_cd_1R:
		return "OpFormatType_cd_1R"
	case OpFormatType_cd_1F:
		return "OpFormatType_cd_1F"
	case OpFormatType_cd_2F:
		return "OpFormatType_cd_2F"
	case OpFormatType_1R_cj:
		return "OpFormatType_1R_cj"
	case OpFormatType_1F_cj:
		return "OpFormatType_1F_cj"
	case OpFormatType_1R_csr:
		return "OpFormatType_1R_csr"
	case OpFormatType_2R_csr:
		return "OpFormatType_2R_csr"
	case OpFormatType_2R_level:
		return "OpFormatType_2R_level"
	case OpFormatType_level:
		return "OpFormatType_level"
	case OpFormatType_0_1R_seq:
		return "OpFormatType_0_1R_seq"
	case OpFormatType_op_2R:
		return "OpFormatType_op_2R"
	case OpFormatType_3F_ca:
		return "OpFormatType_3F_ca"
	case OpFormatType_hint_1R_si12:
		return "OpFormatType_hint_1R_si12"
	case OpFormatType_hint_2R:
		return "OpFormatType_hint_2R"
	case OpFormatType_hint:
		return "OpFormatType_hint"
	case OpFormatType_cj_offset:
		return "OpFormatType_cj_offset"
	case OpFormatType_rj_offset:
		return "OpFormatType_rj_offset"
	case OpFormatType_rj_rd_offset:
		return "OpFormatType_rj_rd_offset"
	case OpFormatType_rd_rj_offset:
		return "OpFormatType_rd_rj_offset"
	case OpFormatType_offset:
		return "OpFormatType_offset"
	default:
		return fmt.Sprintf("loong64.OpFormatType(%d)", int(x))
	}
}
