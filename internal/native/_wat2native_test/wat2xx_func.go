// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

/*
#include <stdint.h>

static int32_t wat2xxI32Add_cgo(int32_t a, int32_t b) {
	return a+b;
}
*/
import "C"

func I32Add(a, b int32) int32 {
	return int32(C.wat2xxI32Add_cgo(C.int32_t(a), C.int32_t(b)))
}
