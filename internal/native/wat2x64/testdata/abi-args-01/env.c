// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#if defined(_WIN64)

#include <stdio.h>
#include <stdint.h>
#include <io.h>

extern int64_t wat2x64_Memory_addr __asm__(".Memory.addr");

// 按照 Win64 调用约定，前四个参数在寄存器，后两个在栈上
int64_t wat2x64_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) {
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

    _write(fd, (void *)(wat2x64_Memory_addr+ptr), size);

    // 返回一个值供 WASM 检查
    return 0;
}

void wat2x64_env_print_i64(int64_t x) {
    printf("printI64: %lld\n", x);
}

#elif defined(__linux__)

#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

extern int64_t wat2x64_Memory_addr __asm__(".Memory.addr");

// 按照 Win64 调用约定，前四个参数在寄存器，后两个在栈上
int64_t wat2x64_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) {
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

    write(fd, (void *)(wat2x64_Memory_addr+ptr), size);

    // 返回一个值供 WASM 检查
    return 0;
}

void wat2x64_env_print_i64(int64_t x) {
    printf("printI64: %ld\n", x);
}

#else

#error "Only Support Windows/Linux"

#endif
