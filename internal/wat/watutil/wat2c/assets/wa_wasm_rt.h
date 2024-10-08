// 版权 @2024 凹语言 作者。保留所有权利。

#pragma once

#ifndef WA_WASM_RT_H_
#define WA_WASM_RT_H_

#include <assert.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct wa_wasm_global_t {
	uint64_t* data;
	size_t size;
} wa_wasm_global_t;

typedef struct wa_wasm_memory_t {
	uint8_t* data;
	size_t size;
} wa_wasm_memory_t;

typedef struct wa_wat2c_table_t {
	uint64_t* data;
	size_t size;
} wa_wasm_table_t;

typedef struct wa_wat2c_module_t {
	wa_wasm_global_t global;
	wa_wasm_memory_t memory;
	wa_wasm_table_t table;
} wa_wasm_module_t;

typedef struct wa_wat2c_stack_t {
	uint64_t* stack;
	size_t size;
	size_t capacity;
} wa_wasm_instance_t;

#ifdef __cplusplus
}
#endif
#endif // WA_WASM_RT_H_
