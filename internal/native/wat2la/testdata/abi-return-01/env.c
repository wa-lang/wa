// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdio.h>
#include <stdint.h>

typedef struct EnvMultiRet EnvMultiRet;

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)        asm(".Wa.Runtime.write");
extern void  _Wa_Runtime_exit(int status)                           asm(".Wa.Runtime.exit");
extern void* _Wa_Runtime_malloc(int size)                           asm(".Wa.Runtime.malloc");
extern void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)  asm(".Wa.Runtime.memcpy");
extern void* _Wa_Runtime_memmove(void* dst, const void* src, int n) asm(".Wa.Runtime.memmove");
extern void* _Wa_Runtime_memset(void* s, int c, int n)              asm(".Wa.Runtime.memset");

extern EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val) asm(".Wa.Import.env.get_multi_values");

extern void _Wa_Import_env_print_i64(int64_t x) asm(".Wa.Import.env.print_i64");

// 定义对应 Wasm 多返回值的结构体
struct EnvMultiRet {
    int64_t v1;
    int64_t v2;
    int64_t v3;
};

// C 函数: 返回结构体
// void _Wa_Import_env_get_multi_values(EnvMultiRet* hidden_ptr, int64_t input_val)
EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val) {
    EnvMultiRet r;
    r.v1 = input_val + 1;
    r.v2 = input_val + 2;
    r.v3 = input_val + 3;
    printf(
        "C side: input=%ld, returning [%ld, %ld, %ld]\n", 
        input_val, r.v1, r.v2, r.v3
    );
    return r;
}

void wat2la_env_print_i64(int64_t x) {
    printf("printI64: %ld\n", x);
}

int _Wa_Runtime_write(int fd, void *buf, int count) {
#if defined(_WIN64)
    return _write(fd, buf, (size_t)(count));
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
