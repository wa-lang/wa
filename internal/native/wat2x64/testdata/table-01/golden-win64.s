# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# 运行时函数
.extern  .Wa.Runtime.write
.extern  .Wa.Runtime.exit
.extern  .Wa.Runtime.malloc
.extern  .Wa.Runtime.memcpy
.extern  .Wa.Runtime.memset
.extern  .Wa.Runtime.memmove

# 导入函数(外部库定义)
.extern .Wa.Import.syscall.write

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
.Wa.Memory.dataPtr.0: .asciz "hello world\n"

# 内存初始化函数
.section .text
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rcx, [rip + .Wa.Memory.maxPages]
    shl  rcx, 16
    call .Wa.Runtime.malloc
    lea  rdx, [rip + .Wa.Memory.addr]
    mov  [rdx], rax

    # 内存清零
    mov  rcx, [rip + .Wa.Memory.addr]
    mov  rdx, 0
    mov  r8, [rip + .Wa.Memory.maxPages]
    shl  r8, 16
    call .Wa.Runtime.memset

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Wa.Memory.addr]
    mov  rcx, [rip + .Wa.Memory.dataOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Wa.Memory.dataPtr.0]
    mov  r8, [rip + .Wa.Memory.dataSize.0]
    call .Wa.Runtime.memcpy

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义表格
.section .data
.align 8
.globl .Wa.Table.addr
.globl .Wa.Table.size
.globl .Wa.Table.maxSize
.Wa.Table.addr: .quad 0
.Wa.Table.size: .quad 1
.Wa.Table.maxSize: .quad 1

# 函数列表
# 保持连续并填充全部函数
.section .data
.align 8
.Wa.Table.funcIndexList:
.Wa.Table.funcIndexList.0: .quad .Wa.Import.syscall.write
.Wa.Table.funcIndexList.1: .quad .Wa.F.main
.Wa.Table.funcIndexList.end: .quad 0

# 表格初始化函数
.section .text
.globl .Wa.Table.initFunc
.Wa.Table.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配表格
    mov  rcx, [rip + .Wa.Table.maxSize]
    shl  rcx, 3 # sizeof(i64) == 8
    call .Wa.Runtime.malloc
    mov  [rip + .Wa.Table.addr], rax

    # 表格填充 0xFF
    mov  rcx, [rip + .Wa.Table.addr]
    mov  rdx, 0xFF
    mov  r8, [rip + .Wa.Table.maxSize]
    shl  r8, 3 # sizeof(i64) == 8
    call .Wa.Runtime.memset

    # 初始化表格

    # 加载表格地址
    mov rax, [rip + .Wa.Table.addr]
    # elem[0]: table[0+0] = syscall.write
    mov qword ptr [rax+0], 0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 汇编程序入口函数
.section .text
.globl _start
_start:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    call .Wa.Memory.initFunc
    call .Wa.Table.initFunc
    call .Wa.F.main

    # runtime.exit(0)
    mov  rcx, 0
    call .Wa.Runtime.exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

.section .data
.align 8
.Wa.Runtime.panic.message: .asciz "panic"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stderr, panicMessage, size)
    mov  rcx, 2 # stderr
    lea  rdx, [rip + .Wa.Runtime.panic.message]
    mov  r8, [rip + .Wa.Runtime.panic.messageLen] # size
    call .Wa.Runtime.write

    # 退出程序
    mov  rcx, 1 # 退出码
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
    sub  rsp, 64

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

    # i32.const 0
    mov rax, 0
    mov [rbp-32], rax

    # 加载函数的地址

    # r10 = table[?]
    mov  rax, [rip+.Wa.Table.addr]
    mov  r10d, [rbp-32] # i32
    mov  r10, [rax+r10*8]

    # r11 = .Table.funcIndexList[r10]
    lea  rax, [rip+.Wa.Table.funcIndexList]
    mov  r11, [rax+r10*8]

    # call_indirect r11(...)
    # type (i64,i64,i64) => i64
    mov rcx, qword ptr [rbp-8] # arg 0
    mov rdx, qword ptr [rbp-16] # arg 1
    mov r8, qword ptr [rbp-24] # arg 2
    call r11
    mov qword ptr [rbp-8], rax
    nop # drop [rbp-8]

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

