// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import "wa-lang.org/wa/internal/native/abi"

// TODO: 补充 opcode 格式注释

// 指令格式类型
type _OpFormatType int

const (
	_ _OpFormatType = iota
	_R
	_I
	_D
	_B
	_CB
	_IW
)

func (x _OpFormatType) Bits() int {
	switch x {
	case _R:
		return 11
	case _I:
		return 10
	case _D:
		return 11
	case _B:
		return 6
	case _CB:
		return 8
	case _IW:
		return 9
	default:
		panic("unreachable")
	}
}

func (x _OpFormatType) HasShamt() bool {
	if x == _R {
		return true
	} else {
		return false
	}
}

func (x _OpFormatType) HasBCond(as abi.As) bool {
	switch as {
	case AB_EQ, // PC = PC+ CondBranchAddr if FLAGS equal
		AB_NE, // PC = PC+ CondBranchAddr if FLAGS not equal
		AB_LT, // PC = PC+ CondBranchAddr if FLAGS less than
		AB_LE, // PC = PC+ CondBranchAddr if FLAGS less than or equal
		AB_GT, // PC = PC+ CondBranchAddr if FLAGS greater than
		AB_GE, // PC = PC+ CondBranchAddr if FLAGS greater than or equal
		AB_LO, // PC = PC+ CondBranchAddr if FLAGS lower
		AB_LS, // PC = PC+ CondBranchAddr if FLAGS lower or same
		AB_HI, // PC = PC+ CondBranchAddr if FLAGS higher
		AB_HS: // PC = PC+ CondBranchAddr if FLAGS higher or same
		return true
	}
	return false
}

// 操作码上下文信息
type _OpContextType struct {
	Format _OpFormatType // 指令类型
	Opcode uint32        // 指令码
	Shamt  uint32        // Shamt
	Cond   uint32        // 条件码
}

