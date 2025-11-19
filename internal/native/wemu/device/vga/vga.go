// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package vga

import (
	"fmt"
)

// 设备的默认映射地址
const VGA_BASE = 0x100000

// 显示器
type VGA struct {
	name    string
	addr    uint64
	width   int
	height  int
	data    []byte  // 内存数据
	palette []byte  // data[:256*4]
	screen  []uint8 // data[256*4:]
}

// 构造 VGA 显示器
func NewVGA(name string, addr uint64, width, height int) *VGA {
	p := &VGA{
		name:   name,
		addr:   addr,
		width:  width,
		height: height,
		data:   make([]byte, 256*4+width*height),
	}

	p.palette = p.data[:256*4]
	p.screen = p.data[256*4:]

	for i := 0; i < 256; i++ {
		p.palette[i+0] = uint8(i)
		p.palette[i+1] = uint8(i)
		p.palette[i+2] = uint8(i)
		p.palette[i+3] = 255
	}

	return p
}

func (p *VGA) Name() string      { return p.name }
func (p *VGA) AddrBegin() uint64 { return p.addr }
func (p *VGA) AddrEnd() uint64   { return p.addr + uint64(len(p.data)) }

// 读颜色表
func (p *VGA) ReadPalette(offset uint64) uint32 {
	v := uint32(p.palette[offset+2])       // B
	v |= uint32(p.palette[offset+1]) << 8  // G
	v |= uint32(p.palette[offset+0]) << 16 // R
	v |= uint32(p.palette[offset+3]) << 24
	return v
}

// 写颜色表
func (p *VGA) WritePalette(offset uint64, value uint32) {
	p.palette[offset+2] = uint8((value >> 0) & 0xFF)  // B
	p.palette[offset+1] = uint8((value >> 8) & 0xFF)  // G
	p.palette[offset+0] = uint8((value >> 16) & 0xFF) // R
	p.palette[offset+3] = 255
}

// 读帧缓存
func (p *VGA) ReadScreen(offset uint64) (value uint32) {
	v := uint32(p.screen[offset+2])       // B
	v |= uint32(p.screen[offset+1]) << 8  // G
	v |= uint32(p.screen[offset+0]) << 16 // R
	v |= uint32(p.screen[offset+3]) << 24
	return v
}

// 写帧缓存(offset不是内存地址, 而是帧缓存内部偏移)
func (vga *VGA) WriteScreen(offset uint64, value uint32) {
	vga.screen[offset+0] = uint8(value & 0xFF)
	vga.screen[offset+1] = uint8((value & 0xFF00) >> 8)
	vga.screen[offset+2] = uint8((value & 0xFF0000) >> 16)
	vga.screen[offset+3] = uint8((value & 0xFF000000) >> 24)
}

func (p *VGA) Read(addr, size uint64) (uint64, error) {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return 0, fmt.Errorf("%s.Read: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch off := addr - p.addr; true {
	case off < 256*4:
		return uint64(p.ReadPalette(off)), nil
	default:
		return uint64(p.ReadScreen(off)), nil
	}
}

func (p *VGA) Write(addr, size, value uint64) error {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return fmt.Errorf("%s.Read: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}

	switch off := addr - p.addr; true {
	case off < 256*4:
		p.WritePalette(off, uint32(value))
		return nil
	default:
		p.WriteScreen(off, uint32(value))
		return nil
	}
}
