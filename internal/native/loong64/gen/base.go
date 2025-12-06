// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

// 指令编码格式
type OpContextType struct {
	mask  uint32   // opcode 掩码
	value uint32   // opcode 值
	name  string   // 操作码定义
	args  instArgs // args[i] 表示结束
}

// 指令格式类型
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
	OpFormatType_2R_ui12
	OpFormatType_2R_si14
	OpFormatType_2R_si16
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

// 指令参数列表打包为数组结构, 可以简化初始化
type instArgs [5]InstArg

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

var OpFormatTypeNameList = []string{
	"OpFormatType_NULL",
	"OpFormatType_2R",
	"OpFormatType_2F",
	"OpFormatType_1F_1R",
	"OpFormatType_1R_1F",
	"OpFormatType_3R",
	"OpFormatType_3F",
	"OpFormatType_1F_2R",
	"OpFormatType_4F",
	"OpFormatType_2R_ui5",
	"OpFormatType_2R_ui6",
	"OpFormatType_2R_si12",
	"OpFormatType_2R_ui12",
	"OpFormatType_2R_si14",
	"OpFormatType_2R_si16",
	"OpFormatType_1R_si20",
	"OpFormatType_0_2R",
	"OpFormatType_3R_sa2",
	"OpFormatType_3R_sa3",
	"OpFormatType_code",
	"OpFormatType_code_1R_si12",
	"OpFormatType_2R_msbw_lsbw",
	"OpFormatType_2R_msbd_lsbd",
	"OpFormatType_fcsr_1R",
	"OpFormatType_1R_fcsr",
	"OpFormatType_cd_1R",
	"OpFormatType_cd_1F",
	"OpFormatType_cd_2F",
	"OpFormatType_1R_cj",
	"OpFormatType_1F_cj",
	"OpFormatType_1R_csr",
	"OpFormatType_2R_csr",
	"OpFormatType_2R_level",
	"OpFormatType_level",
	"OpFormatType_0_1R_seq",
	"OpFormatType_op_2R",
	"OpFormatType_3F_ca",
	"OpFormatType_hint_1R_si12",
	"OpFormatType_hint_2R",
	"OpFormatType_hint",
	"OpFormatType_cj_offset",
	"OpFormatType_rj_offset",
	"OpFormatType_rj_rd_offset",
	"OpFormatType_rd_rj_offset",
	"OpFormatType_offset",
}
