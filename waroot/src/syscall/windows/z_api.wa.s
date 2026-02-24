# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.App.argc
.extern .Wa.App.argv
.extern .Wa.App.envp

.extern .Wa.Memory.addr

.extern MessageBoxA
.extern MessageBoxW

.section .text
.globl .Wa.Import.syscall_windows.GetArgc
.Wa.Import.syscall_windows.GetArgc:
    mov rax, [rip+.Wa.App.argc]
    ret


.section .text
.globl .Wa.Import.syscall_windows.GetArgv
.Wa.Import.syscall_windows.GetArgv:
    lea rax, [rip+.Wa.App.argv]
    ret


.section .text
.globl .Wa.Import.syscall_windows.GetEnvp
.Wa.Import.syscall_windows.GetEnvp:
    lea rax, [rip+.Wa.App.envp]
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

