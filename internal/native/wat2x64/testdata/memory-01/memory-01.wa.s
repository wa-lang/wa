# 源文件: memory-01.wat, 系统: windows/X64
# 自动生成的代码, 不要手动修改!!!

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
.extern _write
.set .Import.syscall._write, _write

# 定义内存
.section .data
.align 8
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
    # 影子空间
    sub rsp, 40

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

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Memory.addr]
    mov  rcx, [rip + .Memory.dataOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Memory.dataPtr.0]
    mov  r8, [rip + .Memory.dataSize.0]
    call .Runtime.memcpy

    # 函数返回
    add rsp, 40
    ret

# 汇编程序入口函数
.section .text
.globl main
main:
    # 影子内存
    sub rsp, 40

    call .Memory.initFunc
    call .F.main

    # runtime.exit(0)
    mov  rcx, 0
    call .Runtime._exit

    # exit 后这里不会被执行, 但是依然保留
    add rsp, 40
    ret

.section .data
.align 8
.Runtime.panic.message: .asciz "panic"
.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Runtime.panic
.Runtime.panic:
    # 影子内存
    sub rsp, 40

    # runtime.write(stderr, panicMessage, size)
    mov  rcx, 2 # stderr
    mov  rdx, [rip + .Runtime.panic.message]
    mov  r8, [rip + .Runtime.panic.messageLen] # size
    call .Runtime.panic

    # 退出程序
    mov  rcx, 1 # 退出码
    call .Runtime._exit

    # return
    add rsp, 40
    ret

.section .text

.F.main:
    # 栈帧开始
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    sub rsp, 0
    # i64.const 1
    movabs rax, 1
    mov qword ptr [rbp -24], rax
    # i64.const 8
    movabs rax, 8
    mov qword ptr [rbp -32], rax
    # i64.const 12
    movabs rax, 12
    mov qword ptr [rbp -40], rax
    mov R4, qword ptr [rbp -24]
    mov R5, qword ptr [rbp -16]
    mov R6, qword ptr [rbp -8]
    call syscall.write    mov R4, qword ptr [rbp -24]
    nop # drop R0
.L.return:

    # 栈帧结束
    add  rsp, 32
    pop  rbp
    ret

