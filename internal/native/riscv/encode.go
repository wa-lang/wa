// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// 指令参数
type AsArgument struct {
	Rd  uint32 // 目标寄存器
	Rs1 uint32 // 原寄存器1
	Rs2 uint32 // 原寄存器2
	Rs3 uint32 // 原寄存器3
	Imm int32  // 立即数
}

// 编码指令
func (as As) Encode(arg *AsArgument) (uint32, error) {
	ctx := &AOpContextTable[as]
	return ctx.encode(as, arg)
}

func (ctx *OpContextType) encode(as As, arg *AsArgument) (uint32, error) {
	switch ctx.Opcode & OpBase_Mask {
	case OpBase_LOAD:
		// imm[11:0], rs1, funct3, rd, opcode
		if arg.Imm < -2048 || arg.Imm > 2047 {
			return 0, fmt.Errorf("imm out of range for LOAD: %d", arg.Imm)
		}
		return ctx.encodeI(arg.Rd, arg.Rs1, arg.Imm), nil

	case OpBase_LOAD_FP:
		// imm[11:0] | rs1 | funct3 | rd | opcode
		imm := uint32(ctx.maskImm(arg.Imm, 12))
		return (imm << 20) | (arg.Rs1 << 15) | (ctx.Funct3 << 12) | (arg.Rd << 7) | 0x07, nil

	case OpBase_CUSTOM_0:
		return 0, fmt.Errorf("riscv.Encode(%v): unsupport", as)

	case OpBase_MISC_MEN:
		// imm[11:0] | rs1=0 | funct3 | rd=0 | opcode
		imm := uint32(ctx.maskImm(arg.Imm, 12))
		return (imm << 20) | (ctx.Funct3 << 12) | 0x0F, nil

	case OpBase_OP_IMM:
		imm := uint32(ctx.maskImm(arg.Imm, 12))
		return (imm << 20) | (arg.Rs1 << 15) | (ctx.Funct3 << 12) | (arg.Rd << 7) | 0x13, nil

	case OpBase_AUIPC:
		imm := uint32(ctx.maskImm(arg.Imm, 20)) // 高 20 位立即数
		return (imm << 12) | (arg.Rd << 7) | 0x17, nil

	case OpBase_OP_IMM_32:
		imm := uint32(ctx.maskImm(arg.Imm, 12))
		return (imm << 20) | (arg.Rs1 << 15) | (ctx.Funct3 << 12) | (arg.Rd << 7) | 0x1B, nil

	case OpBase_STORE:
		if arg.Imm < -2048 || arg.Imm > 2047 {
			return 0, fmt.Errorf("imm out of range for STORE: %d", arg.Imm)
		}
		return ctx.encodeS(arg.Rs1, arg.Rs2, arg.Imm), nil

	case OpBase_STORE_FP:
		return ctx.encodeS(arg.Rs1, arg.Rs2, ctx.maskImm(arg.Imm, 12)), nil

	case OpBase_CUSTOM_1:
		return 0, fmt.Errorf("riscv.Encode(%v): unsupport", as)

	case OpBase_AMO:
		return ctx.encodeR(arg.Rd, arg.Rs1, arg.Rs2), nil

	case OpBase_OP:
		return ctx.encodeR(arg.Rd, arg.Rs1, arg.Rs2), nil

	case OpBase_LUI:
		return ctx.encodeU(arg.Rd, arg.Imm), nil

	case OpBase_OP_32:
		return ctx.encodeR(arg.Rd, arg.Rs1, arg.Rs2), nil

	case OpBase_MADD:
		return ctx.encodeR4(arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3), nil

	case OpBase_MSUB:
		return ctx.encodeR4(arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3), nil

	case OpBase_NMSUB:
		return ctx.encodeR4(arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3), nil

	case OpBase_NMADD:
		return ctx.encodeR4(arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3), nil

	case OpBase_OP_FP:
		return ctx.encodeR(arg.Rd, arg.Rs1, arg.Rs2), nil

	case OpBase_CUSTOM_2:
		return 0, fmt.Errorf("riscv.Encode(%v): unsupport", as)

	case OpBase_BRANCH:
		if arg.Imm%2 != 0 {
			return 0, fmt.Errorf("branch imm must be 2-byte aligned")
		}
		if arg.Imm < -4096 || arg.Imm > 4094 {
			return 0, fmt.Errorf("imm out of range for BRANCH: %d", arg.Imm)
		}
		return ctx.encodeB(arg.Rs1, arg.Rs2, arg.Imm), nil

	case OpBase_JALR:
		// imm[11:0] + rs1 + funct3 + rd + opcode
		imm := ctx.maskImm(arg.Imm, 12)
		ins := ctx.encodeI(arg.Rd, arg.Rs1, imm)
		return ins, nil

	case OpBase_JAL:
		if arg.Imm%2 != 0 {
			return 0, fmt.Errorf("jal imm must be 2-byte aligned")
		}
		return ctx.encodeJ(arg.Rd, arg.Imm), nil

	case OpBase_SYSTEM:
		// SYSTEM 有两类情况：ECALL/EBREAK 和 CSR指令
		switch as {
		case AECALL: // ECALL = imm=0, rs1=0, rd=0, funct3=0
			return ctx.encodeI(0, 0, 0), nil
		case AEBREAK: // EBREAK = imm=1, rs1=0, rd=0, funct3=0
			return ctx.encodeI(0, 0, 1), nil
		default:
			// CSR 指令：imm = CSR 地址
			imm := ctx.maskImm(arg.Imm, 12)
			ins := ctx.encodeI(arg.Rd, arg.Rs1, imm)
			return ins, nil
		}

	case OpBase_CUSTOM_3:
		return 0, fmt.Errorf("riscv.Encode(%v): unsupport", as)

	default:
		return 0, fmt.Errorf("riscv.Encode(%v): unsupport", as)
	}
}

