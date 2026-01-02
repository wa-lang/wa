// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package zparser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
)

func (p *parser) parseInst(fn *ast.Func) (inst *ast.Instruction) {
	switch p.cpu {
	case abi.LOONG64:
		return p.parseInst_loong(fn)
	case abi.RISCV32, abi.RISCV64:
		return p.parseInst_riscv(fn)
	default:
		panic("unreachable")
	}
}
