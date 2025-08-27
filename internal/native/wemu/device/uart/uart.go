// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 串口设备
package uart

import "fmt"

// 默认的地址
const (
	UART_BASE   = 0x10000000
	UART_SIZE   = 0x100
	UART_RHR    = UART_BASE + 0
	UART_THR    = UART_BASE + 0
	UART_LCR    = UART_BASE + 3
	UART_LSR    = UART_BASE + 5
	UART_LSR_RX = 1
	UART_LSR_TX = 1 << 5

	UART_IRQ = 10
)

type UART struct {
	name string
	data []byte

	rx    uint8 // 最近接收到的字符
	tx    uint8 // 最近写出的字符
	hasRX bool  // 是否有待读数据(需要通过其他Goroutine管道传入)
}

func NewUART(name string) *UART {
	u := &UART{name: name}
	u.data = make([]byte, UART_SIZE)
	return u
}

func (p *UART) Name() string { return p.name }
func (p *UART) Size() uint64 { return uint64(len(p.data)) }

func (p *UART) Read(addr, size uint64) (uint64, error) {
	if addr+size >= p.Size() {
		return 0, fmt.Errorf("%s: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch addr {
	case 0: // RHR
		if p.hasRX {
			p.hasRX = false
			return uint64(p.rx), nil
		}
		return 0, nil
	case 3: // LCR
		return 0x03, nil
	case 5: // LSR
		lsr := uint64(UART_LSR_TX)
		if p.hasRX {
			lsr |= UART_LSR_RX
		}
		return lsr, nil
	default:
		return 0, fmt.Errorf("uart: unhandled read offset 0x%x", addr)
	}
}

func (p *UART) Write(addr, size, value uint64) error {
	if addr+size >= p.Size() {
		return fmt.Errorf("%s: bad address [0x%08X, 0x%x08X)", p.name, addr, addr+size)
	}
	switch addr {
	case 0: // THR
		p.tx = uint8(value)
		fmt.Printf("%c", p.tx)
		return nil
	case 3: // LCR
		return nil
	default:
		return fmt.Errorf("uart: unhandled write offset 0x%x", addr)
	}
}
