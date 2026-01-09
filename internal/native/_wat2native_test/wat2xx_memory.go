// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

/*
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

static uint8_t* wat2xx_memory = NULL;
static int32_t  wat2xx_memory_pages = 1;
static int32_t  wat2xx_memory_maxPages = 3;

static void wat2xx_memory_init() {
	wat2xx_memory = malloc(wat2xx_memory_maxPages*(64<<10));
}

static int32_t wat2xx_memory_pages_get() {
	return wat2xx_memory_pages;
}

static int32_t wat2xx_memory_maxPages_get() {
	return wat2xx_memory_maxPages;
}

static int32_t wat2xx_memory_grow(int32_t xPages) {
	int32_t old = wat2xx_memory_pages;
	if(wat2xx_memory_pages+xPages <= wat2xx_memory_maxPages) {
		wat2xx_memory = realloc(wat2xx_memory, (wat2xx_memory_pages+xPages)*(64<<10));
		if(wat2xx_memory == NULL) abort();
		wat2xx_memory_pages += xPages;
		return old;
	}
	return -1;
}

static int32_t wat2xx_I32Load(int32_t addr, int32_t offset) {
	int32_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static float wat2xx_F32Load(int32_t addr, int32_t offset) {
	float result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static double wat2xx_F64Load(int32_t addr, int32_t offset) {
	double result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}

static int32_t wat2xx_I32Load8_s(int32_t addr, int32_t offset) {
	int32_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int32_t wat2xx_I32Load8_u(int32_t addr, int32_t offset) {
	int32_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int32_t wat2xx_I32Load16_s(int32_t addr, int32_t offset) {
	int32_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int32_t wat2xx_I32Load16_u(int32_t addr, int32_t offset) {
	int32_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}

static int64_t wat2xx_I64Load8_s(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load8_u(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load16_s(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load16_u(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load32_s(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}
static int64_t wat2xx_I64Load32_u(int32_t addr, int32_t offset) {
	int64_t result;
	memcpy(&result, wat2xx_memory+addr+offset, sizeof(result));
	return result;
}

static void wat2xx_I32Store(int32_t addr, int32_t offset, int32_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_I64Store(int32_t addr, int32_t offset, int64_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_F32Store(int32_t addr, int32_t offset, float v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_F64Store(int32_t addr, int32_t offset, double v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}

static void wat2xx_I32Store8(int32_t addr, int32_t offset, int8_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_I32Store16(int32_t addr, int32_t offset, int16_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}

static void wat2xx_I64Store8(int32_t addr, int32_t offset, int8_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_I64Store16(int32_t addr, int32_t offset, int16_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
static void wat2xx_I64Store32(int32_t addr, int32_t offset, int32_t v) {
	memcpy(wat2xx_memory+addr+offset, &v, sizeof(v));
}
*/
import "C"

const PageSize = 64 << 10

func init() {
	C.wat2xx_memory_init()
}

func MemorySize() int32 {
	return int32(C.wat2xx_memory_pages_get())
}

func MemoryMaxSize() int32 {
	return int32(C.wat2xx_memory_maxPages_get())
}

func MemoryGrow(n int32) int32 {
	return int32(C.wat2xx_memory_grow(C.int32_t(n)))
}

func I32Load(addr, offset int32) int32 {
	return int32(C.wat2xx_I32Load(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load(C.int32_t(addr), C.int32_t(offset)))
}
func F32Load(addr, offset int32) float32 {
	return float32(C.wat2xx_F32Load(C.int32_t(addr), C.int32_t(offset)))
}
func F64Load(addr, offset int32) float64 {
	return float64(C.wat2xx_F64Load(C.int32_t(addr), C.int32_t(offset)))
}

func I32Load8_s(addr, offset int32) int32 {
	return int32(C.wat2xx_I32Load8_s(C.int32_t(addr), C.int32_t(offset)))
}
func I32Load8_u(addr, offset int32) int32 {
	return int32(C.wat2xx_I32Load8_u(C.int32_t(addr), C.int32_t(offset)))
}
func I32Load16_s(addr, offset int32) int32 {
	return int32(C.wat2xx_I32Load16_s(C.int32_t(addr), C.int32_t(offset)))
}
func I32Load16_u(addr, offset int32) int32 {
	return int32(C.wat2xx_I32Load16_u(C.int32_t(addr), C.int32_t(offset)))
}

func I64Load8_s(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load8_s(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load8_u(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load8_u(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load16_s(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load16_s(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load16_u(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load16_u(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load32_s(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load32_s(C.int32_t(addr), C.int32_t(offset)))
}
func I64Load32_u(addr, offset int32) int64 {
	return int64(C.wat2xx_I64Load32_u(C.int32_t(addr), C.int32_t(offset)))
}

func I32Store(addr, offset int32, value int32) {
	C.wat2xx_I32Store(C.int32_t(addr), C.int32_t(offset), C.int32_t(value))
}
func I64Store(addr, offset int32, value int64) {
	C.wat2xx_I64Store(C.int32_t(addr), C.int32_t(offset), C.int64_t(value))
}
func F32Store(addr, offset int32, value float32) {
	C.wat2xx_F32Store(C.int32_t(addr), C.int32_t(offset), C.float(value))
}
func F64Store(addr, offset int32, value float64) {
	C.wat2xx_F64Store(C.int32_t(addr), C.int32_t(offset), C.double(value))
}

func I32Store8(addr, offset int32, value int8) {
	C.wat2xx_I32Store8(C.int32_t(addr), C.int32_t(offset), C.int8_t(value))
}
func I32Store16(addr, offset int32, value int16) {
	C.wat2xx_I32Store16(C.int32_t(addr), C.int32_t(offset), C.int16_t(value))
}

func I64Store8(addr, offset int32, value int8) {
	C.wat2xx_I64Store8(C.int32_t(addr), C.int32_t(offset), C.int8_t(value))
}
func I64Store16(addr, offset int32, value int16) {
	C.wat2xx_I64Store16(C.int32_t(addr), C.int32_t(offset), C.int16_t(value))
}
func I64Store32(addr, offset int32, value int32) {
	C.wat2xx_I64Store32(C.int32_t(addr), C.int32_t(offset), C.int32_t(value))
}
