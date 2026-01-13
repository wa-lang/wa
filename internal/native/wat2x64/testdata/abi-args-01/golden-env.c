// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdio.h>
#include <stdint.h>
#include <io.h>

// 按照 Win64 调用约定，前四个参数在寄存器，后两个在栈上
int64_t env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6) {
    printf("--- Win64 ABI Verification ---\n");
    printf("Param 1 (RCX): %lld\n", fd);
    printf("Param 2 (RDX): %lld\n", ptr);
    printf("Param 3 (R8):  %lld\n", size);
    printf("Param 4 (R9):  %lld\n", p4);
    printf("Param 5 (Stack RSP+32): %lld\n", p5);
    printf("Param 6 (Stack RSP+40): %lld\n", p6);
    printf("-------------------------------\n");

    _write(fd, ptr, size);

    // 返回一个值供 WASM 检查
    return 0;
}

void env_print_i64(int64_t x) {
    printf("printI64: %lld\n", x);
}
