// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)                asm(".Wa.Runtime.write");

extern void _Wa_Import_syscall_linux_print_str (const char* s, int32_t len) asm(".Wa.Import.syscall_linux.print_str");
extern void _Wa_Import_syscall_linux_print_rune(int32_t c)                  asm(".Wa.Import.syscall_linux.print_rune");
extern void _Wa_Import_syscall_linux_print_i64(int64_t val)                 asm(".Wa.Import.syscall_linux.print_i64");

extern int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) asm(".Wa.Import.env.write");
extern void    _Wa_Import_env_print_i64(int64_t x)                                                                                   asm(".Wa.Import.env.print_i64");

// 打印字符串
static void print_str(const char* s, int32_t len) {
    int i;
    for(i = 0; i < len; i++) {
        _Wa_Import_syscall_linux_print_rune(*s++);
    }
}

// 按照 Win64 调用约定，前四个参数在寄存器，后两个在栈上
int64_t _Wa_Import_env_write(int64_t fd, char* ptr, int64_t size, int64_t p4, int64_t p5, int64_t p6, int64_t p7, int64_t p8) {
    print_str("--- Linux ABI Verification ---\n", sizeof("--- Linux ABI Verification ---\n")-1);

    print_str("Param 1 (RDI): ", sizeof("Param 1 (RDI): ")-1);
    _Wa_Import_syscall_linux_print_i64(fd);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 2 (RSI): ", sizeof("Param 2 (RSI): ")-1);
    _Wa_Import_syscall_linux_print_i64((int64_t)(ptr));
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 3 (RDX): ", sizeof("Param 3 (RDX): ")-1);
    _Wa_Import_syscall_linux_print_i64(size);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 4 (RCX): ", sizeof("Param 4 (RCX): ")-1);
    _Wa_Import_syscall_linux_print_i64(p4);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 5 (R8): ", sizeof("Param 5 (R8): ")-1);
    _Wa_Import_syscall_linux_print_i64(p5);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 6 (R9): ", sizeof("Param 6 (R9): ")-1);
    _Wa_Import_syscall_linux_print_i64(p6);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 7 (Stack RSP+0): ", sizeof("Param 7 (Stack RSP+0): ")-1);
    _Wa_Import_syscall_linux_print_i64(p7);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("Param 8 (Stack RSP+8): ", sizeof("Param 8 (Stack RSP+8): ")-1);
    _Wa_Import_syscall_linux_print_i64(p8);
    _Wa_Import_syscall_linux_print_rune('\n');

    print_str("-------------------------------\n", sizeof("-------------------------------\n")-1);

    _Wa_Runtime_write(fd, (void *)(ptr), size);

    // 返回一个值供 WASM 检查
    return 0;
}

void _Wa_Import_env_print_i64(int64_t x) {
    print_str("printI64: ", sizeof("printI64: "));
    _Wa_Import_syscall_linux_print_i64(x);
    _Wa_Import_syscall_linux_print_rune('\n');
}

