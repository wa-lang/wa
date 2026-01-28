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
	const MaxRegArgs = 4   // 最多4个寄存器参数
	const headSize = 8 * 2 // RIP + RBP

	var (
		// 候选的寄存器列表
		iArgRegList = []abi.RegType{x64.REG_RCX, x64.REG_RDX, x64.REG_R8, x64.REG_R9}
		fArgRegList = []abi.RegType{x64.REG_XMM0, x64.REG_XMM1, x64.REG_XMM2, x64.REG_XMM3}

		iArgRegIdx int
		fArgRegIdx int
		nArgReg    int // 寄存器参数数量

		arg0Base int = 0 // 第一个参数的偏移地址

		sp int = 0 // 局部栈大小
	)

	// 计算输入参数需要的栈大小
	fn.ArgsSize = len(fn.Type.Args) * 8

	// 影子空间是强制保留的
	if fn.ArgsSize < 32 {
		fn.ArgsSize = 32
	}

	// 返回值
	switch len(fn.Type.Return) {
	case 0:
	case 1:
		// 给单个返回值的影子空间
		// 由函数内部保留
		sp -= 8

		// 是基于 rbp 的位置
		switch ret := fn.Type.Return[0]; ret.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			fn.Type.Return[0].Reg = x64.REG_RAX
			fn.Type.Return[0].RBPOff = -8 // rbp-8
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			fn.Type.Return[0].Reg = x64.REG_RAX
			fn.Type.Return[0].RBPOff = -8 // rbp-8
		case token.F32, token.F32_zh:
			fn.Type.Return[0].Reg = x64.REG_XMM0
			fn.Type.Return[0].RBPOff = -8 // rbp-8
		case token.F64, token.F64_zh:
			fn.Type.Return[0].Reg = x64.REG_XMM0
			fn.Type.Return[0].RBPOff = -8 // rbp-8
		default:
			panic("unreachable")
		}
	default:
		// 栈返回值需要跳过输入参数和栈帧头
		base := headSize + 8 + len(fn.Type.Args)*8
		fn.ArgsSize += len(fn.Type.Return)*8 + 8
		for i := 0; i < len(fn.Type.Return); i++ {
			fn.Type.Return[i].RBPOff = base + i*8
			fn.Type.Return[i].Cap = 1
		}

		// 跳过第一个参数寄存器
		// Windows 下参数位置和寄存器是对应的
		arg0Base = 8
		iArgRegIdx++
		fArgRegIdx++
		nArgReg++
	}

	// 输入参数全部在调用方分配空间
	for i, arg := range fn.Type.Args {
		switch arg.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			if nArgReg < MaxRegArgs && iArgRegIdx < len(iArgRegList) {
				arg.Reg = iArgRegList[iArgRegIdx]
				iArgRegIdx++
				nArgReg++
			}
			arg.RBPOff = headSize + arg0Base + i*8 // rbp
			arg.RSPOff = arg0Base + i*8            // 调用前的 rsp
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if nArgReg < MaxRegArgs && iArgRegIdx < len(iArgRegList) {
				arg.Reg = iArgRegList[iArgRegIdx]
				iArgRegIdx++
				nArgReg++
			}
			arg.RBPOff = headSize + arg0Base + i*8 // rbp
			arg.RSPOff = arg0Base + i*8            // 调用前的 rsp
		case token.F32, token.F32_zh:
			if nArgReg < MaxRegArgs && fArgRegIdx < len(fArgRegList) {
				arg.Reg = fArgRegList[fArgRegIdx]
				fArgRegIdx++
				nArgReg++
			}
			arg.RBPOff = headSize + arg0Base + i*8 // rbp
			arg.RSPOff = i * 8                     // 调用前的 rsp
		case token.F64, token.F64_zh:
			if nArgReg < MaxRegArgs && fArgRegIdx < len(fArgRegList) {
				arg.Reg = fArgRegList[fArgRegIdx]
				fArgRegIdx++
				nArgReg++
			}
			arg.RBPOff = headSize + arg0Base + i*8 // rbp
			arg.RSPOff = arg0Base + i*8            // 调用前的 rsp
		default:
			panic("unreachable")
		}
	}

	// 计算局部变量的空间
	// 每个标量在栈上的空间都是8个字节
	for _, x := range fn.Body.Locals {
		sp = sp - x.Cap*8
		x.RBPOff = sp // rbp
	}

	// 当前栈帧大小(不包含调用其他函数的参数)
	// 当前并未做任何对齐操作
	fn.FrameSize = 0 - sp
	return nil
}
