// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

// 编码后指令的长度, 忽略 Symbol 的具体值
func EncodeLen(as abi.As, arg *abi.X64Argument) (size int, err error) {
	panic("TODO")
}

// 指令编码, Symbol 对应的值需要提前解析到 Offset 属性中
func Encode(as abi.As, arg *abi.X64Argument) (code []byte, err error) {
	prog := &p9x86.Prog{}
	prog.To = operand2P9Addr(arg.Dst)
	prog.From = operand2P9Addr(arg.Src)

	switch as {
	case AMOV:
		// TODO: 根据寄存器类型选择 MOV 指令
		prog.As = p9x86.AMOVQ
	case AADD:
		prog.As = p9x86.AADDQ
	case ASUB:
		prog.As = p9x86.ASUBQ
	case APUSH:
		prog.As = p9x86.APUSHQ
	case APOP:
		prog.As = p9x86.APOPQ
	case ARET:
		prog.As = p9x86.ARET
	case ASYSCALL:
		prog.As = p9x86.ASYSCALL
	case ACALL:
		prog.As = p9x86.ACALL
	case AJMP:
		prog.As = p9x86.AJMP

	default:
		panic(fmt.Sprintf("TODO: %v", as))
	}

	code = p9x86.AsmInst(prog)
	return code, nil
}

func operand2P9Addr(op *abi.X64Operand) p9x86.Addr {
	if op == nil {
		return p9x86.Addr{Type: p9x86.TYPE_NONE}
	}

	addr := p9x86.Addr{}

	switch op.Kind {
	case abi.X64Operand_Reg:
		// 映射为寄存器类型
		addr.Type = p9x86.TYPE_REG
		addr.Reg = int16(op.Reg)

	case abi.X64Operand_Imm:
		// 映射为常量（立即数）
		// 例如：MOV $123, RAX 中的 $123
		addr.Type = p9x86.TYPE_CONST
		addr.Offset = op.Imm

	case abi.X64Operand_Mem:
		// 映射为内存引用
		// 这里的逻辑处理 [Reg + Offset]
		addr.Type = p9x86.TYPE_MEM
		addr.Reg = int16(op.Reg) // 基址寄存器 (Base)
		addr.Offset = op.Offset  // 位移 (Displacement)

		// 如果是 RIP 相对寻址 [rip + 0x10]，
		// 你的解析器应确保 op.Reg 是 REG_RIP 的编号
	}

	return addr
}
