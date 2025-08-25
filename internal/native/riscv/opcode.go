// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "wa-lang.org/wa/internal/native/abi"

//
// 0        6 7                   11 12      14 15                  19 20                  24 25                   31
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
// | opcode | |         rd         | | funct3 | |        rs1         | |        rs2         | |      funct7         | R-type, 寄存器类
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
//
// 0        6 7                   11 12      14 15                  19 20                  24 25      26 27        31
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +--------+ +----------+
// | opcode | |         rd         | | funct3 | |        rs1         | |        rs2         | | funct2 | |   rs3    | R4-type, 寄存器类
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +--------+ +----------+
//
// 0        6 7                   11 12      14 15                  19 20                                          31
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------------------------------+
// | opcode | |         rd         | | funct3 | |        rs1         | |        imm[0:11]                           | I-type, 立即数类
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------------------------------+
//
// 0        6 7                   11 12      14 15                  19 20                  24 25                   31
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
// | opcode | |      imm[0:4]      | | funct3 | |        rs1         | |        rs2         | |     imm[5:11]       | S-type, 存储类
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
//
// 0        6 7                   11 12      14 15                  19 20                  24 25                   31
// +--------+ +---------+----------+ +--------+ +--------------------+ +--------------------+ +-----------+---------+
// | opcode | | imm[11] | imm[1:4] | | funct3 | |        rs1         | |        rs2         | | imm[5:10] | imm[12] | B-type, 条件跳转类
// +--------+ +---------+----------+ +--------+ +--------------------+ +--------------------+ +-----------+---------+
//
// 0        6 7                   11 12                                                                            31
// +--------+ +--------------------+ +------------------------------------------------------------------------------+
// | opcode | |         rd         | | imm[12:31], (imm[0:11] = {0})                                                | U-type, 长立即数
// +--------+ +--------------------+ +------------------------------------------------------------------------------+
//
// 0        6 7                   11 12                                                                            31
// +--------+ +--------------------+ +---------------------------+---------+------------------------------+---------+
// | opcode | |         rd         | | imm[12:19]                | imm[11] |    imm[1:10]                 | imm[20] | J-type, 无条件跳转
// +--------+ +--------------------+ +---------------------------+---------+------------------------------+---------+
//

// 指令格式类型
type OpFormatType int

const (
	_ OpFormatType = iota
	R
	R4
	I
	S
	B
	U
	J
)

// Base opecode map (inst[0:1] = 11)
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+
// | inst[2:4] | 000    | 001      | 010      | 011      | 100    | 101      | 110            | 111     |
// | inst[5:6] |        |          |          |          |        |          |                | (> 32b) |
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+
// |       00  | LOAD   | LOAD-FP  | custom-0 | MISC-MEM | OP-IMM | AUIPC    | OP-IMM-32      | 48b     |
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+
// |       01  | STORE  | STORE-FP | custom-1 | AMO      | OP     | LUI      | OP-32          | 64b     |
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+
// |       10  | MADD   | MSUB     | NMSUB    | NMADD    | OP-FP  | reserved | custom-2/rv128 | 48b     |
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+
// |       11  | BRANCH | JALR     | reserved | JAL      | SYSTEM | reserved | custom-3/rv128 | >= 80b  |
// +-----------+--------+----------+----------+----------+--------+----------+----------------+---------+

// 基础 opcode 分类
type OpcodeType uint32

