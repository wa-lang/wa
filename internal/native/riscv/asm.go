// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// TODO: ADDI 打印格式不一致
// TODO: 根据编码后的指令反汇编出汇编程序, 然后通过 riscv 工具链编译链接执行验证

// 汇编语言格式(用别名显式寄存器)
func AsmSyntax(as abi.As, arg *abi.AsArgument) string {
	ctx := &_AOpContextTable[as]
	return ctx.asmSyntax(as, arg, RegAliasString, AsString)
}

// 汇编语言格式, 自定义寄存器和指令名字
func AsmSyntaxEx(
	as abi.As, arg *abi.AsArgument,
	fnRegName func(r abi.RegType) string,
	fnAsName func(x abi.As) string,
) string {
	ctx := &_AOpContextTable[as]
	return ctx.asmSyntax(as, arg, fnRegName, fnAsName)
}

func (ctx *_OpContextType) asmSyntax(
	as abi.As, arg *abi.AsArgument,
	rName func(r abi.RegType) string,
	asName func(x abi.As) string,
) string {
	symbol := arg.Symbol
	if arg.SymbolDecor != 0 {
		symbol = fmt.Sprintf("%v(%s)", arg.SymbolDecor, arg.Symbol)
	}

	switch ctx.Opcode & _OpBase_Mask {
	case _OpBase_OP, _OpBase_OP_32, _OpBase_OP_FP, _OpBase_AMO:
		// XXX rd, rs1, rs2
		return fmt.Sprintf("%s %s, %s, %s", asName(as), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))

	case _OpBase_MADD, _OpBase_MSUB, _OpBase_NMSUB, _OpBase_NMADD:
		// XXX rd, rs1, rs2, rs3
		return fmt.Sprintf("%s %s, %s, %s, %s", asName(as), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), rName(arg.Rs3))

	case _OpBase_OP_IMM, _OpBase_OP_IMM_32:
		// XXX rd, rs1, imm
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s, %v", asName(as), rName(arg.Rd), rName(arg.Rs1), symbol)
		}
		return fmt.Sprintf("%s %s, %s, %d", asName(as), rName(arg.Rd), rName(arg.Rs1), arg.Imm)

	case _OpBase_LOAD, _OpBase_LOAD_FP:
		// XXX rd, offset(rs1)
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asName(as), rName(arg.Rd), symbol, rName(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %d(%s)", asName(as), rName(arg.Rd), arg.Imm, rName(arg.Rs1))

	case _OpBase_MISC_MEN:
		// fence iorw, iorw
		return fmt.Sprintf("%s %s, %s", asName(as), rName(arg.Rd), rName(arg.Rs1))

	case _OpBase_SYSTEM:
		switch as {
		case ACSRRW, ACSRRS, ACSRRC:
			// XXX rd, csr, rs1
			if symbol != "" {
				return fmt.Sprintf("%s %s, %d, %s", asName(as), rName(arg.Rd), arg.Imm, rName(arg.Rs1))
			}
			return fmt.Sprintf("%s %s, %d, %s", asName(as), rName(arg.Rd), arg.Imm, rName(arg.Rs1))

		case ACSRRWI, ACSRRSI, ACSRRCI:
			// XXX rd, csr, uimm
			if symbol != "" {
				return fmt.Sprintf("%s %s, %d, %d", asName(as), rName(arg.Rd), arg.Imm, ctx.regI(arg.Rs1))
			}
			return fmt.Sprintf("%s %s, %d, %d", asName(as), rName(arg.Rd), arg.Imm, ctx.regI(arg.Rs1))

		default:
			panic("unreachable")
		}

	case _OpBase_STORE, _OpBase_STORE_FP:
		// XXX rs2, offset(rs1)
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s(%s)", asName(as), rName(arg.Rs2), symbol, rName(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %d(%s)", asName(as), rName(arg.Rs2), arg.Imm, rName(arg.Rs1))

	case _OpBase_BRANCH:
		// XXX rs1, rs2, label
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s, %s", asName(as), rName(arg.Rs1), rName(arg.Rs2), symbol)
		}
		return fmt.Sprintf("%s %s, %s, %d", asName(as), rName(arg.Rs1), rName(arg.Rs2), int64(arg.Imm))

	case _OpBase_LUI, _OpBase_AUIPC:
		// XXX rd, imm
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s", asName(as), rName(arg.Rd), symbol)
		}
		return fmt.Sprintf("%s %s, 0x%X", asName(as), rName(arg.Rd), arg.Imm)

	case _OpBase_JAL:
		// XXX rd, offset
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s", asName(as), rName(arg.Rd), symbol)
		}
		return fmt.Sprintf("%s %s, %d", asName(as), rName(arg.Rd), int64(arg.Imm))

	case _OpBase_JALR:
		// XXX rd, offset(rs1)
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asName(as), rName(arg.Rd), symbol, rName(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %d(%s)", asName(as), rName(arg.Rd), arg.Imm, rName(arg.Rs1))

	case _OpBase_CUSTOM_0, _OpBase_CUSTOM_1, _OpBase_CUSTOM_2, _OpBase_CUSTOM_3:
		panic("unreachable")

	default:
		panic("unreachable")
	}
}
