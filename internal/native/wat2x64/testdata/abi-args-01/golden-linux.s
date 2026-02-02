# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# 运行时函数
.extern .Wa.Runtime.write
.extern .Wa.Runtime.exit
.extern .Wa.Runtime.malloc
.extern .Wa.Runtime.memcpy
.extern .Wa.Runtime.memset

# 导入函数(外部库定义)
.extern .Wa.Import.env.write
.extern .Wa.Import.env.print_i64

# 定义内存
.section .data
.align 8
.globl .Wa.Memory.addr
.globl .Wa.Memory.pages
.globl .Wa.Memory.maxPages
.Wa.Memory.addr: .quad 0
.Wa.Memory.pages: .quad 1
.Wa.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 8
# memcpy(&Memory[8], data[0], size)
.Wa.Memory.dataOffset.0: .quad 8
.Wa.Memory.dataSize.0: .quad 12
.Wa.Memory.dataPtr.0: .ascii "hello world\n\000"

# 内存初始化函数
.section .text
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rdi, [rip + .Wa.Memory.maxPages]
    shl  rdi, 16
    call .Wa.Runtime.malloc
    mov  [rip + .Wa.Memory.addr], rax

    # 内存清零
    mov  rdi, [rip + .Wa.Memory.addr]
    mov  rsi, 0
    mov  rdx, [rip + .Wa.Memory.maxPages]
    shl  rdx, 16
    call .Wa.Runtime.memset

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Wa.Memory.addr]
    mov  rdi, [rip + .Wa.Memory.dataOffset.0]
    add  rdi, rax
    lea  rsi, [rip + .Wa.Memory.dataPtr.0]
    mov  rdx, [rip + .Wa.Memory.dataSize.0]
    call .Wa.Runtime.memcpy

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

    call .Wa.Memory.initFunc
    call .Wa.F.main

    # runtime.exit(0)
    mov  rdi, 0
    call .Wa.Runtime.exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

.section .data
.align 8
.Wa.Runtime.panic.message: .ascii "panic\000"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stderr, panicMessage, size)
    mov  rdi, 2 # stderr
    lea  rsi, [rip + .Wa.Runtime.panic.message]
    mov  rdx, [rip + .Wa.Runtime.panic.messageLen] # size
    call .Wa.Runtime.write

    # 退出程序
    mov  rdi, 1 # 退出码
    call .Wa.Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.Wa.F.main:
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
    call .Wa.Import.env.write
    mov qword ptr [rbp-8], rax
    nop # drop [rbp-8]

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

