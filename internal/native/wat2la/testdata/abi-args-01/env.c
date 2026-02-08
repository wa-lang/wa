// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)        asm(".Wa.Runtime.write");
extern void  _Wa_Runtime_exit(int status)                           asm(".Wa.Runtime.exit");
extern void* _Wa_Runtime_malloc(int size)                           asm(".Wa.Runtime.malloc");
extern void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)  asm(".Wa.Runtime.memcpy");
extern void* _Wa_Runtime_memmove(void* dst, const void* src, int n) asm(".Wa.Runtime.memmove");
extern void* _Wa_Runtime_memset(void* s, int c, int n)              asm(".Wa.Runtime.memset");

extern int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) asm(".Wa.Import.env.write");
extern void    _Wa_Import_env_print_i64(int64_t x)                                                                                   asm(".Wa.Import.env.print_i64");

// 按照 龙芯 调用约定
int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8, int64_t p9, int64_t p10) {
    printf("--- Linux/Loong64 ABI Verification ---\n");
    printf("Param 1 (a0): %ld\n", fd);
    printf("Param 2 (a1): %ld\n", (int64_t)(ptr));
    printf("Param 3 (a2): %ld\n", size);
    printf("Param 4 (a3): %ld\n", p4);
    printf("Param 5 (a4): %ld\n", p5);
    printf("Param 6 (a5): %ld\n", p6);
    printf("Param 7 (a6): %ld\n", p7);
    printf("Param 8 (a7): %ld\n", p8);
    printf("Param 9 (Stack RSP+0): %ld\n", p9);
    printf("Param10 (Stack RSP+8): %ld\n", p10);
    printf("-------------------------------\n");

    write(fd, (void *)(_Wa_Memory_addr+ptr), size);

    // 返回一个值供 WASM 检查
    return 0;
}

void _Wa_Import_env_print_i64(int64_t x) {
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
