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
	switch as {
	case AMOV:
		return len("bf7b000000") / 2, nil
	case ASYSCALL:
		return len("0f05") / 2, nil
	default:
		panic("TODO")
	}
}

// 指令编码, Symbol 对应的值需要提前解析到 Offset 属性中
func Encode(as abi.As, arg *abi.X64Argument) (code []byte, err error) {
	prog := &p9x86.Prog{}
	prog.To = operand2P9Addr(arg.Dst)
	prog.From = operand2P9Addr(arg.Src)

	switch as {
	case AMOV:
		assert(arg.Src.Kind == abi.X64Operand_Imm)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)

		assert(prog.From.Type == p9x86.TYPE_CONST)
		assert(prog.To.Type == p9x86.TYPE_REG)

		switch {
		case arg.Dst.Reg >= REG_EAX && arg.Dst.Reg <= REG_R15D:
			assert(prog.To.Reg >= p9x86.REG_AX && prog.To.Reg <= p9x86.REG_R15)
			assert(arg.Src.Imm == prog.From.Offset)
			prog.As = p9x86.AMOVL
		case arg.Dst.Reg >= REG_RAX && arg.Dst.Reg <= REG_R15:
			assert(prog.To.Reg >= p9x86.REG_AX && prog.To.Reg <= p9x86.REG_R15)
			assert(arg.Src.Imm == prog.From.Offset)
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
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
		assert(arg.Dst.Kind == 0)
		assert(arg.Src.Kind == 0)
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

func reg2p9Reg(r abi.RegType) int16 {
	switch {
	case r >= REG_EAX && r <= REG_R15D:
		return p9x86.REG_AX + int16(r-REG_EAX)
	case r >= REG_RAX && r <= REG_R15:
		return p9x86.REG_AX + int16(r-REG_RAX)
	default:
		panic("unreachable")
	}
}

func operand2P9Addr(op *abi.X64Operand) p9x86.Addr {
	if op == nil {
		return p9x86.Addr{Type: p9x86.TYPE_NONE}
	}

	addr := p9x86.Addr{}

	switch op.Kind {
	case abi.X64Operand_Reg:
		addr.Type = p9x86.TYPE_REG
		addr.Reg = reg2p9Reg(op.Reg)

	case abi.X64Operand_Imm:
		addr.Type = p9x86.TYPE_CONST
		addr.Offset = op.Imm

	case abi.X64Operand_Mem:
		addr.Type = p9x86.TYPE_MEM // lea/jmp 等需要再次修复
		addr.Reg = reg2p9Reg(op.Reg)
		addr.Offset = op.Offset
	}

	return addr
}
