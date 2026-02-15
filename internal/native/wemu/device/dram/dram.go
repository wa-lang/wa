// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package dram

import (
	"encoding/binary"
	"fmt"

	"wa-lang.org/wa/internal/native/link/pe"
)

// 默认的内存大小
const DRAM_SIZE = 1024 * 1024 * 16

// 龙芯的默认启动地址(4GB+512MB)
const DRAM_BASE_LA64 = 0x120000000

// Same as QEMU virt machine, DRAM starts at 0x80000000.
const DRAM_BASE_RISCV = 0x80000000

// X64/Linux 默认地址
const DRAM_BASE_X64_LINUX = 0x400000

// X64/Linux 默认地址
const DRAM_BASE_X64_WINDOWS = pe.DefaultImageBase

// 内存设备
type DRAM struct {
	name     string
	addr     uint64
	data     []byte
	readonly bool
}

// 长度必须是4的倍数
func NewDRAM(name string, addr, size uint64, readonly bool) *DRAM {
	if size == 0 {
		panic("rom size wa zero")
	}
	if size%4 != 0 {
		panic("rom size must align with 4")
	}
	return &DRAM{name, addr, make([]byte, size), readonly}
}

func (p *DRAM) Name() string { return p.name }

func (p *DRAM) AddrBegin() uint64 { return p.addr }
func (p *DRAM) AddrEnd() uint64   { return p.addr + uint64(len(p.data)) }

// 填充内存数据
func (p *DRAM) Fill(addr uint64, data []byte) error {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return fmt.Errorf("%s.Fill: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+uint64(len(data)))
	}
	copy(p.data[addr-p.AddrBegin():], data)
	return nil
}

func (p *DRAM) Read(addr, size uint64) (uint64, error) {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return 0, fmt.Errorf("%s.Read: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch size {
	case 1:
		return uint64(p.data[addr-p.addr]), nil
	case 2:
		return uint64(binary.LittleEndian.Uint16(p.data[addr-p.addr:])), nil
	case 4:
		return uint64(binary.LittleEndian.Uint32(p.data[addr-p.addr:])), nil
	case 8:
		return binary.LittleEndian.Uint64(p.data[addr-p.addr:]), nil
	default:
		return 0, fmt.Errorf("%s.Read: bad size, %d must one of 1/2/4/8", p.name, size)
	}
}

func (p *DRAM) Write(addr, size, value uint64) error {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return fmt.Errorf("%s.Write: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	if p.readonly {
		return fmt.Errorf("%s.Write: is readonly", p.name)
	}
	switch size {
	case 1:
		p.data[addr-p.addr] = uint8(value)
		return nil
	case 2:
		binary.LittleEndian.PutUint16(p.data[addr-p.addr:], uint16(value))
		return nil
	case 4:
		binary.LittleEndian.PutUint32(p.data[addr-p.addr:], uint32(value))
		return nil
	case 8:
		binary.LittleEndian.PutUint64(p.data[addr-p.addr:], uint64(value))
		return nil
	default:
		return fmt.Errorf("%s.Write: bad size, %d must one of 1/2/4/8", p.name, size)
	}
}
