# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# hello.wa.s

.intel_syntax noprefix

TRUE  = 1
FALSE = 0
NULL  = 0

STD_INPUT_HANDLE  = -10
STD_OUTPUT_HANDLE = -11
STD_ERROR_HANDLE  = -12

MB_OK = 0x0

# --- 外部函数声明 ---
.extern ExitProcess
.extern MessageBoxA

.section .data
msg:
    .asciz "Welcom to Windows World!" # .asciz 自动添加末尾的 \0
cap:
    .asciz "Windows 10 says:"

.section .text
.globl main

main:
    push rbp
    mov  rbp, rsp

    # --- MessageBoxA(NULL, msg, cap, MB_OK) ---
    # 参数 1: HWND (RCX)
    xor  rcx, rcx            # 0 也可以用 xor 优化
    
    # 参数 2: lpText (RDX)
    lea  rdx, [rip + msg]    # 必须显式使用 rip 相对寻址
    
    # 参数 3: lpCaption (R8)
    lea  r8, [rip + cap]
    
    # 参数 4: uType (R9)
    mov  r9d, MB_OK          # MB_OK 是 32 位立即数，使用 r9d 即可
    
    # 必须的影子空间 (Shadow Space)
    sub  rsp, 32
    call MessageBoxA
    add  rsp, 32

    # --- 退出程序 ---
    # 按照 ABI 习惯，通常建议调用 ExitProcess 而不是直接 ret
    # xor rcx, rcx           # 退出码 0
    # sub rsp, 32
    # call ExitProcess

    mov  rax, 0
    leave
    ret
