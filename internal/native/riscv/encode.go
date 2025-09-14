// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 返回寄存器机器码编号
func RegI(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regI(r)
}

// 返回浮点数寄存器机器码编号
func RegF(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regF(r)
}

// 编码RISCV32指令
func Encode(cpu abi.CPUType, as abi.As, arg *abi.AsArgument) (uint32, error) {
	switch cpu {
	case abi.RISCV32:
		return EncodeRV32(as, arg)
	case abi.RISCV64:
		return EncodeRV64(as, arg)
	default:
		return 0, fmt.Errorf("unknonw cpu: %v", cpu)
	}
}

// 编码RISCV32指令
func EncodeRV32(as abi.As, arg *abi.AsArgument) (uint32, error) {
	ctx := &_AOpContextTable[as]
	if ctx.PseudoAs != 0 {
		return ctx.encodePseudo(32, as, arg)
	}
	return ctx.encodeRaw(32, as, arg)
}

// 编码RISCV64指令
func EncodeRV64(as abi.As, arg *abi.AsArgument) (uint32, error) {
	ctx := &_AOpContextTable[as]
	if ctx.PseudoAs != 0 {
		return ctx.encodePseudo(64, as, arg)
	}
	return ctx.encodeRaw(64, as, arg)
}

// 输入一个 32 位有符号立即数 imm, 输出 low(12bit)/high(20bit)
// 满足 imm 约等于 (high << 12) + low, 以便于进行长地址跳转的拆分
func Split32BitImmediate(imm int64) (low12bit, high20bit int64, err error) {
	return split32BitImmediate(imm)
}

func (ctx *_OpContextType) encodeRaw(xlen int, as abi.As, arg *abi.AsArgument) (uint32, error) {
	if ctx.PseudoAs != 0 {
		panic("unreachable")
	}

	// 检查模板和参数
	if err := ctx.checkArgMarks(xlen, as, arg, ctx.ArgMarks); err != nil {
		return 0, err
	}

	// 编码指令
	switch ctx.Opcode.FormatType() {
	case _R:
		return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
	case _R4:
		return ctx.encodeR4(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.regF(arg.Rs3)), nil
	case _I:
		switch as {
		case AECALL:
			return ctx.encodeI(0, 0, 0b_0000_0000_0000), nil
		case AEBREAK:
			return ctx.encodeI(0, 0, 0b_0000_0000_0001), nil
		default:
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil
		}
	case _S:
		return ctx.encodeS(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), uint32(arg.Imm)), nil
	case _B:
		return ctx.encodeB(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), uint32(arg.Imm)), nil
	case _U:
		return ctx.encodeU(ctx.regI(arg.Rd), uint32(arg.Imm)), nil
	case _J:
		// jal 伪指令和基础指令同名, 但是 rd 参数可选
		if ctx.ArgMarks&_ARG_RD_IS_X != 0 {
			if arg.Rd == 0 {
				return ctx.encodeJ(ctx.regI(REG_X0), uint32(arg.Imm)), nil
			}
		}
		return ctx.encodeJ(ctx.regI(arg.Rd), uint32(arg.Imm)), nil

	default:
		return 0, fmt.Errorf("riscv.encodeRaw(%v): no implement", as)
	}
}

// R-type
func (ctx *_OpContextType) encodeR(rd, rs1, rs2 uint32) uint32 {
	return ctx.Funct7<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// R4-type
func (ctx *_OpContextType) encodeR4(rd, rs1, rs2, rs3 uint32) uint32 {
	funct2 := ctx.Funct7
	return rs3<<27 | funct2<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// I-type
func (ctx *_OpContextType) encodeI(rd, rs1 uint32, imm uint32) uint32 {
	return imm<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// S-type
func (ctx *_OpContextType) encodeS(rs1, rs2 uint32, imm uint32) uint32 {
	return (imm>>5)<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | (imm&0b_1_1111)<<7 | uint32(ctx.Opcode)
}

// B-type
func (ctx *_OpContextType) encodeB(rs1, rs2 uint32, imm uint32) uint32 {
	return ctx.encodeB_Imm(imm) | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | uint32(ctx.Opcode)
}
func (ctx *_OpContextType) encodeB_Imm(imm uint32) uint32 {
	return (imm>>12)<<31 | ((imm>>5)&0x3f)<<25 | ((imm>>1)&0xf)<<8 | ((imm>>11)&0x1)<<7
}

// U-type
func (ctx *_OpContextType) encodeU(rd uint32, imm uint32) uint32 {
	return (imm << 12) | rd<<7 | uint32(ctx.Opcode)
}

// J-type
func (ctx *_OpContextType) encodeJ(rd uint32, imm uint32) uint32 {
	return ctx.encodeJ_Imm(imm) | rd<<7 | uint32(ctx.Opcode)
}
func (ctx *_OpContextType) encodeJ_Imm(imm uint32) uint32 {
	return (imm>>20)<<31 | ((imm>>1)&0x3ff)<<21 | ((imm>>11)&0x1)<<20 | ((imm>>12)&0xff)<<12
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regI(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_X0, REG_X31)
}

// 返回浮点数寄存器机器码编号
func (ctx *_OpContextType) regF(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_F0, REG_F31)
}

// 返回寄存器机器码编号
func (ctx *_OpContextType) regVal(r, min, max abi.RegType) uint32 {
	if r < min || r > max {
		panic(fmt.Sprintf("register out of range, want %d <= %d <= %d", min, r, max))
	}
	return uint32(r - min)
}
