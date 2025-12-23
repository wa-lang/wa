// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package framestack

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

type RVFramestack struct {
	iArgReg abi.RegType
	fArgReg abi.RegType
	iRetReg abi.RegType
	fRetReg abi.RegType

	argsOffset  int // 参数和返回值在栈上的偏移
	localOffset int // 局部变量偏移
}

func NewRVFramestack() *RVFramestack {
	return &RVFramestack{
		iArgReg: riscv.REG_A0,
		fArgReg: riscv.REG_FA0,
		iRetReg: riscv.REG_A0,
		fRetReg: riscv.REG_FA0,
	}
}

func (p *RVFramestack) AllocArg(typ token.Token) (reg abi.RegType, off int) {
	switch typ {
	case token.I32, token.U32, token.I32_zh, token.U32_zh:
		if p.iArgReg <= riscv.REG_A7 {
			reg = p.iArgReg
			p.iArgReg++
		} else {
			off = p.argsOffset
			p.argsOffset += 4
		}
	case token.I64, token.U64, token.I64_zh, token.U64_zh:
		if p.iArgReg <= riscv.REG_A7 {
			reg = p.iArgReg
			p.iArgReg++
		} else {
			if p.argsOffset%8 != 0 {
				p.argsOffset += 4
			}
			off = p.argsOffset
			p.argsOffset += 8
		}
	case token.F32, token.F32_zh:
		if p.fArgReg <= riscv.REG_FA7 {
			reg = p.fArgReg
			p.fArgReg++
		} else {
			off = p.argsOffset
			p.argsOffset += 4
		}
	case token.F64, token.F64_zh:
		if p.fArgReg <= riscv.REG_FA7 {
			reg = p.fArgReg
			p.fArgReg++
		} else {
			if p.argsOffset%8 != 0 {
				p.argsOffset += 4
			}
			off = p.argsOffset
			p.argsOffset += 8
		}
	default:
		panic("unreachable")
	}
	return
}

func (p *RVFramestack) AllocRet(typ token.Token) (reg abi.RegType, off int) {
	switch typ {
	case token.I32, token.U32, token.I32_zh, token.U32_zh:
		if p.iRetReg <= riscv.REG_A7 {
			reg = p.iRetReg
			p.iRetReg++
		} else {
			off = p.argsOffset
			p.argsOffset += 4
		}
	case token.I64, token.U64, token.I64_zh, token.U64_zh:
		if p.iRetReg <= riscv.REG_A7 {
			reg = p.iRetReg
			p.iRetReg++
		} else {
			if p.argsOffset%8 != 0 {
				p.argsOffset += 4
			}
			off = p.argsOffset
			p.argsOffset += 8
		}
	case token.F32, token.F32_zh:
		if p.fRetReg <= riscv.REG_FA7 {
			reg = p.fRetReg
			p.fRetReg++
		} else {
			off = p.argsOffset
			p.argsOffset += 4
		}
	case token.F64, token.F64_zh:
		if p.fRetReg <= riscv.REG_FA7 {
			reg = p.fRetReg
			p.fRetReg++
		} else {
			if p.argsOffset%8 != 0 {
				p.argsOffset += 4
			}
			off = p.argsOffset
			p.argsOffset += 8
		}
	default:
		panic("unreachable")
	}
	return
}

func (p *RVFramestack) AllocLocal(typ token.Token, cap int) (off int) {
	switch typ {
	case token.I32, token.U32, token.I32_zh, token.U32_zh:
		off = p.localOffset
		p.localOffset += 4 * cap
	case token.I64, token.U64, token.I64_zh, token.U64_zh:
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
	default:
		panic("unreachable")
	}
	return
}
