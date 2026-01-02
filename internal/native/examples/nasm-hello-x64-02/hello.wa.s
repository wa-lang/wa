; Copyright (C) 2026 武汉凹语言科技有限公司
; SPDX-License-Identifier: AGPL-3.0-or-later

; hello.wa.s

section .data
    msg db "hello, world",0
    NL  db 0xA ; 换行

section .bss

section .text
    global main

main:
    mov rax, 1 ; 1 表示写入
    mov rdi, 1 ; 1 表示标准输出
    mov rsi, msg ; 要显示的字符串
    mov rdx, 12  ; 字符串长度, 不包含 0
    syscall      ; 显示字符串

    mov rax, 1  ; 1 表示写入
    mov rdi, 1  ; 1 表示标准输出
    mov rsi, NL ; 换行
    mov rdx, 1  ; 字符串长度
    syscall     ; 显示字符串

    mov rax, 60 ; 60 表示退出
    mov rdi, 0  ; 0 是退出状态码
    syscall     ; 退出

