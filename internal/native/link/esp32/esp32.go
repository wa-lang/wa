// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// https://github.com/espressif/esp-idf/blob/master/components/bootloader_support/include/esp_app_format.h

package esp32

// 魔数
const Magic = uint8(0xE9)

// 芯片ID
type ChipIDType uint16

const (
	ESP_CHIP_ID_ESP32    ChipIDType = 0x0000
	ESP_CHIP_ID_ESP32S2  ChipIDType = 0x0002
	ESP_CHIP_ID_ESP32C3  ChipIDType = 0x0005
	ESP_CHIP_ID_ESP32S3  ChipIDType = 0x0009
	ESP_CHIP_ID_ESP32C2  ChipIDType = 0x000C
	ESP_CHIP_ID_ESP32C6  ChipIDType = 0x000D
	ESP_CHIP_ID_ESP32H2  ChipIDType = 0x0010
	ESP_CHIP_ID_ESP32P4  ChipIDType = 0x0012
	ESP_CHIP_ID_ESP32C5  ChipIDType = 0x0017
	ESP_CHIP_ID_ESP32C61 ChipIDType = 0x0014
	ESP_CHIP_ID_ESP32H21 ChipIDType = 0x0019
	ESP_CHIP_ID_ESP32H4  ChipIDType = 0x001C
	ESP_CHIP_ID_INVALID  ChipIDType = 0xFFFF
)

// SPI Flash 模式
type SPIMode uint8

const (
	SPIMode_QIO       SPIMode = 0
	SPIMode_QOUT      SPIMode = 1
	SPIMode_DIO       SPIMode = 2
	SPIMode_DOUT      SPIMode = 3
	SPIMode_FAST_READ SPIMode = 4
	SPIMode_SLOW_READ SPIMode = 5
)

// Flash类型(容量和速度)
type FlashSpeedSize uint8

const (
	// 高四位 - Flash大小
	FlashSpeedSize_1MB   FlashSpeedSize = 0x00
	FlashSpeedSize_2MB   FlashSpeedSize = 0x01
	FlashSpeedSize_4MB   FlashSpeedSize = 0x02
	FlashSpeedSize_6MB   FlashSpeedSize = 0x03
	FlashSpeedSize_16MB  FlashSpeedSize = 0x04
	FlashSpeedSize_64MB  FlashSpeedSize = 0x05
	FlashSpeedSize_128MB FlashSpeedSize = 0x06

	// 低四位 - Flash频率
	FlashSpeedSize_40MHz FlashSpeedSize = 0x00 // 1/2 频率
	FlashSpeedSize_26MHz FlashSpeedSize = 0x10 // 1/3 频率
	FlashSpeedSize_20MHz FlashSpeedSize = 0x20 // 1/4 频率
	FlashSpeedSize_80MHz FlashSpeedSize = 0xF0 // 等于时钟频率
)

// 文件头部(Image header)
type ImageHeader struct {
	Magic        uint8          // 固定的 Magic 字节 0xE9
	SegmentCount uint8          // 段的数量
	SpiMode      SPIMode        // Flash 模式
	SpiSpeedSize FlashSpeedSize // Flash 大小和频率
	EntryAddr    uint32         // 程序的入口点地址 (_start)
}

// 扩展头部(Image extended header)
type ImageExtendedHeader struct {
	WPPin          uint8      // 写保护引脚, C3的例子默认是 0xEE
	SpiPinDrv      [3]uint8   // SPI Flash 引脚的驱动设置, 默认值 0x00
	ChipID         ChipIDType // 芯片ID
	MinChipRev     uint8      // 镜像支持的最低芯片版本(已弃用)
	MinChipRevFull uint16     // 镜像支持的最低芯片版本, 格式: major * 100 + minor
	MaxChipRevFull uint16     // 镜像支持的最高芯片版本, 格式: major * 100 + minor
	Reserved       [4]byte    // 保留字节
	HashAppended   uint8      // 如果为1, 则SHA256摘要附加在校验和之后
}

// 段数据头部
type SegmentHeader struct {
	LoadAddr   uint32 // 段应该加载到的内存地址
	DataLength uint32 // 段的二进制数据长度
}
