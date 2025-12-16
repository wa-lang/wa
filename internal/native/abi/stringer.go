// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

import (
	"fmt"
	"strings"
)

func ParseCPUType(name string) CPUType {
	for i, s := range _CPUType_strings {
		if s != "" && strings.EqualFold(s, name) {
			return CPUType(i)
		}
	}
	return CPU_Nil
}

func (cpu CPUType) String() string {
	if cpu >= 0 && int(cpu) < len(_CPUType_strings) {
		return _CPUType_strings[cpu]
	}
	return fmt.Sprintf("abi.CPUType(%d)", int(cpu))
}

func ParseBuiltinFn(cpu CPUType, name string) BuiltinFn {
	for i, s := range _BuiltinFn_strings {
		if s != "" && strings.EqualFold(s, name) {
			if fn := BuiltinFn(i); fn.IsValid(cpu) {
				return fn
			}
		}
	}
	return BuiltinFn_Nil
}

func (fn BuiltinFn) String() string {
	if fn >= 0 && int(fn) < len(_BuiltinFn_strings) {
		return _BuiltinFn_strings[fn]
	}
	return fmt.Sprintf("abi.BuiltinFn(%d)", int(fn))
}

// 不同平台的内置函数不同
func (fn BuiltinFn) IsValid(cpu CPUType) bool {
	if fn <= BuiltinFn_Nil || fn >= BuiltinFn_Max {
		return false
	}

	// sizeof 通用
	if fn == BuiltinFn_SIZEOF || fn == BuiltinFn_SIZEOF_zh {
		return true
	}

	// 少量是 RISCV 特有的
	onlyRiscvEnabled := false
	switch fn {
	case BuiltinFn_HI,
		BuiltinFn_LO,
		BuiltinFn_PCREL_HI,
		BuiltinFn_PCREL_LO:
		onlyRiscvEnabled = false
	case BuiltinFn_HI_zh,
		BuiltinFn_LO_zh,
		BuiltinFn_PCREL_HI_zh,
		BuiltinFn_PCREL_LO_zh:
		onlyRiscvEnabled = false
	}

	switch cpu {
	case LOONG64:
		return !onlyRiscvEnabled
	case RISCV64, RISCV32:
		return onlyRiscvEnabled
	default:
		return false
	}
}
