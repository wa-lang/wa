# 源文件: wa-runtime-01.wat, ABI: x64-Windows
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
.extern _write
.extern _exit
.extern malloc
.extern memcpy
.extern memset
.set .Runtime.write, _write
.set .Runtime.exit, _exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset

# 导入函数(外部库定义)
.extern wat2x64_syscall_write
.set .Import.syscall.write, wat2x64_syscall_write

# 定义内存
.section .data
.align 8
.globl .Memory.addr
.globl .Memory.pages
.globl .Memory.maxPages
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

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Memory.addr]
    mov  rcx, [rip + .Memory.dataOffset.0]
    add  rcx, rax
    lea  rdx, [rip + .Memory.dataPtr.0]
    mov  r8, [rip + .Memory.dataSize.0]
    call .Runtime.memcpy

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义全局变量
.section .data
.align 8
$G.__stack_ptr: .long 1024
$G.__heap_base: .long 1048576
$G.__heap_lfixed_cap: .long 64

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

    # call syscall.write(...)
    mov rcx, qword ptr [rbp-8] # arg 0
    mov rdx, qword ptr [rbp-16] # arg 1
    mov r8, qword ptr [rbp-24] # arg 2
    call .Import.syscall.write
    mov qword ptr [rbp-8], rax
    nop # drop [rbp-8]

    # 根据ABI处理返回值
.L.return.main:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.throw
.section .text
.F.runtime.throw:
    push rbp
    mov  rbp, rsp
    sub  rsp, 0

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    call .Runtime.panic # unreachable

    # 根据ABI处理返回值
.L.return.runtime.throw:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.getStackPtr
.section .text
.F.runtime.getStackPtr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret .F.ret.0 = 0
    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+$G.__stack_ptr]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.runtime.getStackPtr:
    mov rax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.setStackPtr
.section .text
.F.runtime.setStackPtr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], rcx # save arg sp
    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get sp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+$G.__stack_ptr], eax

    # 根据ABI处理返回值
.L.return.runtime.setStackPtr:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.stackAlloc
.section .text
.F.runtime.stackAlloc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], rcx # save arg size
    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret .F.ret.0 = 0
    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+$G.__stack_ptr]
    mov dword ptr [rbp-16], eax
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.sub
    mov eax, dword ptr [rbp-16]
    sub eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rip+$G.__stack_ptr], eax
    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+$G.__stack_ptr]
    mov dword ptr [rbp-16], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.stackAlloc

    # 根据ABI处理返回值
.L.return.runtime.stackAlloc:
    mov rax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.stackFree
.section .text
.F.runtime.stackFree:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], rcx # save arg size
    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+$G.__stack_ptr]
    mov dword ptr [rbp-8], eax
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-8]
    add eax, dword ptr [rbp+0]
    mov dword ptr [rbp+0], eax
    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+$G.__stack_ptr], eax

    # 根据ABI处理返回值
.L.return.runtime.stackFree:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.heapBase
.section .text
.F.runtime.heapBase:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret .F.ret.0 = 0
    # 没有局部变量需要初始化为0

    # global.get __heap_base i32
    mov eax, dword ptr [rip+$G.__heap_base]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.runtime.heapBase:
    mov rax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

