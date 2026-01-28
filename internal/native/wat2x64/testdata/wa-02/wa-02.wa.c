// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <unistd.h>

extern uintptr_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

int _Wa_Runtime_write(int fd, void *buf, int count) {
    return write(fd, buf, (size_t)(count));
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

void _Wa_Import_syscall_linux_print_rune(int32_t c) {
    printf("%c", c);
    fflush(stdout);
}

void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len) {
    const char* s = (const char*)(_Wa_Memory_addr+(uintptr_t)(ptr));
    printf("%.*s", len, s);
    fflush(stdout);
}

void _Wa_Import_syscall_linux_print_position(int32_t pos) {
    printf("{pos:%d}", pos);
    fflush(stdout);
}

void _Wa_Import_syscall_linux_proc_exit(int32_t code) {
    exit(code);
}
