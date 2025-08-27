// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package device

import (
	"fmt"
)

// 设备接口
// 设备无法提前预知被挂载的地址, 这里读写的地址是设备自身的地址
type Device interface {
	Name() string // 设备的名字

	AddrBegin() uint64 // 开始地址
	AddrEnd() uint64   // 介绍地址(开区间)

	Read(addr, size uint64) (uint64, error)
	Write(addr, size, value uint64) error
}

// 总线
// 对应全部的外设抽象
// 地址不得重叠, 名字也不得相同
type Bus struct {
	Devices map[string]Device // 设备列表
}

// 构造总线
func NewBus() *Bus {
	return &Bus{
		Devices: make(map[string]Device),
	}
}

// 内存映射设备
func (bus *Bus) MapDevice(d Device) {
	if d.Name() == "" {
		panic("invalid device")
	}

	// 设备的名字必须唯一
	if _, ok := bus.Devices[d.Name()]; ok {
		panic(fmt.Sprintf("%s: device exists", d.Name()))
	}

	// 判断地址是否重叠
	start := d.AddrBegin()
	end := d.AddrEnd()
	for _, x := range bus.Devices {
		if end < x.AddrBegin() || start >= x.AddrEnd() {
			continue // 没有交集
		}
		panic(fmt.Sprintf("%s: device %s use the same address space([%08X,%08X][%08X,%08X])",
			d.Name(), x.Name(),
			d.AddrBegin(), d.AddrEnd(),
			x.AddrBegin(), x.AddrEnd(),
		))
	}

	// OK
	bus.Devices[d.Name()] = d
}

// 从总线上指定地址读数据
func (bus *Bus) Read(addr, size uint64) (uint64, error) {
	for _, x := range bus.Devices {
		if x.AddrBegin() <= addr && addr < x.AddrEnd() {
			return x.Read(addr, size)
		}
	}
	panic(fmt.Sprintf("bus: no devide at [0x%08X]", addr))
}

// 向总线上指定地址写数据
func (bus *Bus) Write(addr, size, value uint64) error {
	for _, x := range bus.Devices {
		if x.AddrBegin() <= addr && addr < x.AddrEnd() {
			return x.Write(addr, size, value)
		}
	}
	panic(fmt.Sprintf("bus: no devide at [0x%08X]", addr))
}
