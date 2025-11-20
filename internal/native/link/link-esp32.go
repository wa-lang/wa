// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package link

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/link/esp32"
)

// 生成 esp32 bin 文件
func LinkESP32Bin(prog *abi.LinkedProgram) ([]byte, error) {
	switch prog.CPU {
	case abi.RISCV32:
		return _BuildESP32Image(prog)
	default:
		return nil, fmt.Errorf("link.LinkESP32C3Bin: unknown cpu type: %v", prog.CPU)
	}
}

func _BuildESP32Image(prog *abi.LinkedProgram) ([]byte, error) {
	hdr := &esp32.ImageHeader{
		Magic:          esp32.Magic,
		SegmentCount:   2, // 指令+数据
		SpiMode:        esp32.SPIMode_DIO,
		SpiSpeedSize:   esp32.FlashSpeedSize_2MB | esp32.FlashSpeedSize_80MHz,
		EntryAddr:      uint32(prog.TextAddr),
		WPPin:          0xEE, // 参考 C3 例子, 可能是无效端口
		SpiPinDrv:      [3]uint8{},
		ChipID:         esp32.ESP_CHIP_ID_ESP32C3, // TODO: 传入参数
		MinChipRev:     0x03,
		MinChipRevFull: 0x03,
		MaxChipRevFull: 0xC7,
		HashAppended:   0, // 忽略 SHA256
	}

	var checksum = uint8(0)
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.LittleEndian, hdr); err != nil {
		return nil, err
	}

	// 指令段
	if len(prog.TextData) > 0 {
		segHdr := esp32.SegmentHeader{
			LoadAddr:   uint32(prog.TextAddr),
			DataLength: uint32(len(prog.TextData)),
		}
		if err := binary.Write(&buffer, binary.LittleEndian, segHdr); err != nil {
			return nil, err
		}
		if _, err := buffer.Write(prog.TextData); err != nil {
			return nil, err
		}

		for _, b := range prog.TextData {
			checksum ^= b
		}
	}

	// 数据段
	if len(prog.DataData) > 0 {
		segHdr := esp32.SegmentHeader{
			LoadAddr:   uint32(prog.DataAddr),
			DataLength: uint32(len(prog.DataData)),
		}
		if err := binary.Write(&buffer, binary.LittleEndian, segHdr); err != nil {
			return nil, err
		}
		if _, err := buffer.Write(prog.DataData); err != nil {
			return nil, err
		}

		for _, b := range prog.DataData {
			checksum ^= b
		}
	}

	// 填充到16的倍数-1
	paddingNeeded := (16 - ((len(buffer.Bytes()) + 1) % 16)) % 16
	if paddingNeeded > 0 {
		padding := make([]byte, paddingNeeded)
		buffer.Write(padding)
	}

	// 计算校验和
	checksum ^= 0xEF
	if _, err := buffer.Write([]uint8{checksum}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
