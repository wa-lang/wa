// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

/*
#include <stdint.h>

static int32_t wat2xxI32Add_linux_amd64(int32_t a, int32_t b) {
    int32_t result;

    __asm__ volatile (
        ".intel_syntax noprefix;"

        "mov eax, %1;"    // 对应 a
        "add eax, %2;"    // 对应 b (执行 i32.add)

        // 2. 将结果写回输出
        "mov %0, eax;"

        ".att_syntax;"
        : "=r" (result)          // %0: 输出
        : "r" (a), "r" (b)       // %1, %2: 输入
        : "rax", "cc"            // 告知 GCC eax 和 标志位被修改了
    );

    return result;
}
*/
import "C"

func wat2xxI32Add(a, b int32) int32 {
	return int32(C.wat2xxI32Add_linux_amd64(C.int32_t(a), C.int32_t(b)))
}
