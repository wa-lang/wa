package main

// ESP32-C3 固件映像头 (Image Header)
// 结构体字段顺序和大小必须与硬件文档严格匹配
type ImageHeader struct {
	Magic        uint32   // 魔术字：必须是 0xE9
	SegmentCount byte     // 段的数量 (对于单段程序为 1)
	SpiMode      byte     // SPI Flash 模式 (通常为 0x02)
	SpiSpeed     byte     // SPI Flash 速度 (通常为 0x01)
	SpiSize      byte     // SPI Flash 大小 (通常为 0x20)
	EntryAddr    uint32   // 程序入口点地址 (例如: 0x42000000)
	Reserved     [16]byte // 保留字段
	Checksum     byte     // 映像头和所有段头的校验和（通常放在最后计算）
}

// 段头 (Segment Header)
type SegmentHeader struct {
	LoadAddr uint32 // 该段数据将被加载到的 RAM/IRAM 地址
	DataLen  uint32 // 该段数据的字节长度
}

const (
	ESP_IMAGE_MAGIC = 0xE9 // ESP32-C3 固件映像魔术字
	// 假设的程序入口点地址 (Flash 地址)
	// ROM Bootloader 会将它作为入口点执行
	DEFAULT_ENTRY_ADDR = 0x42000000
)

// 假设这是您的 Go 汇编器生成的原始机器指令
var rawMachineCode = []byte{
	// 假设这是 RISC-V 指令的机器码序列
	0x13, 0x05, 0x00, 0x00, // li a0, 0
	0x67, 0x80, 0x00, 0x00, // ret
	// 实际代码会更长，包括 UART 操作
}
