// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv64

import "fmt"

// 总线
// 对应全部的外设抽象
// 地址不得重叠, 名字也不得相同
type Bus struct {
	Devices map[uint64]Device // 设备列表
}

// 设备接口
// 设备无法提前预知被挂载的地址, 这里读写的地址是设备自身的地址
type Device interface {
	Name() string // 设备的名字
	Size() uint64 // 设备的地址范围 [0, device.Size())
	Read(addr uint64) (uint64, error)
	Write(addr, value uint64) error
}

// 内存映射设备
func (bus *Bus) MapDevice(d Device, startAddr uint64) {
	if d.Name() == "" || d.Size() == 0 {
		panic("invalid device")
	}

	// 判断名字和地址是否重叠
	endAddr := startAddr + d.Size()
	for xStart, xDevice := range bus.Devices {
		if d.Name() == xDevice.Name() {
			panic(fmt.Sprintf("device(%s) exists", d.Name()))
		}
		if xStart < endAddr || xStart+xDevice.Size() > startAddr {
			panic(fmt.Sprintf("hte address space([%d,%d]) has been used", startAddr, endAddr-1))
		}
	}

	// OK
	bus.Devices[startAddr] = d
}

// 从总线上指定地址读数据
func (bus *Bus) Read(addr uint64) (uint64, error) {
	for xStart, xDevice := range bus.Devices {
		if xStart <= addr && addr < xStart+xDevice.Size() {
			return xDevice.Read(addr - xStart)
		}
	}
	panic(fmt.Sprintf("bus: no devide for %0x08X", addr))
}

// 向总线上指定地址写数据
func (bus *Bus) Write(addr, value uint64) error {
	for xStart, xDevice := range bus.Devices {
		if xStart <= addr && addr < xStart+xDevice.Size() {
			return xDevice.Write(addr-xStart, value)
		}
	}
	panic(fmt.Sprintf("bus: no devide for %0x08X", addr))
}