// 指令编码信息表
// 保持和《计算机组成与设计：硬件/软件接口》的参考表格顺序一致
var _AOpContextTable = []_OpContextType{
	AB:      {Format: _B, Opcode: 0b_000101, Shamt: 0b_0},           // R[Rd] = R[Rn] + R[Rm]
	AFMULS:  {Format: _R, Opcode: 0b_00011110001, Shamt: 0b_000010}, // R[Rd] = R[Rn] + ALUImm
	AFDIVS:  {Format: _R, Opcode: 0b_00011110001, Shamt: 0b_000110}, // R[Rd], FLAGS = R[Rn] + ALUImm
	AFCMPS:  {Format: _R, Opcode: 0b_00011110001, Shamt: 0b_001000}, // R[Rd], FLAGS = R[Rn] +R[Rm]
	AFADDS:  {Format: _R, Opcode: 0b_00011110001, Shamt: 0b_001010}, // R[Rd] = R[Rn] & R[Rm]
	AFSUBS:  {Format: _R, Opcode: 0b_00011110001, Shamt: 0b_001110}, // R[Rd] = R[Rn] & ALUImm
	AFMULD:  {Format: _R, Opcode: 0b_00011110011, Shamt: 0b_000010}, // R[Rd], FLAGS = R[Rn] & ALUImm
	AFDIVD:  {Format: _R, Opcode: 0b_00011110011, Shamt: 0b_000110}, // R[Rd], FLAGS=R[Rn] & R[Rm]
	AFCMPD:  {Format: _R, Opcode: 0b_00011110011, Shamt: 0b_001000}, // PC = PC + BranchAddr
	AFADDD:  {Format: _R, Opcode: 0b_00011110011, Shamt: 0b_001010}, // PC = PC+ CondBranchAddr if FLAGS equal
	AFSUBD:  {Format: _R, Opcode: 0b_00011110011, Shamt: 0b_001110}, // PC = PC+ CondBranchAddr if FLAGS not equal
	ASTURB:  {Format: _D, Opcode: 0b_00111000000, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS less than
	ALDURB:  {Format: _D, Opcode: 0b_00111000010, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS less than or equal
	AB_EQ:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_00000},     // PC = PC+ CondBranchAddr if FLAGS equal
	AB_NE:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_00001},     // PC = PC+ CondBranchAddr if FLAGS not equal
	AB_LT:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01011},     // PC = PC+ CondBranchAddr if FLAGS less than
	AB_LE:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01101},     // PC = PC+ CondBranchAddr if FLAGS less than or equal
	AB_GT:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01100},     // PC = PC+ CondBranchAddr if FLAGS greater than
	AB_GE:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01010},     // PC = PC+ CondBranchAddr if FLAGS greater than or equal
	AB_LO:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_00011},     // PC = PC+ CondBranchAddr if FLAGS lower
	AB_LS:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01001},     // PC = PC+ CondBranchAddr if FLAGS lower or same
	AB_HI:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_01000},     // PC = PC+ CondBranchAddr if FLAGS higher
	AB_HS:   {Format: _CB, Opcode: 0b_01010100, Cond: 0b_00010},     // PC = PC+ CondBranchAddr if FLAGS higher or same
	ASTURH:  {Format: _D, Opcode: 0b_01111000000, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS greater than or equal
	ALDURH:  {Format: _D, Opcode: 0b_01111000010, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS lower
	AAND:    {Format: _R, Opcode: 0b_10001010000, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS lower or same
	AADD:    {Format: _R, Opcode: 0b_10001011000, Shamt: 0b_0},      // PC = PC+ CondBranchAddr if FLAGS higher
	AADDI:   {Format: _I, Opcode: 0b_1001000100, Shamt: 0b_0},       // PC = PC+ CondBranchAddr if FLAGS higher or same
	AANDI:   {Format: _I, Opcode: 0b_1001001000, Shamt: 0b_0},       // R[30] = PC+4; PC = PC+ BranchAddr
	ABL:     {Format: _B, Opcode: 0b_100101, Shamt: 0b_0},           // PC= R[Rt]
	ASDIV:   {Format: _R, Opcode: 0b_10011010110, Shamt: 0b_000010}, // PC = PC + CondBranchAddr if(R[Rt]!=0)
	AUDIV:   {Format: _R, Opcode: 0b_10011010110, Shamt: 0b_000011}, // PC = PC + CondBranchAddr if(R[Rt]==0)
	AMUL:    {Format: _R, Opcode: 0b_10011011000, Shamt: 0b_011111}, // R[Rd] = R[Rn] ^ R[Rm]
	ASMULH:  {Format: _R, Opcode: 0b_10011011010, Shamt: 0b_0},      // R[Rd] = R[Rn] ^ ALUImm
	AUMULH:  {Format: _R, Opcode: 0b_10011011110, Shamt: 0b_0},      // R[Rt] = M[R[Rn] + DTAddr]
	AORR:    {Format: _R, Opcode: 0b_10101010000, Shamt: 0b_0},      // R[Rt]={56'bо, M[R[Rn] + DTAddr](7:0)}
	AADDS:   {Format: _R, Opcode: 0b_10101011000, Shamt: 0b_0},      // R[Rt]={48'bо, M[R[Rn] + DTAddr](15:0)}
	AADDIS:  {Format: _I, Opcode: 0b_1011000100, Shamt: 0b_0},       // R[Rt]={32{M[R[Rn] + DTAddr][31]}, M[R[Rn] + DTAddr] (31:0)}
	AORRI:   {Format: _I, Opcode: 0b_1011001000, Shamt: 0b_0},       // R[Rd] = M[R[Rn] + DTAddr]
	ACBZ:    {Format: _CB, Opcode: 0b_10110100, Shamt: 0b_0},        // R[Rd] = R[Rn] << shamt
	ACBNZ:   {Format: _CB, Opcode: 0b_10110101, Shamt: 0b_0},        // R[Rd] = R[Rn]>>> shamt
	ASTURW:  {Format: _D, Opcode: 0b_10111000000, Shamt: 0b_0},      // R[Rd](Instruction[22:21]*16: Instruction[22:21]*16-15)=MOVImm
	ALDURSW: {Format: _D, Opcode: 0b_10111000100, Shamt: 0b_0},      // R[Rd] = {MOVImm <<(Instruction[22:21]*16)}
	ASTURS:  {Format: _R, Opcode: 0b_10111100000, Shamt: 0b_0},      // R[Rd] = R[Rn] | R[Rm]
	ALDURS:  {Format: _R, Opcode: 0b_10111100010, Shamt: 0b_0},      // R[Rd] = R[Rn] | ALUImm
	ASTXR:   {Format: _D, Opcode: 0b_11001000000, Shamt: 0b_0},      // M[R[Rn] + DTAddr] = R[Rt]
	ALDXR:   {Format: _D, Opcode: 0b_11001000010, Shamt: 0b_0},      // M[R[Rn]+DTAddr](7:0) = R[Rt](7:0)
	AEOR:    {Format: _R, Opcode: 0b_11001010000, Shamt: 0b_0},      // M[R[Rn] + DTAddr](15:0) = R[Rt](15:0)
	ASUB:    {Format: _R, Opcode: 0b_11001011000, Shamt: 0b_0},      // M[R[Rn] + DTAddr](31:0) = R[Rt](31:0)
	ASUBI:   {Format: _I, Opcode: 0b_1101000100, Shamt: 0b_0},       // M[R[Rn] + DTAddr] = R[Rt]; R[Rm] = (atomic)?0:1
	AEORI:   {Format: _I, Opcode: 0b_1101001000, Shamt: 0b_0},       // R[Rd] = R[Rn] - R[Rm]
	AMOVZ:   {Format: _IW, Opcode: 0b_110100101, Shamt: 0b_0},       // R[Rd] = R[Rn] - ALUImm
	ALSR:    {Format: _R, Opcode: 0b_11010011010, Shamt: 0b_0},      // R[Rd], FLAGS = R[Rn] - ALUImm
	ALSL:    {Format: _R, Opcode: 0b_11010011011, Shamt: 0b_0},      // R[Rd], FLAGS = R[Rn] - R[Rm]
	ABR:     {Format: _R, Opcode: 0b_11010110000, Shamt: 0b_0},      // S[Rd] = S[Rn] + S[Rm]
	AANDS:   {Format: _R, Opcode: 0b_11101010000, Shamt: 0b_0},      // D[Rd] = D[Rn] + D[Rm]
	ASUBS:   {Format: _R, Opcode: 0b_11101011000, Shamt: 0b_0},      // FLAGS = (S[Rn] vs S[Rm])
	ASUBIS:  {Format: _I, Opcode: 0b_1111000100, Shamt: 0b_0},       // FLAGS = (D[Rn] vs D[Rm])
	AANDIS:  {Format: _I, Opcode: 0b_1111001000, Shamt: 0b_0},       // S[Rd] = S[Rn]/S[Rm]
	AMOVK:   {Format: _IW, Opcode: 0b_111100101, Shamt: 0b_0},       // D[Rd] = D[Rn]/D[Rm]
	ASTUR:   {Format: _D, Opcode: 0b_11111000000, Shamt: 0b_0},      // S[Rd] = S[Rn] * S[Rm]
	ALDUR:   {Format: _D, Opcode: 0b_11111000010, Shamt: 0b_0},      // D[Rd] = D[Rn] * D[Rm]
	ASTURD:  {Format: _R, Opcode: 0b_11111100000, Shamt: 0b_0},      // S[Rd] = S[Rn] – S[Rm]
	ALDURD:  {Format: _R, Opcode: 0b_11111100010, Shamt: 0b_0},      // D[Rd] = D[Rn] - D[Rm]
}
