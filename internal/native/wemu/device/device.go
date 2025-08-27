// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package device

// 设备接口
// 设备无法提前预知被挂载的地址, 这里读写的地址是设备自身的地址
type Device interface {
	Name() string // 设备的名字
	Size() uint64 // 设备的地址范围 [0, device.Size())
	Read(addr uint64) (uint64, error)
	Write(addr, value uint64) error
}
