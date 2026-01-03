; Copyright (C) 2026 武汉凹语言科技有限公司
; SPDX-License-Identifier: AGPL-3.0-or-later

; hello.wa.s

TRUE equ 1
FALSE equ 0
NULL equ 0

STD_INPUT_HANDLE equ -10
STD_OUTPUT_HANDLE equ -11
STD_ERROR_HANDLE equ -12

extern WriteFile
extern GetStdHandle

section .data
    msg db "Hello, World!!",10,0
    msgLen EQU $-msg-1

section .bss
    hFile resq 1
    lpNumberOfByteWritten resq 1

default rel
section .text
    global main

main:
    push rbp
    mov  rbp, rsp

    ; hFile = GetStdHandle(STD_OUTPUT_HANDLE)
    mov rcx, STD_OUTPUT_HANDLE
    sub rsp, 32
    call GetStdHandle
    add rsp, 32
    mov qword[hFile], rax

    ; WriteFile(hFile, msg, msgLen, lpNumberOfByteWritten, NULL)
    mov rcx, qword[hFile]
    lea rdx, [msg]
    mov r8, msgLen
    lea r9, [lpNumberOfByteWritten]
    push NULL
    sub rsp, 32
    call WriteFile

    mov rax, 0
    leave
    ret

