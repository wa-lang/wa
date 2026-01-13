# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# 导入 系统调用
.extern _write
.extern _exit
.set .Import.syscall.write, _write
.set .Import.syscall.exit,  _exit

# 导入 runtime 函数
.extern malloc
.extern memcpy
.extern memset
.set .Import.runtime.malloc, malloc
.set .Import.runtime.memcpy, memcpy
.set .Import.runtime.memset, memset

# 导入 env 函数
.extern env_write
.extern env_print_i64
.set .Import.env_write, env_write
.set .Import.env_print_i64, env_print_i64

# 定义内存
.section .data
.align 8
.Memory.addr: .quad 0
.Memory.pages: .quad 1
.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 8
# Memory[8]: hello worl...
.Memory.dataOffset.0: .quad 8
.Memory.dataSize.0: .quad 12
.Memory.dataPtr.0: .ascii "hello world\n"

# 内存初始化函数
.section .text
.globl .Memory.initFunc
.Memory.initFunc:
    # 影子空间
    sub rsp, 40

    # 分配内存
    mov  rcx, [rip + .Memory.maxPages]
    shl  rcx, 16
    call malloc
    lea  rdx, [rip + .Memory.addr]
    mov  [rdx], rax

    # 内存清零
    mov  rcx, [rip + .Memory.addr]
    mov  rdx, 0
    mov  r8, [rip + .Memory.maxPages]
    shl  r8, 16
    call .Import.runtime.memset

    # 初始化内存

    # Memory[8]: hello worl...
    mov  rax, [rip + .Memory.addr]
    mov  rcx, [rip + .Memory.dataOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Memory.dataPtr.0]
    mov  r8, [rip + .Memory.dataSize.0]
    call .Import.runtime.memcpy

    # 函数返回
    add rsp, 40
    ret

# Wasm 入口函数, 后续是编译器自动生成
.section .text
.global .F.main
.F.main:
    # 栈对齐并分配空间:
    # 32 (影子空间) + 16 (参数5和6) + 8 (对齐) = 56 或更高
    sub rsp, 56

    # syscall.write(stdout, &memory[ptr], size)
    mov  rcx, 1 # stdout
    mov  rax, [rip + .Memory.addr]
    mov  rdx, [rip + .Memory.dataOffset.0]
    add  rdx, rax # rdx = &memory[ptr]
    mov  r8, [rip + .Memory.dataSize.0] # size
    call .Import.syscall.write

    # env_write(stdout, &memory[ptr], size, 100, 200, 300)
    mov  rcx, 1 # stdout
    mov  rax, [rip + .Memory.addr]
    mov  rdx, [rip + .Memory.dataOffset.0]
    add  rdx, rax # rdx = &memory[ptr]
    mov  r8, [rip + .Memory.dataSize.0] # size
    mov r9, 100 # p4
    mov rax, 200 # 处理超出影子空间的参数
    mov [rsp + 32], rax  # 第5个参数紧跟在影子空间后面
    mov rax, 300
    mov [rsp + 40], rax  # 第6个参数
    call .Import.env_write

    # 函数返回
    add rsp, 56
    ret

# 汇编程序入口函数
.section .text
.global main
main:
    # rsp%16 == 0
    # (rsp-8-40)%16 == 0
    sub rsp, 40

    call .Memory.initFunc
    call .F.main

    # syscall.exit(0)
    mov  rcx, 0 # exit code
    call .Import.syscall.exit

    # exit 后这里不会被执行, 但是依然保留
    add rsp, 40
    ret
