// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build cgo && !inline_asm

package main

/*
#include <stdint.h>

int32_t wat2xx_global_i32_v = 0;
int64_t wat2xx_global_i64_v = 0;
float   wat2xx_global_f32_v = 0;
double  wat2xx_global_f64_v = 0;
*/
import "C"

func GlobalGetI32() int32 {
	return int32(C.wat2xx_global_i32_v)
}

func GlobalGetI64() int64 {
	return int64(C.wat2xx_global_i64_v)
}

func GlobalGetF32() float32 {
	return float32(C.wat2xx_global_f32_v)
}

func GlobalGetF64() float64 {
	return float64(C.wat2xx_global_f64_v)
}

func GlobalSetI32(v int32) {
	C.wat2xx_global_i32_v = C.int32_t(v)
}

func GlobalSetI64(v int64) {
	C.wat2xx_global_i64_v = C.int64_t(v)
}

func GlobalSetF32(v float32) {
	C.wat2xx_global_f32_v = C.float(v)
}

func GlobalSetF64(v float64) {
	C.wat2xx_global_f64_v = C.double(v)
}
