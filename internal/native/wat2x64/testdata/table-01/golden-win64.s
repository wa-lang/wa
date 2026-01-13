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

# 定义表格
.section .data
.align 8
.Table.addr: .quad 0
.Table.size: .quad 1
.Table.maxSize: .quad 1

# 表格元素
.section .data
.align 8
.Table.elemOffset.0: .quad 0
.Table.elemFuncIndex.0: .quad 0 # .Import.syscall.write

# 函数列表
# 保持连续并填充全部函数
.section .data
.align 8
.Func.indexList:
.Func.indexList.0: .quad .Import.syscall.write
.Func.indexList.end: .quad 0

# 内存初始化函数
.section .text
.globl .Memory.initFunc
.Memory.initFunc:
    # 影子空间
    sub rsp, 40

    # 分配内存
    mov  rcx, [rip + .Memory.maxPages]
    shl  rcx, 16 # sizeof(page) == 65536
    call malloc
    lea  rdx, [rip + .Memory.addr]
    mov  [rdx], rax

    # 内存清零
    mov  rcx, [rip + .Memory.addr]
    mov  rdx, 0
    mov  r8, [rip + .Memory.maxPages]
    shl  r8, 16 # sizeof(page) == 65536
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


# 表格初始化函数
.section .text
.globl .Table.initFunc
.Table.initFunc:
    # 影子空间
    sub rsp, 40

    # 分配表格
    mov  rcx, [rip + .Table.maxSize]
    shl  rcx, 3 # sizeof(int) == 8
    call malloc
    lea  rdx, [rip + .Table.addr]
    mov  [rdx], rax

    # 表格清零
    mov  rcx, [rip + .Table.addr]
    mov  rdx, 0
    mov  r8, [rip + .Table.maxSize]
    shl  r8, 3 # sizeof(int) == 8
    call .Import.runtime.memset

    # 初始化表格
    
    # data[0]: ...
    mov  rax, [rip + .Table.addr]
    mov  rcx, [rip + .Table.elemOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Table.elemFuncIndex.0]
    mov  r8, 8
    call .Import.runtime.memcpy

    # 函数返回
    add rsp, 40
    ret

# Wasm 入口函数, 后续是编译器自动生成
.section .text
.global .F.main
.F.main:
    # rsp%16 == 0
    # (rsp-8-40)%16 == 0
    sub rsp, 40

    # r10 = table[0]
    mov  rax, [rip + .Table.addr]
    mov  rbx, 0
    mov  r10, [rax + rbx*8]

    # r11 = .Func.indexList[r10]
    lea  r11, [rip + .Func.indexList]
    mov  rax, [r11 + r10*8]
    mov  r11, rax

    # table[0](stdout, &memory[ptr], size)
    mov  rcx, 1 # stdout
    mov  rax, [rip + .Memory.addr]
    mov  rdx, [rip + .Memory.dataOffset.0]
    add  rdx, rax # rdx = &memory[ptr]
    mov  r8, [rip + .Memory.dataSize.0] # size
    call r11

    # 函数返回
    add rsp, 40
    ret

# 汇编程序入口函数
.section .text
.global main
main:
    # rsp%16 == 0
    # (rsp-8-40)%16 == 0
    sub rsp, 40

    call .Memory.initFunc
    call .Table.initFunc
    call .F.main

    # syscall.exit(0)
    mov  rcx, 0 # exit code
    call .Import.syscall.exit

    # exit 后这里不会被执行, 但是依然保留
    add rsp, 40
    ret
