// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern void _Wa_Import_syscall_linux_print_str (const char* s, int32_t len) asm(".Wa.Import.syscall_linux.print_str");
extern void _Wa_Import_syscall_linux_print_rune(int32_t c)                  asm(".Wa.Import.syscall_linux.print_rune");
extern void _Wa_Import_syscall_linux_print_i64(int64_t val)                 asm(".Wa.Import.syscall_linux.print_i64");

typedef struct EnvMultiRet EnvMultiRet;

extern EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val) asm(".Wa.Import.env.get_multi_values");
extern void        _Wa_Import_env_print_i64(int64_t x)                asm(".Wa.Import.env.print_i64");

// 定义对应 Wasm 多返回值的结构体
struct EnvMultiRet {
    int64_t v1;
    int64_t v2;
    int64_t v3;
};

// 打印字符串
static void print_str(const char* s, int32_t len) {
    int i;
    for(i = 0; i < len; i++) {
        _Wa_Import_syscall_linux_print_rune(*s++);
    }
}

// C 函数: 返回结构体
// 注意: 根据 Win64 ABI，这会被编译为:
// void _Wa_Import_env_get_multi_values(EnvMultiRet* hidden_ptr, int64_t input_val)
EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val) {
    EnvMultiRet r;
    r.v1 = input_val + 1;
    r.v2 = input_val + 2;
    r.v3 = input_val + 3;

    print_str("C side: input=", sizeof("C side: input=")-1);
    _Wa_Import_syscall_linux_print_i64(input_val);
    print_str(", returning [", sizeof(", returning [")-1);
    _Wa_Import_syscall_linux_print_i64(r.v1);
    _Wa_Import_syscall_linux_print_rune(',');
    _Wa_Import_syscall_linux_print_rune(' ');
    _Wa_Import_syscall_linux_print_i64(r.v2);
    _Wa_Import_syscall_linux_print_rune(',');
    _Wa_Import_syscall_linux_print_rune(' ');
    _Wa_Import_syscall_linux_print_i64(r.v3);
    _Wa_Import_syscall_linux_print_rune(']');
    _Wa_Import_syscall_linux_print_rune('\n');

    return r;
}

void _Wa_Import_env_print_i64(int64_t x) {
    print_str("printI64: ", sizeof("printI64: "));
    _Wa_Import_syscall_linux_print_i64(x);
    _Wa_Import_syscall_linux_print_rune('\n');
}

