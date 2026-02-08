# 源文件: abi-return-01.wat, ABI: x64-Windows
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
# .extern .Wa.Import.env.get_multi_values
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
    # local input: i64

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local input = 0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 100
    movabs rax, 100
    mov    [rbp-16], rax

    # local.set input i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-8], rax

    # local.get input i64
    mov rax, qword ptr [rbp-8]
    mov qword ptr [rbp-16], rax

    # call env.get_multi_values(...)
    lea rcx, [rsp+40] # return address
    mov rdx, qword ptr [rbp-16] # arg 0
    call .Wa.Import.env.get_multi_values
    mov r10, qword ptr [rax+0]
    mov qword ptr [rbp-16], r10
    mov r10, qword ptr [rax+8]
    mov qword ptr [rbp-24], r10
    mov r10, qword ptr [rax+16]
    mov qword ptr [rbp-32], r10

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-32] # arg 0
    call .Wa.Import.env.print_i64

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.Import.env.print_i64

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-16] # arg 0
    call .Wa.Import.env.print_i64

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

