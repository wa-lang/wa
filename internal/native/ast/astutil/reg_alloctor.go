// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package astutil

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/x64"
)

type _RegAlloctor struct {
	iRegList []abi.RegType
	fRegList []abi.RegType

	iRegIdx int
	fRegIdx int
}

// Linux/LA64 ABI 输入参数寄存器
func newArgsRegAlloctor_LA64Linux() *_RegAlloctor {
	return &_RegAlloctor{
		iRegList: []abi.RegType{
			loong64.REG_A0,
			loong64.REG_A1,
			loong64.REG_A2,
			loong64.REG_A3,
			loong64.REG_A4,
			loong64.REG_A5,
			loong64.REG_A6,
			loong64.REG_A7,
		},
		fRegList: []abi.RegType{
			loong64.REG_FA0,
			loong64.REG_FA1,
			loong64.REG_FA2,
			loong64.REG_FA3,
			loong64.REG_FA4,
			loong64.REG_FA5,
			loong64.REG_FA6,
			loong64.REG_FA7,
		},
	}
}

// Linux/LA64 ABI 返回值寄存器
func newRetRegAlloctor_LA64Linux() *_RegAlloctor {
	return &_RegAlloctor{
		iRegList: []abi.RegType{
			loong64.REG_A0,
			loong64.REG_A1,
		},
		fRegList: []abi.RegType{
			loong64.REG_FA0,
			loong64.REG_FA1,
		},
	}
}

// Windows/X64 ABI 输入参数寄存器
func newArgsRegAlloctor_X64Windows() *_RegAlloctor {
	return &_RegAlloctor{
		iRegList: []abi.RegType{
			x64.REG_RCX,
			x64.REG_RDX,
			x64.REG_R8,
			x64.REG_R9,
		},
		fRegList: []abi.RegType{
			x64.REG_XMM0,
			x64.REG_XMM1,
			x64.REG_XMM2,
			x64.REG_XMM3,
		},
	}
}

// Linux/X64 ABI 输入参数寄存器
func newArgsRegAlloctor_X64Linux() *_RegAlloctor {
	return &_RegAlloctor{
		iRegList: []abi.RegType{
			x64.REG_RDI,
			x64.REG_RSI,
			x64.REG_RDX,
			x64.REG_RCX,
			x64.REG_R8,
			x64.REG_R9,
		},
		fRegList: []abi.RegType{
			x64.REG_XMM0,
			x64.REG_XMM1,
			x64.REG_XMM2,
			x64.REG_XMM3,
			x64.REG_XMM4,
			x64.REG_XMM5,
			x64.REG_XMM6,
			x64.REG_XMM7,
		},
	}
}

// Linux/X64 ABI 返回值寄存器
func newRetRegAlloctor_X64Linux() *_RegAlloctor {
	return &_RegAlloctor{
		iRegList: []abi.RegType{
			x64.REG_RAX,
			x64.REG_RDX,
		},
		fRegList: []abi.RegType{
			x64.REG_XMM0,
			x64.REG_XMM1,
		},
	}
}

// 已经使用的总数(Windows ABI需要)
func (p *_RegAlloctor) UsedNum() int {
	return p.iRegIdx + p.fRegIdx
}

// 分配整数寄存器, 失败返回 0
func (p *_RegAlloctor) GetInt() (r abi.RegType) {
	if p.iRegIdx < len(p.iRegList) {
		r = p.iRegList[p.iRegIdx]
		p.iRegIdx++
	}
	return
}

// 分配浮点数寄存器, 失败返回 0
func (p *_RegAlloctor) GetFloat() (r abi.RegType) {
	if p.fRegIdx < len(p.fRegList) {
		r = p.fRegList[p.fRegIdx]
		p.fRegIdx++
	}
	return
}
