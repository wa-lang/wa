// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

/*
#include <stdint.h>

// i32 constants
#define W2_I32_MAX   2147483647
#define W2_I32_MIN   -2147483648
#define W2_I32_HEX   0x12345678

// i64 constants
#define W2_I64_MAX   9223372036854775807LL
#define W2_I64_MIN   ((int64_t)0x8000000000000000ULL)
#define W2_I64_HEX   0x1122334455667788LL

// f32 bit patterns (IEEE 754)
#define W2_F32_PI_BITS    0x40490fdb          // 3.1415926
#define W2_F32_INF_BITS   0x7f800000          // +Infinity
#define W2_F32_NEG0_BITS  0x80000000          // -0.0

// f64 bit patterns (IEEE 754)
#define W2_F64_PI_BITS    0x400921fb54442d18ULL // 3.141592653589793
#define W2_F64_NAN_BITS   0x7ff8000000000000ULL // Quiet NaN

static int32_t wat2xx_I32Const_W2_I32_MAX() {
	return W2_I32_MAX;
}
static int32_t wat2xx_I32Const_W2_I32_MIN() {
	return W2_I32_MIN;
}
static int32_t wat2xx_I32Const_W2_I32_HEX() {
	return W2_I32_HEX;
}

static int64_t wat2xx_I64Const_W2_I64_MAX() {
	return W2_I64_MAX;
}
static int64_t wat2xx_I64Const_W2_I64_MIN() {
	return W2_I64_MIN;
}
static int64_t wat2xx_I64Const_W2_I64_HEX() {
	return W2_I64_HEX;
}

static float wat2xx_F32Const_W2_F32_PI() {
	union { float f; uint32_t u32; } v;
	v.u32 = W2_F32_PI_BITS;
	return v.f;
}
static float wat2xx_F32Const_W2_F32_INF() {
	union { float f; uint32_t u32; } v;
	v.u32 = W2_F32_INF_BITS;
	return v.f;
}
static float wat2xx_F32Const_W2_F32_NEG0() {
	union { float f; uint32_t u32; } v;
	v.u32 = W2_F32_NEG0_BITS;
	return v.f;
}

static double wat2xx_F64Const_W2_F64_PI() {
	union { double f; uint64_t u64; } v;
	v.u64 = W2_F64_PI_BITS;
	return v.f;
}
static double wat2xx_F64Const_W2_F64_NAN() {
	union { double f; uint64_t u64; } v;
	v.u64 = W2_F64_NAN_BITS;
	return v.f;
}
*/
import "C"
import "math"

const (
	W2_I32_MAX int32 = C.W2_I32_MAX
	W2_I32_MIN int32 = C.W2_I32_MIN
	W2_I32_HEX int32 = C.W2_I32_HEX

	W2_I64_MAX int64 = C.W2_I64_MAX
	W2_I64_MIN int64 = C.W2_I64_MIN
	W2_I64_HEX int64 = C.W2_I64_HEX

	W2_F32_PI_BITS   uint32 = C.W2_F32_PI_BITS
	W2_F32_INF_BITS  uint32 = C.W2_F32_INF_BITS
	W2_F32_NEG0_BITS uint32 = C.W2_F32_NEG0_BITS

	W2_F64_PI_BITS  uint64 = C.W2_F64_PI_BITS
	W2_F64_NAN_BITS uint64 = C.W2_F64_NAN_BITS
)

func I32Const(x int32) int32 {
	switch x {
	case W2_I32_MAX:
		return int32(C.wat2xx_I32Const_W2_I32_MAX())
	case W2_I32_MIN:
		return int32(C.wat2xx_I32Const_W2_I32_MIN())
	case W2_I32_HEX:
		return int32(C.wat2xx_I32Const_W2_I32_HEX())
	default:
		panic("unreachable")
	}
}
func I64Const(x int64) int64 {
	switch x {
	case W2_I64_MAX:
		return int64(C.wat2xx_I64Const_W2_I64_MAX())
	case W2_I64_MIN:
		return int64(C.wat2xx_I64Const_W2_I64_MIN())
	case W2_I64_HEX:
		return int64(C.wat2xx_I64Const_W2_I64_HEX())
	default:
		panic("unreachable")
	}
}
func F32Const(x float32) float32 {
	switch math.Float32bits(x) {
	case W2_F32_PI_BITS:
		return float32(C.wat2xx_F32Const_W2_F32_PI())
	case W2_F32_INF_BITS:
		return float32(C.wat2xx_F32Const_W2_F32_INF())
	case W2_F32_NEG0_BITS:
		return float32(C.wat2xx_F32Const_W2_F32_NEG0())
	default:
		panic("unreachable")
	}
}
func F64Const(x float64) float64 {
	switch math.Float64bits(x) {
	case W2_F64_PI_BITS:
		return float64(C.wat2xx_F64Const_W2_F64_PI())
	case W2_F64_NAN_BITS:
		return float64(C.wat2xx_F64Const_W2_F64_NAN())
	default:
		panic("unreachable")
	}
}
