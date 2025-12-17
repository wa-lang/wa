// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pcrel

// MakeLa64PCRel 计算龙芯64目标地址相对于 PC 地址的相对偏移量
// TODO: hi20 结果有问题
func MakeLa64PCRel(targetAddress, pc int64) (pc_hi20, pc_lo12 int32) {
	delta := targetAddress - pc

	pc_hi20 = int32((delta + 0x800) >> 12)
	pc_lo12 = int32(targetAddress & 0xFFF)

	if pc_lo12 >= 0x800 {
		pc_lo12 -= 0x1000
	}

	return
}
