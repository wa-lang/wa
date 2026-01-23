# 源文件: fib-01.wat, ABI: x64-Windows
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
.extern _write
.extern _exit
.extern malloc
.extern memcpy
.extern memset
.extern memmove
.set .Runtime.write, _write
.set .Runtime.exit, _exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset
.set .Runtime.memmove, memmove

# 导入函数(外部库定义)
.extern wat2x64_env_print_i64
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

# 内存初始化函数
.section .text
.globl .Memory.initFunc
.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rcx, [rip + .Memory.maxPages]
    shl  rcx, 16
    call .Runtime.malloc
    mov  [rip + .Memory.addr], rax

    # 内存清零
    mov  rcx, [rip + .Memory.addr]
    mov  rdx, 0
    mov  r8, [rip + .Memory.maxPages]
    shl  r8, 16
    call .Runtime.memset

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义表格
.section .data
.align 8
.globl .Table.addr
.globl .Table.size
.globl .Table.maxSize
.Table.addr: .quad 0
.Table.size: .quad 1
.Table.maxSize: .quad 1

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
    mov  rcx, 0
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
    mov  rcx, 2 # stderr
    lea  rdx, [rip + .Runtime.panic.message]
    mov  r8, [rip + .Runtime.panic.messageLen] # size
    call .Runtime.write

    # 退出程序
    mov  rcx, 1 # 退出码
    call .Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.F.main:
    # local N: i64
    # local i: i64

    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local N = 0
    mov dword ptr [rbp-16], 0 # local i = 0

    # i64.const 10
    movabs rax, 10
    mov    [rbp-24], rax

    # local.set N i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-8], rax

    # i64.const 1
    movabs rax, 1
    mov    [rbp-24], rax

    # local.set i i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # loop.begin(name=my_loop, suffix=00000000)
.L.brNext.my_loop.00000000:
    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .F.fib
    mov qword ptr [rbp-24], rax

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Import.env.print_i64

    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # i64.const 1
    movabs rax, 1
    mov    [rbp-32], rax

    # i64.add
    mov rax, qword ptr [rbp-24]
    add rax, qword ptr [rbp-32]
    mov qword ptr [rbp-24], rax

    # local.set i i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # local.get N i64
    mov rax, qword ptr [rbp-8]
    mov qword ptr [rbp-32], rax

    # i64.lt_u
    mov   r10, qword ptr [rbp-24]
    mov   r11, qword ptr [rbp-32]
    cmp   r10, r11
    setb  al
    movzx eax, al
    mov   dword ptr [rbp-24], eax

    # br_if my_loop00000000
    mov eax, dword ptr [rbp-24]
    cmp eax, 0
    je  .L.brFallthrough.my_loop.00000000
    jmp .L.brNext.my_loop.00000000
.L.brFallthrough.my_loop.00000000:

    # loop.end(name=my_loop, suffix=00000000)


    # 根据ABI处理返回值
.L.return.main:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func fib(n:i64) => i64
.section .text
.globl .F.fib
.F.fib:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov [rbp+16], rcx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-16], rax

    # i64.const 2
    movabs rax, 2
    mov    [rbp-24], rax

    # i64.le_u
    mov   r10, qword ptr [rbp-16]
    mov   r11, qword ptr [rbp-24]
    cmp   r10, r11
    setbe al
    movzx eax, al
    mov   dword ptr [rbp-16], eax

    # if.begin(name=, suffix=00000001)
    mov eax, [rbp-16]
    cmp eax, 0
    je  .L.else.00000001

    # if.body(name=, suffix=00000001)
    # i64.const 1
    movabs rax, 1
    mov    [rbp-16], rax

    jmp .L.brNext.00000001

    # if.else(name=, suffix=00000001)
.L.else.00000001:
    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-16], rax

    # i64.const 1
    movabs rax, 1
    mov    [rbp-24], rax

    # i64.sub
    mov rax, qword ptr [rbp-16]
    sub rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-16] # arg 0
    call .F.fib
    mov qword ptr [rbp-16], rax

    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-24], rax

    # i64.const 2
    movabs rax, 2
    mov    [rbp-32], rax

    # i64.sub
    mov rax, qword ptr [rbp-24]
    sub rax, qword ptr [rbp-32]
    mov qword ptr [rbp-24], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .F.fib
    mov qword ptr [rbp-24], rax

    # i64.add
    mov rax, qword ptr [rbp-16]
    add rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

.L.brNext.00000001:
    # if.end(name=, suffix=00000001)

    # return
    # copy result from stack
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-8], rax # ret.0

    # 根据ABI处理返回值
.L.return.fib:
    mov rax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

