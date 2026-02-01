// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"log"

	"wa-lang.org/wa/internal/native/link/elf"
)

func main() {
	// 1. 打开动态库文件
	f, err := elf.Open("libhello.so")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer f.Close()

	// 2. 获取动态符号表 (Dynamic Symbols)
	// 对于 .so 文件，导出的函数通常在 DynamicSymbols 中
	symbols, err := f.DynamicSymbols()
	if err != nil {
		log.Fatalf("无法读取动态符号表: %v", err)
	}

	// 3. 遍历符号寻找 "add"
	found := false
	for _, sym := range symbols {
		if sym.Name == "add" {
			fmt.Printf("找到符号: %s\n", sym.Name)
			fmt.Printf("偏移量 (Value): 0x%x\n", sym.Value)
			fmt.Printf("大小 (Size): %d bytes\n", sym.Size)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("未在 libhello.so 中找到 'add' 函数")
	}
}
