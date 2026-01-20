// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#if defined(_WIN64)

#include <stdint.h>
#include <stdio.h>
#include <io.h>

extern int64_t wat2x64_Memory_addr __asm__(".Memory.addr");

int64_t wat2x64_syscall_write(int64_t fd, int64_t ptr, int64_t size) {
    printf("wat2x64_syscall_write: %lld, %lld, %lld\n", fd, ptr, size);
    return _write(fd, (void *)(wat2x64_Memory_addr+ptr), size);
}

#elif defined(__linux__)

#include <stdint.h>
#include <stdio.h>
#include <unistd.h>

extern int64_t wat2x64_Memory_addr __asm__(".Memory.addr");

int64_t wat2x64_syscall_write(int64_t fd, int64_t ptr, int64_t size) {
    printf("wat2x64_syscall_write: %ld, %ld, %ld\n", fd, ptr, size);
    return write(fd, (void *)(wat2x64_Memory_addr+ptr), size);
}

#else

#error "Only Support Windows/Linux"

#endif
