// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#if defined(_WIN64)
#   include <io.h>
#else
#   include <unistd.h>
#endif

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)        asm(".Wa.Runtime.write");
extern void  _Wa_Runtime_exit(int status)                           asm(".Wa.Runtime.exit");
extern void* _Wa_Runtime_malloc(int size)                           asm(".Wa.Runtime.malloc");
extern void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)  asm(".Wa.Runtime.memcpy");
extern void* _Wa_Runtime_memmove(void* dst, const void* src, int n) asm(".Wa.Runtime.memmove");
extern void* _Wa_Runtime_memset(void* s, int c, int n)              asm(".Wa.Runtime.memset");

extern int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) asm(".Wa.Import.env.write");
extern void    _Wa_Import_env_print_i64(int64_t x)                                                                                   asm(".Wa.Import.env.print_i64");

// 按照 Win64 调用约定，前四个参数在寄存器，后两个在栈上
int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) {
#if defined(_WIN64)
    printf("--- Win64 ABI Verification ---\n");
    printf("Param 1 (RCX): %lld\n", fd);
    printf("Param 2 (RDX): %lld\n", ptr);
    printf("Param 3 (R8):  %lld\n", size);
    printf("Param 4 (R9):  %lld\n", p4);
    printf("Param 5 (Stack RSP+32): %lld\n", p5);
    printf("Param 6 (Stack RSP+40): %lld\n", p6);
    printf("Param 7 (Stack RSP+48): %lld\n", p7);
    printf("Param 8 (Stack RSP+56): %lld\n", p8);
    printf("-------------------------------\n");

    _write(fd, (void *)(_Wa_Memory_addr+ptr), size);
#else
    printf("--- Linux ABI Verification ---\n");
    printf("Param 1 (RDI): %ld\n", fd);
    printf("Param 2 (RSI): %ld\n", (int64_t)(ptr));
    printf("Param 3 (RDX):  %ld\n", size);
    printf("Param 4 (RCX):  %ld\n", p4);
    printf("Param 5 (R8):   %ld\n", p5);
    printf("Param 6 (R9):   %ld\n", p6);
    printf("Param 7 (Stack RSP+0): %ld\n", p7);
    printf("Param 8 (Stack RSP+8): %ld\n", p8);
    printf("-------------------------------\n");

    write(fd, (void *)(_Wa_Memory_addr+ptr), size);
#endif

    // 返回一个值供 WASM 检查
    return 0;
}

void _Wa_Import_env_print_i64(int64_t x) {
    printf("printI64: %d\n", (int)(x));
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
