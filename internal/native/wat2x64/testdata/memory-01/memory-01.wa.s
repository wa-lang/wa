# 源文件: memory-01.wat, ABI: x64-Windows
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
.extern .Wa.Runtime.write
.extern .Wa.Runtime.exit
.extern .Wa.Runtime.malloc
.extern .Wa.Runtime.memcpy
.extern .Wa.Runtime.memset
.extern .Wa.Runtime.memmove

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
.Wa.Memory.dataPtr.0: .asciz "hello world\012"

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
    mov  [rip + .Wa.Memory.addr], rax

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

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 1
    movabs rax, 1
    mov    [rbp-8], rax

    # i64.const 8
    movabs rax, 8
    mov    [rbp-16], rax

    # i64.const 12
    movabs rax, 12
    mov    [rbp-24], rax

    # call syscall.write(...)
    mov rcx, qword ptr [rbp-8] # arg 0
    mov rdx, qword ptr [rbp-16] # arg 1
    mov r8, qword ptr [rbp-24] # arg 2
    call .Wa.Import.syscall.write
    mov qword ptr [rbp-8], rax

    nop # drop [rbp-8]

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

