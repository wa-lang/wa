# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern ExitProcess
.extern MessageBoxA

.section .data
.align 8
msg_title: .ascii "MessageBox\0"
msg_body:  .ascii "Hello from Assembly!\0"

.section .text
.align 8
.global _start
_start:
    sub rsp, 40

    call hello_msg

    # return 0
    xor ecx, ecx
    call ExitProcess


.section .text
.align 8
.global hello_msg
hello_msg:
    sub rsp, 40

    # int MessageBoxA(
    #   [in, optional] HWND   hWnd,
    #   [in, optional] LPCSTR lpText,
    #   [in, optional] LPCSTR lpCaption,
    #   [in]           UINT   uType
    # );

    xor rcx, rcx            # arg.1: hWnd = 0
    lea rdx, [rip+msg_body] # arg.2: lpText
    lea r8, [rip+msg_title] # arg.3: lpCaption
    mov r9d, 0x40           # arg.4: uType (MB_OK | MB_ICONINFORMATION)

    # GBK 编码
    # 中文要改为 UTF16编码, 用 MessageBoxW 输出
    call MessageBoxA

    add rsp, 40
    ret

