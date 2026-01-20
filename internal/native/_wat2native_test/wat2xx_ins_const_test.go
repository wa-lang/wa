// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI32Const(t *testing.T) {
	tests := []struct {
		name string
		val  int32
	}{
		{"Max", W2_I32_MAX},
		{"Min", W2_I32_MIN},
		{"Hex", W2_I32_HEX},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := I32Const(tt.val)
			if got != tt.val {
				t.Errorf("I32Const(%v) = %v; 想要 %v", tt.name, got, tt.val)
			}
		})
	}
}

func TestI64Const(t *testing.T) {
	tests := []struct {
		name string
		val  int64
	}{
		{"Max", W2_I64_MAX},
		{"Min", W2_I64_MIN},
		{"Hex", W2_I64_HEX},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := I64Const(tt.val)
			if got != tt.val {
				t.Errorf("I64Const(%v) = %v; expect %v", tt.name, got, tt.val)
			}
		})
	}
}

func TestF32Const(t *testing.T) {
	tests := []struct {
		name string
		bits uint32
	}{
		{"Pi", W2_F32_PI_BITS},
		{"Infinity", W2_F32_INF_BITS},
		{"NegativeZero", W2_F32_NEG0_BITS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := math.Float32frombits(tt.bits)
			got := F32Const(input)

			gotBits := math.Float32bits(got)
			if gotBits != tt.bits {
				t.Errorf("F32Const(%v) ERR: got 0x%x; expect 0x%x", tt.name, gotBits, tt.bits)
			}
		})
	}
}

func TestF64Const(t *testing.T) {
	tests := []struct {
		name string
		bits uint64
	}{
		{"Pi", W2_F64_PI_BITS},
		{"NaN", W2_F64_NAN_BITS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := math.Float64frombits(tt.bits)
			got := F64Const(input)

			gotBits := math.Float64bits(got)
			if gotBits != tt.bits {
				if math.IsNaN(input) && math.IsNaN(got) {
					t.Logf("F64Const(%v): bits mismatch, but remains a valid NaN (got 0x%x)", tt.name, gotBits)
				} else {
					t.Errorf("F64Const(%v) ERR: got 0x%x; expect 0x%x", tt.name, gotBits, tt.bits)
				}
			}
		})
	}
}
