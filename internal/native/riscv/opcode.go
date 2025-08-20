// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

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

// 操作码上下文信息
type OpContextType struct {
	Opcode OpcodeType
	Funct3 uint32
	Rs1    uint32
	Rs2    uint32
	Csr    int64
	Funct7 uint32 // 和 Funct2 共用
}

// 指令编码信息表
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
var AOpContextTable = []OpContextType{
	// RV32I Base Instruction Set

	ALUI:    {Opcode: OpBase_LUI},
	AAUIPC:  {Opcode: OpBase_AUIPC},
	AJAL:    {Opcode: OpBase_JAL},
	AJALR:   {Opcode: OpBase_JALR, Funct3: 0b000},
	ABEQ:    {Opcode: OpBase_BRANCH, Funct3: 0b000},
	ABNE:    {Opcode: OpBase_BRANCH, Funct3: 0b001},
	ABLT:    {Opcode: OpBase_BRANCH, Funct3: 0b100},
	ABGE:    {Opcode: OpBase_BRANCH, Funct3: 0b101},
	ABLTU:   {Opcode: OpBase_BRANCH, Funct3: 0b110},
	ABGEU:   {Opcode: OpBase_BRANCH, Funct3: 0b111},
	ALB:     {Opcode: OpBase_LOAD, Funct3: 0b000},
	ALH:     {Opcode: OpBase_LOAD, Funct3: 0b001},
	ALW:     {Opcode: OpBase_LOAD, Funct3: 0b010},
	ALBU:    {Opcode: OpBase_LOAD, Funct3: 0b100},
	ALHU:    {Opcode: OpBase_LOAD, Funct3: 0b101},
	ASB:     {Opcode: OpBase_STORE, Funct3: 0b000},
	ASH:     {Opcode: OpBase_STORE, Funct3: 0b001},
	ASW:     {Opcode: OpBase_STORE, Funct3: 0b010},
	AADDI:   {Opcode: OpBase_OP_IMM, Funct3: 0b000},
	ASLTI:   {Opcode: OpBase_OP_IMM, Funct3: 0b010},
	ASLTIU:  {Opcode: OpBase_OP_IMM, Funct3: 0b011},
	AXORI:   {Opcode: OpBase_OP_IMM, Funct3: 0b100},
	AORI:    {Opcode: OpBase_OP_IMM, Funct3: 0b110},
	AANDI:   {Opcode: OpBase_OP_IMM, Funct3: 0b111},
	ASLLI:   {Opcode: OpBase_OP_IMM, Funct3: 0b001, Funct7: 0b0000000},
	ASRLI:   {Opcode: OpBase_OP_IMM, Funct3: 0b101, Funct7: 0b0000000},
	ASRAI:   {Opcode: OpBase_OP_IMM, Funct3: 0b101, Funct7: 0b0100000},
	AADD:    {Opcode: OpBase_OP, Funct3: 0b000, Funct7: 0b0000000},
	ASUB:    {Opcode: OpBase_OP, Funct3: 0b000, Funct7: 0b0100000},
	ASLL:    {Opcode: OpBase_OP, Funct3: 0b001, Funct7: 0b0000000},
	ASLT:    {Opcode: OpBase_OP, Funct3: 0b010, Funct7: 0b0000000},
	ASLTU:   {Opcode: OpBase_OP, Funct3: 0b011, Funct7: 0b0000000},
	AXOR:    {Opcode: OpBase_OP, Funct3: 0b100, Funct7: 0b0000000},
	ASRL:    {Opcode: OpBase_OP, Funct3: 0b101, Funct7: 0b0000000},
	ASRA:    {Opcode: OpBase_OP, Funct3: 0b101, Funct7: 0b0100000},
	AOR:     {Opcode: OpBase_OP, Funct3: 0b110, Funct7: 0b0000000},
	AAND:    {Opcode: OpBase_OP, Funct3: 0b111, Funct7: 0b0000000},
	AFENCE:  {Opcode: OpBase_MISC_MEN, Funct3: 0b000},
	AECALL:  {Opcode: OpBase_SYSTEM, Rs2: 0, Funct7: 0}, // imm[11:0] = 0b000000000000
	AEBREAK: {Opcode: OpBase_SYSTEM, Rs2: 1, Funct7: 0}, // imm[11:0] = 0b000000000001

	// RV64I Base Instruction Set (in addition to RV32I)

	ALWU:   {Opcode: OpBase_LOAD, Funct3: 0b110},
	ALD:    {Opcode: OpBase_LOAD, Funct3: 0b011},
	ASD:    {Opcode: OpBase_STORE},
	AADDIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b000},
	ASLLIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b001, Funct7: 0b0000000},
	ASRLIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b101, Funct7: 0b0000000},
	ASRAIW: {Opcode: OpBase_OP_IMM_32, Funct3: 0b101, Funct7: 0b0100000},
	AADDW:  {Opcode: OpBase_OP_32, Funct3: 0b000, Funct7: 0b0000000},
	ASUBW:  {Opcode: OpBase_OP_32, Funct3: 0b000, Funct7: 0b0100000},
	ASLLW:  {Opcode: OpBase_OP_32, Funct3: 0b001, Funct7: 0b0000000},
	ASRLW:  {Opcode: OpBase_OP_32, Funct3: 0b101, Funct7: 0b0000000},
	ASRAW:  {Opcode: OpBase_OP_32, Funct3: 0b101, Funct7: 0b0100000},

	// RV32/RV64 Zicsr Standard Extension

	ACSRRW:  {Opcode: OpBase_SYSTEM, Funct3: 0b001},
	ACSRRS:  {Opcode: OpBase_SYSTEM, Funct3: 0b010},
	ACSRRC:  {Opcode: OpBase_SYSTEM, Funct3: 0b011},
	ACSRRWI: {Opcode: OpBase_SYSTEM, Funct3: 0b101},
	ACSRRSI: {Opcode: OpBase_SYSTEM, Funct3: 0b110},
	ACSRRCI: {Opcode: OpBase_SYSTEM, Funct3: 0b111},

	// RV32M Standard Extension

	AMUL:    {Opcode: OpBase_OP, Funct3: 0b000, Funct7: 0b0000001},
	AMULH:   {Opcode: OpBase_OP, Funct3: 0b001, Funct7: 0b0000001},
	AMULHSU: {Opcode: OpBase_OP, Funct3: 0b010, Funct7: 0b0000001},
	AMULHU:  {Opcode: OpBase_OP, Funct3: 0b011, Funct7: 0b0000001},
	ADIV:    {Opcode: OpBase_OP, Funct3: 0b100, Funct7: 0b0000001},
	ADIVU:   {Opcode: OpBase_OP, Funct3: 0b101, Funct7: 0b0000001},
	AREM:    {Opcode: OpBase_OP, Funct3: 0b110, Funct7: 0b0000001},
	AREMU:   {Opcode: OpBase_OP, Funct3: 0b111, Funct7: 0b0000001},

	// RV64M Standard Extension (in addition to RV32M)

	AMULW:  {Opcode: OpBase_OP_32, Funct3: 0b000, Funct7: 0b0000001},
	ADIVW:  {Opcode: OpBase_OP_32, Funct3: 0b100, Funct7: 0b0000001},
	ADIVUW: {Opcode: OpBase_OP_32, Funct3: 0b101, Funct7: 0b0000001},
	AREMW:  {Opcode: OpBase_OP_32, Funct3: 0b110, Funct7: 0b0000001},
	AREMUW: {Opcode: OpBase_OP_32, Funct3: 0b111, Funct7: 0b0000001},

	// TODO: RV32F Standard Extension

	// TODO: RV64F Standard Extension (in addition to RV32F)

	// TODO: RV32D Standard Extension

	// TODO: RV64D Standard Extension (in addition to RV32D)

	// TODO: 验证遗漏的指令
}
