// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

type U32Slice []byte

// 编码RISCV32指令
func EncodeRV32(as abi.As, arg *abi.AsArgument) (uint32, error) {
	ctx := &AOpContextTable[as]
	if ctx.PseudoAs != 0 {
		return ctx.encodePseudo(32, as, arg)
	}
	return ctx.encodeRaw(32, as, arg)
}

// 编码RISCV64指令
func EncodeRV64(as abi.As, arg *abi.AsArgument) (uint32, error) {
	ctx := &AOpContextTable[as]
	if ctx.PseudoAs != 0 {
		return ctx.encodePseudo(64, as, arg)
	}
	return ctx.encodeRaw(64, as, arg)
}

func (ctx *OpContextType) encodeRaw(xlen int, as abi.As, arg *abi.AsArgument) (uint32, error) {
	if ctx.PseudoAs != 0 {
		panic("unreachable")
	}
	switch ctx.Opcode.FormatType() {
	case R:
		if ctx.Funct3&0b_111 != 0 {
			panic(fmt.Sprintf("encodeR: %d is invalid funct3", ctx.Funct3))
		}
		if ctx.Funct7&0b_111_1111 != 0 {
			panic(fmt.Sprintf("encodeR: %d is invalid funct7", ctx.Funct7))
		}
		if arg.Rs3 != 0 {
			panic("encodeR: rs3 was nonzero")
		}
		if arg.Imm != 0 {
			panic("encodeR: imm was nonzero")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_OP:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		case OpBase_OP_32:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		case OpBase_OP_FP:
			if ctx.Rs2 != nil {
				if arg.Rs2 != 0 {
					panic("encodeR: rs2 was nonzero")
				}
				return ctx.encodeR(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), *ctx.Rs2), nil
			}
			return ctx.encodeR(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2)), nil
		case OpBase_AMO:
			return ctx.encodeR(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.regI(arg.Rs2)), nil
		default:
			panic("unreachable")
		}
	case R4:
		if ctx.Funct3&0b_111 != 0 {
			panic(fmt.Sprintf("encodeR4: %d is invalid funct3", ctx.Funct3))
		}
		if funct2 := ctx.Funct7; funct2&0b_11 != 0 {
			panic(fmt.Sprintf("encodeR4: %d is invalid funct2", funct2))
		}
		if arg.Imm != 0 {
			panic("encodeR4: imm was nonzero")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_MADD:
			return ctx.encodeR4(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.regF(arg.Rs3)), nil
		case OpBase_MSUB:
			return ctx.encodeR4(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.regF(arg.Rs3)), nil
		case OpBase_NMSUB:
			return ctx.encodeR4(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.regF(arg.Rs3)), nil
		case OpBase_NMADD:
			return ctx.encodeR4(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), ctx.regF(arg.Rs3)), nil
		default:
			panic("unreachable")
		}

	case I:
		if arg.Rs2 != 0 {
			panic("encodeI: rs2 was nonzero")
		}
		if arg.Rs3 != 0 {
			panic("encodeI: rs3 was nonzero")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_OP_IMM:
			// 是否是移位命令
			// 移位是 I 格式的变化版本
			if ctx.HasShamt {
				switch as {
				// SLLI/SRLI/SRAI 只有这3个指令在 RV32/RV64 中不同
				case ASLLI, ASRLI, ASRAI:
					switch xlen {
					case 32:
						if arg.Imm < 0 || arg.Imm > 0b_1_1111 {
							panic("encodeI: imm(shamt5bit) overflow")
						}
					case 64:
						if arg.Imm < 0 || arg.Imm > 0b_11_1111 {
							panic("encodeI: imm(shamt6bit) overflow")
						}
					default:
						panic("encodeI: xlen must be 32 or 64")
					}
				default:
					// 其他都是 5bit
					if arg.Imm < 0 || arg.Imm > 0b_1_1111 {
						panic("encodeI: imm(shamt5bit) overflow")
					}
				}
			}
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.Funct7<<5|uint32(arg.Imm)), nil
		case OpBase_OP_IMM_32:
			switch as {
			case AADDIW:
				if ctx.Funct7 != 0 {
					panic("encodeI: funct7 was nonzero")
				}
				if arg.Imm < -2048 || arg.Imm > 2047 {
					panic("encodeI: imm must be in [-2048, 2047]")
				}

				return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil

			case ASRAIW:
				if ctx.Funct7 != 0b_010_0000 {
					panic("encodeI: funct7 != 0b_010_0000")
				}
				if arg.Imm < 0 || arg.Imm > 0b_1_1111 {
					panic("encodeI: imm(shamt5bit) overflow")
				}
			case ASLLIW, ASRLIW:
				if arg.Imm < 0 || arg.Imm > 0b_1_1111 {
					panic("encodeI: imm(shamt5bit) overflow")
				}
			default:
				panic("unreachable")
			}
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), ctx.Funct7<<5|uint32(arg.Imm)), nil
		case OpBase_JALR:
			if ctx.Funct3 != 0 {
				panic("encodeI: funct3 was nonzero")
			}
			if ctx.Funct7 != 0 {
				panic("encodeI: funct7 was nonzero")
			}
			if arg.Imm < -2048 || arg.Imm > 2047 {
				panic("encodeI: imm must be in [-2048, 2047]")
			}
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil
		case OpBase_LOAD:
			if ctx.Funct7 != 0 {
				panic("encodeI: funct7 was nonzero")
			}
			if arg.Imm < -2048 || arg.Imm > 2047 {
				panic("encodeI: imm must be in [-2048, 2047]")
			}
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil
		case OpBase_LOAD_FP:
			if ctx.Funct7 != 0 {
				panic("encodeI: funct7 was nonzero")
			}
			return ctx.encodeI(ctx.regF(arg.Rd), ctx.regF(arg.Rs1), uint32(arg.Imm)), nil
		case OpBase_MISC_MEN:
			if ctx.Funct7 != 0 {
				panic("encodeI: funct7 was nonzero")
			}
			switch as {
			case AFENCE:
				if ctx.Funct3 != 0 {
					panic("encodeI: funct3 was nonzero")
				}
				// 检查 fm, pred, succ
			default:
				panic("unreachable")
			}
			return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil
		case OpBase_SYSTEM:
			if ctx.Funct7 != 0 {
				panic("encodeI: funct7 was nonzero")
			}
			switch as {
			case AECALL, AEBREAK:
				if ctx.Funct3 != 0 {
					panic("encodeI: funct3 was nonzero")
				}
				if arg.Rd != 0 {
					panic("encodeI: rd must be zero")
				}
				if arg.Rs1 != 0 {
					panic("encodeI: rs1 must be zero")
				}
			}

			switch as {
			case AECALL:
				if arg.Imm != 0b_0000_0000_0000 {
					panic("encodeI: imm must be 0b_0000_0000_0000")
				}
				return ctx.encodeI(0, 0, 0b_0000_0000_0000), nil
			case AEBREAK:
				if arg.Imm != 0b_0000_0000_0001 {
					panic("encodeI: imm must be 0b_0000_0000_0001")
				}
				return ctx.encodeI(0, 0, 0b_0000_0000_0001), nil
			default:
				if arg.Imm < 0 || arg.Imm > 4096 {
					panic("encodeI: imm(csr) must be in [0, 4096]")
				}
				return ctx.encodeI(ctx.regI(arg.Rd), ctx.regI(arg.Rs1), uint32(arg.Imm)), nil
			}
		default:
			panic("unreachable")
		}
	case S:
		if ctx.Funct7 != 0 {
			panic("encodeS: funct7 was nonzero")
		}
		if arg.Rd != 0 {
			panic("encodeS: rd was nonzero")
		}
		if arg.Rs3 != 0 {
			panic("encodeS: rs3 was nonzero")
		}
		if arg.Imm < -2048 || arg.Imm > 2047 {
			panic("encodeI: imm must be in [-2048, 2047]")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_STORE:
			return ctx.encodeS(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), uint32(arg.Imm)), nil
		case OpBase_STORE_FP:
			return ctx.encodeS(ctx.regF(arg.Rs1), ctx.regF(arg.Rs2), uint32(arg.Imm)), nil
		default:
			panic("unreachable")
		}
	case B:
		if ctx.Funct7 != 0 {
			panic("encodeB: funct7 was nonzero")
		}
		if arg.Rs3 != 0 {
			panic("encodeB: rs3 was nonzero")
		}
		if arg.Imm < -4096 || arg.Imm > 4095 {
			panic("encodeB: imm must be in [-4096, 4095]")
		}
		if arg.Imm%2 != 0 {
			panic("encodeB: imm must align 2bytes")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_BRANCH:
			return ctx.encodeB(ctx.regI(arg.Rs1), ctx.regI(arg.Rs2), uint32(arg.Imm)), nil
		default:
			panic("unreachable")
		}
	case U:
		if ctx.Funct3 != 0 {
			panic("encodeU: funct3 was nonzero")
		}
		if ctx.Funct7 != 0 {
			panic("encodeU: funct7 was nonzero")
		}
		if arg.Rs1 != 0 {
			panic("encodeU: rs1 was nonzero")
		}
		if arg.Rs2 != 0 {
			panic("encodeU: rs2 was nonzero")
		}
		if arg.Rs3 != 0 {
			panic("encodeU: rs3 was nonzero")
		}

		// imm 对应 imm20, 不包含低 12bit
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_LUI:
			return ctx.encodeU(ctx.regI(arg.Rd), uint32(arg.Imm)), nil
		case OpBase_AUIPC:
			return ctx.encodeU(ctx.regI(arg.Rd), uint32(arg.Imm)), nil
		default:
			panic("unreachable")
		}
	case J:
		if ctx.Funct3 != 0 {
			panic("encodeJ: funct3 was nonzero")
		}
		if ctx.Funct7 != 0 {
			panic("encodeJ: funct7 was nonzero")
		}
		if arg.Rs1 != 0 {
			panic("encodeJ: rs1 was nonzero")
		}
		if arg.Rs2 != 0 {
			panic("encodeJ: rs2 was nonzero")
		}
		if arg.Rs3 != 0 {
			panic("encodeJ: rs3 was nonzero")
		}
		if min, max := -(1 << 21), 1<<21-1; arg.Imm < int32(min) || arg.Imm > int32(max) {
			panic(fmt.Sprintf("encodeU: imm must be in [%d, %d]", min, max))
		}
		if arg.Imm&0b1 != 0 {
			panic("encodeU: imm must be align 2 bytes")
		}
		switch ctx.Opcode & OpBase_Mask {
		case OpBase_JAL:
			return ctx.encodeJ(ctx.regI(arg.Rd), uint32(arg.Imm)), nil
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
	funct2 := ctx.Funct7
	return rs3<<27 | funct2<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// I-type
func (ctx *OpContextType) encodeI(rd, rs1 uint32, imm uint32) uint32 {
	return imm<<20 | rs1<<15 | ctx.Funct3<<12 | rd<<7 | uint32(ctx.Opcode)
}

// S-type
func (ctx *OpContextType) encodeS(rs1, rs2 uint32, imm uint32) uint32 {
	return (imm>>5)<<25 | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | (imm&0b_1_1111)<<7 | uint32(ctx.Opcode)
}

// B-type
func (ctx *OpContextType) encodeB(rs1, rs2 uint32, imm uint32) uint32 {
	return ctx.encodeB_Imm(imm) | rs2<<20 | rs1<<15 | ctx.Funct3<<12 | uint32(ctx.Opcode)
}
func (ctx *OpContextType) encodeB_Imm(imm uint32) uint32 {
	return (imm>>12)<<31 | ((imm>>5)&0x3f)<<25 | ((imm>>1)&0xf)<<8 | ((imm>>11)&0x1)<<7
}

// U-type
func (ctx *OpContextType) encodeU(rd uint32, imm uint32) uint32 {
	return (imm << 12) | rd<<7 | uint32(ctx.Opcode)
}

// J-type
func (ctx *OpContextType) encodeJ(rd uint32, imm uint32) uint32 {
	return ctx.encodeJ_Imm(imm) | rd<<7 | uint32(ctx.Opcode)
}
func (ctx *OpContextType) encodeJ_Imm(imm uint32) uint32 {
	return (imm>>20)<<31 | ((imm>>1)&0x3ff)<<21 | ((imm>>11)&0x1)<<20 | ((imm>>12)&0xff)<<12
}

// 返回寄存器机器码编号
func (ctx *OpContextType) regI(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_X0, REG_X31)
}

// 返回浮点数寄存器机器码编号
func (ctx *OpContextType) regF(r abi.RegType) uint32 {
	return ctx.regVal(r, REG_F0, REG_F31)
}

// 返回寄存器机器码编号
func (ctx *OpContextType) regVal(r, min, max abi.RegType) uint32 {
	if r < min || r > max {
		panic(fmt.Sprintf("register out of range, want %d <= %d <= %d", min, r, max))
	}
	return uint32(r - min)
}
