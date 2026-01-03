; Copyright (C) 2026 武汉凹语言科技有限公司
; SPDX-License-Identifier: AGPL-3.0-or-later

; hello.wa.s

extern printf

section .data
    msg db "Welcom to Windows World!",0
    fmt db "Windows 10 says: %s",10,0

section .text
    global main

main:
    push rbp
    mov  rbp, rsp

    mov  rcx, fmt
    mov  rdx, msg
    sub  rsp, 32
    call printf
    add  rsp, 32

    mov  rax, 0
    leave
    ret

