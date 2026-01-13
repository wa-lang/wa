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

# 导入函数(由导入文件定义)
.extern $Import.syscall.write
.set $F.syscall.write, $Import.syscall.write

# 定义内存
.section .data
.align 8
$Memory.addr: .quad 0
$Memory.pages: .quad 1
$Memory.maxPages: .quad 1

# Memory[8]: hello worl...
$Memory.dataOffset.0: .quad 8
$Memory.dataSize.0: .quad 12
$Memory.dataPtr.0: .asciz "hello world\n"

# 内存初始化函数
.section .text
.globl $Memory.initFunc
$Memory.initFunc:
    # 影子空间
    sub rsp, 40

    # 分配内存
    mov  rcx, 65536 # 1 pages
    call malloc
    lea  rdx, [rip + $Memory.addr]
    mov  [rdx], rax

    # 内存清零
    lea  rcx, [rip + $Memory.addr]
    mov  rdx, 0
    mov  r8, 65536
    call .Runtime.memset

    # 初始化内存

    # Memory[8]: hello worl...
    lea  rcx, [rip + $Memory.addr]
    add  rcx, 8
    mov  rdx, [rip + $Memory.dataOffset.0]
    mov  r8, 12
    call .Runtime.memcpy

    # 函数返回
    add rsp, 40
    ret

.section .text
.globl main
main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    call $Memory.initFunc
    call $Table.initFunc
    call $F.main

    # return 0
    xor  eax, eax
    add  rsp, 32
    pop  rbp
    ret

.section .text

$F.main:
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

    # return
    add rsp, 40
    ret

