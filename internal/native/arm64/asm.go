// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import "wa-lang.org/wa/internal/native/abi"

// 汇编语言格式(用别名显式寄存器)
func AsmSyntax(as abi.As, asName string, arg *abi.AsArgument) string {
	//ctx := &_AOpContextTable[as]
	//return ctx.asmSyntax(as, asName, arg, RegAliasString, AsString)
	panic("TODO")
}

// 汇编语言格式, 自定义寄存器和指令名字
func AsmSyntaxEx(
	as abi.As, asName string, arg *abi.AsArgument,
	fnRegName func(r abi.RegType) string,
	fnAsName func(x abi.As, xName string) string,
) string {
	//ctx := &_AOpContextTable[as]
	//return ctx.asmSyntax(as, asName, arg, fnRegName, fnAsName)
	panic("TODO")
}
