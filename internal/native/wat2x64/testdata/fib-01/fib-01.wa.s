# 源文件: fib-01.wat, ABI: x64-Windows
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
# .extern .Wa.Runtime.write
# .extern .Wa.Runtime.exit
# .extern .Wa.Runtime.malloc
# .extern .Wa.Runtime.memcpy
# .extern .Wa.Runtime.memset
# .extern .Wa.Runtime.memmove

# 导入函数(外部库定义)
# .extern .Wa.Import.env.print_i64

# 定义内存
.section .data
.align 8
.globl .Wa.Memory.addr
.Wa.Memory.addr: .quad 0
.globl .Wa.Memory.pages
.Wa.Memory.pages: .quad 1
.globl .Wa.Memory.maxPages
.Wa.Memory.maxPages: .quad 1

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
.Wa.Runtime.panic.message: .ascii "panic"
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

    # fn.body.begin(name=main, suffix=00000000)

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

    # loop.begin(name=my_loop, suffix=00000001)
.Wa.L.brNext.my_loop.00000001:
    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.F.fib
    mov qword ptr [rbp-24], rax

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.Import.env.print_i64

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

    # br_if my_loop.00000001
    mov eax, dword ptr [rbp-24]
    cmp eax, 0
    je  .Wa.L.brFallthrough.my_loop.00000001
    jmp .Wa.L.brNext.my_loop.00000001
.Wa.L.brFallthrough.my_loop.00000001:

    # loop.end(name=my_loop, suffix=00000001)

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func fib(n:i64) => i64
.section .text
.globl .Wa.F.fib
.Wa.F.fib:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov [rbp+16], rcx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=fib, suffix=00000002)

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

    # if.begin(name=, suffix=00000003)
    mov eax, [rbp-16]
    cmp eax, 0
    je  .Wa.L.else.00000003

    # if.body(name=, suffix=00000003)
    # i64.const 1
    movabs rax, 1
    mov    [rbp-16], rax

    jmp .Wa.L.brNext.00000003

.Wa.L.else.00000003:
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
    call .Wa.F.fib
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
    call .Wa.F.fib
    mov qword ptr [rbp-24], rax

    # i64.add
    mov rax, qword ptr [rbp-16]
    add rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

.Wa.L.brNext.00000003:
    # if.end(name=, suffix=00000003)

.Wa.L.brNext.fib.00000002:
    # fn.body.end(name=fib, suffix=00000002)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-8], rax # ret.0

    # 将返回值变量复制到寄存器
    mov rax, [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

