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
	}
	return fmt.Sprintf("abi.CPUType(%d)", int(cpu))
}
