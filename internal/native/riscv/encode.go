// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// 指令参数
type AsArguments struct {
	Rd     uint32 // 目标寄存器
	Rs1    uint32 // 原寄存器1
	Rs2    uint32 // 原寄存器2
	Rs3    uint32 // 原寄存器3
	Imm    int32  // 立即数
	ImmMin int64  // 立即数范围
	ImmMax int64  // 立即数范围
	Funct3 uint32 // Function 3
	Funct7 uint32 // Function 7 (or Function 2)
}

// 编码指令
func (as As) Encode(arg *AsArguments) (uint32, error) {
	switch AOpTab[as].Opcode & OpBase_Mask {
	case OpBase_LOAD:
		return as.encode_OpBase_LOAD(arg)
	case OpBase_LOAD_FP:
		return as.encode_OpBase_LOAD_FP(arg)
	case OpBase_CUSTOM_0:
		return 0, fmt.Errorf("riscv.As.Encode(%v): TODO", as)
	case OpBase_MISC_MEN:
		return as.encode_OpBase_MISC_MEN(arg)
	case OpBase_OP_IMM:
		return as.encode_OpBase_OP_IMM(arg)
	case OpBase_AUIPC:
		return as.encode_OpBase_AUIPC(arg)
	case OpBase_OP_IMM_32:
		return as.encode_OpBase_OP_IMM_32(arg)

	case OpBase_STORE:
		return as.encode_OpBase_STORE(arg)
	case OpBase_STORE_FP:
		return as.encode_OpBase_STORE_FP(arg)
	case OpBase_CUSTOM_1:
		return 0, fmt.Errorf("riscv.As.Encode(%v): TODO", as)
	case OpBase_AMO:
		return as.encode_OpBase_AMO(arg)
	case OpBase_OP:
		return as.encode_OpBase_OP(arg)
	case OpBase_LUI:
		return as.encode_OpBase_LUI(arg)
	case OpBase_OP_32:
		return as.encode_OpBase_OP_32(arg)

	case OpBase_MADD:
		return as.encode_OpBase_MADD(arg)
	case OpBase_MSUB:
		return as.encode_OpBase_MSUB(arg)
	case OpBase_NMSUB:
		return as.encode_OpBase_NMSUB(arg)
	case OpBase_NMADD:
		return as.encode_OpBase_NMADD(arg)
	case OpBase_OP_FP:
		return as.encode_OpBase_OP_FP(arg)
	case OpBase_CUSTOM_2:
		return 0, fmt.Errorf("riscv.As.Encode(%v): TODO", as)

	case OpBase_BRANCH:
		return as.encode_OpBase_BRANCH(arg)
	case OpBase_JALR:
		return as.encode_OpBase_JALR(arg)
	case OpBase_JAL:
		return as.encode_OpBase_JAL(arg)
	case OpBase_SYSTEM:
		return as.encode_OpBase_SYSTEM(arg)
	case OpBase_CUSTOM_3:
		return 0, fmt.Errorf("riscv.As.Encode(%v): TODO", as)

	default:
		return 0, fmt.Errorf("riscv.As.Encode(%v): TODO", as)
	}
}

func (as As) encode_OpBase_LOAD(arg *AsArguments) (uint32, error) {
	// imm[11:0], rs1, funct3, rd, opcode
	if arg.Imm < -2048 || arg.Imm > 2047 {
		return 0, fmt.Errorf("imm out of range for LOAD: %d", arg.Imm)
	}
	op := AOpTab[as]
	return op.Opcode.encodeI(op.Funct3, arg.Rd, arg.Rs1, arg.Imm), nil
}

