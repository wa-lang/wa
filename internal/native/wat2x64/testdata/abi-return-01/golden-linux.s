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
.extern .Wa.Import.env.get_multi_values
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
.global .Wa.F.main
.Wa.F.main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # env_get_multi_values()
    # 根据 Linux ABI，第一个参数 RDI 必须指向存储返回值的内存地址
    # 原本 WAT 里的第一个参数 (param i64) 顺延到 rsi
    lea  rdi, [rsp + 0]
    mov  rsi, 100
    call .Wa.Import.env.get_multi_values
    mov  rdi, [rsp + 16] # v3
    call .Wa.Import.env.print_i64
    mov  rdi, [rsp + 8] # v2
    call .Wa.Import.env.print_i64
    mov  rdi, [rsp + 0] # v1
    call .Wa.Import.env.print_i64

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

