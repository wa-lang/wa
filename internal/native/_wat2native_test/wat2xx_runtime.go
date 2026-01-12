// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

/*
#include <stdint.h>
#include <string.h>
#include <stdio.h>

static void wat2xx_memcpy(uintptr_t dst, uintptr_t src, int32_t size) {
	memcpy((void*)dst, (void*)src, (size_t)size);
}
static void wat2xx_memset(uintptr_t dst, uint8_t value, int32_t size) {
	memset((void*)dst, value, size);
}

static void wat2xx_printString(uintptr_t str, int32_t len) {
	fwrite((void*)str, 1, len, stdout);
	fflush(stdout);
}
*/
import "C"

func Memcpy(dst, src uintptr, size int32) {
	C.wat2xx_memcpy(C.uintptr_t(dst), C.uintptr_t(src), C.int32_t(size))
}

func Memset(dst uintptr, value byte, size int32) {
	C.wat2xx_memset(C.uintptr_t(dst), C.uint8_t(value), C.int32_t(size))
}

func PrintString(s uintptr, size int32) {
	C.wat2xx_printString(C.uintptr_t(s), C.int32_t(size))
}
