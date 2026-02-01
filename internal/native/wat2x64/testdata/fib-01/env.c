// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#if defined(_WIN64)
#   include <io.h>
#else
#   include <unistd.h>
#endif

extern int64_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

void _Wa_Import_env_print_i64(int64_t x) {
    printf("printI64: %d\n", (int)(x));
}

int _Wa_Import_syscall_write(int fd, int ptr, int size) {
    printf("_Wa_Import_syscall_write: %d, %d, %d\n", fd, ptr, size);
#if defined(_WIN64)
    return _write(fd, (void*)(_Wa_Memory_addr+ptr), size);
#else
    return write(fd, (void*)(_Wa_Memory_addr+ptr), size);
#endif
}

int _Wa_Runtime_write(int fd, void *buf, int count) {
#if defined(_WIN64)
    return _write(fd, buf, (size_t)(count));
#else
    return write(fd, buf, (size_t)(count));
#endif
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
