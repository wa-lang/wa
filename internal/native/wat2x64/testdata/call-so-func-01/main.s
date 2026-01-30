# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.section .data
lib_path: .asciz "./libhello.so"

.section .text
.globl _start

_start:
    # 1. sys_open: 打开动态库文件
    mov eax, 2              # sys_open
    lea rdi, [rip + lib_path]
    xor rsi, rsi            # O_RDONLY
    syscall
    mov r12, rax            # 保存文件描述符 fd 到 r12

    # 2. sys_mmap: 将整个 .so 映射到内存
    # void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset)
    mov eax, 9              # sys_mmap
    xor rdi, rdi            # addr = NULL
    mov rsi, 8192           # length = 8KB (假设库很小)
    mov rdx, 5              # prot = PROT_READ | PROT_EXEC (1 | 4)
    mov r10, 2              # flags = MAP_PRIVATE
    mov r8, r12             # fd
    xor r9, r9              # offset = 0
    syscall
    mov r13, rax            # r13 现在是动态库在内存的基址

    # 3. 调用动态库函数
    # 假设通过预先分析得知 'add' 函数相对于基址的偏移是 0x10f9
    mov rdi, 10             # 参数 1: a = 10
    mov rsi, 20             # 参数 2: b = 20
    
    mov rbx, 0x10f9         # 这里的偏移量需要你根据实际编译的 .so 确定
    add rbx, r13            # 绝对地址 = 基址 + 偏移
    call rbx                # 调用 add(10, 20)

    # 此时 rax = 30
    
    # 4. sys_exit: 退出程序，返回码设为 rax 的结果
    mov rdi, rax            # 将 30 作为退出码以便验证 (echo $?)
    mov eax, 60             # sys_exit
    syscall
