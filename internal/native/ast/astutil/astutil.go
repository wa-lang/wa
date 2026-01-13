// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package astutil

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
)

// 构建函数的栈帧
// fn 需要包含 输入参数/返回值/局部变量 信息
func BuildFuncFrame(cpu abi.CPUType, fn *ast.Func) error {
	switch cpu {
	case abi.LOONG64:
		return buildFuncFrame_loong64(cpu, fn)
	case abi.RISCV32:
		return buildFuncFrame_riscv(cpu, fn)
	case abi.RISCV64:
		return buildFuncFrame_riscv(cpu, fn)
	case abi.X64:
		return buildFuncFrame_x64_windows(fn)
	case abi.X64Unix:
		return buildFuncFrame_x64_unix(cpu, fn)
	default:
		panic("unreachable")
	}
}
