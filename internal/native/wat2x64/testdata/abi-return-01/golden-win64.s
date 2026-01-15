# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# 运行时函数
.extern _write
.extern _exit
.extern malloc
.extern memcpy
.extern memset
.set .Runtime._write, _write
.set .Runtime._exit, _exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset

# 导入函数(外部库定义)
.extern wat2x64_env_get_multi_values
.extern wat2x64_env_print_i64
.set .Import.env.get_multi_values, wat2x64_env_get_multi_values
.set .Import.env.print_i64, wat2x64_env_print_i64

# 定义内存
.section .data
.align 8
.globl .Memory.addr
.globl .Memory.pages
.globl .Memory.maxPages
.Memory.addr: .quad 0
.Memory.pages: .quad 1
.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 8
# memcpy(&Memory[8], data[0], size)
.Memory.dataOffset.0: .quad 8
.Memory.dataSize.0: .quad 12
.Memory.dataPtr.0: .ascii "hello world\n"

# 内存初始化函数
.section .text
.globl .Memory.initFunc
.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rcx, [rip + .Memory.maxPages]
    shl  rcx, 16
    call .Runtime.malloc
    lea  rdx, [rip + .Memory.addr]
    mov  [rdx], rax

    # 内存清零
    mov  rcx, [rip + .Memory.addr]
    mov  rdx, 0
    mov  r8, [rip + .Memory.maxPages]
    shl  r8, 16
    call .Runtime.memset

    # 初始化内存

    # memcpy(&Memory[0], data[0], size)
    mov  rax, [rip + .Memory.addr]
    mov  rcx, [rip + .Memory.dataOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Memory.dataPtr.0]
    mov  r8, [rip + .Memory.dataSize.0]
    call .Runtime.memcpy

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 汇编程序入口函数
.section .text
.globl main
main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    call .Memory.initFunc
    call .F.main

    # runtime.exit(0)
    mov  rcx, 0
    call .Runtime._exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

.section .data
.align 8
.Runtime.panic.message: .asciz "panic"
.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Runtime.panic
.Runtime.panic:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stderr, panicMessage, size)
    mov  rcx, 2 # stderr
    mov  rdx, [rip + .Runtime.panic.message]
    mov  r8, [rip + .Runtime.panic.messageLen] # size
    call .Runtime.panic

    # 退出程序
    mov  rcx, 1 # 退出码
    call .Runtime._exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# Wasm 入口函数, 后续是编译器自动生成
.section .text
.global .F.main
.F.main:
    # 栈对齐并分配空间:
    # 影子空间(32字节) + 结构体空间(24字节) + 16字节对齐保护 = 64 字节
    sub rsp, 64

    # env_get_multi_values()
    # 根据 Win64 ABI，第一个参数 RCX 必须指向存储返回值的内存地址
    # 我们使用刚分配的栈空间：[rsp + 32] (跳过影子空间)
    # 原本 WAT 里的第一个参数 (param i64) 顺延到 RDX
    lea  rcx, [rsp + 32]
    mov  rdx, 100
    call .Import.env.get_multi_values
    mov  rcx, [rsp + 48] # v3
    call .Import.env.print_i64
    mov  rcx, [rsp + 40] # v2
    call .Import.env.print_i64
    mov  rcx, [rsp + 32] # v1
    call .Import.env.print_i64

    # 函数返回
    add rsp, 64
    ret
