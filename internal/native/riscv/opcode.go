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

func (opFormatType OpFormatType) ArgMarks() ArgMarks {
	switch opFormatType {
	case R:
		return ARG_RType
	case R4:
		return ARG_R4Type
	case I:
		return ARG_IType
	case S:
		return ARG_SType
	case B:
		return ARG_BType
	case U:
		return ARG_UType
	case J:
		return ARG_JType
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

	ARG_RType  = ARG_RD | ARG_RS1 | ARG_RS2 | ARG_FUNCT3
	ARG_R4Type = ARG_RType | ARG_RS3 | ARG_FUNCT2
	ARG_IType  = ARG_RD | ARG_RS1 | ARG_IMM | ARG_FUNCT3
	ARG_SType  = ARG_RS1 | ARG_RS2 | ARG_IMM | ARG_FUNCT3
	ARG_BType  = ARG_RS1 | ARG_RS2 | ARG_IMM | ARG_FUNCT3
	ARG_UType  = ARG_RD | ARG_IMM
	ARG_JType  = ARG_RD | ARG_IMM
)

// 操作码上下文信息
type OpContextType struct {
	Opcode         OpcodeType
	Funct3         uint32
	Funct7         uint32   // 和 Funct2 共用
	Rs2            *uint32  // 是否用到 Rs2 模板
	HasShamt       bool     // 是否有 shamt 参数. SLLI/SRLI/SRAI 在 RV32 是 5bit, RV64 是 6bit
	PseudoAs       abi.As   // 伪指令对应的原生指令
	PseudoArgMarks ArgMarks // 伪指令的参数
}

// 指令编码信息表
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
var AOpContextTable = []OpContextType{
	// RV32I Base Instruction Set

	ALUI:    {Opcode: OpBase_LUI},
	AAUIPC:  {Opcode: OpBase_AUIPC},
	AJAL:    {Opcode: OpBase_JAL},
	AJALR:   {Opcode: OpBase_JALR, Funct3: 0b_000},
	ABEQ:    {Opcode: OpBase_BRANCH, Funct3: 0b_000},
	ABNE:    {Opcode: OpBase_BRANCH, Funct3: 0b_001},
	ABLT:    {Opcode: OpBase_BRANCH, Funct3: 0b_100},
	ABGE:    {Opcode: OpBase_BRANCH, Funct3: 0b_101},
	ABLTU:   {Opcode: OpBase_BRANCH, Funct3: 0b_110},
	ABGEU:   {Opcode: OpBase_BRANCH, Funct3: 0b_111},
	ALB:     {Opcode: OpBase_LOAD, Funct3: 0b_000},
	ALH:     {Opcode: OpBase_LOAD, Funct3: 0b_001},
	ALW:     {Opcode: OpBase_LOAD, Funct3: 0b_010},
	ALBU:    {Opcode: OpBase_LOAD, Funct3: 0b_100},
	ALHU:    {Opcode: OpBase_LOAD, Funct3: 0b_101},
	ASB:     {Opcode: OpBase_STORE, Funct3: 0b_000},
	ASH:     {Opcode: OpBase_STORE, Funct3: 0b_001},
	ASW:     {Opcode: OpBase_STORE, Funct3: 0b_010},
	AADDI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_000},
	ASLTI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_010},
	ASLTIU:  {Opcode: OpBase_OP_IMM, Funct3: 0b_011},
	AXORI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_100},
	AORI:    {Opcode: OpBase_OP_IMM, Funct3: 0b_110},
	AANDI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_111},
	ASLLI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAI:   {Opcode: OpBase_OP_IMM, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADD:    {Opcode: OpBase_OP, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUB:    {Opcode: OpBase_OP, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLL:    {Opcode: OpBase_OP, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASLT:    {Opcode: OpBase_OP, Funct3: 0b_010, Funct7: 0b_000_0000},
	ASLTU:   {Opcode: OpBase_OP, Funct3: 0b_011, Funct7: 0b_000_0000},
	AXOR:    {Opcode: OpBase_OP, Funct3: 0b_100, Funct7: 0b_000_0000},
	ASRL:    {Opcode: OpBase_OP, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRA:    {Opcode: OpBase_OP, Funct3: 0b_101, Funct7: 0b_010_0000},
	AOR:     {Opcode: OpBase_OP, Funct3: 0b_110, Funct7: 0b_000_0000},
	AAND:    {Opcode: OpBase_OP, Funct3: 0b_111, Funct7: 0b_000_0000},
	AFENCE:  {Opcode: OpBase_MISC_MEN, Funct3: 0b_000},
	AECALL:  {Opcode: OpBase_SYSTEM, Funct3: 0b_000}, // imm[11:0] = 0b000000000000
	AEBREAK: {Opcode: OpBase_SYSTEM, Funct3: 0b_000}, // imm[11:0] = 0b000000000001

	// RV64I Base Instruction Set (in addition to RV32I)

	ALWU:   {Opcode: OpBase_LOAD, Funct3: 0b_110},
	ALD:    {Opcode: OpBase_LOAD, Funct3: 0b_011},
	ASD:    {Opcode: OpBase_STORE, Funct3: 0b_011},
	AADDIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b_000},
	ASLLIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADDW:  {Opcode: OpBase_OP_32, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUBW:  {Opcode: OpBase_OP_32, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLLW:  {Opcode: OpBase_OP_32, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASRLW:  {Opcode: OpBase_OP_32, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRAW:  {Opcode: OpBase_OP_32, Funct3: 0b_101, Funct7: 0b_010_0000},

	// RV32/RV64 Zicsr Standard Extension

	ACSRRW:  {Opcode: OpBase_SYSTEM, Funct3: 0b_001},
	ACSRRS:  {Opcode: OpBase_SYSTEM, Funct3: 0b_010},
	ACSRRC:  {Opcode: OpBase_SYSTEM, Funct3: 0b_011},
	ACSRRWI: {Opcode: OpBase_SYSTEM, Funct3: 0b_101},
	ACSRRSI: {Opcode: OpBase_SYSTEM, Funct3: 0b_110},
	ACSRRCI: {Opcode: OpBase_SYSTEM, Funct3: 0b_111},

	// RV32M Standard Extension

	AMUL:    {Opcode: OpBase_OP, Funct3: 0b_000, Funct7: 0b_000_0001},
	AMULH:   {Opcode: OpBase_OP, Funct3: 0b_001, Funct7: 0b_000_0001},
	AMULHSU: {Opcode: OpBase_OP, Funct3: 0b_010, Funct7: 0b_000_0001},
	AMULHU:  {Opcode: OpBase_OP, Funct3: 0b_011, Funct7: 0b_000_0001},
	ADIV:    {Opcode: OpBase_OP, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVU:   {Opcode: OpBase_OP, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREM:    {Opcode: OpBase_OP, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMU:   {Opcode: OpBase_OP, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV64M Standard Extension (in addition to RV32M)

	AMULW:  {Opcode: OpBase_OP_32, Funct3: 0b_000, Funct7: 0b_000_0001},
	ADIVW:  {Opcode: OpBase_OP_32, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVUW: {Opcode: OpBase_OP_32, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREMW:  {Opcode: OpBase_OP_32, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMUW: {Opcode: OpBase_OP_32, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV32F Standard Extension

	AFLW:     {Opcode: OpBase_LOAD_FP, Funct3: 0b_010},
	AFSW:     {Opcode: OpBase_STORE_FP, Funct3: 0b_010},
	AFMADDS:  {Opcode: OpBase_MADD, Funct7: 0b_00},  // funct2
	AFMSUBS:  {Opcode: OpBase_MSUB, Funct7: 0b_00},  // funct2
	AFNMSUBS: {Opcode: OpBase_NMADD, Funct7: 0b_00}, // funct2
	AFNMADDS: {Opcode: OpBase_NMSUB, Funct7: 0b_00}, // funct2
	AFADDS:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_0000},
	AFSUBS:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_0100},
	AFMULS:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_1000},
	AFDIVS:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_1100},
	AFSQRTS:  {Opcode: OpBase_OP_FP, Funct7: 0b_000_1100, Rs2: newU32(0b_0_0000)},
	AFSGNJS:  {Opcode: OpBase_OP_FP, Funct7: 0b_001_0000},
	AFSGNJNS: {Opcode: OpBase_OP_FP, Funct7: 0b_001_0000},
	AFSGNJXS: {Opcode: OpBase_OP_FP, Funct7: 0b_001_0000},
	AFMINS:   {Opcode: OpBase_OP_FP, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0000)},
	AFMAXS:   {Opcode: OpBase_OP_FP, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0001)},
	AFCVTWS:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0000)},
	AFCVTWUS: {Opcode: OpBase_OP_FP, Funct7: 0b_110_0000},
	AFMVXW:   {Opcode: OpBase_OP_FP, Funct7: 0b_111_0000},
	AFEQS:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0000},
	AFLTS:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0000},
	AFLES:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0000},
	AFCLASSS: {Opcode: OpBase_OP_FP, Funct7: 0b_111_0000, Rs2: newU32(0b_0_0000)},
	AFCVTSW:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0000)},
	AFCVTSWU: {Opcode: OpBase_OP_FP, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0001)},
	AFMVWX:   {Opcode: OpBase_OP_FP, Funct7: 0b_111_1000, Rs2: newU32(0b_0_0000)},

	// RV64F Standard Extension (in addition to RV32F)

	AFCVTLS:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0010)},
	AFCVTLUS: {Opcode: OpBase_OP_FP, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0011)},
	AFCVTSL:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0010)},
	AFCVTSLU: {Opcode: OpBase_OP_FP, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0011)},

	// RV32D Standard Extension

	AFLD:     {Opcode: OpBase_LOAD_FP, Funct3: 0b_011},
	AFSD:     {Opcode: OpBase_STORE_FP, Funct3: 0b_011},
	AFMADDD:  {Opcode: OpBase_MADD, Funct7: 0b_00},  // funct2
	AFMSUBD:  {Opcode: OpBase_MSUB, Funct7: 0b_00},  // funct2
	AFNMSUBD: {Opcode: OpBase_NMADD, Funct7: 0b_00}, // funct2
	AFNMADDD: {Opcode: OpBase_NMSUB, Funct7: 0b_00}, // funct2
	AFADDD:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_0001},
	AFSUBD:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_0101},
	AFMULD:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_1001},
	AFDIVD:   {Opcode: OpBase_OP_FP, Funct7: 0b_000_1101},
	AFSQRTD:  {Opcode: OpBase_OP_FP, Funct7: 0b_010_1101, Rs2: newU32(0b_0_0000)},
	AFSGNJD:  {Opcode: OpBase_OP_FP, Funct7: 0b_001_0001},
	AFSGNJND: {Opcode: OpBase_OP_FP, Funct7: 0b_001_0001},
	AFSGNJXD: {Opcode: OpBase_OP_FP, Funct7: 0b_001_0001},
	AFMIND:   {Opcode: OpBase_OP_FP, Funct7: 0b_001_0101},
	AFMAXD:   {Opcode: OpBase_OP_FP, Funct7: 0b_001_0101},
	AFCVTSD:  {Opcode: OpBase_OP_FP, Funct7: 0b_010_0000, Rs2: newU32(0b_0_0001)},
	AFCVTDS:  {Opcode: OpBase_OP_FP, Funct7: 0b_010_0001, Rs2: newU32(0b_0_0000)},
	AFEQD:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0001},
	AFLTD:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0001},
	AFLED:    {Opcode: OpBase_OP_FP, Funct7: 0b_101_0001},
	AFCLASSD: {Opcode: OpBase_OP_FP, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWD:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWUD: {Opcode: OpBase_OP_FP, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0001)},
	AFCVTDW:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0000)},
	AFCVTDWU: {Opcode: OpBase_OP_FP, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0001)},

	// RV64D Standard Extension (in addition to RV32D)

	AFCVTLD:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0010)},
	AFCVTLUD: {Opcode: OpBase_OP_FP, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0011)},
	AFMVXD:   {Opcode: OpBase_OP_FP, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTDL:  {Opcode: OpBase_OP_FP, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0010)},
	AFCVTDLU: {Opcode: OpBase_OP_FP, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0011)},
	AFMVDX:   {Opcode: OpBase_OP_FP, Funct7: 0b_111_1001, Rs2: newU32(0b_0_0000)},

	// 伪指令
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook

	// TODO: 补充 PseudoArgMarks

	ANOP:       {PseudoAs: AADDI, PseudoArgMarks: 0},    // nop                      => addi     x0, x0, 0
	AMV:        {PseudoAs: AADDI, PseudoArgMarks: 0},    // mv        rd, rs         => addi     rd, rs, 0
	ANOT:       {PseudoAs: AXORI, PseudoArgMarks: 0},    // not       rd, rs         => xori     rd, rs, -1
	ANEG:       {PseudoAs: ASUB, PseudoArgMarks: 0},     // neg       rd, rs         => sub      rd, x0, rs
	ANEGW:      {PseudoAs: ASUBW, PseudoArgMarks: 0},    // negw      rd, rs         => subw     rd, x0, rs
	ASEXT_W:    {PseudoAs: AADDIW, PseudoArgMarks: 0},   // sext.w    rs, rs         => addiw    rd, rs, 0
	ASEQZ:      {PseudoAs: ASLTIU, PseudoArgMarks: 0},   // seqz      rd, rs         => sltiu    rd, rs, 1
	ASNEZ:      {PseudoAs: ASLTU, PseudoArgMarks: 0},    // snez      rd, rs         => sltu     rd, x0, rs
	ASLTZ:      {PseudoAs: ASLT, PseudoArgMarks: 0},     // sltz      rd, rs         => slt      rd, rs, x0
	ASGTZ:      {PseudoAs: ASLT, PseudoArgMarks: 0},     // sgtz      rd, rs         => slt      rd, x0, rs
	AFMV_S:     {PseudoAs: AFSGNJS, PseudoArgMarks: 0},  // fmv.s     rd, rs         => fsgnj.s  rd, rs, rs
	AFABS_S:    {PseudoAs: AFSGNJXS, PseudoArgMarks: 0}, // fabc.s    rd, rs         => fsgnjx.s rd, rs, rs
	AFNEG_S:    {PseudoAs: AFSGNJNS, PseudoArgMarks: 0}, // fneg.s    rd, rs         => fsgnjn.s rd, rs, rs
	AFMV_D:     {PseudoAs: AFSGNJD, PseudoArgMarks: 0},  // fmv.d     rd, rs         => fsgnj.d  rd, rs, rs
	AFABS_D:    {PseudoAs: AFSGNJXD, PseudoArgMarks: 0}, // fabs.d    rd, rs         => fsgnjx.d rd, rs, rs
	AFNEG_D:    {PseudoAs: AFSGNJND, PseudoArgMarks: 0}, // fneg.d    rd, rs         => fsgnjn.d rd, rs, rs
	ABEQZ:      {PseudoAs: ABEQ, PseudoArgMarks: 0},     // beqz      rs, offset     => beq      rs, x0, pffset
	ABNEZ:      {PseudoAs: ABNE, PseudoArgMarks: 0},     // bnez      rs, offset     => bne      rs, x0, pffset
	ABLEZ:      {PseudoAs: ABGE, PseudoArgMarks: 0},     // blez      rs, offset     => bge      x0, rs, pffset
	ABGEZ:      {PseudoAs: ABGE, PseudoArgMarks: 0},     // bgez      rs, offset     => bge      rs, x0, pffset
	ABLTZ:      {PseudoAs: ABLT, PseudoArgMarks: 0},     // bltz      rs, offset     => blt      rs, x0, pffset
	ABGTZ:      {PseudoAs: ABLT, PseudoArgMarks: 0},     // bgtz      rs, offset     => blt      x0, rs, pffset
	ABGT:       {PseudoAs: ABLT, PseudoArgMarks: 0},     // bgt       rs, rt, offset => blt      rt, rs, offset
	ABLE:       {PseudoAs: ABGE, PseudoArgMarks: 0},     // ble       rs, rt, offset => bge      rt, rs, offset
	ABGTU:      {PseudoAs: ABLTU, PseudoArgMarks: 0},    // bgtu      rs, rt, offset => bltu     rt, rs, offset
	ABLEU:      {PseudoAs: ABGEU, PseudoArgMarks: 0},    // bleu      rs, rt, offset => bgeu     rt, rs, offset
	AJ:         {PseudoAs: AJAL, PseudoArgMarks: 0},     // j         offset         => jal      x0, offset
	AJR:        {PseudoAs: AJALR, PseudoArgMarks: 0},    // jr        rs             => jalr     x0, 0(rs)
	ARET:       {PseudoAs: AJALR, PseudoArgMarks: 0},    // ret                      => jalr     x0, 0(x1)
	ARDINSTRET: {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // rdinstret rd             => csrrs    rd, instret, x0
	ARDCYCLE:   {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // rdcyle    rd             => csrrs    rd, cycle, x0
	ARDTIME:    {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // rdtime    rd             => csrrs    rd, time, x0
	ACSRR:      {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // csrr      rd, csr        => csrrs    rd, csr, x0
	ACSRW:      {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // csrr      csr, rd        => csrrs    x0, csr, rs
	ACSRS:      {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // csrr      csr, rd        => csrrs    x0, csr, rs
	ACSRC:      {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // csrr      csr, rd        => csrrs    x0, csr, rs
	ACSRWI:     {PseudoAs: ACSRRWI, PseudoArgMarks: 0},  // csrwi     csr, imm       => csrrwi   x0 csr, imm
	ACSRSI:     {PseudoAs: ACSRRSI, PseudoArgMarks: 0},  // csrsi     csr, imm       => csrrsi   x0 csr, imm
	ACSRCI:     {PseudoAs: ACSRRCI, PseudoArgMarks: 0},  // csrci     csr, imm       => csrrci   x0 csr, imm
	AFRCSR:     {PseudoAs: ACSRRS, PseudoArgMarks: 0},   // frcsr     rd             => csrrs    rd, fcsr
	AFSCSR:     {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // fscsr     rd, rs         => csrrw    rd, fcsr, rs # rd 可省略
	AFRRM:      {PseudoAs: ACSRRS, PseudoArgMarks: 0},   // frrm      rd             => csrrs    rd, frm, x0
	AFSRM:      {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // fsrm      rd, rs         => csrrw    rd, frm, rs # rd 可省略
	AFRFLAGS:   {PseudoAs: ACSRRS, PseudoArgMarks: 0},   // frflags   rd             => csrrs    rd, fflags, x0
	AFSFLAGS:   {PseudoAs: ACSRRW, PseudoArgMarks: 0},   // fsflags   rd, rs         => csrrw    rd, fflags, rs # rd 可省略

	// End marker

	ALAST: {},
}
