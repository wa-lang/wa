# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# hello.wa.s - GAS (Intel Syntax) for Windows x64

.intel_syntax noprefix

.extern _write
.globl main

.section .data
msg: 
    .ascii "hello, world\n"
    msg_len = . - msg

.section .text
main:
    push rbp
    mov  rbp, rsp

    # _write(STDOUT, msg, count)
    mov  rcx, 1
    lea  rdx, [rip + msg]
    mov  r8d, msg_len

    sub  rsp, 32
    call _write
    add  rsp, 32

    # return 0
    mov  eax, 0
    leave
    ret

