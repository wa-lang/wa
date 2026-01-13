// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package astutil

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/native/x64"
)

// ABI: $(WA_REPO)/docs/asm_abi_x64.md

// 构建栈帧中参数和返回值的位置
func buildFuncFrame_x64_windows(fn *ast.Func) error {
	const headSize = 8 * 2 // RIP + RBP

	var (
		// 候选的寄存器列表
		iArgRegList = []abi.RegType{x64.REG_RCX, x64.REG_RDX, x64.REG_R8, x64.REG_R9}
		fArgRegList = []abi.RegType{x64.REG_XMM0, x64.REG_XMM1, x64.REG_XMM2, x64.REG_XMM3}

		iArgRegIdx int
		fArgRegIdx int

		sp int = 0 // 局部栈大小
	)

	// 计算输入参数需要的栈大小
	argsStackSize := len(fn.Type.Return) * 8

	// 影子空间是强制保留的
	if argsStackSize < 32 {
		argsStackSize = 32
	}

	// 调用其他函数的参数也需要 16 字节对齐
	if x := argsStackSize; x%16 != 0 {
		x = ((x + 15) / 16) * 16
		argsStackSize = x
	}

	// 返回值
	switch len(fn.Type.Return) {
	case 0:
	case 1:
		// 给返回值的影子空间
		sp -= 8

		// 是基于 rbp 的位置
		switch ret := fn.Type.Return[0]; ret.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			fn.Type.Return[0].Reg = x64.REG_RAX
			fn.Type.Return[0].Off = -8 // rbp-8
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			fn.Type.Return[0].Reg = x64.REG_RAX
			fn.Type.Return[0].Off = -8 // rbp-8
		case token.F32, token.F32_zh:
			fn.Type.Return[0].Reg = x64.REG_XMM0
			fn.Type.Return[0].Off = -8 // rbp-8
		case token.F64, token.F64_zh:
			fn.Type.Return[0].Reg = x64.REG_XMM0
			fn.Type.Return[0].Off = -8 // rbp-8
		default:
			panic("unreachable")
		}
	default:
		// 栈返回值需要跳过输入参数和栈帧头
		for i := 0; i < len(fn.Type.Return); i++ {
			fn.Type.Return[i].Off = headSize + argsStackSize + i*8
			fn.Type.Return[i].Cap = 1
		}
	}

	// 输入参数全部在调用方分配空间
	for i, arg := range fn.Type.Args {
		switch arg.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			if iArgRegIdx < len(iArgRegList) {
				arg.Reg = iArgRegList[iArgRegIdx]
				iArgRegIdx++
			}
			arg.Off = headSize + i*8 // rbp
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if iArgRegIdx < len(iArgRegList) {
				arg.Reg = iArgRegList[iArgRegIdx]
				iArgRegIdx++
			}
			arg.Off = headSize + i*8 // rbp
		case token.F32, token.F32_zh:
			if fArgRegIdx < len(fArgRegList) {
				arg.Reg = fArgRegList[fArgRegIdx]
				fArgRegIdx++
			}
			arg.Off = headSize + i*8 // rbp
		case token.F64, token.F64_zh:
			if fArgRegIdx < len(fArgRegList) {
				arg.Reg = fArgRegList[fArgRegIdx]
				fArgRegIdx++
			}
			arg.Off = headSize + i*8 // rbp
		default:
			panic("unreachable")
		}
	}

	// 计算局部变量的空间
	for i := 0; i < len(fn.Body.Locals); i++ {
		switch x := fn.Body.Locals[i]; x.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			sp = sp - 4*x.Cap
			x.Off = sp // rbp
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8*x.Cap
			x.Off = sp // rbp
		case token.F32, token.F32_zh:
			sp = sp - 4*x.Cap
			x.Off = sp // rbp
		case token.F64, token.F64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8*x.Cap
			x.Off = sp // rbp
		default:
			panic("unreachable")
		}
	}

	// rsp 需要 16 字节对齐
	if x := 0 - sp; x%16 != 0 {
		x = ((x + 15) / 16) * 16
		sp = 0 - x
	}

	// 当前栈帧大小(不包含调用其他函数的参数)
	fn.FrameSize = 0 - sp
	return nil
}
