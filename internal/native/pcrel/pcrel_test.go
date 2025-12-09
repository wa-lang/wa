// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pcrel

import (
	"math"
	"testing"
)

// TestData 结构用于定义测试用例
type TestData struct {
	name          string
	delta         int32  // 原始偏移量 (targetAddress - pc)
	expectedHi    int32  // 期望的 %pcrel_hi
	expectedLo    int32  // 期望的 %pcrel_lo (必须在 [-2048, 2047] 范围内)
	targetAddress uint32 // 用于 MakePCRel 测试 (A)
	pc            uint32 // 用于 MakePCRel 测试 (PC)
}

// 完整的测试用例列表 (已修复 uint32 溢出问题)
var testCases = []TestData{
	// --- I. 零点和普通偏移量 ---
	{"1_Zero Delta", 0, 0, 0, 0x100000, 0x100000},

	// Delta = 8483. targetAddress 设大以避免 pc 溢出。
	{"2_Simple Positive", 8483, 2, 291, 0x100000, 0x100000 - 8483},

	// Delta = -8483. pc = targetAddress + 8483 (正数，无需修改)
	{"3_Simple Negative", -8483, -2, -291, 0x2000, 0x2000 + 8483},

	// --- II. Lo 边界测试 (±2048 附近) ---

	// Delta = 6143. targetAddress 设大以避免 pc 溢出。
	{"4_Boundary_LoMaxPos", 6143, 1, 2047, 0x100000, 0x100000 - 6143},

	// Delta = 6144. targetAddress 设大以避免 pc 溢出。
	{"5_Boundary_LoMaxNegAbs", 6144, 2, -2048, 0x100000, 0x100000 - 6144},

	// Delta = 4095. targetAddress 设大以避免 pc 溢出。
	{"6_Boundary_LoNeg1", 4095, 1, -1, 0x100000, 0x100000 - 4095},

	// Delta = 4096. targetAddress 设大以避免 pc 溢出。
	{"7_Boundary_LoIsZero", 4096, 1, 0, 0x100000, 0x100000 - 4096},

	// Delta = -2049. pc = targetAddress + 2049 (正数，无需修改)
	{"8_Boundary_LoMinNegAbs_minus1", -2049, -1, 2047, 0x1000, 0x1000 + 2049},

	// Delta = -2048. pc = targetAddress + 2048 (正数，无需修改)
	{"9_Boundary_LoMaxNeg", -2048, 0, -2048, 0x1000, 0x1000 + 2048},

	// --- III. Hi 极端值测试 (32位 Int32 边界) ---

	// Delta = MaxInt32. pc = targetAddress - MaxInt32.
	// targetAddress 必须至少为 MaxInt32 (0x7FFFFFFF) 才能保持 pc 为正。
	// 使用 0x80000000 (2^31) 作为 targetAddress，pc = 1。
	{"10_MaxInt32", math.MaxInt32, 0x80000, -1, 0x80000000, 0x80000000 - math.MaxInt32},

	// Delta = MinInt32. pc = targetAddress + 2^31.
	// targetAddress 必须使用 uint32(MaxInt32 + 1) 以上的值才能确保 pc 为正。
	// 使用 0x100000， pc = 0x100000 + 2^31，仍在 uint32 范围内。
	{"11_MinInt32", math.MinInt32, -0x80000, 0, 0x100000, 0x100000 - math.MinInt32},
}

func TestSplitOffset(t *testing.T) {
	for _, tc := range testCases {
		hi, lo := SplitOffset(tc.delta)

		// 1. 检查拆分结果是否正确
		if hi != tc.expectedHi || lo != tc.expectedLo {
			t.Errorf("SplitOffset(%d) (0x%X) failed.\nGot hi=%d (0x%X), lo=%d (0x%X);\nWant hi=%d (0x%X), lo=%d (0x%X)",
				tc.delta, tc.delta, hi, hi, lo, lo, tc.expectedHi, tc.expectedHi, tc.expectedLo, tc.expectedLo)
		}

		// 2. 检查结果的重构是否等于原始 Delta
		reconstructed := (hi << 12) + lo
		if reconstructed != tc.delta {
			t.Errorf("Reconstruction failed: (%d << 12) + %d = %d; want %d",
				hi, lo, reconstructed, tc.delta)
		}

		// 3. 检查 Lo 是否在 [-2048, 2047] 范围内 (12位带符号立即数)
		if lo < -2048 || lo > 2047 {
			t.Errorf("Lo value (%d) is outside the valid range [-2048, 2047]", lo)
		}

	}
}

func TestMakePCRel(t *testing.T) {
	for _, tc := range testCases {
		// 验证 Delta 是否匹配测试用例的 Delta
		// MakePCRel 使用 int64 减法来计算 Delta，这是安全的。
		expectedDelta := int32(int64(tc.targetAddress) - int64(tc.pc))

		if expectedDelta != tc.delta {
			t.Fatalf("Test Case Error: target/pc addresses resulted in Delta=%d, but expected Delta is %d", expectedDelta, tc.delta)
		}

		// 调用 MakePCRel
		hi, lo := MakePCRel(tc.targetAddress, tc.pc)

		// 检查 SplitOffset 的结果
		if hi != tc.expectedHi || lo != tc.expectedLo {
			t.Errorf("MakePCRel(0x%X, 0x%X) [Delta %d] failed.\nGot hi=%d, lo=%d;\nWant hi=%d, lo=%d",
				tc.targetAddress, tc.pc, tc.delta, hi, lo, tc.expectedHi, tc.expectedLo)
		}
	}
}
