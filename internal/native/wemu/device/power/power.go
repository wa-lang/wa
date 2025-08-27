// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package power

import (
	"encoding/binary"
	"fmt"
)

// 退出设备的寄存器地址
const POWER_BASE = 0x100000

// 退出的状态
type Status uint32

const (
	ExitOK   = Status(0x5555) // 正常退出
	ExitFail = Status(0x3333) // 异常退出
)

// 电源设备(4个字节)
type Power struct {
	name string
	data []byte
}

func NewPower(name string) *Power {
	return &Power{name, make([]byte, 4)}
}

func (p *Power) Name() string { return p.name }
func (p *Power) Size() uint64 { return uint64(len(p.data)) }

// 电源状态
func (p *Power) Status() Status {
	return Status(binary.LittleEndian.Uint32(p.data))
}

func (p *Power) Read(addr, size uint64) (uint64, error) {
	if addr+size >= p.Size() {
		return 0, fmt.Errorf("%s: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch size {
	case 1:
		return uint64(p.data[addr]), nil
	case 2:
		return uint64(binary.LittleEndian.Uint16(p.data[addr:])), nil
	case 4:
		return uint64(binary.LittleEndian.Uint32(p.data[addr:])), nil
	case 8:
		return binary.LittleEndian.Uint64(p.data[addr:]), nil
	default:
		panic(fmt.Sprintf("rom: size %d must 1/2/4/8", size))
	}
}

func (p *Power) Write(addr, size, value uint64) error {
	if addr+size >= p.Size() {
		return fmt.Errorf("%s: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch size {
	case 1:
		p.data[addr] = uint8(value)
		return nil
	case 2:
		binary.LittleEndian.PutUint16(p.data[addr:], uint16(value))
		return nil
	case 4:
		binary.LittleEndian.PutUint32(p.data[addr:], uint32(value))
		return nil
	case 8:
		binary.LittleEndian.PutUint64(p.data[addr:], uint64(value))
		return nil
	default:
		return fmt.Errorf("%s: bad size, %d must one of 1/2/4/8", p.name, size)
	}
}
