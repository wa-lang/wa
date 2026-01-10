// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestWat2xxI32Add(t *testing.T) {
	tests := []struct {
		name string
		a, b int32
		want int32
	}{
		{"基础加法", 1, 2, 3},
		{"负数测试", -10, 5, -5},
		{"零值测试", 0, 100, 100},
		{"最大值溢出", math.MaxInt32, 1, math.MinInt32}, // Wasm 要求 i32 自动截断/溢出
		{"最小值相加", math.MinInt32, math.MinInt32, 0},
		{"大数碰撞", 0x7FFFFFFF, 0x7FFFFFFF, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := I32Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("%s 失败: %d + %d = 预期 %d, 实际 %d",
					tt.name, tt.a, tt.b, tt.want, got)
			}
		})
	}
}
