# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# 运行时函数
.extern write
.extern exit
.extern malloc
.extern memcpy
.extern memset
.set .Runtime.write, write
.set .Runtime.exit, exit
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
    mov  rdi, [rip + .Memory.maxPages]
    shl  rdi, 16
    call .Runtime.malloc
    mov  [rip + .Memory.addr], rax

    # 内存清零
    mov  rdi, [rip + .Memory.addr]
    mov  rsi, 0
    mov  rdx, [rip + .Memory.maxPages]
    shl  rdx, 16
    call .Runtime.memset

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Memory.addr]
    mov  rdi, [rip + .Memory.dataOffset.0]
    add  rdi, rax
    lea  rsi, [rip + .Memory.dataPtr.0]
    mov  rdx, [rip + .Memory.dataSize.0]
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
    mov  rdi, 0
    call .Runtime.exit

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
    mov  rdi, 2 # stderr
    lea  rsi, [rip + .Runtime.panic.message]
    mov  rdx, [rip + .Runtime.panic.messageLen] # size
    call .Runtime.write

    # 退出程序
    mov  rdi, 1 # 退出码
    call .Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.global .F.main
.F.main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # env_get_multi_values()
    # 根据 Linux ABI，第一个参数 RDI 必须指向存储返回值的内存地址
    # 原本 WAT 里的第一个参数 (param i64) 顺延到 rsi
    lea  rdi, [rsp + 0]
    mov  rsi, 100
    call .Import.env.get_multi_values
    mov  rdi, [rsp + 16] # v3
    call .Import.env.print_i64
    mov  rdi, [rsp + 8] # v2
    call .Import.env.print_i64
    mov  rdi, [rsp + 0] # v1
    call .Import.env.print_i64

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

