// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

/*
#include <stdint.h>

static void wat2xx_func_foo(){}
static void wat2xx_func_bar(){}

static int64_t wat2xx_GetFuncPtr_foo() {
	return (int64_t)(wat2xx_func_foo);
}
static int64_t wat2xx_GetFuncPtr_bar() {
	return (int64_t)(wat2xx_func_bar);
}

static int64_t wat2xx_table_array[10];

static int32_t wat2xx_table_size() {
	return sizeof(wat2xx_table_array)/sizeof(wat2xx_table_array[0]);
}
static int64_t wat2xx_TabelGet(int i) {
	return wat2xx_table_array[i];
}
static void wat2xx_TabelSet(int i, int64_t v) {
	wat2xx_table_array[i] = v;
}
*/
import "C"

func GetFuncPtr_foo() int64 {
	return int64(C.wat2xx_GetFuncPtr_foo())
}
func GetFuncPtr_bar() int64 {
	return int64(C.wat2xx_GetFuncPtr_bar())
}

func TableSize() int32 {
	return int32(C.wat2xx_table_size())
}

func TabelGet(i int32) int64 {
	return int64(C.wat2xx_TabelGet(C.int32_t(i)))
}

func TabelSet(i int32, v int64) {
	C.wat2xx_TabelSet(C.int32_t(i), C.int64_t(v))
}
