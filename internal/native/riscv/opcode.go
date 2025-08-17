// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

//
// 0        6 7                   11 12      14 15                  19 20                  24 25                   31
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
// | opcode | |         rd         | | funct3 | |        rs1         | |        rs2         | |      funct7         | R-type, 寄存器类
// +--------+ +--------------------+ +--------+ +--------------------+ +--------------------+ +---------------------+
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

// 指令编码类型
type OpType int

const (
	R OpType = iota + 1
	I
	S
	B
	U
	J
)

type OptabType struct {
	Optype OpType
	Opcode uint32
	Funct3 uint32
	Rs1    uint32
	Rs2    uint32
	Csr    int64
	Funct7 uint32
}

// 指令编码信息表
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rv32-64g
var AOpTab = []OptabType{
	// RV32I Base Instruction Set
	ALUI:    {Optype: U, Opcode: 0b0110111},
	AAUIPC:  {Optype: U, Opcode: 0b0010111},
	AJAL:    {Optype: J, Opcode: 0b1101111},
	AJALR:   {Optype: I, Opcode: 0b1100111, Funct3: 0b000},
	ABEQ:    {Optype: B, Opcode: 0b1100011, Funct3: 0b000},
	ABNE:    {Optype: B, Opcode: 0b1100011, Funct3: 0b001},
	ABLT:    {Optype: B, Opcode: 0b1100011, Funct3: 0b100},
	ABGE:    {Optype: B, Opcode: 0b1100011, Funct3: 0b101},
	ABLTU:   {Optype: B, Opcode: 0b1100011, Funct3: 0b110},
	ABGEU:   {Optype: B, Opcode: 0b1100011, Funct3: 0b111},
	ALB:     {Optype: I, Opcode: 0b0000011, Funct3: 0b000},
	ALH:     {Optype: I, Opcode: 0b0000011, Funct3: 0b001},
	ALW:     {Optype: I, Opcode: 0b0000011, Funct3: 0b010},
	ALBU:    {Optype: I, Opcode: 0b0000011, Funct3: 0b100},
	ALHU:    {Optype: I, Opcode: 0b0000011, Funct3: 0b101},
	ASB:     {Optype: S, Opcode: 0b0100011, Funct3: 0b000},
	ASH:     {Optype: S, Opcode: 0b0100011, Funct3: 0b001},
	ASW:     {Optype: S, Opcode: 0b0100011, Funct3: 0b010},
	AADDI:   {Optype: I, Opcode: 0b0010011, Funct3: 0b000},
	ASLTI:   {Optype: I, Opcode: 0b0010011, Funct3: 0b010},
	ASLTIU:  {Optype: I, Opcode: 0b0010011, Funct3: 0b011},
	AXORI:   {Optype: I, Opcode: 0b0010011, Funct3: 0b100},
	AORI:    {Optype: I, Opcode: 0b0010011, Funct3: 0b110},
	AANDI:   {Optype: I, Opcode: 0b0010011, Funct3: 0b111},
	ASLLI:   {Optype: R, Opcode: 0b0010011, Funct3: 0b001, Funct7: 0b0000000},
	ASRLI:   {Optype: R, Opcode: 0b0010011, Funct3: 0b101, Funct7: 0b0000000},
	ASRAI:   {Optype: R, Opcode: 0b0010011, Funct3: 0b101, Funct7: 0b0100000},
	AADD:    {Optype: R, Opcode: 0b0110011, Funct3: 0b000, Funct7: 0b0000000},
	ASUB:    {Optype: R, Opcode: 0b0110011, Funct3: 0b000, Funct7: 0b0100000},
	ASLL:    {Optype: R, Opcode: 0b0110011, Funct3: 0b001, Funct7: 0b0000000},
	ASLT:    {Optype: R, Opcode: 0b0110011, Funct3: 0b010, Funct7: 0b0000000},
	ASLTU:   {Optype: R, Opcode: 0b0110011, Funct3: 0b011, Funct7: 0b0000000},
	AXOR:    {Optype: R, Opcode: 0b0110011, Funct3: 0b100, Funct7: 0b0000000},
	ASRL:    {Optype: R, Opcode: 0b0110011, Funct3: 0b101, Funct7: 0b0000000},
	ASRA:    {Optype: R, Opcode: 0b0110011, Funct3: 0b101, Funct7: 0b0100000},
	AOR:     {Optype: R, Opcode: 0b0110011, Funct3: 0b110, Funct7: 0b0000000},
	AAND:    {Optype: R, Opcode: 0b0110011, Funct3: 0b111, Funct7: 0b0000000},
	AECALL:  {Optype: I, Opcode: 0b1110011, Rs2: 0, Funct7: 0}, // imm[11:0] = 0b000000000000
	AEBREAK: {Optype: I, Opcode: 0b1110011, Rs2: 1, Funct7: 0}, // imm[11:0] = 0b000000000001

	// TODO: RV64I Base Instruction Set (in addition to RV32I)

	// TODO: RV32/RV64 Zicsr Standard Extension

	// TODO: RV32M Standard Extension

	// TODO: RV64M Standard Extension (in addition to RV32M)

	// TODO: RV32F Standard Extension

	// TODO: RV64F Standard Extension (in addition to RV32F)

	// TODO: RV32D Standard Extension

	// TODO: RV64D Standard Extension (in addition to RV32D)

	// TODO: 验证遗漏的指令
}
