// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build inline_asm

package main

/*
#include <stdint.h>

static int32_t wat2xxI32Add_linux_loong64(int32_t a, int32_t b) {
    int32_t result;

    __asm__ volatile (
        // 1. 将输入 a 加载到临时寄存器 $t0 (通过 %1 映射)
        // 2. 将输入 b 加载到临时寄存器 $t1 (通过 %2 映射)
        // 3. 执行 32 位加法运算：$t2 = $t0 + $t1
        "add.w %0, %1, %2;"

        : "=r" (result)          // %0: 输出寄存器
        : "r" (a), "r" (b)       // %1, %2: 输入寄存器
        : "memory"               // 龙芯上 add.w 不修改标志位, 只需声明内存屏障
    );

    return result;
}
*/
import "C"

func I32Add(a, b int32) int32 {
	return int32(C.wat2xxI32Add_linux_loong64(C.int32_t(a), C.int32_t(b)))
}
