// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"testing"
)

func TestSplit32BitImmediate(t *testing.T) {
	tests := []struct {
		imm       int64
		shouldErr bool
	}{
		// 12bit 以内
		{0, false},
		{2047, false},
		{-2048, false},

		// 超过 12bit
		{0x12345, false},
		{-0x12345, false},
	}

	for _, tt := range tests {
		low, high, err := split32BitImmediate(tt.imm)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("imm=%d: expected error, got none", tt.imm)
			}
			continue
		}
		if err != nil {
			t.Errorf("imm=%d: unexpected error: %v", tt.imm, err)
			continue
		}
		reconstructed := (high << 12) + low
		if reconstructed != tt.imm {
			t.Errorf("imm=%d: reconstructed=%d (low=%d, high=%d)", tt.imm, reconstructed, low, high)
		}
		if low < -2048 || low > 2047 {
			t.Errorf("imm=%d: low=%d out of 12-bit signed range", tt.imm, low)
		}
	}
}
