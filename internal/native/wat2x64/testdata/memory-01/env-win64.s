# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.Memory.addr

# kernel32.dll

.extern ExitProcess
.extern GetStdHandle
.extern VirtualAlloc
.extern WriteFile

# int _Wa_Import_syscall_write(int fd, int ptr, int size) {
.section .text
.globl .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    push rbp
    mov  r10, [rip+.Wa.Memory.addr]
    add  rdx, r10 # buf
    call .Wa.Runtime.write
    pop  rbp
    ret


# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    // 保存寄存器
    mov [rbp-8], r12
    mov [rbp-16], r13
    mov [rbp-24], r14

    # 保存参数
    mov r12, rcx # fd
    mov r13, rdx # buf
    mov r14, r8  # count

    # 获取标准输出句柄
    cmp rcx, 2
    je  .Wa.L.Runtime.write.stderr
    jmp .Wa.L.Runtime.write.stdout

.Wa.L.Runtime.write.stdout:
    mov ecx, -11 # STD_OUTPUT_HANDLE
    jmp .Wa.L.Runtime.write.gethandle

.Wa.L.Runtime.write.stderr:
    mov ecx, -12 # STD_ERROR_HANDLE
    jmp .Wa.L.Runtime.write.gethandle

.Wa.L.Runtime.write.gethandle:
    # rax = GetStdHandle(nStdHandle)
    sub  rsp, 32
    call GetStdHandle
    add  rsp, 32

    # BOOL WriteFile(
    #   [in]                HANDLE       hFile,
    #   [in]                LPCVOID      lpBuffer,
    #   [in]                DWORD        nNumberOfBytesToWrite,
    #   [out, optional]     LPDWORD      lpNumberOfBytesWritten,
    #   [in, out, optional] LPOVERLAPPED lpOverlapped
    # );

    sub  rsp, 48
    mov  rcx, rax               # arg.0: fd
    mov  rdx, r13               # arg.1: buf
    mov  r8, r14                # arg2: count
    mov  qword ptr [rsp+40], 0  # arg.3: (*lpNumberOfBytesWritten) = 0
    lea  r9, qword ptr [rsp+40] # arg.3: lpNumberOfBytesWritten
    mov  qword ptr [rsp+32], 0  # arg.4: lpOverlapped = NULL
    call WriteFile
    add  rsp, 48

    # return nWrite
    mov eax, dword ptr [rsp+40]
    
    // 恢复寄存器
    mov r12, [rbp-8]
    mov r13, [rbp-16]
    mov r14, [rbp-24]

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

