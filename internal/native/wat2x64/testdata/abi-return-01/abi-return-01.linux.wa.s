# 源文件: abi-return-01.wat, ABI: X64-Unix
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
.extern write
.extern exit
.extern malloc
.extern memcpy
.extern memset
.set .Runtime.write, write
.set .Runtime.exit, exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset

# 导入函数(外部库定义)
.extern wat2x64_env_get_multi_values
.extern wat2x64_env_print_i64
.set .Import.env.get_multi_values, wat2x64_env_get_multi_values
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
    mov  rdi, [rip + .Memory.maxPages]
    shl  rdi, 16
    call .Runtime.malloc
    mov  [rip + .Memory.addr], rax

    # 内存清零
    mov  rdi, [rip + .Memory.addr]
    mov  rsi, 0
    mov  rdx, [rip + .Memory.maxPages]
    shl  rdx, 16
    call .Runtime.memset

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

    call .Memory.initFunc
    call .F.main

    # runtime.exit(0)
    mov  rdi, 0
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
    mov  rdi, 2 # stderr
    lea  rsi, [rip + .Runtime.panic.message]
    mov  rdx, [rip + .Runtime.panic.messageLen] # size
    call .Runtime.write

    # 退出程序
    mov  rdi, 1 # 退出码
    call .Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.F.main:
    # local input: i64

    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local input = 0

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
    lea rdi, [rsp+8] # return address
    mov rsi, qword ptr [rbp-16] # arg 0
    call .Import.env.get_multi_values
    mov r10, qword ptr [rax+0]
    mov qword ptr [rbp-16], r10
    mov r10, qword ptr [rax+8]
    mov qword ptr [rbp-24], r10
    mov r10, qword ptr [rax+16]
    mov qword ptr [rbp-32], r10
    # call env.print_i64(...)
    mov rdi, qword ptr [rbp-32] # arg 0
    call .Import.env.print_i64
    # call env.print_i64(...)
    mov rdi, qword ptr [rbp-24] # arg 0
    call .Import.env.print_i64
    # call env.print_i64(...)
    mov rdi, qword ptr [rbp-16] # arg 0
    call .Import.env.print_i64

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

