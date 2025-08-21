// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "testing"

func TestAOpContextTable(t *testing.T) {
	for i := AXXX; i < A_ARCHSPECIFIC; i++ {
		if AOpContextTable[i].Opcode != 0 {
			t.Fatalf("invalid AOpContextTable[%v]", i)
		}
	}
	for i := A_ARCHSPECIFIC; i < ALAST; i++ {
		if AOpContextTable[i].Opcode == 0 {
			t.Fatalf("invalid AOpContextTable[%s]", AsString(i))
		}
	}
}
