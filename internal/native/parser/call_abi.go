// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/token"
)

// 参数和返回值寄存器分配
type _AbiRegAlloctor struct {
	cpu abi.CPUType

	iArgRegEnd abi.RegType
	fArgRegEnd abi.RegType
	iRetRegEnd abi.RegType
	fRetRegEnd abi.RegType

	iArgReg abi.RegType
	fArgReg abi.RegType
	iRetReg abi.RegType
	fRetReg abi.RegType

	argRetOffset int // 栈上的偏移
	localOffset  int // 局部变量偏移
}

func (p *_AbiRegAlloctor) Reset() {
	switch p.cpu {
	case abi.LOONG64:
		p.iArgRegEnd = loong64.REG_A7
		p.fArgRegEnd = loong64.REG_F7
		p.iRetRegEnd = loong64.REG_A1
		p.fRetRegEnd = loong64.REG_F1

		p.iArgReg = loong64.REG_A0
		p.fArgReg = loong64.REG_F0
		p.iRetReg = loong64.REG_A0
		p.fRetReg = loong64.REG_F0

		p.argRetOffset = 0
		p.localOffset = 0

	case abi.RISCV32:
		// TODO
	case abi.RISCV64:
		// TODO
	default:
		panic("unreachable")
	}
}

func (p *_AbiRegAlloctor) AllocArg(typ token.Token) (argReg abi.RegType, argOff int) {
	switch p.cpu {
	case abi.LOONG64:
		switch typ {
		case token.I32, token.I32_zh:
			if p.iArgReg <= p.iArgRegEnd {
				argReg = p.iArgReg
				p.iArgReg++
			} else {
				argOff = p.argRetOffset
				p.argRetOffset += 4
			}
		case token.I64, token.I64_zh:
			if p.iArgReg <= p.iArgRegEnd {
				argReg = p.iArgReg
				p.iArgReg++
			} else {
				if p.argRetOffset%8 != 0 {
					p.argRetOffset += 4
				}
				argOff = p.argRetOffset
				p.argRetOffset += 8
			}
		case token.F32, token.F32_zh:
			if p.fArgReg <= p.fArgRegEnd {
				argReg = p.fArgReg
				p.fArgReg++
			} else {
				argOff = p.argRetOffset
				p.argRetOffset += 4
			}
		case token.F64, token.F64_zh:
			if p.fArgReg <= p.fArgRegEnd {
				argReg = p.fArgReg
				p.fArgReg++
			} else {
				if p.argRetOffset%8 != 0 {
					p.argRetOffset += 4
				}
				argOff = p.argRetOffset
				p.argRetOffset += 8
			}
		}
		return

	case abi.RISCV32:
		// TODO
		return 0, 0
	case abi.RISCV64:
		// TODO
		return 0, 0
	default:
		panic("unreachable")
	}
}

func (p *_AbiRegAlloctor) AllocRet(typ token.Token) (retReg abi.RegType, retOff int) {
	switch p.cpu {
	case abi.LOONG64:
		switch typ {
		case token.I32, token.I32_zh:
			if p.iRetReg <= p.iRetRegEnd {
				retReg = p.iRetReg
				p.iRetReg++
			} else {
				retOff = p.argRetOffset
				p.argRetOffset += 4
			}
		case token.I64, token.I64_zh:
			if p.iRetReg <= p.iRetRegEnd {
				retReg = p.iRetReg
				p.iRetReg++
			} else {
				if p.argRetOffset%8 != 0 {
					p.argRetOffset += 4
				}
				retOff = p.argRetOffset
				p.argRetOffset += 8
			}
		case token.F32, token.F32_zh:
			if p.fRetReg <= p.fRetRegEnd {
				retReg = p.fRetReg
				p.fArgReg++
			} else {
				retOff = p.argRetOffset
				p.argRetOffset += 4
			}
		case token.F64, token.F64_zh:
			if p.fRetReg <= p.fRetRegEnd {
				retReg = p.fRetReg
				p.fRetReg++
			} else {
				if p.argRetOffset%8 != 0 {
					p.argRetOffset += 4
				}
				retOff = p.argRetOffset
				p.argRetOffset += 8
			}
		}
		return

	case abi.RISCV32:
		// TODO
		return 0, 0
	case abi.RISCV64:
		// TODO
		return 0, 0
	default:
		panic("unreachable")
	}
}

func (p *_AbiRegAlloctor) AllocLocal(typ token.Token, cap int) (off int) {
	switch p.cpu {
	case abi.LOONG64:
		switch typ {
		case token.I32, token.I32_zh:
			off = p.localOffset
			p.localOffset += 4 * cap
		case token.I64, token.I64_zh:
			if p.localOffset%8 != 0 {
				p.localOffset += 4
			}
			off = p.localOffset
			p.localOffset += 8 * cap
		case token.F32, token.F32_zh:
			off = p.localOffset
			p.localOffset += 4 * cap
		case token.F64, token.F64_zh:
			if p.localOffset%8 != 0 {
				p.localOffset += 4
			}
			off = p.localOffset
			p.localOffset += 8 * cap
		}
		return

	case abi.RISCV32:
		// TODO
		return 0
	case abi.RISCV64:
		// TODO
		return 0
	default:
		panic("unreachable")
	}
}
