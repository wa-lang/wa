# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern write
.extern exit

.extern malloc
.extern memcpy
.extern memmove
.extern memset

.extern .Wa.Memory.addr

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.global _Wa_Runtime_write
_Wa_Runtime_write:
    jmp write

# void _Wa_Runtime_exit(int status)
.section .text
.global _Wa_Runtime_exit
_Wa_Runtime_exit:
    jmp exit

# void* _Wa_Runtime_malloc(int size)
.section .text
.global _Wa_Runtime_malloc
_Wa_Runtime_malloc:
    jmp malloc

# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.global _Wa_Runtime_memcpy
_Wa_Runtime_memcpy:
    jmp memcpy

# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.global _Wa_Runtime_memmove
_Wa_Runtime_memmove:
    jmp memmove

# void* _Wa_Runtime_memset(void* s, int c, int n)
.global _Wa_Runtime_memset
_Wa_Runtime_memset:
    jmp memset

# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
.global _Wa_Import_syscall_linux_print_str
_Wa_Import_syscall_linux_print_str:
    # rax = base + ptr
    mov rax, [rip + .Wa.Memory.addr]
    add rax, rdi

    mov rdx, rsi # arg.2: len
    mov rsi, rax # arg.1: base + ptr
    mov rdi, 1   # arg.0: stdout
    jmp write

# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
.global _Wa_Import_syscall_linux_proc_exit
_Wa_Import_syscall_linux_proc_exit:
    jmp exit

