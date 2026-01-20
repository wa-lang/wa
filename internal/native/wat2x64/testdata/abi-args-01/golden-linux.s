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
.extern wat2x64_env_write
.extern wat2x64_env_print_i64
.set .Import.env.write, wat2x64_env_write
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
.Memory.dataPtr.0: .asciz "hello world\n"

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
.F.main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 80

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # i64.const 1
    movabs rax, 1
    mov    [rbp-8], rax

    # i64.const 8
    movabs rax, 8
    mov    [rbp-16], rax

    # i64.const 12
    movabs rax, 12
    mov    [rbp-24], rax

    # i64.const 100
    movabs rax, 100
    mov    [rbp-32], rax

    # i64.const 200
    movabs rax, 200
    mov    [rbp-40], rax

    # i64.const 300
    movabs rax, 300
    mov    [rbp-48], rax

    # i64.const 400
    movabs rax, 400
    mov    [rbp-56], rax
    
    # i64.const 500
    movabs rax, 500
    mov    [rbp-64], rax

    # call env.write(...)
    mov rdi, qword ptr [rbp-8] # arg 0
    mov rsi, qword ptr [rbp-16] # arg 1
    mov rdx, qword ptr [rbp-24] # arg 2
    mov rcx, qword ptr [rbp-32] # arg 3
    mov r8, qword ptr [rbp-40] # arg 4
    mov r9, qword ptr [rbp-48] # arg 5
    mov rax, qword ptr [rbp-56]
    mov qword ptr [rsp+0], rax
    mov rax, qword ptr [rbp-64]
    mov qword ptr [rsp+8], rax
    call .Import.env.write
    mov qword ptr [rbp-8], rax
    nop # drop [rbp-8]

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

