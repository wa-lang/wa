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

	switch ctx.FormatType() {
	case OpFormatType_2R:
		return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1))
	case OpFormatType_3R:
		return fmt.Sprintf("%s %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))
	case OpFormatType_4R:
		return fmt.Sprintf("%s %s, %s, %s, %s", asNameFn(as, asName), rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), rName(arg.Rs3))

	case OpFormatType_2RI8:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asNameFn(as, asName), rName(arg.Rd), symbol, rName(arg.Rs1))
		} else {
			return fmt.Sprintf("%s %s, %d(%s)", asNameFn(as, asName), rName(arg.Rd), arg.Imm, rName(arg.Rs1))
		}
	case OpFormatType_2RI12:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asNameFn(as, asName), rName(arg.Rd), symbol, rName(arg.Rs1))
		} else {
			return fmt.Sprintf("%s %s, %d(%s)", asNameFn(as, asName), rName(arg.Rd), arg.Imm, rName(arg.Rs1))
		}
	case OpFormatType_2RI14:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asNameFn(as, asName), rName(arg.Rd), symbol, rName(arg.Rs1))
		} else {
			return fmt.Sprintf("%s %s, %d(%s)", asNameFn(as, asName), rName(arg.Rd), arg.Imm, rName(arg.Rs1))
		}
	case OpFormatType_2RI16:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %v(%s)", asNameFn(as, asName), rName(arg.Rd), symbol, rName(arg.Rs1))
		} else {
			return fmt.Sprintf("%s %s, %d(%s)", asNameFn(as, asName), rName(arg.Rd), arg.Imm, rName(arg.Rs1))
		}

	case OpFormatType_1RI20:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), symbol)
		} else {
			return fmt.Sprintf("%s %s, 0x%X", asNameFn(as, asName), rName(arg.Rd), arg.Imm)
		}

	case OpFormatType_1RI21:
		if symbol != "" {
			return fmt.Sprintf("%s %s, %s", asNameFn(as, asName), rName(arg.Rd), symbol)
		} else {
			return fmt.Sprintf("%s %s, 0x%X", asNameFn(as, asName), rName(arg.Rd), arg.Imm)
		}

	case OpFormatType_I26:
		if symbol != "" {
			return fmt.Sprintf("%s  %s", asNameFn(as, asName), symbol)
		} else {
			return fmt.Sprintf("%s %d", asNameFn(as, asName), int64(arg.Imm))
		}

	default:
		panic(fmt.Sprintf("TODO: %v", ctx.FormatType()))
	}
}