func (as As) encode_OpBase_LOAD_FP(arg *AsArguments) (uint32, error) {
	// imm[11:0] | rs1 | funct3 | rd | opcode
	imm := uint32(maskImm(arg.Imm, 12))
	return (imm << 20) | (arg.Rs1 << 15) | (arg.Funct3 << 12) | (arg.Rd << 7) | 0x07, nil
}
func (as As) encode_OpBase_MISC_MEN(arg *AsArguments) (uint32, error) {
	// imm[11:0] | rs1=0 | funct3 | rd=0 | opcode
	imm := uint32(maskImm(arg.Imm, 12))
	return (imm << 20) | (arg.Funct3 << 12) | 0x0F, nil
}
func (as As) encode_OpBase_OP_IMM(arg *AsArguments) (uint32, error) {
	imm := uint32(arg.Imm) // TODO: maskImm(arg.Imm, 12)
	return (imm << 20) | (arg.Rs1 << 15) | (arg.Funct3 << 12) | (arg.Rd << 7) | 0x13, nil
}
func (as As) encode_OpBase_AUIPC(arg *AsArguments) (uint32, error) {
	imm := uint32(maskImm(arg.Imm, 20)) // 高 20 位立即数
	return (imm << 12) | (arg.Rd << 7) | 0x17, nil
}
func (as As) encode_OpBase_OP_IMM_32(arg *AsArguments) (uint32, error) {
	imm := uint32(maskImm(arg.Imm, 12))
	return (imm << 20) | (arg.Rs1 << 15) | (arg.Funct3 << 12) | (arg.Rd << 7) | 0x1B, nil
}

func (as As) encode_OpBase_STORE(arg *AsArguments) (uint32, error) {
	if arg.Imm < -2048 || arg.Imm > 2047 {
		return 0, fmt.Errorf("imm out of range for STORE: %d", arg.Imm)
	}
	op := AOpTab[as]
	return op.Opcode.encodeS(op.Funct3, arg.Rs1, arg.Rs2, arg.Imm), nil
}
func (as As) encode_OpBase_STORE_FP(arg *AsArguments) (uint32, error) {
	imm := maskImm(arg.Imm, 12) // imm[11:0]
	op := AOpTab[as]
	return op.Opcode.encodeS(arg.Funct3, arg.Rs1, arg.Rs2, imm), nil
}
func (as As) encode_OpBase_AMO(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	// AMO 类使用 R-type，funct7 存储 funct5|aq|rl
	return op.Opcode.encodeR(arg.Funct3, arg.Rd, arg.Rs1, arg.Rs2, arg.Funct7), nil
}
func (as As) encode_OpBase_OP(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR(arg.Funct3, arg.Rd, arg.Rs1, arg.Rs2, arg.Funct7), nil
}
func (as As) encode_OpBase_LUI(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeU(arg.Rd, arg.Imm), nil
}
func (as As) encode_OpBase_OP_32(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR(arg.Funct3, arg.Rd, arg.Rs1, arg.Rs2, arg.Funct7), nil
}

func (as As) encode_OpBase_MADD(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR4(arg.Funct7, arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3, arg.Funct3), nil
}
func (as As) encode_OpBase_MSUB(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR4(arg.Funct7, arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3, arg.Funct3), nil
}
func (as As) encode_OpBase_NMSUB(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR4(arg.Funct7, arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3, arg.Funct3), nil
}
func (as As) encode_OpBase_NMADD(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR4(arg.Funct7, arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3, arg.Funct3), nil
}
func (as As) encode_OpBase_OP_FP(arg *AsArguments) (uint32, error) {
	op := AOpTab[as]
	return op.Opcode.encodeR(arg.Funct3, arg.Rd, arg.Rs1, arg.Rs2, arg.Funct7), nil
}

