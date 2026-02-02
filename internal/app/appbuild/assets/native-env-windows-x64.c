// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

#if defined(_WIN64)
#   include <io.h>
#else
#   include <unistd.h>
#endif

extern uintptr_t _Wa_Memory_addr __asm__(".Wa.Memory.addr");

extern int   _Wa_Runtime_write(int fd, void *buf, int count)        asm(".Wa.Runtime.write");
extern void  _Wa_Runtime_exit(int status)                           asm(".Wa.Runtime.exit");
extern void* _Wa_Runtime_malloc(int size)                           asm(".Wa.Runtime.malloc");
extern void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)  asm(".Wa.Runtime.memcpy");
extern void* _Wa_Runtime_memmove(void* dst, const void* src, int n) asm(".Wa.Runtime.memmove");
extern void* _Wa_Runtime_memset(void* s, int c, int n)              asm(".Wa.Runtime.memset");

extern void _Wa_Import_syscall_js_print_bool(_Bool v)                   asm(".Wa.Import.syscall_js.print_bool");
extern void _Wa_Import_syscall_js_print_f32(float v)                    asm(".Wa.Import.syscall_js.print_f32");
extern void _Wa_Import_syscall_js_print_f64(double v)                   asm(".Wa.Import.syscall_js.print_f64");
extern void _Wa_Import_syscall_js_print_i32(int32_t v)                  asm(".Wa.Import.syscall_js.print_i32");
extern void _Wa_Import_syscall_js_print_i64(int64_t v)                  asm(".Wa.Import.syscall_js.print_i64");
extern void _Wa_Import_syscall_js_print_u32(uint32_t v)                 asm(".Wa.Import.syscall_js.print_u32");
extern void _Wa_Import_syscall_js_print_u64(uint64_t v)                 asm(".Wa.Import.syscall_js.print_u64");
extern void _Wa_Import_syscall_js_print_rune(int32_t c)                 asm(".Wa.Import.syscall_js.print_rune");
extern void _Wa_Import_syscall_js_print_ptr (uint32_t ptr)              asm(".Wa.Import.syscall_js.print_ptr");
extern void _Wa_Import_syscall_js_print_str (uint32_t ptr, int32_t len) asm(".Wa.Import.syscall_js.print_str");
extern void _Wa_Import_syscall_js_print_position(int32_t pos)           asm(".Wa.Import.syscall_js.print_position");
extern void _Wa_Import_syscall_js_proc_exit(int32_t code)               asm(".Wa.Import.syscall_js.proc_exit");

int _Wa_Runtime_write(int fd, void *buf, int count) {
#if defined(_WIN64)
    return _write(fd, buf, (unsigned int)(count));
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

void _Wa_Import_syscall_js_print_bool(_Bool v) {
    printf("%s", v? "true": "false");
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_f32(float v) {
    printf("%f", v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_f64(double v) {
    printf("%f", v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_i32(int32_t v) {
    printf("%d", v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_i64(int64_t v) {
    printf("%" PRId64, v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_u32(uint32_t v) {
    printf("%u", v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_u64(uint64_t v) {
    printf("%" PRIu64, v);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_rune(int32_t c) {
    printf("%c", c);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_ptr (uint32_t ptr) {
    const void* p = (const void*)(_Wa_Memory_addr+(uintptr_t)(ptr));
    printf("%p", p);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_str (uint32_t ptr, int32_t len) {
    const char* s = (const char*)(_Wa_Memory_addr+(uintptr_t)(ptr));
    printf("%.*s", len, s);
    fflush(stdout);
}

void _Wa_Import_syscall_js_print_position(int32_t pos) {
    printf("{pos:%d}", pos);
    fflush(stdout);
}

void _Wa_Import_syscall_js_proc_exit(int32_t code) {
    exit(code);
}
