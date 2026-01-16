// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package astutil

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// ABI: $(WA_REPO)/docs/asm_abi_x64.md

// 构建栈帧中参数和返回值的位置
func buildFuncFrame_x64_unix(fn *ast.Func) error {
	const headSize = 8 * 2 // RIP + RBP

	var (
		argRegAlloctor = newArgsRegAlloctor_X64Linux()
		retRegAlloctor = newRetRegAlloctor_X64Linux()

		retOnStack bool // 返回值在栈上

		argsSize  int // 调用方输入参数大小
		frameSize int // 被调用方局部栈大小
	)

	// 尝试为返回值分配寄存器
	// 返回值是否走栈要看寄存器分配结果
	{
		for _, reti := range fn.Type.Return {
			switch reti.Type {
			case token.I32, token.U32, token.I32_zh, token.U32_zh:
				if reti.Reg = retRegAlloctor.GetInt(); reti.Reg == 0 {
					retOnStack = true
				}
			case token.I64, token.U64, token.I64_zh, token.U64_zh:
				if reti.Reg = retRegAlloctor.GetInt(); reti.Reg == 0 {
					retOnStack = true
				}
			case token.F32, token.F32_zh:
				if reti.Reg = retRegAlloctor.GetFloat(); reti.Reg == 0 {
					retOnStack = true
				}
			case token.F64, token.F64_zh:
				if reti.Reg = retRegAlloctor.GetFloat(); reti.Reg == 0 {
					retOnStack = true
				}
			default:
				panic("unreachable")
			}
		}
		if retOnStack {
			for _, reti := range fn.Type.Return {
				reti.Reg = 0
			}
		}
	}

	// 为输入参数分配寄存器和对应的内存
	// 输入参数可能有寄存器和栈内存混合的情况
	for _, arg := range fn.Type.Args {
		switch arg.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			if r := argRegAlloctor.GetInt(); r != 0 {
				arg.Reg = r
				arg.RBPOff = 0 - frameSize - 8
				frameSize += 8
			} else {
				arg.RSPOff = argsSize
				arg.RBPOff = argsSize + headSize
				argsSize += 8
			}
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if r := argRegAlloctor.GetInt(); r != 0 {
				arg.Reg = r
				arg.RBPOff = 0 - frameSize - 8
				frameSize += 8
			} else {
				arg.RSPOff = argsSize
				arg.RBPOff = argsSize + headSize
				argsSize += 8
			}
		case token.F32, token.F32_zh:
			if r := argRegAlloctor.GetFloat(); r != 0 {
				arg.Reg = r
				arg.RBPOff = 0 - frameSize - 8
				frameSize += 8
			} else {
				arg.RSPOff = argsSize
				arg.RBPOff = argsSize + headSize
				argsSize += 8
			}
		case token.F64, token.F64_zh:
			if r := argRegAlloctor.GetFloat(); r != 0 {
				arg.Reg = r
				arg.RBPOff = 0 - frameSize - 8
				frameSize += 8
			} else {
				arg.RSPOff = argsSize
				arg.RBPOff = argsSize + headSize
				argsSize += 8
			}
		default:
			panic("unreachable")
		}
	}

	// 为返回值分配对应的内存
	// 返回值只能是全部走寄存器或者全部走栈内存
	if retOnStack {
		// 返回值全部走栈, 需要输入方分配
		for _, ret := range fn.Type.Return {
			ret.RSPOff = argsSize
			ret.RBPOff = argsSize + headSize
			argsSize += 8
		}
	} else {
		// 返回值全部走寄存器
		for _, ret := range fn.Type.Return {
			ret.RBPOff = 0 - frameSize - 8
			frameSize += 8
		}
	}

	// 局部变量
	for _, local := range fn.Body.Locals {
		local.RBPOff = 0 - frameSize - 8*local.Cap
		frameSize += 8 * local.Cap
	}

	fn.ArgsSize = argsSize
	fn.FrameSize = frameSize
	return nil
}
