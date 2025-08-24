// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "testing"

func TestAOpContextTable(t *testing.T) {
	for i := AADDI; i < ALAST; i++ {
		ctx := AOpContextTable[i]
		if ctx.Pseudo {
			continue
		}
		if ctx.Opcode == 0 {
			t.Fatalf("invalid AOpContextTable[%s]", AsString(i))
		}
	}
}
