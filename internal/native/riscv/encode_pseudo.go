// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "wa-lang.org/wa/internal/native/abi"

// 伪指令只产生1个机器指令
func (ctx *_OpContextType) encodePseudo(xlen int, as abi.As, arg *abi.AsArgument) (uint32, error) {
	if ctx.PseudoAs == 0 {
		panic("unreachable")
	}

	// 检查模板和参数
	if err := ctx.checkArgMarks(xlen, as, arg, ctx.ArgMarks); err != nil {
		return 0, err
	}

	// 避免在 case 直接使用 ctx.PseudoAs
	// 便于在单元测试交叉验证两组数据一致性
	switch as {
	default:
		panic("unreachable")
	case A_NOP:
		// nop => addi x0, x0, 0
		return ctx.encodeRaw(xlen, AADDI, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: REG_X0,
			Imm: 0,
		})
	case A_MV:
		// mv rd, rs1 => addi rd, rs1, 0
		return ctx.encodeRaw(xlen, AADDI, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Imm: 0,
		})
	case A_NOT:
		// not rd, rs1 => xori rd, rs1, -1
		return ctx.encodeRaw(xlen, AXORI, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Imm: -1,
		})
	case A_NEG:
		// neg rd, rs1 => sub rd, x0, rs1
		return ctx.encodeRaw(xlen, ASUB, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Rs2: arg.Rs1,
		})
	case A_NEGW:
		// negw rd, rs1 => subw rd, x0, rs1
		return ctx.encodeRaw(xlen, ASUBW, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Rs2: arg.Rs1,
		})
	case A_SEXT_W:
		// sext.w rd, rs1 => addiw rd, rs1, 0
		return ctx.encodeRaw(xlen, AADDIW, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Imm: 0,
		})
	case A_SEQZ:
		// seqz rd, rs1 => sltiu rd, rs1, 1
		return ctx.encodeRaw(xlen, ASLTIU, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Imm: 1,
		})
	case A_SNEZ:
		// snez rd, rs1 => sltu rd, x0, rs1
		return ctx.encodeRaw(xlen, ASLTU, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Rs2: arg.Rs1,
		})
	case A_SLTZ:
		// sltz rd, rs1 => slt rd, rs1, x0
		return ctx.encodeRaw(xlen, ASLT, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: REG_X0,
		})
	case A_SGTZ:
		// sgtz rd, rs1 => slt rd, x0, rs1
		return ctx.encodeRaw(xlen, ASLT, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Rs2: arg.Rs1,
		})
	case A_FMV_S:
		// fmv.s rd, rs1 => fsgnj.s rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJ_S, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_FABS_S:
		// fabs.s rd, rs1 => fsgnjx.s rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJX_S, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_FNEG_S:
		// fneg.s rd, rs1 => fsgnjn.s rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJN_S, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_FMV_D:
		// fmv.d rd, rs1 => fsgnj.d rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJ_D, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_FABS_D:
		// fabs.d rd, rs1 => fsgnjx.d rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJX_D, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_FNEG_D:
		// fneg.d rd, rs1 => fsgnjn.d rd, rs1, rs1
		return ctx.encodeRaw(xlen, AFSGNJN_D, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: arg.Rs1,
			Rs2: arg.Rs1,
		})
	case A_BEQZ:
		// beqz rs1, offset => beq rs1, x0, offset
		return ctx.encodeRaw(xlen, ABEQ, &abi.AsArgument{
			Rs1: arg.Rs1,
			Rs2: REG_X0,
			Imm: arg.Imm,
		})
	case A_BNEZ:
		// bnez rs1, offset => bne rs1, x0, offset
		return ctx.encodeRaw(xlen, ABNE, &abi.AsArgument{
			Rs1: arg.Rs1,
			Rs2: REG_X0,
			Imm: arg.Imm,
		})
	case A_BLEZ:
		// blez rs1, offset => bge x0, rs1, offset
		return ctx.encodeRaw(xlen, ABGE, &abi.AsArgument{
			Rs1: REG_X0,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_BGEZ:
		// bgez rs1, offset => bge rs1, x0, offset
		return ctx.encodeRaw(xlen, ABGE, &abi.AsArgument{
			Rs1: arg.Rs1,
			Rs2: REG_X0,
			Imm: arg.Imm,
		})
	case A_BLTZ:
		// bltz rs1, offset => blt rs1, x0, offset
		return ctx.encodeRaw(xlen, ABLT, &abi.AsArgument{
			Rs1: arg.Rs1,
			Rs2: REG_X0,
			Imm: arg.Imm,
		})
	case A_BGTZ:
		// bgtz rs1, offset => blt x0, rs1, offset
		return ctx.encodeRaw(xlen, ABLT, &abi.AsArgument{
			Rs1: REG_X0,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_BGT:
		// bgt rs1, rs2, offset => blt rs2, rs1, offset
		return ctx.encodeRaw(xlen, ABLT, &abi.AsArgument{
			Rs1: arg.Rs2,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_BLE:
		// ble rs1, rs2, offset => bge rs2, rs1, offset
		return ctx.encodeRaw(xlen, ABGE, &abi.AsArgument{
			Rs1: arg.Rs2,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_BGTU:
		// bgtu rs1, rs2, offset => bltu rs2, rs1, offset
		return ctx.encodeRaw(xlen, ABLTU, &abi.AsArgument{
			Rs1: arg.Rs2,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_BLEU:
		// bleu rs1, rs2, offset => bgeu rs2, rs1, offset
		return ctx.encodeRaw(xlen, ABGEU, &abi.AsArgument{
			Rs1: arg.Rs2,
			Rs2: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_J:
		// j offset => jal x0, offset
		return ctx.encodeRaw(xlen, AJAL, &abi.AsArgument{
			Rd:  REG_X0,
			Imm: arg.Imm,
		})
	case A_JR:
		// jr rs1 => jalr x0, 0(rs1)
		return ctx.encodeRaw(xlen, AJALR, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_RET:
		// ret => jalr x0, 0(x1)
		return ctx.encodeRaw(xlen, AJALR, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: REG_X1,
			Imm: arg.Imm,
		})
	case A_RDINSTRET:
		// rdinstret rd => csrrs rd, instret, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: -1022,
		})
	case A_RDCYCLE:
		// rdcyle rd => csrrs rd, cycle, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: -1024,
		})
	case A_RDTIME:
		// rdtime rd => csrrs rd, time, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: -1023,
		})
	case A_CSRR:
		// csrr rd, csr => csrrs rd, csr, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: arg.Imm,
		})
	case A_CSRW:
		// csrw csr, rs1 => csrrw x0, csr, rs1
		return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_CSRS:
		// csrw csr, rs1 => csrrs x0, csr, rs1
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_CSRC:
		// csrc csr, rs1 => csrrc x0, csr, rs1
		return ctx.encodeRaw(xlen, ACSRRC, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_CSRWI:
		// csrwi csr, imm => csrrwi x0 csr, imm
		return ctx.encodeRaw(xlen, ACSRRWI, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_CSRSI:
		// csrsi csr, imm => csrrsi x0 csr, imm
		return ctx.encodeRaw(xlen, ACSRRSI, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_CSRCI:
		// csrci csr, imm => csrrci x0 csr, imm
		return ctx.encodeRaw(xlen, ACSRRCI, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: arg.Rs1,
			Imm: arg.Imm,
		})
	case A_FRCSR:
		// frcsr rd => csrrs rd, fcsr, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: 0x003,
		})
	case A_FSCSR:
		// fscsr rd, rs1 => csrrw rd, fcsr, rs1 # rd 可省略
		if arg.Rd == 0 {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  REG_X0,
				Rs1: arg.Rs1,
				Imm: 0x003,
			})
		} else {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  arg.Rd,
				Rs1: arg.Rs1,
				Imm: 0x003,
			})
		}
	case A_FRRM:
		// frrm rd => csrrs rd, frm, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: 0x002,
		})
	case A_FSRM:
		// fsrm rd, rs1 => csrrw rd, frm, rs1 # rd 可省略
		if arg.Rd == 0 {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  REG_X0,
				Rs1: arg.Rs1,
				Imm: 0x002,
			})
		} else {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  arg.Rd,
				Rs1: arg.Rs1,
				Imm: 0x002,
			})
		}
	case A_FRFLAGS:
		// frflags rd => csrrs rd, fflags, x0
		return ctx.encodeRaw(xlen, ACSRRS, &abi.AsArgument{
			Rd:  arg.Rd,
			Rs1: REG_X0,
			Imm: 0x001,
		})
	case A_FSFLAGS:
		// fsflags rd, rs1 => csrrw rd, fflags, rs1 # rd 可省略
		if arg.Rd == 0 {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  REG_X0,
				Rs1: arg.Rs1,
				Imm: 0x001,
			})
		} else {
			return ctx.encodeRaw(xlen, ACSRRW, &abi.AsArgument{
				Rd:  arg.Rd,
				Rs1: arg.Rs1,
				Imm: 0x001,
			})
		}
	}
}
