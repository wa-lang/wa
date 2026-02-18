# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.Memory.addr

.extern .Wa.Import.syscall_windows.print_rune
.extern .Wa.Import.syscall_windows.print_i64

# void _Wa_Import_syscall_windows_print_i64(int64_t val)
.section .text
.globl .Wa.Import.env.print_i64
.Wa.Import.env.print_i64:
    call .Wa.Import.syscall_windows.print_i64

    mov  rcx, 10 # '\n'
    call .Wa.Import.syscall_windows.print_rune

    ret
