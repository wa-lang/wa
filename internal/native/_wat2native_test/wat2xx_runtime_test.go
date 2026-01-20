// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bytes"
	"testing"
	"unsafe"
)

func TestMemset(t *testing.T) {
	// 1. 准备一段 Go 的字节切片
	buf := make([]byte, 10)

	// 2. 获取指针并调用 C 函数
	ptr := uintptr(unsafe.Pointer(&buf[0]))
	var val byte = 0xAA
	var size int32 = 5

	Memset(ptr, val, size)

	// 3. 验证前 5 个字节是否被修改，后 5 个是否保持原样
	for i := 0; i < 5; i++ {
		if buf[i] != 0xAA {
			t.Errorf("Index %d should be 0xAA, got 0x%X", i, buf[i])
		}
	}
	for i := 5; i < 10; i++ {
		if buf[i] != 0x00 {
			t.Errorf("Index %d should be 0x00, got 0x%X", i, buf[i])
		}
	}
}

func TestMemcpy(t *testing.T) {
	// 1. 准备源数据和目标缓冲区
	src := []byte("SourceData") // 10 bytes
	dst := make([]byte, 10)

	srcPtr := uintptr(unsafe.Pointer(&src[0]))
	dstPtr := uintptr(unsafe.Pointer(&dst[0]))

	// 2. 拷贝前 6 个字节 "Source"
	Memcpy(dstPtr, srcPtr, 6)

	// 3. 验证结果
	expected := []byte("Source")
	if !bytes.Equal(dst[:6], expected) {
		t.Errorf("Expected %s, got %s", expected, dst[:6])
	}

	// 验证剩余部分是否未被干扰
	if dst[6] != 0x00 {
		t.Error("Memory beyond size should remain untouched")
	}
}