const (
	OpBase_LOAD      OpcodeType = 0b_00_000_11
	OpBase_LOAD_FP   OpcodeType = 0b_00_001_11
	OpBase_CUSTOM_0  OpcodeType = 0b_00_010_11
	OpBase_MISC_MEN  OpcodeType = 0b_00_011_11
	OpBase_OP_IMM    OpcodeType = 0b_00_100_11
	OpBase_AUIPC     OpcodeType = 0b_00_101_11
	OpBase_OP_IMM_32 OpcodeType = 0b_00_110_11

	OpBase_STORE    OpcodeType = 0b_01_000_11
	OpBase_STORE_FP OpcodeType = 0b_01_001_11
	OpBase_CUSTOM_1 OpcodeType = 0b_01_010_11
	OpBase_AMO      OpcodeType = 0b_01_011_11
	OpBase_OP       OpcodeType = 0b_01_100_11
	OpBase_LUI      OpcodeType = 0b_01_101_11
	OpBase_OP_32    OpcodeType = 0b_01_110_11

	OpBase_MADD     OpcodeType = 0b_10_000_11
	OpBase_MSUB     OpcodeType = 0b_10_001_11
	OpBase_NMSUB    OpcodeType = 0b_10_010_11
	OpBase_NMADD    OpcodeType = 0b_10_011_11
	OpBase_OP_FP    OpcodeType = 0b_10_100_11
	_               OpcodeType = 0b_10_101_11
	OpBase_CUSTOM_2 OpcodeType = 0b_10_110_11

	OpBase_BRANCH   OpcodeType = 0b_11_000_11
	OpBase_JALR     OpcodeType = 0b_11_001_11
	_               OpcodeType = 0b_11_010_11
	OpBase_JAL      OpcodeType = 0b_11_011_11
	OpBase_SYSTEM   OpcodeType = 0b_11_100_11
	_               OpcodeType = 0b_11_101_11
	OpBase_CUSTOM_3 OpcodeType = 0b_11_110_11

	OpBase_Mask OpcodeType = 0b_11_111_11
)

// 指令的编码格式
func (opcode OpcodeType) FormatType() OpFormatType {
	switch opcode & OpBase_Mask {
	case OpBase_OP,
		OpBase_OP_32,
		OpBase_OP_FP,
		OpBase_AMO:
		return R
	case OpBase_MADD,
		OpBase_MSUB,
		OpBase_NMSUB,
		OpBase_NMADD:
		return R4
	case OpBase_OP_IMM,
		OpBase_OP_IMM_32,
		OpBase_JALR,
		OpBase_LOAD,
		OpBase_LOAD_FP,
		OpBase_MISC_MEN,
		OpBase_SYSTEM:
		return I
	case OpBase_STORE,
		OpBase_STORE_FP:
		return S
	case OpBase_BRANCH:
		return B
	case OpBase_LUI,
		OpBase_AUIPC:
		return U
	case OpBase_JAL:
		return J

	case OpBase_CUSTOM_0,
		OpBase_CUSTOM_1,
		OpBase_CUSTOM_2,
		OpBase_CUSTOM_3:
		return 0

	default:
		return 0
	}
}

// 参数标志
type ArgMarks uint16

const (
	ARG_RD ArgMarks = 1 << iota
	ARG_RS1
	ARG_RS2
	ARG_RS3
	ARG_IMM
	ARG_FUNCT3
	ARG_FUNCT7
	ARG_FUNCT2
	ARG_SHAMT

	ARG_RType  = ARG_RD | ARG_RS1 | ARG_RS2 | ARG_FUNCT3 | ARG_FUNCT7
	ARG_R4Type = ARG_RD | ARG_RS1 | ARG_RS2 | ARG_RS3 | ARG_FUNCT2
	ARG_IType  = ARG_RD | ARG_RS1 | ARG_IMM | ARG_FUNCT3
	ARG_SType  = ARG_RS1 | ARG_RS2 | ARG_IMM | ARG_FUNCT3
	ARG_BType  = ARG_RS1 | ARG_RS2 | ARG_IMM | ARG_FUNCT3
	ARG_UType  = ARG_RD | ARG_IMM
	ARG_JType  = ARG_RD | ARG_IMM
)

// 操作码上下文信息
type OpContextType struct {
	Opcode   OpcodeType // 基础指令码
	ArgMarks ArgMarks   // 指令的参数标志
	Funct3   uint32     // funct3
	Funct7   uint32     // 和 Funct2 共用
	Rs2      *uint32    // 是否用到 Rs2 模板
	HasShamt bool       // 是否有 shamt 参数. SLLI/SRLI/SRAI 在 RV32 是 5bit, RV64 是 6bit
	PseudoAs abi.As     // 伪指令对应的原生指令
}

// Imm 取值范围
type ImmRange struct {
	Align int64
	Min   int64
	Max   int64
}

