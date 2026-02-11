// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 汇编语言格式(用别名显式寄存器)
func AsmSyntax(as abi.As, asName string, arg *abi.X64Argument) string {
	switch {
	case arg.Dst == nil && arg.Src == nil:
		return AsString(as, asName)
	case arg.Dst != nil && arg.Src == nil:
		return fmt.Sprintf("%s %s", AsString(as, asName), AsmOperand(arg.Dst))
	case arg.Dst != nil && arg.Src != nil:
		return fmt.Sprintf("%s %s, %s", AsString(as, asName), AsmOperand(arg.Dst), AsmOperand(arg.Src))
	default:
		panic("unreachable")
	}
}

// 格式化操作数
func AsmOperand(op *abi.X64Operand) string {
	switch op.Kind {
	case abi.X64X64Operand_Reg:
		return RegString(op.Reg)
	case abi.X64X64Operand_Imm:
		return fmt.Sprintf("%d", op.Imm)
	case abi.X64X64Operand_Mem:
		// 处理指针前缀
		ptrPrefix := ""
		switch op.PtrTyp {
		case abi.X64BytePtr:
			ptrPrefix = "byte ptr "
		case abi.X64DWordPtr:
			ptrPrefix = "dword ptr "
		case abi.X64QWordPtr:
			ptrPrefix = "qword ptr "
		}

		// 处理地址部分 [base + symbol + offset]
		addr := ""
		if op.Reg != 0 {
			addr = RegString(op.Reg)
		}
		if op.Symbol != "" {
			if addr != "" {
				addr += " + "
			}
			addr += op.Symbol
		}
		if op.Offset != 0 {
			if addr != "" && op.Offset > 0 {
				addr += " + "
			}
			addr += fmt.Sprintf("%d", op.Offset)
		}
		return fmt.Sprintf("%s[%s]", ptrPrefix, addr)
	}
	return ""
}
