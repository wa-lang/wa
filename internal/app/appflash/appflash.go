// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appflash

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/3rdparty/serial"
)

var CmdFlash = &cli.Command{
	Hidden:    true,
	Name:      "flash",
	Usage:     "flash a bare-metal firmware image to an ESP32-C3 device.",
	ArgsUsage: "[firmware.bin]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port, p",
			Usage:   "Serial port device to use for flashing (e.g., /dev/ttyUSB0 or COM3)",
			Value:   "/dev/ttyUSB0",
			Aliases: []string{"p"},
		},
		&cli.IntFlag{
			Name:    "baud, b",
			Usage:   "Baud rate for flashing communication",
			Value:   115200,
			Aliases: []string{"b"},
		},
		&cli.StringFlag{
			Name:    "addr, a",
			Usage:   "Flash memory address to start writing the firmware (e.g., 0x42000000)",
			Value:   "0x42000000",
			Aliases: []string{"a"},
		},
	},
	Action: func(c *cli.Context) error {
		if false {
			cmdMain() // TODO
		}
		return nil
	},
}

const (
	ESP_SYNC        = 0x08
	ESP_FLASH_BEGIN = 0x02
	ESP_FLASH_DATA  = 0x03
	ESP_FLASH_END   = 0x04
)

func cmdMain() {
	port := openPort("/dev/ttyUSB0", 115200)

	// step 1: handshake
	if err := espSync(port); err != nil {
		panic(err)
	}

	// step 2: load firmware
	data, _ := os.ReadFile("firmware.bin")

	// step 3: flash
	flash(port, data, 0x0) // 地址 0x0
}

func openPort(name string, baud int) io.ReadWriteCloser {
	cfg := &serial.Config{
		Name:        name,
		Baud:        baud,
		ReadTimeout: time.Millisecond * 500,
	}
	p, err := serial.OpenPort(cfg)
	if err != nil {
		panic(err)
	}
	return p
}

func slipEncode(data []byte) []byte {
	out := []byte{0xC0}
	for _, b := range data {
		switch b {
		case 0xC0:
			out = append(out, 0xDB, 0xDC)
		case 0xDB:
			out = append(out, 0xDB, 0xDD)
		default:
			out = append(out, b)
		}
	}
	out = append(out, 0xC0)
	return out
}

func espSync(port io.ReadWriteCloser) error {
	// SYNC payload（固定）
	payload := []byte{
		0x07, 0x07, 0x12, 0x20,
	}
	for i := 0; i < 32; i++ {
		payload = append(payload, 0x55)
	}

	// 组包
	packet := buildPacket(ESP_SYNC, payload)
	port.Write(packet)

	// 等待设备回应
	buf := make([]byte, 100)
	_, err := port.Read(buf)
	return err
}

func buildPacket(cmd byte, data []byte) []byte {
	header := make([]byte, 8)
	header[0] = cmd
	header[1] = byte(len(data))
	header[2] = byte(len(data) >> 8)
	header[3] = byte(len(data) >> 16)
	header[4] = 0 // checksum; XOR
	header[5] = 0
	header[6] = 0
	header[7] = 0

	// 简单 XOR 校验
	chk := byte(0)
	for _, b := range data {
		chk ^= b
	}
	header[4] = chk

	return slipEncode(append(header, data...))
}

