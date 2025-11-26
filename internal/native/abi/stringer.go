// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

import "fmt"

func (cpu CPUType) String() string {
	switch cpu {
	case RISCV64:
		return "RISCV64"
	case RISCV32:
		return "RISCV32"
	case LOONG64:
		return "LOONG64"
	}
	return fmt.Sprintf("abi.CPUType(%d)", int(cpu))
}

func (fn BuiltinFn) String() string {
	switch fn {
	case BuiltinFn_HI:
		return "%hi"
	case BuiltinFn_LO:
		return "%lo"
	case BuiltinFn_PCREL_HI:
		return "%pcrel_hi"
	case BuiltinFn_PCREL_LO:
		return "%pcrel_lo"

	case BuiltinFn_HI_zh:
		return "%高位"
	case BuiltinFn_LO_zh:
		return "%低位"
	case BuiltinFn_PCREL_HI_zh:
		return "%相对高位"
	case BuiltinFn_PCREL_LO_zh:
		return "%相对低位"
	}
	return fmt.Sprintf("abi.BuiltinFn(%d)", int(fn))
}
