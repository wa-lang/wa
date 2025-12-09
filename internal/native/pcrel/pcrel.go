// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pcrel

// GetTargetAddress 根据 pc 和参数计算出原始的地址
func GetTargetAddress(pc uint32, pcrel_hi, pcrel_lo int32) (targetAddress uint32) {
	delta := CombineOffset(pcrel_hi, pcrel_lo)
	targetAddress = uint32(int64(pc) + int64(delta))
	return
}

// CombineOffset 合并出原始的偏移量
func CombineOffset(pcrel_hi, pcrel_lo int32) (delta int32) {
	delta = (pcrel_hi << 12) + pcrel_lo
	return delta
}

// MakePCRel 计算目标地址相对于 PC 地址的相对偏移量
func MakePCRel(targetAddress, pc uint32) (pcrel_hi, pcrel_lo int32) {
	delta := int32(int64(targetAddress) - int64(pc))
	return SplitOffset(delta)
}

// SplitOffset 将一个 32 位的 PC 相对偏移量拆分为 20 位的校正高位和 12 位的带符号低位.
//
// 这个逻辑适用于 LoongArch/RISC-V 等架构中 PC 相对寻址所需的地址校正.
// 在通过 PC 相对位置计算地址时, 范围是 PC 附近的 2GB 地址空间
func SplitOffset(delta int32) (hiCorr int32, loSigned int32) {
	// 0. 检查 LO 部分的符号位 (第 11 位, 即 0x800)
	const lowBitMask = 0x800

	// 1. 计算校正后的 HI 值 (HI_corr)
	// 如果第 11 位为 1, 说明 LO 部分是负数, HI 需要向上取整 (+1)
	if (delta & lowBitMask) != 0 {
		// 如果 LO 部分是负数, HI = (Delta >> 12) + 1
		// 这样是为了抵消 LO 部分的负值
		hiCorr = (delta >> 12) + 1
	} else {
		// 如果 LO 部分是正数或零, HI = Delta >> 12
		hiCorr = delta >> 12
	}

	// 2. 计算 LO 值 (LO_signed)
	// LO_signed = Delta - (HI_corr << 12)
	// 这确保了 LO_signed 必然落在了 [-2048, 2047] 的 12 位带符号立即数范围内
	loSigned = delta - (hiCorr << 12)

	return hiCorr, loSigned
}