func flash(port io.ReadWriteCloser, fw []byte, offset int) {
	blockSize := 0x400
	numBlocks := (len(fw) + blockSize - 1) / blockSize

	// FLASH_BEGIN
	beginData := make([]byte, 16)
	binary.LittleEndian.PutUint32(beginData[0:], uint32(len(fw)))
	binary.LittleEndian.PutUint32(beginData[4:], uint32(numBlocks))
	binary.LittleEndian.PutUint32(beginData[8:], uint32(blockSize))
	binary.LittleEndian.PutUint32(beginData[12:], uint32(offset))

	p := buildPacket(ESP_FLASH_BEGIN, beginData)
	port.Write(p)
	if err := waitAckOK(port); err != nil {
		panic(err)
	}

	// FLASH_DATA
	for seq := 0; seq < numBlocks; seq++ {
		start := seq * blockSize
		end := start + blockSize
		if end > len(fw) {
			end = len(fw)
		}
		buf := fw[start:end]

		data := make([]byte, 16)
		binary.LittleEndian.PutUint32(data[0:], uint32(len(buf)))
		binary.LittleEndian.PutUint32(data[4:], uint32(seq))
		binary.LittleEndian.PutUint32(data[8:], uint32(0)) // flash offset; unused
		binary.LittleEndian.PutUint32(data[12:], uint32(0))

		p := buildPacket(ESP_FLASH_DATA, append(data, buf...))
		port.Write(p)
		waitAck(port)
	}

	// FLASH_END
	endData := make([]byte, 4)
	binary.LittleEndian.PutUint32(endData, 1) // reboot = 1

	p = buildPacket(ESP_FLASH_END, endData)
	port.Write(p)
	if err := waitAckOK(port); err != nil {
		panic(err)
	}
}

// ErrTimeout 表示等待响应超时
var ErrTimeout = errors.New("timeout waiting for response")

// waitAck 从串口读取一个 SLIP 包并返回解码后的 payload。超时返回 ErrTimeout。
// 注意：上层需要根据协议检查 payload 内容（比如 header 或状态码）。
func waitAck(port io.ReadWriteCloser) ([]byte, error) {
	// 使用短超时避免永久阻塞；如果你在高波特率下需要更短或更长，可调整
	return readSLIPPacket(port, 3*time.Second)
}

// readSLIPPacket 读取一个以 0xC0 开始并以 0xC0 结束的 SLIP 包，处理转义字节（0xDB）。
// 返回的是 SLIP 解码后的内部 payload（不含外层 0xC0）。
func readSLIPPacket(port io.ReadWriteCloser, timeout time.Duration) ([]byte, error) {
	r := bufio.NewReader(port)
	deadline := time.After(timeout)

	// 等待起始 0xC0
	for {
		select {
		case <-deadline:
			return nil, ErrTimeout
		default:
		}

		b, err := r.ReadByte()
		if err != nil {
			// 如果串口设置了 ReadTimeout，会因为超时返回错误；把它包装成超时
			if err == io.EOF {
				// EOF from underlying stream, treat as timeout-like
				return nil, ErrTimeout
			}
			return nil, err
		}
		if b == 0xC0 {
			break
		}
		// otherwise keep skipping until start marker
	}

	// 读取直到下一个 0xC0，处理转义
	out := make([]byte, 0, 256)
	for {
		select {
		case <-deadline:
			return nil, ErrTimeout
		default:
		}

		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil, ErrTimeout
			}
			return nil, err
		}

		if b == 0xC0 {
			// end of SLIP packet
			break
		}
		if b == 0xDB {
			// escape sequence: next byte should be 0xDC or 0xDD
			nb, err := r.ReadByte()
			if err != nil {
				if err == io.EOF {
					return nil, ErrTimeout
				}
				return nil, err
			}
			if nb == 0xDC {
				out = append(out, 0xC0)
			} else if nb == 0xDD {
				out = append(out, 0xDB)
			} else {
				// 非标准转义：把两字节都放回（尽量恢复）
				out = append(out, nb)
			}
		} else {
			out = append(out, b)
		}
		// 简单保护：避免无穷增长
		if len(out) > 10*1024*1024 {
			return nil, fmt.Errorf("frame too large")
		}
	}

	return out, nil
}

func waitAckOK(port io.ReadWriteCloser) error {
	payload, err := waitAck(port)
	if err != nil {
		return err
	}
	// 下面是示例性的“成功判定”。按你需要替换为真实的协议字段判断。
	if len(payload) >= 4 {
		// 假设 header[0] 为 status（0 == OK）
		if payload[0] == 0x00 {
			return nil
		}
		return fmt.Errorf("device returned error status: 0x%02x", payload[0])
	}
	return fmt.Errorf("unexpected response length: %d", len(payload))
}
