// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdio.h>
#include <stdint.h>

// 定义对应 Wasm 多返回值的结构体
typedef struct {
    int64_t v1;
    int64_t v2;
    int64_t v3;
} EnvMultiRet;

// C 函数: 返回结构体
// 注意: 根据 Win64 ABI，这会被编译为:
// void env_get_multi_values(EnvMultiRet* hidden_ptr, int64_t input_val)
EnvMultiRet env_get_multi_values(int64_t input_val) {
    EnvMultiRet r;
    r.v1 = input_val + 1;
    r.v2 = input_val + 2;
    r.v3 = input_val + 3;
    printf(
        "C side: input=%llu, returning [%llu, %llu, %llu]\n", 
        input_val, r.v1, r.v2, r.v3
    );
    return r;
}

void env_print_i64(int64_t x) {
    printf("printI64: %lld\n", x);
}
