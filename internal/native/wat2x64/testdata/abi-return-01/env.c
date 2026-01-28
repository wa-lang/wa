// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

#if defined(_WIN64)
#   include <io.h>
#else
#   include <unistd.h>
#endif

extern uintptr_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

int _Wa_Runtime_write(int fd, void *buf, int count) {
#if defined(_WIN64)
    return _write(fd, buf, (unsigned int)(count));
#else
    return write(fd, buf, (size_t)(count));
#endif
}

void _Wa_Runtime_exit(int status) {
    exit(status);
}

void* _Wa_Runtime_malloc(int size) {
    return malloc((size_t)(size));
}

void* _Wa_Runtime_memcpy(void* dst, const void* src, int n) {
    return memcpy(dst, src, (size_t)(n));
}

void* _Wa_Runtime_memmove(void* dst, const void* src, int n) {
    return memmove(dst, src, (size_t)(n));
}

void* _Wa_Runtime_memset(void* s, int c, int n) {
    return memset(s, c, (size_t)(n));
}

// 定义对应 Wasm 多返回值的结构体
typedef struct {
    int64_t v1;
    int64_t v2;
    int64_t v3;
} EnvMultiRet;

// C 函数: 返回结构体
// 注意: 根据 Win64 ABI，这会被编译为:
// void _Wa_Import_env_get_multi_values(EnvMultiRet* hidden_ptr, int64_t input_val)
EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val) {
    EnvMultiRet r;
    r.v1 = input_val + 1;
    r.v2 = input_val + 2;
    r.v3 = input_val + 3;
    printf(
        "C side: input=%lld, returning [%lld, %lld, %lld]\n", 
        input_val, r.v1, r.v2, r.v3
    );
    return r;
}

void _Wa_Import_env_print_i64(int64_t x) {
    printf("printI64: %lld\n", x);
}

