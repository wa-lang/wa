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

    # 保存寄存器
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
    
    # 恢复寄存器
    mov r12, [rbp-8]
    mov r13, [rbp-16]
    mov r14, [rbp-24]

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret


# void _Wa_Runtime_exit(int status)
.section .text
.globl .Wa.Runtime.exit
.Wa.Runtime.exit:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # void ExitProcess(
    #   [in] UINT uExitCode
    # );

    call ExitProcess

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret


# void* _Wa_Runtime_malloc(int size)
.section .text
.globl .Wa.Runtime.malloc
.Wa.Runtime.malloc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # LPVOID VirtualAlloc(
    #   [in, optional] LPVOID lpAddress,
    #   [in]           SIZE_T dwSize,
    #   [in]           DWORD  flAllocationType,
    #   [in]           DWORD  flProtect
    # );

    mov  rdx, rcx      # dwSize
    xor  rcx, rcx      # lpAddress = NULL
    mov  r8,  0x3000   # MEM_COMMIT | MEM_RESERVE
    mov  r9,  0x04     # PAGE_READWRITE
    call VirtualAlloc # rax = allocated memory

    mov rsp, rbp
    pop rbp
    ret


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
    mov  rax, rcx
    test r8, r8
    jz   .Wa.L.memcpy.done

.Wa.L.memcpy.loop:
    mov r9b, byte ptr [rdx]
    mov byte ptr [rcx], r9b
    inc rcx
    inc rdx
    dec r8
    jnz .Wa.L.memcpy.loop

.Wa.L.memcpy.done:
    ret


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memmove
.Wa.Runtime.memmove:
    mov  rax, rcx # 备份 dst 用于返回
    test r8, r8   # n == 0 ?
    jz   .Wa.L.memmove.done

    cmp rcx, rdx
    je  .Wa.L.memmove.done
    jb  .Wa.L.memmove.forward  # dst < src → 前向拷贝

    # =========================
    # 后向拷贝 (dst > src)
    # =========================

    push rdi
    push rsi

    mov rdi, rcx
    mov rsi, rdx

    add rdi, r8
    dec rdi
    add rsi, r8
    dec rsi

    mov rcx, r8 # 计数器

.Wa.L.memmove.back_loop:
    mov r9b, byte ptr [rsi]
    mov byte ptr [rdi], r9b
    dec rdi
    dec rsi
    dec rcx
    jnz .Wa.L.memmove.back_loop

    pop rsi
    pop rdi
    jmp .Wa.L.memmove.done

.Wa.L.memmove.forward:

    # =========================
    # 前向拷贝 (dst < src)
    # =========================

    push rdi
    push rsi

    mov rdi, rcx
    mov rsi, rdx
    mov rcx, r8

.Wa.L.memmove.fwd_loop:
    mov r9b, byte ptr [rsi]
    mov byte ptr [rdi], r9b
    inc rdi
    inc rsi
    dec rcx
    jnz .Wa.L.memmove.fwd_loop

    pop rsi
    pop rdi

.Wa.L.memmove.done:
    ret


# void* _Wa_Runtime_memset(void* s, int c, int n)
.section .text
.globl .Wa.Runtime.memset
.Wa.Runtime.memset:
    mov  rax, rcx # 返回 s
    test r8, r8
    jz   .Wa.L.memset.done

.Wa.L.memset.loop:
    mov byte ptr [rcx], dl # c 的低 8 位
    inc rcx
    dec r8
    jnz .Wa.L.memset.loop

.Wa.L.memset.done:
    ret
