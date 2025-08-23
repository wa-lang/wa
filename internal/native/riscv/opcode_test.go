// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "testing"

func TestAOpContextTable(t *testing.T) {
	for i := AXXX + 1; i < ALAST; i++ {
		if AOpContextTable[i].Opcode == 0 {
			t.Fatalf("invalid AOpContextTable[%s]", AsString(i))
		}
	}
}
