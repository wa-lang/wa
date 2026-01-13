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
.extern env_get_multi_values
.extern env_print_i64
.set .Import.env_get_multi_values, env_get_multi_values
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
    # 影子空间(32字节) + 结构体空间(24字节) + 16字节对齐保护 = 64 字节
    sub rsp, 64

    # syscall.write(stdout, &memory[ptr], size)
    mov  rcx, 1 # stdout
    mov  rax, [rip + .Memory.addr]
    mov  rdx, [rip + .Memory.dataOffset.0]
    add  rdx, rax # rdx = &memory[ptr]
    mov  r8, [rip + .Memory.dataSize.0] # size
    call .Import.syscall.write

    # env_get_multi_values()
    # 根据 Win64 ABI，第一个参数 RCX 必须指向存储返回值的内存地址
    # 我们使用刚分配的栈空间：[rsp + 32] (跳过影子空间)
    # 原本 WAT 里的第一个参数 (param i64) 顺延到 RDX
    lea  rcx, [rsp + 32]
    mov  rdx, 100
    call env_get_multi_values
    mov  rcx, [rsp + 32] # v1
    call env_print_i64
    mov  rcx, [rsp + 40] # v2
    call env_print_i64
    mov  rcx, [rsp + 48] # v3
    call env_print_i64

    # 函数返回
    add rsp, 64
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
