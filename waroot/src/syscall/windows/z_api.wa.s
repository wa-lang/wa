# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.App.argc
.extern .Wa.App.argv
.extern .Wa.App.envp

.extern .Wa.Memory.addr

.extern .Wa.Runtime.memcpy

.extern MessageBoxA
.extern MessageBoxW

# BUG: windows 命令行参数需要通过函数获取
# LPWSTR GetCommandLineW(); // 返回程序名字字符串
#
# LPWSTR * CommandLineToArgvW(
#   [in]  LPCWSTR lpCmdLine,
#   [out] int     *pNumArgs
# );

# func GetArgc => int
.section .text
.globl .Wa.Import.syscall_windows.GetArgc
.Wa.Import.syscall_windows.GetArgc:
    mov rax, [rip+.Wa.App.argc]
    ret


# func GetArgvLen(idx: int) => int
.section .text
.globl .Wa.Import.syscall_windows.GetArgvLen
.Wa.Import.syscall_windows.GetArgvLen:
    # 取第 i 个参数的地址
    mov rax, [rip+.Wa.App.argv]
    shl rcx, 3     # rcx = rcx*8
    add rax, rcx   # rax = .Wa.App.argv + rcx
    mov rax, [rax] # rax = argv[idx]

    # 保存开始地址
    mov rdx, rax   # start = rax

.Wa.L.syscall_windows.strlen_loop:
    mov r8b, byte ptr [rax]
    cmp r8b, 0
    je  .Wa.L.syscall_windows.strlen_done
    inc rax
    jmp .Wa.L.syscall_windows.strlen_loop

.Wa.L.syscall_windows.strlen_done:
    sub rax, rdx
    ret


# func GetArgvData(dst: uintptr, srcIdx, n: int)
.section .text
.globl .Wa.Import.syscall_windows.GetArgvData
.Wa.Import.syscall_windows.GetArgvData:
    # dst = dst + .Wa.Memory.addr
    mov rax, [rip+.Wa.Memory.addr]
    add rcx, rax

    # src = argv[idx]
    mov rax, [rip+.Wa.App.argv]
    shl rdx, 3     # rdx = rcx*8
    add rax, rdx   # rax = .Wa.App.argv + rdx
    mov rdx, [rax] # rdx = argv[idx]

    sub  rsp, 32
    call .Wa.Runtime.memcpy
    add  rsp, 32
    ret


.section .text
.globl .Wa.Import.syscall_windows.MessageBoxA
.Wa.Import.syscall_windows.MessageBoxA:
    sub  rsp, 40
    add  rdx, [rip+.Wa.Memory.addr]
    add  r8, [rip+.Wa.Memory.addr]
    call MessageBoxA
    add  rsp, 40
    ret

.section .text
.globl .Wa.Import.syscall_windows.MessageBoxW
.Wa.Import.syscall_windows.MessageBoxW:
    sub  rsp, 40
    add  rdx, [rip+.Wa.Memory.addr]
    add  r8, [rip+.Wa.Memory.addr]
    call MessageBoxW
    add  rsp, 40
    ret

