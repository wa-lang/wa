// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

/*
#include <stdint.h>

static int32_t wat2xxI32Add_windows_amd64(int32_t a, int32_t b) {
    int32_t result;

    __asm__ volatile (
        ".intel_syntax noprefix;"  // 切换到 Intel 语法, 不带前缀

        // 1. 建立模拟栈环境 (rbp 已由 C 调用约定保护或建立)
        "push rbp;"
        "mov rbp, rsp;"
        "sub rsp, 16;"

        // 2. 存入输入数据
        "mov dword ptr [rbp-4], %1;"
        "mov dword ptr [rbp-8], %2;"

        // 3. 执行 i32.add 逻辑
        "mov eax, dword ptr [rbp-4];"
        "add eax, dword ptr [rbp-8];"
        "mov %0, eax;"             // 将结果存入输出变量

        // 4. 清理栈环境
        "add rsp, 16;"
        "pop rbp;"

        ".att_syntax;"             // 重要: 切回 AT&T 语法, 避免污染后续编译
        : "=r" (result)            // %0
        : "r" (a), "r" (b)         // %1, %2
        : "rax", "memory", "cc"
    );

    return result;
}
*/
import "C"

func wat2xxI32Add(a, b int32) int32 {
	return int32(C.wat2xxI32Add_windows_amd64(C.int32_t(a), C.int32_t(b)))
}
