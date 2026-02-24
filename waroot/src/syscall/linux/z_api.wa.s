# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.App.argc
.extern .Wa.App.argv
.extern .Wa.App.envp

.extern .Wa.Memory.addr

.extern .Wa.Runtime.memcpy


# func GetArgc => int
.section .text
.globl .Wa.Import.syscall_linux.GetArgc
.Wa.Import.syscall_linux.GetArgc:
    mov rax, [rip+.Wa.App.argc]
    ret


# func GetArgvLen(idx: int) => int
.section .text
.globl .Wa.Import.syscall_linux.GetArgvLen
.Wa.Import.syscall_linux.GetArgvLen:
    # 取第 i 个参数的地址
    mov rax, [rip+.Wa.App.argv]
    shl rdi, 3     # rdi = rdi*8
    add rax, rdi   # rax = .Wa.App.argv + rdi
    mov rax, [rax] # rax = argv[idx]

    # 保存开始地址
    mov rdx, rax   # start = rax

.Wa.L.syscall_linux.strlen_loop:
    mov r8b, byte ptr [rax]
    cmp r8b, 0
    je  .Wa.L.syscall_linux.strlen_done
    inc rax
    jmp .Wa.L.syscall_linux.strlen_loop

.Wa.L.syscall_linux.strlen_done:
    sub rax, rdx
    ret


# func GetArgvData(dst: uintptr, srcIdx, n: int)
.section .text
.globl .Wa.Import.syscall_linux.GetArgvData
.Wa.Import.syscall_linux.GetArgvData:
    # dst = dst + .Wa.Memory.addr
    mov rax, [rip+.Wa.Memory.addr]
    add rdi, rax

    # src = argv[idx]
    mov rax, [rip+.Wa.App.argv]
    shl rsi, 3     # rsi = rsi*8
    add rax, rsi   # rax = .Wa.App.argv + rsi
    mov rsi, [rax] # rsi = argv[idx]

    sub  rsp, 32
    call .Wa.Runtime.memcpy
    add  rsp, 32
    ret


