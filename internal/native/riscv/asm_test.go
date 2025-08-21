// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv_test

import (
	"fmt"

	"wa-lang.org/wa/internal/native/riscv"
)

func ExampleAsmSyntax() {
	fmt.Println(riscv.AsmSyntax(riscv.AADD, &riscv.AsArgument{
		Rd:  riscv.REG_X1,
		Rs1: riscv.REG_X2,
		Rs2: riscv.REG_X3,
	}))

	// Output:
	// ADD X1, X2, X3
}