var (
	ImmRanges_IType   = ImmRange{Align: 1, Min: -(1 << 11), Max: (1 << 11) - 1}
	ImmRanges_SType   = ImmRange{Align: 1, Min: -(1 << 11), Max: (1 << 11) - 1}
	ImmRanges_BType   = ImmRange{Align: 1 << 1, Min: -(1 << 12), Max: (1 << 12) - 2}
	ImmRanges_UType   = ImmRange{Align: 1, Min: -(1 << 20), Max: (1 << 20) - 1}
	ImmRanges_JType   = ImmRange{Align: 1, Min: -(1 << 20), Max: (1 << 20) - 2}
	ImmRanges_Shamt32 = ImmRange{Align: 1, Min: 0, Max: (1 << 5) - 1}
	ImmRanges_Shamt64 = ImmRange{Align: 1, Min: 0, Max: (1 << 6) - 1}
)

// 指令编码信息表
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
var AOpContextTable = []OpContextType{
	// RV32I Base Instruction Set

	ALUI:    {Opcode: OpBase_LUI, ArgMarks: ARG_UType},
	AAUIPC:  {Opcode: OpBase_AUIPC, ArgMarks: ARG_UType},
	AJAL:    {Opcode: OpBase_JAL, ArgMarks: ARG_JType},
	AJALR:   {Opcode: OpBase_JALR, ArgMarks: ARG_IType, Funct3: 0b_000},
	ABEQ:    {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_000},
	ABNE:    {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_001},
	ABLT:    {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_100},
	ABGE:    {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_101},
	ABLTU:   {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_110},
	ABGEU:   {Opcode: OpBase_BRANCH, ArgMarks: ARG_BType, Funct3: 0b_111},
	ALB:     {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_000},
	ALH:     {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_001},
	ALW:     {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_010},
	ALBU:    {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_100},
	ALHU:    {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_101},
	ASB:     {Opcode: OpBase_STORE, ArgMarks: ARG_SType, Funct3: 0b_000},
	ASH:     {Opcode: OpBase_STORE, ArgMarks: ARG_SType, Funct3: 0b_001},
	ASW:     {Opcode: OpBase_STORE, ArgMarks: ARG_SType, Funct3: 0b_010},
	AADDI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_000},
	ASLTI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_010},
	ASLTIU:  {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_011},
	AXORI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_100},
	AORI:    {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_110},
	AANDI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_111},
	ASLLI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAI:   {Opcode: OpBase_OP_IMM, ArgMarks: ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADD:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUB:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLL:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASLT:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_010, Funct7: 0b_000_0000},
	ASLTU:   {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_011, Funct7: 0b_000_0000},
	AXOR:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0000},
	ASRL:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRA:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_010_0000},
	AOR:     {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0000},
	AAND:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0000},
	AFENCE:  {Opcode: OpBase_MISC_MEN, ArgMarks: ARG_IType, Funct3: 0b_000},
	AECALL:  {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_000}, // imm[11:0] = 0b000000000000
	AEBREAK: {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_000}, // imm[11:0] = 0b000000000001

	// RV64I Base Instruction Set (in addition to RV32I)

	ALWU:   {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_110},
	ALD:    {Opcode: OpBase_LOAD, ArgMarks: ARG_IType, Funct3: 0b_011},
	ASD:    {Opcode: OpBase_STORE, ArgMarks: ARG_SType, Funct3: 0b_011},
	AADDIW: {Opcode: OpBase_OP_IMM_32, ArgMarks: ARG_IType, Funct3: 0b_000},
	ASLLIW: {Opcode: OpBase_OP_IMM_32, ArgMarks: ARG_IType, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLIW: {Opcode: OpBase_OP_IMM_32, ArgMarks: ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAIW: {Opcode: OpBase_OP_IMM_32, ArgMarks: ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADDW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUBW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLLW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASRLW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRAW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_010_0000},

	// RV32/RV64 Zicsr Standard Extension

	ACSRRW:  {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_001},
	ACSRRS:  {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_010},
	ACSRRC:  {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_011},
	ACSRRWI: {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_101},
	ACSRRSI: {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_110},
	ACSRRCI: {Opcode: OpBase_SYSTEM, ArgMarks: ARG_IType, Funct3: 0b_111},

	// RV32M Standard Extension

	AMUL:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0001},
	AMULH:   {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0001},
	AMULHSU: {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_010, Funct7: 0b_000_0001},
	AMULHU:  {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_011, Funct7: 0b_000_0001},
	ADIV:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVU:   {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREM:    {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMU:   {Opcode: OpBase_OP, ArgMarks: ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV64M Standard Extension (in addition to RV32M)

	AMULW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0001},
	ADIVW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVUW: {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREMW:  {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMUW: {Opcode: OpBase_OP_32, ArgMarks: ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV32F Standard Extension

	AFLW:     {Opcode: OpBase_LOAD_FP, ArgMarks: ARG_IType, Funct3: 0b_010},
	AFSW:     {Opcode: OpBase_STORE_FP, ArgMarks: ARG_SType, Funct3: 0b_010},
	AFMADDS:  {Opcode: OpBase_MADD, ArgMarks: ARG_R4Type, Funct7: 0b_00},  // funct2
	AFMSUBS:  {Opcode: OpBase_MSUB, ArgMarks: ARG_R4Type, Funct7: 0b_00},  // funct2
	AFNMSUBS: {Opcode: OpBase_NMADD, ArgMarks: ARG_R4Type, Funct7: 0b_00}, // funct2
	AFNMADDS: {Opcode: OpBase_NMSUB, ArgMarks: ARG_R4Type, Funct7: 0b_00}, // funct2
	AFADDS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_0000},
	AFSUBS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_0100},
	AFMULS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_1000},
	AFDIVS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_1100},
	AFSQRTS:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_1100, Rs2: newU32(0b_0_0000)},
	AFSGNJS:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0000},
	AFSGNJNS: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0000},
	AFSGNJXS: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0000},
	AFMINS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0000)},
	AFMAXS:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0001)},
	AFCVTWS:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0000)},
	AFCVTWUS: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0000},
	AFMVXW:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_0000},
	AFEQS:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0000},
	AFLTS:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0000},
	AFLES:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0000},
	AFCLASSS: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_0000, Rs2: newU32(0b_0_0000)},
	AFCVTSW:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0000)},
	AFCVTSWU: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0001)},
	AFMVWX:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_1000, Rs2: newU32(0b_0_0000)},

	// RV64F Standard Extension (in addition to RV32F)

	AFCVTLS:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0010)},
	AFCVTLUS: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0011)},
	AFCVTSL:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0010)},
	AFCVTSLU: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0011)},

	// RV32D Standard Extension

	AFLD:     {Opcode: OpBase_LOAD_FP, ArgMarks: ARG_IType, Funct3: 0b_011},
	AFSD:     {Opcode: OpBase_STORE_FP, ArgMarks: ARG_SType, Funct3: 0b_011},
	AFMADDD:  {Opcode: OpBase_MADD, ArgMarks: ARG_R4Type, Funct7: 0b_00},  // funct2
	AFMSUBD:  {Opcode: OpBase_MSUB, ArgMarks: ARG_R4Type, Funct7: 0b_00},  // funct2
	AFNMSUBD: {Opcode: OpBase_NMADD, ArgMarks: ARG_R4Type, Funct7: 0b_00}, // funct2
	AFNMADDD: {Opcode: OpBase_NMSUB, ArgMarks: ARG_R4Type, Funct7: 0b_00}, // funct2
	AFADDD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_0001},
	AFSUBD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_0101},
	AFMULD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_1001},
	AFDIVD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_000_1101},
	AFSQRTD:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_010_1101, Rs2: newU32(0b_0_0000)},
	AFSGNJD:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0001},
	AFSGNJND: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0001},
	AFSGNJXD: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0001},
	AFMIND:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0101},
	AFMAXD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_001_0101},
	AFCVTSD:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_010_0000, Rs2: newU32(0b_0_0001)},
	AFCVTDS:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_010_0001, Rs2: newU32(0b_0_0000)},
	AFEQD:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0001},
	AFLTD:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0001},
	AFLED:    {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_101_0001},
	AFCLASSD: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWD:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWUD: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0001)},
	AFCVTDW:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0000)},
	AFCVTDWU: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0001)},

	// RV64D Standard Extension (in addition to RV32D)

	AFCVTLD:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0010)},
	AFCVTLUD: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0011)},
	AFMVXD:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTDL:  {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0010)},
	AFCVTDLU: {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0011)},
	AFMVDX:   {Opcode: OpBase_OP_FP, ArgMarks: ARG_RType, Funct7: 0b_111_1001, Rs2: newU32(0b_0_0000)},

	// 伪指令
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook

	ANOP:       {PseudoAs: AADDI, ArgMarks: 0},                           // nop                        => addi     x0, x0, 0
	AMV:        {PseudoAs: AADDI, ArgMarks: ARG_RD | ARG_RS1},            // mv        rd, rs1          => addi     rd, rs1, 0
	ANOT:       {PseudoAs: AXORI, ArgMarks: ARG_RD | ARG_RS1},            // not       rd, rs1          => xori     rd, rs1, -1
	ANEG:       {PseudoAs: ASUB, ArgMarks: ARG_RD | ARG_RS1},             // neg       rd, rs1          => sub      rd, x0, rs1
	ANEGW:      {PseudoAs: ASUBW, ArgMarks: ARG_RD | ARG_RS1},            // negw      rd, rs1          => subw     rd, x0, rs1
	ASEXT_W:    {PseudoAs: AADDIW, ArgMarks: ARG_RD | ARG_RS1},           // sext.w    rs, rs1          => addiw    rd, rs1, 0
	ASEQZ:      {PseudoAs: ASLTIU, ArgMarks: ARG_RD | ARG_RS1},           // seqz      rd, rs1          => sltiu    rd, rs1, 1
	ASNEZ:      {PseudoAs: ASLTU, ArgMarks: ARG_RD | ARG_RS1},            // snez      rd, rs1          => sltu     rd, x0, rs1
	ASLTZ:      {PseudoAs: ASLT, ArgMarks: ARG_RD | ARG_RS1},             // sltz      rd, rs1          => slt      rd, rs1, x0
	ASGTZ:      {PseudoAs: ASLT, ArgMarks: ARG_RD | ARG_RS1},             // sgtz      rd, rs1          => slt      rd, x0, rs1
	AFMV_S:     {PseudoAs: AFSGNJS, ArgMarks: ARG_RD | ARG_RS1},          // fmv.s     rd, rs1          => fsgnj.s  rd, rs1, rs1
	AFABS_S:    {PseudoAs: AFSGNJXS, ArgMarks: ARG_RD | ARG_RS1},         // fabc.s    rd, rs1          => fsgnjx.s rd, rs1, rs1
	AFNEG_S:    {PseudoAs: AFSGNJNS, ArgMarks: ARG_RD | ARG_RS1},         // fneg.s    rd, rs1          => fsgnjn.s rd, rs1, rs1
	AFMV_D:     {PseudoAs: AFSGNJD, ArgMarks: ARG_RD | ARG_RS1},          // fmv.d     rd, rs1          => fsgnj.d  rd, rs1, rs1
	AFABS_D:    {PseudoAs: AFSGNJXD, ArgMarks: ARG_RD | ARG_RS1},         // fabs.d    rd, rs1          => fsgnjx.d rd, rs1, rs1
	AFNEG_D:    {PseudoAs: AFSGNJND, ArgMarks: ARG_RD | ARG_RS1},         // fneg.d    rd, rs1          => fsgnjn.d rd, rs1, rs1
	ABEQZ:      {PseudoAs: ABEQ, ArgMarks: ARG_RS1 | ARG_IMM},            // beqz      rs1, offset      => beq      rs1, x0, offset
	ABNEZ:      {PseudoAs: ABNE, ArgMarks: ARG_RS1 | ARG_IMM},            // bnez      rs1, offset      => bne      rs1, x0, offset
	ABLEZ:      {PseudoAs: ABGE, ArgMarks: ARG_RS1 | ARG_IMM},            // blez      rs1, offset      => bge      x0, rs1, offset
	ABGEZ:      {PseudoAs: ABGE, ArgMarks: ARG_RS1 | ARG_IMM},            // bgez      rs1, offset      => bge      rs1, x0, offset
	ABLTZ:      {PseudoAs: ABLT, ArgMarks: ARG_RS1 | ARG_IMM},            // bltz      rs1, offset      => blt      rs1, x0, offset
	ABGTZ:      {PseudoAs: ABLT, ArgMarks: ARG_RS1 | ARG_IMM},            // bgtz      rs1, offset      => blt      x0, rs1, offset
	ABGT:       {PseudoAs: ABLT, ArgMarks: ARG_RS1 | ARG_RS2 | ARG_IMM},  // bgt       rs1, rs2, offset => blt      rs2, rs1, offset
	ABLE:       {PseudoAs: ABGE, ArgMarks: ARG_RS1 | ARG_RS2 | ARG_IMM},  // ble       rs1, rs2, offset => bge      rs2, rs1, offset
	ABGTU:      {PseudoAs: ABLTU, ArgMarks: ARG_RS1 | ARG_RS2 | ARG_IMM}, // bgtu      rs1, rs2, offset => bltu     rs2, rs1, offset
	ABLEU:      {PseudoAs: ABGEU, ArgMarks: ARG_RS1 | ARG_RS2 | ARG_IMM}, // bleu      rs1, rs2, offset => bgeu     rs2, rs1, offset
	AJ:         {PseudoAs: AJAL, ArgMarks: ARG_IMM},                      // j         offset           => jal      x0, offset
	AJR:        {PseudoAs: AJALR, ArgMarks: ARG_RS1},                     // jr        rs1              => jalr     x0, 0(rs1)
	ARET:       {PseudoAs: AJALR, ArgMarks: 0},                           // ret                        => jalr     x0, 0(x1)
	ARDINSTRET: {PseudoAs: ACSRRW, ArgMarks: ARG_RD},                     // rdinstret rd               => csrrs    rd, instret, x0
	ARDCYCLE:   {PseudoAs: ACSRRW, ArgMarks: ARG_RD},                     // rdcyle    rd               => csrrs    rd, cycle, x0
	ARDTIME:    {PseudoAs: ACSRRW, ArgMarks: ARG_RD},                     // rdtime    rd               => csrrs    rd, time, x0
	ACSRR:      {PseudoAs: ACSRRW, ArgMarks: ARG_RD | ARG_IMM},           // csrr      rd, csr          => csrrs    rd, csr, x0
	ACSRW:      {PseudoAs: ACSRRW, ArgMarks: ARG_RS1 | ARG_IMM},          // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRS:      {PseudoAs: ACSRRW, ArgMarks: ARG_RS1 | ARG_IMM},          // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRC:      {PseudoAs: ACSRRW, ArgMarks: ARG_RS1 | ARG_IMM},          // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRWI:     {PseudoAs: ACSRRWI, ArgMarks: ARG_RS1 | ARG_IMM},         // csrwi     csr, imm         => csrrwi   x0 csr, imm
	ACSRSI:     {PseudoAs: ACSRRSI, ArgMarks: ARG_RS1 | ARG_IMM},         // csrsi     csr, imm         => csrrsi   x0 csr, imm
	ACSRCI:     {PseudoAs: ACSRRCI, ArgMarks: ARG_RS1 | ARG_IMM},         // csrci     csr, imm         => csrrci   x0 csr, imm
	AFRCSR:     {PseudoAs: ACSRRS, ArgMarks: ARG_RD},                     // frcsr     rd               => csrrs    rd, fcsr
	AFSCSR:     {PseudoAs: ACSRRW, ArgMarks: ARG_RD | ARG_RS1},           // fscsr     rd, rs1          => csrrw    rd, fcsr, rs1 # rd 可省略
	AFRRM:      {PseudoAs: ACSRRS, ArgMarks: ARG_RD},                     // frrm      rd               => csrrs    rd, frm, x0
	AFSRM:      {PseudoAs: ACSRRW, ArgMarks: ARG_RD | ARG_RS1},           // fsrm      rd, rs1          => csrrw    rd, frm, rs1 # rd 可省略
	AFRFLAGS:   {PseudoAs: ACSRRS, ArgMarks: ARG_RD},                     // frflags   rd               => csrrs    rd, fflags, x0
	AFSFLAGS:   {PseudoAs: ACSRRW, ArgMarks: ARG_RD | ARG_RS1},           // fsflags   rd, rs1          => csrrw    rd, fflags, rs1 # rd 可省略

	// End marker

	ALAST: {},
}
