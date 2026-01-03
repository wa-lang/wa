; Copyright (C) 2026 武汉凹语言科技有限公司
; SPDX-License-Identifier: AGPL-3.0-or-later

; hello.wa.s

%include "win32n.inc"

extern ExitProcess
extern MessageBoxA

section .data
    msg db "Welcom to Windows World!",0
    cap db "Windows 10 says:",0

default rel
section .text
    global main

main:
    push rbp
    mov  rbp, rsp

    mov  rcx, 0
    lea  rdx, [msg]
    lea  r8, [cap]
    mov  r9d, MB_OK
    sub  rsp, 32
    call MessageBoxA
    add  rsp, 32
    leave
    ret

