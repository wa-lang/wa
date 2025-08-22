// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// TODO: ADDI 打印格式不一致
// TODO: 根据编码后的指令反汇编出汇编程序, 然后通过 riscv 工具链编译链接执行验证

// 汇编语言格式
func AsmSyntax(pc int64, as As, arg *AsArgument) string {
	ctx := &AOpContextTable[as]
	return ctx.asmSyntax(pc, as, arg, RegString)
}

// 汇编语言格式, 自定义寄存器名字
func AsmSyntaxEx(pc int64, as As, arg *AsArgument, fnRegName func(r RegType) string) string {
	ctx := &AOpContextTable[as]
	return ctx.asmSyntax(pc, as, arg, fnRegName)
}

func (ctx *OpContextType) asmSyntax(pc int64, as As, arg *AsArgument, rName func(r RegType) string) string {
	switch ctx.Opcode.FormatType() {
	case R:
		return fmt.Sprintf("%s %s, %s, %s", as, rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2))
	case R4:
		return fmt.Sprintf("%s %s, %s, %s, %s", as, rName(arg.Rd), rName(arg.Rs1), rName(arg.Rs2), rName(arg.Rs3))
	case I:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s(%s)", as, rName(arg.Rd), arg.ImmName, rName(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %d(%s)", as, rName(arg.Rd), arg.Imm, rName(arg.Rs1))
	case S:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s(%s)", as, rName(arg.Rs2), arg.ImmName, rName(arg.Rs1))
		}
		return fmt.Sprintf("%s %s, %d(%s)", as, rName(arg.Rs2), arg.Imm, rName(arg.Rs1))
	case B:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s, %s", as, rName(arg.Rs1), rName(arg.Rs2), arg.ImmName)
		}
		return fmt.Sprintf("%s %s, %s, 0x%X", as, rName(arg.Rs1), rName(arg.Rs2), pc+int64(arg.Imm))
	case U:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s", as, rName(arg.Rd), arg.ImmName)
		}
		return fmt.Sprintf("%s %s, 0x%X", as, rName(arg.Rd), arg.Imm)
	case J:
		if arg.ImmName != "" {
			return fmt.Sprintf("%s %s, %s", as, rName(arg.Rd), arg.ImmName)
		}
		return fmt.Sprintf("%s %s, 0x%X", as, rName(arg.Rd), pc+int64(arg.Imm))
	default:
		return AsString(as)
	}
}