func (as As) encode_OpBase_BRANCH(arg *AsArguments) (uint32, error) {
	if arg.Imm%2 != 0 {
		return 0, fmt.Errorf("branch imm must be 2-byte aligned")
	}
	if arg.Imm < -4096 || arg.Imm > 4094 {
		return 0, fmt.Errorf("imm out of range for BRANCH: %d", arg.Imm)
	}
	op := AOpTab[as]
	return op.Opcode.encodeB(op.Funct3, arg.Rs1, arg.Rs2, arg.Imm), nil
}
func (as As) encode_OpBase_JALR(arg *AsArguments) (uint32, error) {
	// imm[11:0] + rs1 + funct3 + rd + opcode
	imm := maskImm(arg.Imm, 12)
	op := AOpTab[as].Opcode
	ins := op.encodeI(arg.Funct3, arg.Rd, arg.Rs1, imm)
	return ins, nil
}
func (as As) encode_OpBase_JAL(arg *AsArguments) (uint32, error) {
	if arg.Imm%2 != 0 {
		return 0, fmt.Errorf("jal imm must be 2-byte aligned")
	}
	op := AOpTab[as]
	return op.Opcode.encodeJ(arg.Rd, arg.Imm), nil
}
func (as As) encode_OpBase_SYSTEM(arg *AsArguments) (uint32, error) {
	op := AOpTab[as].Opcode

	// SYSTEM 有两类情况：ECALL/EBREAK 和 CSR指令
	switch as {
	case AECALL: // ECALL = imm=0, rs1=0, rd=0, funct3=0
		return op.encodeI(0, 0, 0, 0), nil
	case AEBREAK: // EBREAK = imm=1, rs1=0, rd=0, funct3=0
		return op.encodeI(0, 0, 0, 1), nil
	default:
		// CSR 指令：imm = CSR 地址
		imm := maskImm(arg.Imm, 12)
		ins := op.encodeI(arg.Funct3, arg.Rd, arg.Rs1, imm)
		return ins, nil
	}
}

// maskImm 立即数掩码
func maskImm(imm int32, bits int) int32 {
	mask := int32(1<<bits - 1)
	return int32(imm & mask)
}

// R-type
func (opcode OpcodeType) encodeR(funct3, funct7, rd, rs1, rs2 uint32) uint32 {
	return funct7<<25 | rs2<<20 | rs1<<15 | funct3<<12 | rd<<7 | uint32(opcode)
}

// R4-type (浮点四寄存器, MADD/MSUB/NMSUB/NMADD)
func (opcode OpcodeType) encodeR4(funct2, rd, rs1, rs2, rs3, funct3 uint32) uint32 {
	return (rs3 << 27) |
		(funct2 << 25) |
		(rs2 << 20) |
		(rs1 << 15) |
		(funct3 << 12) |
		(rd << 7) |
		uint32(opcode)
}

// I-type
func (opcode OpcodeType) encodeI(funct3, rd, rs1 uint32, imm int32) uint32 {
	imm12 := uint32(imm) & 0xFFF
	return imm12<<20 | rs1<<15 | funct3<<12 | rd<<7 | uint32(opcode)
}

// S-type
func (opcode OpcodeType) encodeS(funct3, rs1, rs2 uint32, imm int32) uint32 {
	imm12 := uint32(imm) & 0xFFF
	return (imm12>>5)<<25 | rs2<<20 | rs1<<15 | funct3<<12 | (imm12&0x1F)<<7 | uint32(opcode)
}

// B-type
func (opcode OpcodeType) encodeB(funct3, rs1, rs2 uint32, imm int32) uint32 {
	imm13 := uint32(imm) & 0x1FFF
	return ((imm13>>12)&1)<<31 | ((imm13>>5)&0x3F)<<25 |
		rs2<<20 | rs1<<15 | funct3<<12 |
		((imm13>>1)&0xF)<<8 | ((imm13>>11)&1)<<7 | uint32(opcode)
}

// U-type
func (opcode OpcodeType) encodeU(rd uint32, imm int32) uint32 {
	imm20 := uint32(imm) & 0xFFFFF000
	return imm20 | rd<<7 | uint32(opcode)
}

// J-type
func (opcode OpcodeType) encodeJ(rd uint32, imm int32) uint32 {
	imm21 := uint32(imm) & 0x1FFFFF
	return ((imm21>>20)&1)<<31 | ((imm21>>1)&0x3FF)<<21 |
		((imm21>>11)&1)<<20 | ((imm21>>12)&0xFF)<<12 |
		rd<<7 | uint32(opcode)
}
