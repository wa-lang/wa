# 源文件: table-01.wat, 系统: windows/X64
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 系统调用
.extern malloc
.globl $builtin.memcpy
.globl $builtin.memset

# 导入函数(由导入文件定义)
.extern $Import.syscall.write
.set $F.syscall.write, $Import.syscall.write

# 定义表格
.section .data
.align 8
$Table.size: .quad 1
$Table.maxSize: .quad 0
$Table.addr: .fill 8, 1, 0
$Table.funcIndexList: .fill 8, 2, 0

# 表格初始化函数
.section .text
.globl $Table.initFunc
$Table.initFunc:
    # 影子空间
    sub rsp, 40

    # 初始化全部函数索引列表

    # 加载函数索引列表地址
    lea rax, [rip + $Table.funcIndexList]
    # 导入[0] syscall.write
    lea rcx, [rip + $F.syscall.write]
    mov [rax + 0], rcx
    # 函数[0] main
    lea rcx, [rip + $F.main]
    mov [rax + 8], rcx

    # 初始化表格元素

    # 加载表格地址
    lea rax, [rip + $Table.addr]
    # 表格[0] = syscall.write
    mov qword ptr [rax + 0], 0

    # 函数返回
    add rsp, 40
    ret

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
    call $builtin.memset

    # 初始化内存

    # Memory[8]: hello worl...
    lea  rcx, [rip + $Memory.addr]
    add  rcx, 8
    mov  rdx, [rip + $Memory.dataOffset.0]
    mov  r8, 12
    call $builtin.memcpy

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
    # i32.const 0
    movabs rax, 0
    mov qword ptr [rbp -48], rax
    # 根据函数索引编号从列表查询函数地址
    lea rax, [rip + $Table.funcIndexList]
    mov r10, [rbp -48]
    add rax, r10
    lea r11, [rax]
    mov R4, qword ptr [rbp -24]
    mov R5, qword ptr [rbp -16]
    mov R6, qword ptr [rbp -8]
    call r11
    mov R4, qword ptr [rbp -24]
    nop # drop R0
.L.return:

    # 栈帧结束
    add  rsp, 32
    pop  rbp
    ret

