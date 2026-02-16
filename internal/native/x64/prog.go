// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

// 底层的指令别名
type Prog p9x86.Prog

// 初始化b编码表格
func init() {
	p9x86.Init()
}

// 从解析器获得的参数构建底层的指令
// 如果有标识符, 默认均以转化为了具体的数值
func BuildProg(as abi.As, arg *abi.X64Argument) (inst *Prog, err error) {
	return new(Prog).buildProg(as, arg)
}

func (p *Prog) Encode() []byte {
	return ((*p9x86.Prog)(p)).Encode()
}

func (p *Prog) nArg(arg *abi.X64Argument) int {
	switch {
	case arg.Dst != nil && arg.Src != nil:
		return 2
	case arg.Dst != nil && arg.Src == nil:
		return 1
	case arg.Dst == nil && arg.Src == nil:
		return 0
	default:
		panic("unreachable")
	}
}

func (p *Prog) xLen(arg *abi.X64Argument) int {
	dstXlen := p.operandXlen(arg.Dst)
	srcXlen := p.operandXlen(arg.Src)
	if dstXlen > srcXlen {
		return dstXlen
	} else {
		return srcXlen
	}
}

func (p *Prog) reg2p9Reg(r abi.RegType) int16 {
	switch {
	case r >= REG_EAX && r <= REG_R15D:
		return p9x86.REG_AX + int16(r-REG_EAX)
	case r >= REG_RAX && r <= REG_R15:
		return p9x86.REG_AX + int16(r-REG_RAX)
	default:
		panic("unreachable")
	}
}

func (p *Prog) operandXlen(op *abi.X64Operand) int {
	if op == nil {
		return 1
	}
	switch op.Kind {
	case abi.X64Operand_Reg:
		return RegXLen(op.Reg)
	case abi.X64Operand_Mem:
		switch op.PtrTyp {
		case abi.X64BytePtr:
			return 1
		case abi.X64WordPtr:
			return 2
		case abi.X64DWordPtr:
			return 4
		case abi.X64QWordPtr:
			return 8
		default:
			panic("unreachable")
		}
	case abi.X64Operand_Imm:
		return 1 // 立即数按照单字节计算, 后续再修复
	default:
		panic("unreachable")
	}
}

func (p *Prog) operand2P9Addr(op *abi.X64Operand) p9x86.Addr {
	if op == nil {
		return p9x86.Addr{Type: p9x86.TYPE_NONE}
	}

	addr := p9x86.Addr{}

	switch op.Kind {
	case abi.X64Operand_Reg:
		addr.Type = p9x86.TYPE_REG
		addr.Reg = p.reg2p9Reg(op.Reg)

	case abi.X64Operand_Imm:
		addr.Type = p9x86.TYPE_CONST
		addr.Offset = op.Imm

	case abi.X64Operand_Mem:
		addr.Type = p9x86.TYPE_MEM // lea/jmp 等需要再次修复
		addr.Reg = p.reg2p9Reg(op.Reg)
		addr.Offset = op.Offset
	}

	return addr
}
