# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# hello.wa.s

.intel_syntax noprefix

# --- 常量定义 ---
TRUE  = 1
FALSE = 0
NULL  = 0

STD_INPUT_HANDLE  = -10
STD_OUTPUT_HANDLE = -11
STD_ERROR_HANDLE  = -12

# --- 外部函数声明 ---
.extern WriteFile
.extern GetStdHandle

.section .data
msg:
    .ascii "Hello, World!!\n"      # GAS 中 \n 代表换行 (10)
    .byte 0                       # 末尾补 0
msgLen = . - msg - 1              # . 代表当前地址，相当于 NASM 的 $

.section .bss
.align 8
hFile:
    .quad 0                       # GAS 用 .quad 预留 8 字节空间
lpNumberOfByteWritten:
    .quad 0

.section .text
.globl main

main:
    push rbp
    mov  rbp, rsp

    # --- hFile = GetStdHandle(STD_OUTPUT_HANDLE) ---
    mov rcx, STD_OUTPUT_HANDLE
    sub rsp, 32                   # 开辟影子空间
    call GetStdHandle
    add rsp, 32
    # 注意: GAS Intel 模式下访问变量地址建议显式配合 RIP
    mov [rip + hFile], rax

    # --- WriteFile(hFile, msg, msgLen, lpNumberOfByteWritten, NULL) ---
    mov rcx, [rip + hFile]
    lea rdx, [rip + msg]
    mov r8, msgLen
    lea r9, [rip + lpNumberOfByteWritten]
    
    # 第 5 个参数走栈
    push NULL                     
    sub rsp, 32                   # 影子空间（必须在 push 参数之后，call 之前）
    call WriteFile
    # 这里的清理需要同时处理影子空间和 push 的那个 NULL 参数 (32 + 8)
    add rsp, 40

    mov rax, 0
    leave
    ret
