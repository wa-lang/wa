// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package astutil

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

// 构建栈帧中参数和返回值的位置
func buildFuncFrame_riscv(cpu abi.CPUType, fn *ast.Func) error {
	if err := buildFuncArgReturn_riscv(cpu, fn); err != nil {
		return err
	}
	if err := buildFuncLocals_riscv(cpu, fn); err != nil {
		return err
	}
	return nil
}

func buildFuncArgReturn_riscv(cpu abi.CPUType, fn *ast.Func) error {
	var headSize = 4 * 2 // RA + FP
	if cpu == abi.RISCV64 {
		headSize = 8 * 2
	}

	var (
		iArgReg abi.RegType = riscv.REG_A0
		fArgReg abi.RegType = riscv.REG_FA0
		iRetReg abi.RegType = riscv.REG_A0
		fRetReg abi.RegType = riscv.REG_FA0

		sp int = 0 - headSize // 栈顶位置
	)

	// 返回值
	for i := len(fn.Type.Return) - 1; i >= 0; i-- {
		switch ret := fn.Type.Return[i]; ret.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			if iRetReg <= riscv.REG_A1 {
				ret.Reg = iRetReg
				iRetReg++
			}
			sp = sp - 4
			ret.Off = sp

		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if iRetReg <= riscv.REG_A1 {
				ret.Reg = iRetReg
				iRetReg++
			}
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8
			ret.Off = sp

		case token.F32, token.F32_zh:
			if fRetReg <= riscv.REG_FA1 {
				ret.Reg = fRetReg
				fRetReg++
			}
			sp = sp - 4
			ret.Off = sp

		case token.F64, token.F64_zh:
			if fRetReg <= riscv.REG_FA1 {
				ret.Reg = fRetReg
				fRetReg++
			}
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8
			ret.Off = sp

		default:
			panic("unreachable")
		}
	}

	// 输入参数
	for _, arg := range fn.Type.Args {
		switch arg.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			if iArgReg <= riscv.REG_A7 {
				arg.Reg = iArgReg
				iArgReg++
			}
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if iArgReg <= riscv.REG_A7 {
				arg.Reg = iArgReg
				iArgReg++
			}
		case token.F32, token.F32_zh:
			if fArgReg <= riscv.REG_FA7 {
				arg.Reg = fArgReg
				fArgReg++
			}
		case token.F64, token.F64_zh:
			if fArgReg <= riscv.REG_FA7 {
				arg.Reg = fArgReg
				fArgReg++
			}
		default:
			panic("unreachable")
		}
	}

	// 修复走寄存器的输入参数位置
	for i := len(fn.Type.Args) - 1; i >= 0; i-- {
		arg := fn.Type.Args[i]

		// 跳过走栈的输入参数(已经在调用方处理)
		if arg.Reg == 0 {
			continue
		}

		switch arg.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			sp = sp - 4
			arg.Off = sp

		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8
			arg.Off = sp

		case token.F32, token.F32_zh:
			sp = sp - 4
			arg.Off = sp

		case token.F64, token.F64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8
			arg.Off = sp

		default:
			panic("unreachable")
		}
	}

	fn.FrameSize = 0 - sp
	return nil
}

// 构造局部遍历在栈帧的位置
func buildFuncLocals_riscv(cpu abi.CPUType, fn *ast.Func) error {
	var headSize = 4 * 2 // RA + FP
	if cpu == abi.RISCV64 {
		headSize = 8 * 2
	}

	var sp int = 0 - headSize // 栈顶位置

	if len(fn.Body.Locals) == 0 {
		return nil
	}

	// 对齐到输入参数和返回值位置的底部
	if len(fn.Type.Return) > 0 {
		if off := fn.Type.Return[0].Off; off < sp {
			sp = off
		}
	}
	if len(fn.Type.Args) > 0 {
		if off := fn.Type.Args[0].Off; off < sp {
			sp = off
		}
	}

	// 局部变量
	for i := len(fn.Body.Locals) - 1; i >= 0; i-- {
		switch x := fn.Body.Locals[i]; x.Type {
		case token.I32, token.U32, token.I32_zh, token.U32_zh:
			sp = sp - 4*x.Cap
			x.Off = sp
		case token.I64, token.U64, token.I64_zh, token.U64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8*x.Cap
			x.Off = sp
		case token.F32, token.F32_zh:
			sp = sp - 4*x.Cap
			x.Off = sp
		case token.F64, token.F64_zh:
			if sp%8 != 0 {
				sp -= 4
			}
			sp = sp - 8*x.Cap
			x.Off = sp
		default:
			panic("unreachable")
		}
	}

	fn.FrameSize = 0 - sp
	return nil
}
