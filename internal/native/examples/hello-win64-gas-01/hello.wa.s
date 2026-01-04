# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# hello.wa.s - GAS (Intel Syntax) for Windows x64

.intel_syntax noprefix

.section .data

msg:
    .asciz "Welcom to Windows World!"
fmt:
    .asciz "Windows 10 says: %s\n"

.section .text

.globl main
.extern printf

main:
    push rbp
    mov  rbp, rsp

    # Windows x64 ABI 传参: RCX, RDX
    # 注意: GAS 环境下通常使用 lea 来获取符号地址以保证兼容性
    lea  rcx, [rip + fmt]     # 相对 RIP 寻址，比直接加载更稳妥
    lea  rdx, [rip + msg]
    
    sub  rsp, 32              # 关键: Windows 影子空间
    call printf
    add  rsp, 32

    mov  rax, 0
    leave
    ret
