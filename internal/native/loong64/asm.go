// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	immValue := arg.Symbol
	if arg.SymbolDecor != 0 {
		immValue = fmt.Sprintf("%v(%s)", arg.SymbolDecor, arg.Symbol)
	}
	if arg.Symbol == "" {
		immValue = fmt.Sprint(arg.Imm)
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
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_2R_ui6:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_2R_si12:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_2R_ui12:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_2R_si14:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_1R_si20:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), immValue)
	case OpFormatType_0_2R:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_3R_sa2:
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), immValue)
	case OpFormatType_3R_sa3:
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), immValue)
	case OpFormatType_code:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), immValue)
	case OpFormatType_code_1R_si12:
		codeString := arg.RdName
		if codeString == "" {
			codeString = fmt.Sprint(int(arg.Rd))
		}
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), codeString, rName(arg.Rs1), immValue)
	case OpFormatType_2R_msbw_lsbw:
		msbw := arg.Rs2Name
		lsbw := arg.Rs3Name
		if msbw == "" {
			msbw = fmt.Sprint(int(arg.Rs2))
		}
		if lsbw == "" {
			lsbw = fmt.Sprint(int(arg.Rs3))
		}
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), msbw, lsbw)
	case OpFormatType_2R_msbd_lsbd:
		msbd := arg.Rs2Name
		lsbd := arg.Rs3Name
		if msbd == "" {
			msbd = fmt.Sprint(int(arg.Rs2))
		}
		if lsbd == "" {
			lsbd = fmt.Sprint(int(arg.Rs3))
		}
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), msbd, lsbd)
	case OpFormatType_fcsr_1R:
		fcsrSymbol := arg.RdName
		if fcsrSymbol == "" {
			fcsrSymbol = fmt.Sprint(int(arg.Rd))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), fcsrSymbol, rName(arg.Rs1))
	case OpFormatType_1R_fcsr:
		fcsrSymbol := arg.Rs1Name
		if fcsrSymbol == "" {
			fcsrSymbol = fmt.Sprint(int(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), fcsrSymbol)
	case OpFormatType_cd_1R:
		cdSymbol := arg.RdName
		if cdSymbol == "" {
			cdSymbol = fmt.Sprint(int(arg.Imm))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), cdSymbol, rName(arg.Rs1))
	case OpFormatType_cd_1F:
		cdSymbol := arg.RdName
		if cdSymbol == "" {
			cdSymbol = fmt.Sprint(int(arg.Imm))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), cdSymbol, rName(arg.Rs1))
	case OpFormatType_cd_2F:
		cdSymbol := arg.RdName
		if cdSymbol == "" {
			cdSymbol = fmt.Sprint(int(arg.Imm))
		}
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), cdSymbol, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_1R_cj:
		cjSymbol := arg.Rs1Name
		if cjSymbol == "" {
			cjSymbol = fmt.Sprint(int(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), cjSymbol)
	case OpFormatType_1F_cj:
		cjSymbol := arg.Rs1Name
		if cjSymbol == "" {
			cjSymbol = fmt.Sprint(int(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), cjSymbol)
	case OpFormatType_1R_csr:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), immValue)
	case OpFormatType_2R_csr:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_2R_level:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_level:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), immValue)
	case OpFormatType_0_1R_seq:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rs1), immValue)
	case OpFormatType_op_2R:
		opSymbol := arg.RdName
		if opSymbol == "" {
			opSymbol = fmt.Sprint(int(arg.Rd))
		}
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), opSymbol, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_3F_ca:
		return fmt.Sprintf("%s %s, %s, %s, %d", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), arg.Imm)
	case OpFormatType_hint_1R_si12:
		hintSymbol := arg.RdName
		if hintSymbol == "" {
			hintSymbol = fmt.Sprint(int(arg.Rd))
		}
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), hintSymbol, rName(arg.Rs1), immValue)
	case OpFormatType_hint_2R:
		hintSymbol := arg.RdName
		if hintSymbol == "" {
			hintSymbol = fmt.Sprint(int(arg.Rd))
		}
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), hintSymbol, rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_hint:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), immValue)
	case OpFormatType_cj_offset:
		cjSymbol := arg.Rs1Name
		if cjSymbol == "" {
			cjSymbol = fmt.Sprint(int(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), cjSymbol, immValue)
	case OpFormatType_rj_offset:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rs1), immValue)
	case OpFormatType_rj_rd_offset:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rs1), rName(arg.Rd), immValue)
	case OpFormatType_rd_rj_offset:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), immValue)
	case OpFormatType_offset:
		return fmt.Sprintf("%s %s", asNameFn(as, asName), immValue)

	default:
		panic("unreachable")
	}
}
