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
type _OpFormatType int

const (
	_ _OpFormatType = iota
	_R
	_R4
	_I
	_S
	_B
	_U
	_J
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
type _OpcodeType uint32

const (
	_OpBase_LOAD      _OpcodeType = 0b_00_000_11
	_OpBase_LOAD_FP   _OpcodeType = 0b_00_001_11
	_OpBase_CUSTOM_0  _OpcodeType = 0b_00_010_11
	_OpBase_MISC_MEN  _OpcodeType = 0b_00_011_11
	_OpBase_OP_IMM    _OpcodeType = 0b_00_100_11
	_OpBase_AUIPC     _OpcodeType = 0b_00_101_11
	_OpBase_OP_IMM_32 _OpcodeType = 0b_00_110_11

	_OpBase_STORE    _OpcodeType = 0b_01_000_11
	_OpBase_STORE_FP _OpcodeType = 0b_01_001_11
	_OpBase_CUSTOM_1 _OpcodeType = 0b_01_010_11
	_OpBase_AMO      _OpcodeType = 0b_01_011_11
	_OpBase_OP       _OpcodeType = 0b_01_100_11
	_OpBase_LUI      _OpcodeType = 0b_01_101_11
	_OpBase_OP_32    _OpcodeType = 0b_01_110_11

	_OpBase_MADD     _OpcodeType = 0b_10_000_11
	_OpBase_MSUB     _OpcodeType = 0b_10_001_11
	_OpBase_NMSUB    _OpcodeType = 0b_10_010_11
	_OpBase_NMADD    _OpcodeType = 0b_10_011_11
	_OpBase_OP_FP    _OpcodeType = 0b_10_100_11
	_                _OpcodeType = 0b_10_101_11
	_OpBase_CUSTOM_2 _OpcodeType = 0b_10_110_11

	_OpBase_BRANCH   _OpcodeType = 0b_11_000_11
	_OpBase_JALR     _OpcodeType = 0b_11_001_11
	_                _OpcodeType = 0b_11_010_11
	_OpBase_JAL      _OpcodeType = 0b_11_011_11
	_OpBase_SYSTEM   _OpcodeType = 0b_11_100_11
	_                _OpcodeType = 0b_11_101_11
	_OpBase_CUSTOM_3 _OpcodeType = 0b_11_110_11

	_OpBase_Mask _OpcodeType = 0b_11_111_11
)

// 指令的编码格式
func (opcode _OpcodeType) FormatType() _OpFormatType {
	switch opcode & _OpBase_Mask {
	case _OpBase_OP,
		_OpBase_OP_32,
		_OpBase_OP_FP,
		_OpBase_AMO:
		return _R
	case _OpBase_MADD,
		_OpBase_MSUB,
		_OpBase_NMSUB,
		_OpBase_NMADD:
		return _R4
	case _OpBase_OP_IMM,
		_OpBase_OP_IMM_32,
		_OpBase_JALR,
		_OpBase_LOAD,
		_OpBase_LOAD_FP,
		_OpBase_MISC_MEN,
		_OpBase_SYSTEM:
		return _I
	case _OpBase_STORE,
		_OpBase_STORE_FP:
		return _S
	case _OpBase_BRANCH:
		return _B
	case _OpBase_LUI,
		_OpBase_AUIPC:
		return _U
	case _OpBase_JAL:
		return _J

	case _OpBase_CUSTOM_0,
		_OpBase_CUSTOM_1,
		_OpBase_CUSTOM_2,
		_OpBase_CUSTOM_3:
		return 0

	default:
		return 0
	}
}

// 参数标志
type _ArgMarks uint16

const (
	_ARG_RD _ArgMarks = 1 << iota
	_ARG_RS1
	_ARG_RS2
	_ARG_RS3
	_ARG_IMM
	_ARG_FUNCT3
	_ARG_FUNCT7
	_ARG_FUNCT2
	_ARG_SHAMT

	_ARG_RType  = _ARG_RD | _ARG_RS1 | _ARG_RS2 | _ARG_FUNCT3 | _ARG_FUNCT7
	_ARG_R4Type = _ARG_RD | _ARG_RS1 | _ARG_RS2 | _ARG_RS3 | _ARG_FUNCT2
	_ARG_IType  = _ARG_RD | _ARG_RS1 | _ARG_IMM | _ARG_FUNCT3
	_ARG_SType  = _ARG_RS1 | _ARG_RS2 | _ARG_IMM | _ARG_FUNCT3
	_ARG_BType  = _ARG_RS1 | _ARG_RS2 | _ARG_IMM | _ARG_FUNCT3
	_ARG_UType  = _ARG_RD | _ARG_IMM
	_ARG_JType  = _ARG_RD | _ARG_IMM
)

// 操作码上下文信息
type _OpContextType struct {
	Opcode   _OpcodeType // 基础指令码
	ArgMarks _ArgMarks   // 指令的参数标志
	Funct3   uint32      // funct3
	Funct7   uint32      // 和 Funct2 共用
	Rs2      *uint32     // 是否用到 Rs2 模板
	HasShamt bool        // 是否有 shamt 参数. SLLI/SRLI/SRAI 在 RV32 是 5bit, RV64 是 6bit
	PseudoAs abi.As      // 伪指令对应的原生指令
}

// Imm 取值范围
type _ImmRange struct {
	Align int64
	Min   int64
	Max   int64
}

var (
	_ImmRanges_IType   = _ImmRange{Align: 1, Min: -(1 << 11), Max: (1 << 11) - 1}
	_ImmRanges_SType   = _ImmRange{Align: 1, Min: -(1 << 11), Max: (1 << 11) - 1}
	_ImmRanges_BType   = _ImmRange{Align: 1 << 1, Min: -(1 << 12), Max: (1 << 12) - 2}
	_ImmRanges_UType   = _ImmRange{Align: 1, Min: -(1 << 20), Max: (1 << 20) - 1}
	_ImmRanges_JType   = _ImmRange{Align: 1, Min: -(1 << 20), Max: (1 << 20) - 2}
	_ImmRanges_Shamt32 = _ImmRange{Align: 1, Min: 0, Max: (1 << 5) - 1}
	_ImmRanges_Shamt64 = _ImmRange{Align: 1, Min: 0, Max: (1 << 6) - 1}
)

// 指令编码信息表
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
var _AOpContextTable = []_OpContextType{
	// RV32I Base Instruction Set

	ALUI:    {Opcode: _OpBase_LUI, ArgMarks: _ARG_UType},
	AAUIPC:  {Opcode: _OpBase_AUIPC, ArgMarks: _ARG_UType},
	AJAL:    {Opcode: _OpBase_JAL, ArgMarks: _ARG_JType},
	AJALR:   {Opcode: _OpBase_JALR, ArgMarks: _ARG_IType, Funct3: 0b_000},
	ABEQ:    {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_000},
	ABNE:    {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_001},
	ABLT:    {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_100},
	ABGE:    {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_101},
	ABLTU:   {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_110},
	ABGEU:   {Opcode: _OpBase_BRANCH, ArgMarks: _ARG_BType, Funct3: 0b_111},
	ALB:     {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_000},
	ALH:     {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_001},
	ALW:     {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_010},
	ALBU:    {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_100},
	ALHU:    {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_101},
	ASB:     {Opcode: _OpBase_STORE, ArgMarks: _ARG_SType, Funct3: 0b_000},
	ASH:     {Opcode: _OpBase_STORE, ArgMarks: _ARG_SType, Funct3: 0b_001},
	ASW:     {Opcode: _OpBase_STORE, ArgMarks: _ARG_SType, Funct3: 0b_010},
	AADDI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_000},
	ASLTI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_010},
	ASLTIU:  {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_011},
	AXORI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_100},
	AORI:    {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_110},
	AANDI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_111},
	ASLLI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAI:   {Opcode: _OpBase_OP_IMM, ArgMarks: _ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADD:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUB:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLL:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASLT:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_010, Funct7: 0b_000_0000},
	ASLTU:   {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_011, Funct7: 0b_000_0000},
	AXOR:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0000},
	ASRL:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRA:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_010_0000},
	AOR:     {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0000},
	AAND:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0000},
	AFENCE:  {Opcode: _OpBase_MISC_MEN, ArgMarks: _ARG_IType, Funct3: 0b_000},
	AECALL:  {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_000}, // imm[11:0] = 0b000000000000
	AEBREAK: {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_000}, // imm[11:0] = 0b000000000001

	// RV64I Base Instruction Set (in addition to RV32I)

	ALWU:   {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_110},
	ALD:    {Opcode: _OpBase_LOAD, ArgMarks: _ARG_IType, Funct3: 0b_011},
	ASD:    {Opcode: _OpBase_STORE, ArgMarks: _ARG_SType, Funct3: 0b_011},
	AADDIW: {Opcode: _OpBase_OP_IMM_32, ArgMarks: _ARG_IType, Funct3: 0b_000},
	ASLLIW: {Opcode: _OpBase_OP_IMM_32, ArgMarks: _ARG_IType, Funct3: 0b_001, HasShamt: true, Funct7: 0b_000_0000},
	ASRLIW: {Opcode: _OpBase_OP_IMM_32, ArgMarks: _ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_000_0000},
	ASRAIW: {Opcode: _OpBase_OP_IMM_32, ArgMarks: _ARG_IType, Funct3: 0b_101, HasShamt: true, Funct7: 0b_010_0000},
	AADDW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0000},
	ASUBW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_010_0000},
	ASLLW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0000},
	ASRLW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0000},
	ASRAW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_010_0000},

	// RV32/RV64 Zicsr Standard Extension

	ACSRRW:  {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_001},
	ACSRRS:  {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_010},
	ACSRRC:  {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_011},
	ACSRRWI: {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_101},
	ACSRRSI: {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_110},
	ACSRRCI: {Opcode: _OpBase_SYSTEM, ArgMarks: _ARG_IType, Funct3: 0b_111},

	// RV32M Standard Extension

	AMUL:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0001},
	AMULH:   {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_001, Funct7: 0b_000_0001},
	AMULHSU: {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_010, Funct7: 0b_000_0001},
	AMULHU:  {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_011, Funct7: 0b_000_0001},
	ADIV:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVU:   {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREM:    {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMU:   {Opcode: _OpBase_OP, ArgMarks: _ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV64M Standard Extension (in addition to RV32M)

	AMULW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_000, Funct7: 0b_000_0001},
	ADIVW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_100, Funct7: 0b_000_0001},
	ADIVUW: {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_101, Funct7: 0b_000_0001},
	AREMW:  {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_110, Funct7: 0b_000_0001},
	AREMUW: {Opcode: _OpBase_OP_32, ArgMarks: _ARG_RType, Funct3: 0b_111, Funct7: 0b_000_0001},

	// RV32F Standard Extension

	AFLW:     {Opcode: _OpBase_LOAD_FP, ArgMarks: _ARG_IType, Funct3: 0b_010},
	AFSW:     {Opcode: _OpBase_STORE_FP, ArgMarks: _ARG_SType, Funct3: 0b_010},
	AFMADDS:  {Opcode: _OpBase_MADD, ArgMarks: _ARG_R4Type, Funct7: 0b_00},  // funct2
	AFMSUBS:  {Opcode: _OpBase_MSUB, ArgMarks: _ARG_R4Type, Funct7: 0b_00},  // funct2
	AFNMSUBS: {Opcode: _OpBase_NMADD, ArgMarks: _ARG_R4Type, Funct7: 0b_00}, // funct2
	AFNMADDS: {Opcode: _OpBase_NMSUB, ArgMarks: _ARG_R4Type, Funct7: 0b_00}, // funct2
	AFADDS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_0000},
	AFSUBS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_0100},
	AFMULS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_1000},
	AFDIVS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_1100},
	AFSQRTS:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_1100, Rs2: newU32(0b_0_0000)},
	AFSGNJS:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0000},
	AFSGNJNS: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0000},
	AFSGNJXS: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0000},
	AFMINS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0000)},
	AFMAXS:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0100, Rs2: newU32(0b_0_0001)},
	AFCVTWS:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0000)},
	AFCVTWUS: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0000},
	AFMVXW:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_0000},
	AFEQS:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0000},
	AFLTS:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0000},
	AFLES:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0000},
	AFCLASSS: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_0000, Rs2: newU32(0b_0_0000)},
	AFCVTSW:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0000)},
	AFCVTSWU: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0001)},
	AFMVWX:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_1000, Rs2: newU32(0b_0_0000)},

	// RV64F Standard Extension (in addition to RV32F)

	AFCVTLS:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0010)},
	AFCVTLUS: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0000, Rs2: newU32(0b_0_0011)},
	AFCVTSL:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0010)},
	AFCVTSLU: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1000, Rs2: newU32(0b_0_0011)},

	// RV32D Standard Extension

	AFLD:     {Opcode: _OpBase_LOAD_FP, ArgMarks: _ARG_IType, Funct3: 0b_011},
	AFSD:     {Opcode: _OpBase_STORE_FP, ArgMarks: _ARG_SType, Funct3: 0b_011},
	AFMADDD:  {Opcode: _OpBase_MADD, ArgMarks: _ARG_R4Type, Funct7: 0b_00},  // funct2
	AFMSUBD:  {Opcode: _OpBase_MSUB, ArgMarks: _ARG_R4Type, Funct7: 0b_00},  // funct2
	AFNMSUBD: {Opcode: _OpBase_NMADD, ArgMarks: _ARG_R4Type, Funct7: 0b_00}, // funct2
	AFNMADDD: {Opcode: _OpBase_NMSUB, ArgMarks: _ARG_R4Type, Funct7: 0b_00}, // funct2
	AFADDD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_0001},
	AFSUBD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_0101},
	AFMULD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_1001},
	AFDIVD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_000_1101},
	AFSQRTD:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_010_1101, Rs2: newU32(0b_0_0000)},
	AFSGNJD:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0001},
	AFSGNJND: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0001},
	AFSGNJXD: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0001},
	AFMIND:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0101},
	AFMAXD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_001_0101},
	AFCVTSD:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_010_0000, Rs2: newU32(0b_0_0001)},
	AFCVTDS:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_010_0001, Rs2: newU32(0b_0_0000)},
	AFEQD:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0001},
	AFLTD:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0001},
	AFLED:    {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_101_0001},
	AFCLASSD: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWD:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0000)},
	AFCVTWUD: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0001)},
	AFCVTDW:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0000)},
	AFCVTDWU: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0001)},

	// RV64D Standard Extension (in addition to RV32D)

	AFCVTLD:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0010)},
	AFCVTLUD: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_0001, Rs2: newU32(0b_0_0011)},
	AFMVXD:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_0001, Rs2: newU32(0b_0_0000)},
	AFCVTDL:  {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0010)},
	AFCVTDLU: {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_110_1001, Rs2: newU32(0b_0_0011)},
	AFMVDX:   {Opcode: _OpBase_OP_FP, ArgMarks: _ARG_RType, Funct7: 0b_111_1001, Rs2: newU32(0b_0_0000)},

	// 伪指令
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook

	ANOP:       {PseudoAs: AADDI, ArgMarks: 0},                              // nop                        => addi     x0, x0, 0
	AMV:        {PseudoAs: AADDI, ArgMarks: _ARG_RD | _ARG_RS1},             // mv        rd, rs1          => addi     rd, rs1, 0
	ANOT:       {PseudoAs: AXORI, ArgMarks: _ARG_RD | _ARG_RS1},             // not       rd, rs1          => xori     rd, rs1, -1
	ANEG:       {PseudoAs: ASUB, ArgMarks: _ARG_RD | _ARG_RS1},              // neg       rd, rs1          => sub      rd, x0, rs1
	ANEGW:      {PseudoAs: ASUBW, ArgMarks: _ARG_RD | _ARG_RS1},             // negw      rd, rs1          => subw     rd, x0, rs1
	ASEXT_W:    {PseudoAs: AADDIW, ArgMarks: _ARG_RD | _ARG_RS1},            // sext.w    rs, rs1          => addiw    rd, rs1, 0
	ASEQZ:      {PseudoAs: ASLTIU, ArgMarks: _ARG_RD | _ARG_RS1},            // seqz      rd, rs1          => sltiu    rd, rs1, 1
	ASNEZ:      {PseudoAs: ASLTU, ArgMarks: _ARG_RD | _ARG_RS1},             // snez      rd, rs1          => sltu     rd, x0, rs1
	ASLTZ:      {PseudoAs: ASLT, ArgMarks: _ARG_RD | _ARG_RS1},              // sltz      rd, rs1          => slt      rd, rs1, x0
	ASGTZ:      {PseudoAs: ASLT, ArgMarks: _ARG_RD | _ARG_RS1},              // sgtz      rd, rs1          => slt      rd, x0, rs1
	AFMV_S:     {PseudoAs: AFSGNJS, ArgMarks: _ARG_RD | _ARG_RS1},           // fmv.s     rd, rs1          => fsgnj.s  rd, rs1, rs1
	AFABS_S:    {PseudoAs: AFSGNJXS, ArgMarks: _ARG_RD | _ARG_RS1},          // fabc.s    rd, rs1          => fsgnjx.s rd, rs1, rs1
	AFNEG_S:    {PseudoAs: AFSGNJNS, ArgMarks: _ARG_RD | _ARG_RS1},          // fneg.s    rd, rs1          => fsgnjn.s rd, rs1, rs1
	AFMV_D:     {PseudoAs: AFSGNJD, ArgMarks: _ARG_RD | _ARG_RS1},           // fmv.d     rd, rs1          => fsgnj.d  rd, rs1, rs1
	AFABS_D:    {PseudoAs: AFSGNJXD, ArgMarks: _ARG_RD | _ARG_RS1},          // fabs.d    rd, rs1          => fsgnjx.d rd, rs1, rs1
	AFNEG_D:    {PseudoAs: AFSGNJND, ArgMarks: _ARG_RD | _ARG_RS1},          // fneg.d    rd, rs1          => fsgnjn.d rd, rs1, rs1
	ABEQZ:      {PseudoAs: ABEQ, ArgMarks: _ARG_RS1 | _ARG_IMM},             // beqz      rs1, offset      => beq      rs1, x0, offset
	ABNEZ:      {PseudoAs: ABNE, ArgMarks: _ARG_RS1 | _ARG_IMM},             // bnez      rs1, offset      => bne      rs1, x0, offset
	ABLEZ:      {PseudoAs: ABGE, ArgMarks: _ARG_RS1 | _ARG_IMM},             // blez      rs1, offset      => bge      x0, rs1, offset
	ABGEZ:      {PseudoAs: ABGE, ArgMarks: _ARG_RS1 | _ARG_IMM},             // bgez      rs1, offset      => bge      rs1, x0, offset
	ABLTZ:      {PseudoAs: ABLT, ArgMarks: _ARG_RS1 | _ARG_IMM},             // bltz      rs1, offset      => blt      rs1, x0, offset
	ABGTZ:      {PseudoAs: ABLT, ArgMarks: _ARG_RS1 | _ARG_IMM},             // bgtz      rs1, offset      => blt      x0, rs1, offset
	ABGT:       {PseudoAs: ABLT, ArgMarks: _ARG_RS1 | _ARG_RS2 | _ARG_IMM},  // bgt       rs1, rs2, offset => blt      rs2, rs1, offset
	ABLE:       {PseudoAs: ABGE, ArgMarks: _ARG_RS1 | _ARG_RS2 | _ARG_IMM},  // ble       rs1, rs2, offset => bge      rs2, rs1, offset
	ABGTU:      {PseudoAs: ABLTU, ArgMarks: _ARG_RS1 | _ARG_RS2 | _ARG_IMM}, // bgtu      rs1, rs2, offset => bltu     rs2, rs1, offset
	ABLEU:      {PseudoAs: ABGEU, ArgMarks: _ARG_RS1 | _ARG_RS2 | _ARG_IMM}, // bleu      rs1, rs2, offset => bgeu     rs2, rs1, offset
	AJ:         {PseudoAs: AJAL, ArgMarks: _ARG_IMM},                        // j         offset           => jal      x0, offset
	AJR:        {PseudoAs: AJALR, ArgMarks: _ARG_RS1},                       // jr        rs1              => jalr     x0, 0(rs1)
	ARET:       {PseudoAs: AJALR, ArgMarks: 0},                              // ret                        => jalr     x0, 0(x1)
	ARDINSTRET: {PseudoAs: ACSRRW, ArgMarks: _ARG_RD},                       // rdinstret rd               => csrrs    rd, instret, x0
	ARDCYCLE:   {PseudoAs: ACSRRW, ArgMarks: _ARG_RD},                       // rdcyle    rd               => csrrs    rd, cycle, x0
	ARDTIME:    {PseudoAs: ACSRRW, ArgMarks: _ARG_RD},                       // rdtime    rd               => csrrs    rd, time, x0
	ACSRR:      {PseudoAs: ACSRRW, ArgMarks: _ARG_RD | _ARG_IMM},            // csrr      rd, csr          => csrrs    rd, csr, x0
	ACSRW:      {PseudoAs: ACSRRW, ArgMarks: _ARG_RS1 | _ARG_IMM},           // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRS:      {PseudoAs: ACSRRW, ArgMarks: _ARG_RS1 | _ARG_IMM},           // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRC:      {PseudoAs: ACSRRW, ArgMarks: _ARG_RS1 | _ARG_IMM},           // csrr      csr, rd          => csrrs    x0, csr, rs1
	ACSRWI:     {PseudoAs: ACSRRWI, ArgMarks: _ARG_RS1 | _ARG_IMM},          // csrwi     csr, imm         => csrrwi   x0 csr, imm
	ACSRSI:     {PseudoAs: ACSRRSI, ArgMarks: _ARG_RS1 | _ARG_IMM},          // csrsi     csr, imm         => csrrsi   x0 csr, imm
	ACSRCI:     {PseudoAs: ACSRRCI, ArgMarks: _ARG_RS1 | _ARG_IMM},          // csrci     csr, imm         => csrrci   x0 csr, imm
	AFRCSR:     {PseudoAs: ACSRRS, ArgMarks: _ARG_RD},                       // frcsr     rd               => csrrs    rd, fcsr
	AFSCSR:     {PseudoAs: ACSRRW, ArgMarks: _ARG_RD | _ARG_RS1},            // fscsr     rd, rs1          => csrrw    rd, fcsr, rs1 # rd 可省略
	AFRRM:      {PseudoAs: ACSRRS, ArgMarks: _ARG_RD},                       // frrm      rd               => csrrs    rd, frm, x0
	AFSRM:      {PseudoAs: ACSRRW, ArgMarks: _ARG_RD | _ARG_RS1},            // fsrm      rd, rs1          => csrrw    rd, frm, rs1 # rd 可省略
	AFRFLAGS:   {PseudoAs: ACSRRS, ArgMarks: _ARG_RD},                       // frflags   rd               => csrrs    rd, fflags, x0
	AFSFLAGS:   {PseudoAs: ACSRRW, ArgMarks: _ARG_RD | _ARG_RS1},            // fsflags   rd, rs1          => csrrw    rd, fflags, rs1 # rd 可省略

	// End marker

	ALAST: {},
}
