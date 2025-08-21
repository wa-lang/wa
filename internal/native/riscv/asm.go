// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// 汇编语言格式
func AsmSyntax(as As, arg *AsArgument) string {
	ctx := &AOpContextTable[as]
	return ctx.asmSyntax(as, arg)
}

func (ctx *OpContextType) asmSyntax(as As, arg *AsArgument) string {
	switch ctx.Opcode.FormatType() {
	case R:
		return fmt.Sprintf("%s %s, %s, %s", as, arg.Rd, arg.Rs1, arg.Rs2)
	case R4:
		return fmt.Sprintf("%s %s, %s, %s, %s", as, arg.Rd, arg.Rs1, arg.Rs2, arg.Rs3)
	case I:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s(%s)", as, arg.Rd, arg.ImmName, arg.Rs1)
		}
		if arg.Imm != 0 {
			return fmt.Sprintf("%s %s, %d(%s)", as, arg.Rd, arg.Imm, arg.Rs1)
		}
		return fmt.Sprintf("%s %s, (%s)", as, arg.Rd, arg.Rs1)
	case S:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s(%s)", as, arg.Rs2, arg.ImmName, arg.Rs1)
		}
		if arg.Imm != 0 {
			return fmt.Sprintf("%s %s, %d(%s)", as, arg.Rs2, arg.Imm, arg.Rs1)
		}
		return fmt.Sprintf("%s %s, (%s)", as, arg.Rs2, arg.Rs1)
	case B:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s, %s", as, arg.Rs1, arg.Rs2, arg.ImmName)
		}
		return fmt.Sprintf("%s %s, %s, %d", as, arg.Rs1, arg.Rs2, arg.Imm)
	case U:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s, %s", as, arg.Rs1, arg.Rs2, arg.ImmName)
		}
		return fmt.Sprintf("%s %s, %s, %d", as, arg.Rs1, arg.Rs2, arg.Imm)
	case J:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s", as, arg.Rd, arg.ImmName)
		}
		return fmt.Sprintf("%s %s, %d", as, arg.Rd, arg.Imm)
	default:
		return AsString(as)
	}
}
