// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("test wat2xx asm code!")
	fmt.Println("1+2=", I32Add(1, 2))

	data := []byte{'H', 'e', 'l', 'l', 'o', 'W', 'a'}
	ptr := uintptr(unsafe.Pointer(&data[0]))
	PrintString(ptr, 5)
}
