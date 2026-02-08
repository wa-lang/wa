// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)        asm(".Wa.Runtime.write");
extern void  _Wa_Runtime_exit(int status)                           asm(".Wa.Runtime.exit");
extern void* _Wa_Runtime_malloc(int size)                           asm(".Wa.Runtime.malloc");
extern void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)  asm(".Wa.Runtime.memcpy");
extern void* _Wa_Runtime_memmove(void* dst, const void* src, int n) asm(".Wa.Runtime.memmove");
extern void* _Wa_Runtime_memset(void* s, int c, int n)              asm(".Wa.Runtime.memset");

extern void _Wa_Import_syscall_linux_print_str (const char* s, int32_t len) asm(".Wa.Import.syscall_linux.print_str");
extern void _Wa_Import_syscall_linux_print_rune(int32_t c)                  asm(".Wa.Import.syscall_linux.print_rune");
extern void _Wa_Import_syscall_linux_print_i64(int64_t val)                 asm(".Wa.Import.syscall_linux.print_i64");

extern void _Wa_Import_env_print_i64(int64_t x) asm(".Wa.Import.env.print_i64");

// 打印字符串
static void print_str(const char* s, int32_t len) {
    int i;
    for(i = 0; i < len; i++) {
        _Wa_Import_syscall_linux_print_rune(*s++);
    }
}

void _Wa_Import_env_print_i64(int64_t x) {
    print_str("printI64: ", sizeof("printI64: "));
    _Wa_Import_syscall_linux_print_i64(x);
    _Wa_Import_syscall_linux_print_rune('\n');
}
