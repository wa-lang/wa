# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.section .data
    msg:
        .ascii "Hello from Kernel32!\12"
    msg_len = . - msg

.section .bss
    .lcomm written, 8           # 预留 8 字节存储实际写入的字符数

.section .text
.globl _start

_start:
    # 1. 栈对齐与影子空间 (Shadow Space)
    # Windows 要求调用函数前栈必须 16 字节对齐，且预留 32 字节空间
    sub rsp, 40                 # 32 (影子空间) + 8 (为了对齐 RSP)

    # 2. 获取标准输出句柄: GetStdHandle(-11)
    mov rcx, -11                # 参数1: nStdHandle = STD_OUTPUT_HANDLE
    mov rax, [rip + __imp_GetStdHandle] # 从导入表获取函数地址
    call rax
    mov r12, rax                # 将返回的句柄保存在 r12

    # 3. 打印字符串: WriteConsoleA(handle, buf, len, &written, 0)
    mov rcx, r12                # 参数1: hConsoleOutput
    lea rdx, [rip + msg]        # 参数2: *lpBuffer
    mov r8, msg_len             # 参数3: nNumberOfCharsToWrite
    lea r9, [rip + written]     # 参数4: *lpNumberOfCharsWritten
    
    # WriteConsoleA 有第 5 个参数，必须放在栈上 (影子空间上方)
    mov qword ptr [rsp + 32], 0 # 参数5: lpReserved = NULL
    
    mov rax, [rip + __imp_WriteConsoleA]
    call rax

    # 4. 退出程序: ExitProcess(0)
    xor rcx, rcx                # 参数1: uExitCode = 0
    mov rax, [rip + __imp_ExitProcess]
    call rax

    # 正常情况下不会执行到这里
    add rsp, 40
    ret

