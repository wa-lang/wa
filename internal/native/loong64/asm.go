package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 汇编语言格式(用别名显式寄存器)
func AsmSyntax(as abi.As, asName string, arg *abi.AsArgument) string {
	ctx := &_AOpContextTable[as]
	return ctx.asmSyntax(as, asName, arg, RegAliasString, AsString)
}

// 汇编语言格式, 自定义寄存器和指令名字
func AsmSyntaxEx(
	as abi.As, asName string, arg *abi.AsArgument,
	fnRegName func(r abi.RegType) string,
	fnAsName func(x abi.As, xName string) string,
) string {
	ctx := &_AOpContextTable[as]
	return ctx.asmSyntax(as, asName, arg, fnRegName, fnAsName)
}

func (ctx *_OpContextType) asmSyntax(
	as abi.As, asName string, arg *abi.AsArgument,
	rName func(r abi.RegType) string,
	asNameFn func(x abi.As, xName string) string,
) string {
	symbol := arg.Symbol
	if arg.SymbolDecor != 0 {
		symbol = fmt.Sprintf("%v(%s)", arg.SymbolDecor, arg.Symbol)
	}

	// 辅助函数：处理 Imm/Symbol 作为偏移量 (Offset) 的格式，如 Imm(Rs1)
	formatOffset := func(rs1 abi.RegType) string {
		if symbol != "" {
			return fmt.Sprintf("%v(%s)", symbol, rName(rs1))
		} else {
			// 假设 Arg.Imm 是一个整数偏移量
			return fmt.Sprintf("%d(%s)", arg.Imm, rName(rs1))
		}
	}

	// 辅助函数：处理 Imm/Symbol 作为立即数 (Immediate) 或目标 (Target) 的格式
	formatImm := func() string {
		if symbol != "" {
			return symbol
		} else {
			// 默认以十进制格式化立即数，如需十六进制请自行调整
			return fmt.Sprintf("%d", arg.Imm)
		}
	}

	// 辅助函数：处理 Imm/Symbol 作为地址 (Address) 的格式
	formatAddr := func() string {
		if symbol != "" {
			return symbol
		} else {
			// 地址或大立即数通常用十六进制表示
			return fmt.Sprintf("0x%X", arg.Imm)
		}
	}

	switch ctx.fmt {
	case OpFormatType_NULL:
		return asNameFn(as, asName)
	case OpFormatType_2R:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1))
	case OpFormatType_2F:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1))
	case OpFormatType_1F_1R:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1))
	case OpFormatType_1R_1F:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1))
	case OpFormatType_3R:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_3F:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_1F_2R:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_4F:
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), rName(arg.Rs3))
	case OpFormatType_2R_ui5:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), formatOffset(arg.Rs1))
	case OpFormatType_2R_ui6:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatImm())
	case OpFormatType_2R_si12:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatImm())
	case OpFormatType_2R_ui12:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatImm())
	case OpFormatType_2R_si14:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatImm())
	case OpFormatType_2R_si16:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatImm())
	case OpFormatType_1R_si20:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), formatAddr())
	case OpFormatType_0_2R:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_3R_sa2:
		return fmt.Sprintf("%s %s, %s, %s, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), arg.Imm)
	case OpFormatType_3R_sa3:
		return fmt.Sprintf("%s %s, %s, %s, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), arg.Imm)
	case OpFormatType_code:
		return fmt.Sprintf("%s %d", asNameFn(as, asName), arg.Imm)
	case OpFormatType_code_1R_si12:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), formatImm())
	case OpFormatType_2R_msbw_lsbw:
		return fmt.Sprintf("%s %s, %s, %d, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), arg.Imm, arg.Imm2)
	case OpFormatType_2R_msbd_lsbd:
		return fmt.Sprintf("%s %s, %s, %d, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), arg.Imm, arg.Imm2)
	case OpFormatType_fcsr_1R:
		return fmt.Sprintf("%s %d, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1))
	case OpFormatType_1R_fcsr:
		return fmt.Sprintf("%s %s, %d", asNameFn(as, asName), rName(arg.Rd), arg.Imm)
	case OpFormatType_cd_1R:
		return fmt.Sprintf("%s %d, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1))
	case OpFormatType_cd_1F:
		return fmt.Sprintf("%s %d, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1))
	case OpFormatType_cd_2R:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_cd_2F:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_1R_cj:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), formatAddr())
	case OpFormatType_1F_cj:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), formatAddr())
	case OpFormatType_1R_csr:
		return fmt.Sprintf("%s %s, %d", asNameFn(as, asName), rName(arg.Rd), arg.Imm)
	case OpFormatType_2R_csr:
		return fmt.Sprintf("%s %s, %s, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), arg.Imm)
	case OpFormatType_2R_level:
		return fmt.Sprintf("%s %s, %s, %d", asNameFn(as, asName), rName(arg.Rs1), rName(arg.Rs2), arg.Imm)
	case OpFormatType_level:
		return fmt.Sprintf("%s %d", asNameFn(as, asName), arg.Imm)
	case OpFormatType_0_1R_seq:
		return fmt.Sprintf("%s %s, %d", asNameFn(as, asName), rName(arg.Rs1), arg.Imm)
	case OpFormatType_op_2R:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_3R_ca:
		return fmt.Sprintf("%s %s, %s, %s, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), arg.Imm)
	case OpFormatType_hint_1R_si12:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), formatImm())
	case OpFormatType_hint_2R:
		return fmt.Sprintf("%s %d, %s, %s", asNameFn(as, asName), arg.Imm, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_hint:
		return fmt.Sprintf("%s %d", asNameFn(as, asName), arg.Imm)
	case OpFormatType_cj_offset:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), formatAddr())
	case OpFormatType_rj_offset:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rs1), formatOffset(arg.Rs1))
	case OpFormatType_rj_rd_offset:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatOffset(arg.Rs1))
	case OpFormatType_rd_rj_offset:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), formatOffset(arg.Rs1))
	case OpFormatType_offset:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), formatAddr())

	default:
		panic("unreachable")
	}
}
