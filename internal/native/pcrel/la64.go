// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pcrel

// GetTargetAddressLa64 根据 LoongArch PC 和指令参数计算出原始的目标地址
func GetTargetAddressLa64(pc int64, pc_hi20, pc_lo12 int32) int64 {
	aluOut := (pc + (int64(pc_hi20) << 12)) &^ 0xFFF
	targetAddress := aluOut + int64(pc_lo12)
	return targetAddress
}

// MakeLa64PCRel 计算龙芯64目标地址相对于 PC 地址的相对偏移量
func MakeLa64PCRel(targetAddress, pc int64) (pc_hi20, pc_lo12 int32) {
	pcPage := pc &^ 0xFFF
	delta := targetAddress - pcPage

	pc_hi20 = int32((delta) >> 12)
	pc_lo12 = int32(delta & 0xFFF)

	if pc_lo12 >= 0x800 {
		pc_hi20++
	}
	pc_hi20 &= 0xFFFFF
	return
}
