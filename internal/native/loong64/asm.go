package loong64

import "wa-lang.org/wa/internal/native/abi"

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
	panic("TODO")
}
