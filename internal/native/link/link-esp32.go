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
		Magic:        esp32.Magic,
		SegmentCount: 0,
		SpiMode:      esp32.SPIMode_DIO,
		SpiSpeedSize: esp32.FlashSpeedSize_2MB | esp32.FlashSpeedSize_80MHz,
		EntryAddr:    uint32(prog.TextAddr),
	}

	extHdr := &esp32.ImageExtendedHeader{
		WPPin:          0xEE, // 参考 C3 例子, 可能是无效端口
		SpiPinDrv:      [3]uint8{},
		ChipID:         esp32.ESP_CHIP_ID_ESP32C3, // TODO: 传入参数
		MinChipRev:     0x03,
		MinChipRevFull: 0x03,
		MaxChipRevFull: 0xC7,
		HashAppended:   0, // 忽略 SHA256
	}

	// 段数据
	var segBuffer bytes.Buffer

	// 指令段
	if len(prog.TextData) > 0 {
		hdr.SegmentCount++
		segHdr := esp32.SegmentHeader{
			LoadAddr:   uint32(prog.TextAddr),
			DataLength: uint32(len(prog.TextData)),
		}
		if err := binary.Write(&segBuffer, binary.LittleEndian, segHdr); err != nil {
			return nil, err
		}
		if _, err := segBuffer.Write(prog.TextData); err != nil {
			return nil, err
		}
	}

	// 数据段
	if len(prog.DataData) > 0 {
		hdr.SegmentCount++
		segHdr := esp32.SegmentHeader{
			LoadAddr:   uint32(prog.DataAddr),
			DataLength: uint32(len(prog.DataData)),
		}
		if err := binary.Write(&segBuffer, binary.LittleEndian, segHdr); err != nil {
			return nil, err
		}
		if _, err := segBuffer.Write(prog.DataData); err != nil {
			return nil, err
		}
	}

	// 写入缓存
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.LittleEndian, hdr); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.LittleEndian, extHdr); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(segBuffer.Bytes()); err != nil {
		return nil, err
	}

	paddingNeeded := (16 - ((len(buffer.Bytes()) + 1) % 16)) % 16
	if paddingNeeded > 0 {
		padding := make([]byte, paddingNeeded)
		buffer.Write(padding)
	}

	// 计算校验和(可能有BUG)
	checksum := esp32CalculateChecksum(segBuffer.Bytes())
	if _, err := buffer.Write([]uint8{checksum}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// 计算校验和
// TODO: 如果不是4字节倍数, 怎么处理?
func esp32CalculateChecksum(data []byte) uint8 {
	const ESP_ROM_CHECKSUM_INITIAL = uint32(0xEF)

	r := bytes.NewReader(data)
	checksum_word := ESP_ROM_CHECKSUM_INITIAL

	for {
		var w uint32
		err := binary.Read(r, binary.LittleEndian, &w)
		if err != nil {
			break
		}
		checksum_word ^= w
	}

	checksum := (checksum_word >> 24) ^ (checksum_word >> 16) ^ (checksum_word >> 8) ^ (checksum_word >> 0)
	return uint8(checksum)
}
