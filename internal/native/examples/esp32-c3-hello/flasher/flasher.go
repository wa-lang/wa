package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"wa-lang.org/wa/internal/3rdparty/serial"
)

// 烧写配置
const (
	portName  = "/dev/ttyUSB0" // 根据您的系统修改为正确的串口名称 (e.g., COM3 on Windows)
	baudRate  = 115200
	dataFile  = "your_firmware.bin" // 假设这是你的汇编器生成的固件
	flashAddr = 0x42000000          // ESP32-C3 Flash的典型起始地址之一
)

func main() {
	// 1. 打开串口
	c := &serial.Config{Name: portName, Baud: baudRate}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("无法打开串口 %s: %v", portName, err)
	}
	defer port.Close()

	log.Printf("成功打开串口 %s", portName)

	// 2. 模拟 ROM Bootloader 握手 (SYNC)
	// 实际操作中，你需要发送一个复杂的命令序列，这里仅模拟发送
	fmt.Println("--- 模拟 ROM Bootloader 握手 ---")

	// 假设的同步命令（实际请查阅文档）
	syncCmd := []byte{0x07, 0x07, 0x12, 0x20}

	n, err := port.Write(syncCmd)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("发送了 %d 字节的同步命令。", n)

	// 等待设备响应 (实际需要解析响应，这里只是等待)
	time.Sleep(100 * time.Millisecond)

	// 3. 模拟发送固件
	fmt.Println("--- 模拟发送固件数据 ---")

	firmware, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("无法读取固件文件 %s: %v", dataFile, err)
	}

	// 假设的写入 Flash 命令头 (实际会包含地址、长度、校验和等)
	// 在实际烧写中，会是 (Command Header + Data Segment Header + Firmware Chunk)

	dataChunkSize := 256 // 假设每次发送256字节

	for i := 0; i < len(firmware); i += dataChunkSize {
		end := i + dataChunkSize
		if end > len(firmware) {
			end = len(firmware)
		}

		chunk := firmware[i:end]

		// 实际烧写时，需要在发送数据块前发送一个带有目标地址和长度的命令包

		// 模拟发送数据块
		n, err := port.Write(chunk)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("发送数据块: 地址 0x%X, 长度 %d 字节", flashAddr+i, n)

		// 实际烧写中，每发送一个块都需要等待设备的 ACK 确认
		time.Sleep(1 * time.Millisecond)
	}

	// 4. 模拟运行命令 (RUN)
	fmt.Println("--- 模拟 RUN 命令 ---")
	// 实际操作中，你需要发送一个 RUN 命令，指定程序入口地址

	// 5. 切换到串口监视模式
	fmt.Println("--- 切换到串口监视模式，等待设备输出 ---")

	// 启动一个 goroutine 来持续读取串口数据
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := port.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("\n串口读取错误: %v\n", err)
				}
				return
			}
			if n > 0 {
				fmt.Printf("%s", buf[:n])
			}
		}
	}()

	// 保持主程序运行，直到用户按下 Enter
	fmt.Println("烧写和监视器已启动。按下 Enter 键退出...")
	fmt.Scanln()
}
