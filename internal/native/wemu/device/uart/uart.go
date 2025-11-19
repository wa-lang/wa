// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 串口设备
package uart

import (
	"fmt"
	"unicode/utf8"
)

// 默认的地址
const (
	UART_BASE_QEMU     = 0x10000000
	UART_BASE_ESP32_C3 = 0x60000000

	UART_BASE_DEFAULT = UART_BASE_QEMU
	UART_SIZE         = 0x100

	// 相对于 Base 的偏移
	UART_RHR    = 0
	UART_THR    = 0
	UART_LCR    = 3
	UART_LSR    = 5
	UART_LSR_RX = 1
	UART_LSR_TX = 1 << 5

	UART_IRQ = 10
)

type UART struct {
	name string
	addr uint64
	data []byte

	rx      uint8 // 最近接收到的字符
	tx      uint8 // 最近写出的字符
	hasRX   bool  // 是否有待读数据(需要通过其他Goroutine管道传入)
	utf8Buf []byte
}

func NewUART(name string, addr uint64) *UART {
	u := &UART{name: name, addr: addr}
	u.data = make([]byte, UART_SIZE)
	u.utf8Buf = make([]byte, 0, 4)
	return u
}

func (p *UART) Name() string      { return p.name }
func (p *UART) AddrBegin() uint64 { return p.addr }
func (p *UART) AddrEnd() uint64   { return p.addr + uint64(len(p.data)) }

func (p *UART) Read(addr, size uint64) (uint64, error) {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return 0, fmt.Errorf("%s.Read: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch addr - p.addr {
	case UART_RHR: // RHR
		// TODO: 输入需要通过 Goroutine 支持
		if p.hasRX {
			p.hasRX = false
			return uint64(p.rx), nil
		}
		return 0, nil
	case UART_LCR: // LCR
		return 0x03, nil
	case UART_LSR: // LSR
		lsr := uint64(UART_LSR_TX)
		if p.hasRX {
			lsr |= UART_LSR_RX
		}
		return lsr, nil
	default:
		return 0, fmt.Errorf("%s.Read: unhandled read offset 0x%x", p.name, addr)
	}
}

func (p *UART) Write(addr, size, value uint64) error {
	if addr < p.AddrBegin() || addr >= p.AddrEnd() {
		return fmt.Errorf("%s.Write: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch addr - p.addr {
	case UART_RHR: // THR
		p.tx = uint8(value)
		p.utf8Buf = append(p.utf8Buf, p.tx)
		if utf8.FullRune(p.utf8Buf) {
			r, _ := utf8.DecodeRune(p.utf8Buf)
			p.utf8Buf = p.utf8Buf[:0]
			fmt.Printf("%c", r)
		}
		return nil
	case UART_LCR: // LCR
		return nil
	default:
		return fmt.Errorf("%s.Write: unhandled write offset 0x%x", p.name, addr)
	}
}