// R-type
func (ctx *OpContextType) encodeR(rd, rs1, rs2 uint32) uint32 {
	return ctx.Funct7<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// R4-type (浮点四寄存器, MADD/MSUB/NMSUB/NMADD)
func (ctx *OpContextType) encodeR4(rd, rs1, rs2, rs3 uint32) uint32 {
	return (rs3 << 27) |
		(ctx.Funct7 << 25) |
		(rs2 << 20) |
		(rs1 << 15) |
		(ctx.Funct3 << 12) |
		(rd << 7) |
		uint32(ctx.Opcode)
}

// I-type
func (ctx *OpContextType) encodeI(rd, rs1 uint32, imm int32) uint32 {
	imm12 := uint32(imm) & 0xFFF
	return imm12<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// S-type
func (ctx *OpContextType) encodeS(rs1, rs2 uint32, imm int32) uint32 {
	imm12 := uint32(imm) & 0xFFF
	return (imm12>>5)<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | (imm12&0x1F)<<7 | uint32(ctx.Opcode)
}

// B-type
func (ctx *OpContextType) encodeB(rs1, rs2 uint32, imm int32) uint32 {
	imm13 := uint32(imm) & 0x1FFF
	return ((imm13>>12)&1)<<31 | ((imm13>>5)&0x3F)<<25 |
		rs2<<20 | rs1<<15 | ctx.Funct3<<12 |
		((imm13>>1)&0xF)<<8 | ((imm13>>11)&1)<<7 | uint32(ctx.Opcode)
}

// U-type
func (ctx *OpContextType) encodeU(rd uint32, imm int32) uint32 {
	imm20 := uint32(imm) & 0xFFFFF000
	return imm20 | rd<<7 | uint32(ctx.Opcode)
}

// J-type
func (ctx *OpContextType) encodeJ(rd uint32, imm int32) uint32 {
	imm21 := uint32(imm) & 0x1FFFFF
	return ((imm21>>20)&1)<<31 | ((imm21>>1)&0x3FF)<<21 |
		((imm21>>11)&1)<<20 | ((imm21>>12)&0xFF)<<12 |
		rd<<7 | uint32(ctx.Opcode)
}

// maskImm 立即数掩码
func (ctx *OpContextType) maskImm(imm int32, bits int) int32 {
	mask := int32(1<<bits - 1)
	return int32(imm & mask)
}
