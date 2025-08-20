// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// 指令参数
type AsArgument struct {
	Rd  RegType // 目标寄存器
	Rs1 RegType // 原寄存器1
	Rs2 RegType // 原寄存器2
	Rs3 RegType // 原寄存器3
	Imm int32   // 立即数
}

// 编码指令
func (as As) Encode(arg *AsArgument) (uint32, error) {
	ctx := &AOpContextTable[as]
	return ctx.encode(as, arg)
}

func (ctx *OpContextType) encode(as As, arg *AsArgument) (uint32, error) {
	switch ctx.Opcode.FormatType() {
	case R:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_OP:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		case OpBase_OP_32:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		case OpBase_OP_FP:
			return ctx.encodeR(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2)), nil
		case OpBase_AMO:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		default:
			panic("unreachable")
		}
	case R4:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_MADD:
			return ctx.encodeR4(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), ctx.regI(arg.Rs3)), nil
		case OpBase_MSUB:
			return ctx.encodeR4(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), ctx.regI(arg.Rs3)), nil
		case OpBase_NMSUB:
			return ctx.encodeR4(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), ctx.regI(arg.Rs3)), nil
		case OpBase_NMADD:
			return ctx.encodeR4(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), ctx.regI(arg.Rs3)), nil
		default:
			panic("unreachable")
		}

	case I:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_OP_IMM:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_OP_IMM_32:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_JALR:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_LOAD:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_LOAD_FP:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_MISC_MEN:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
		case OpBase_SYSTEM:
			// SYSTEM 有两类情况：ECALL/EBREAK 和 CSR指令
			switch as {
			case AECALL: // ECALL = imm=0, rs1=0, rd=0, funct3=0
				return ctx.encodeI(0, 0, 0), nil
			case AEBREAK: // EBREAK = imm=1, rs1=0, rd=0, funct3=0
				return ctx.encodeI(0, 0, 1), nil
			default:
				return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.maskImm(arg.Imm, 12)), nil
			}
		default:
			panic("unreachable")
		}
	case S:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_STORE:
			return ctx.encodeS(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), ctx.rangeImm(arg.Imm, -2048, 2047)), nil
		case OpBase_STORE_FP:
			return ctx.encodeS(ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.rangeImm(arg.Imm, -2048, 2047)), nil
		default:
			panic("unreachable")
		}
	case B:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_BRANCH:
			imm := ctx.alignImm(arg.Imm, 2)
			imm = ctx.rangeImm(imm, -4096, 4095)
			return ctx.encodeB(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), imm), nil
		default:
			panic("unreachable")
		}
	case U:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_LUI:
			return ctx.encodeU(ctx.regI(arg.Rd), arg.Imm), nil
		case OpBase_AUIPC:
			return ctx.encodeU(ctx.regI(arg.Rd), ctx.maskImm(arg.Imm, 20)), nil
		default:
			panic("unreachable")
		}
	case J:
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_JAL:
			return ctx.encodeJ(ctx.regI(arg.Rd), ctx.alignImm(arg.Imm, 2)), nil
		default:
			panic("unreachable")
		}

	default:
		return 0, fmt.Errorf("riscv.Encode(%v): no implement", as)
	}
}

// R-type
func (ctx *OpContextType) encodeR(rd, rs1, rs2 uint32) uint32 {
	return ctx.Funct7<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// R4-type
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

// 立即数取低N个bit
func (ctx *OpContextType) maskImm(imm int32, bits int) int32 {
	mask := int32(1<<bits - 1)
	return int32(imm & mask)
}

// 立即数掩码校验范围
func (ctx *OpContextType) rangeImm(imm int32, min, max int32) int32 {
	if imm < min || imm > max {
		panic(fmt.Sprintf("imm out of range, want %d <= %d <= %d", min, imm, max))
	}
	return imm
}

// 立即数必须是n字节对齐
func (ctx *OpContextType) alignImm(imm int32, n int32) int32 {
	if (imm % n) != 0 {
		panic(fmt.Sprintf("imm must be 2-byte aligned, want %d%%%d == 0", imm, n))
	}
	return imm
}

// 返回寄存器机器码编号
func (ctx *OpContextType) regI(r RegType) uint32 {
	return ctx.regVal(r, REG_X0, REG_X31)
}

// 返回浮点数寄存器机器码编号
func (ctx *OpContextType) regF(r RegType) uint32 {
	return ctx.regVal(r, REG_F0, REG_F31)
}

// 返回寄存器机器码编号
func (ctx *OpContextType) regVal(r, min, max RegType) uint32 {
	if r < min || r > max {
		panic(fmt.Sprintf("register out of range, want %d <= %d <= %d", min, r, max))
	}
	return uint32(r - min)
}
