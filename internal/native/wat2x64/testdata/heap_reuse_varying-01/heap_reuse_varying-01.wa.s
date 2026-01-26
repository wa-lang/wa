# 源文件: heap_reuse_varying-01.wat, ABI: x64-Windows
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

# 定义全局变量
.section .data
.align 8
.G.__heap_l128_freep: .long 0

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
    sub  rsp, 48

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 123
    movabs rax, 123
    mov    [rbp-8], rax

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-8] # arg 0
    call .Import.env.print_i64

.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.size(ptr:i32) => i32
.section .text
.F.heap_block.size:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.size, suffix=00000001)

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.load
    mov rax, [rip + .Memory.addr]
    mov r10d, dword ptr [rbp-16]
    add r10, rax
    mov eax, dword ptr [r10+0]
    mov dword ptr [rbp-16], eax

.L.brNext.heap_block.size.00000001:
    # fn.body.end(name=heap_block.size, suffix=00000001)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov eax,  dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax # ret.0

    # 将返回值变量复制到寄存器
    mov eax, [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.set_size(ptr:i32,size:i32)
.section .text
.F.heap_block.set_size:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.set_size, suffix=00000002)

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get size i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    mov rax, [rip + .Memory.addr]
    mov r10d, dword ptr [rbp-8]
    add r10, rax
    mov eax, dword ptr [rbp-16]
    mov dword ptr [r10+0], eax

.L.brNext.heap_block.set_size.00000002:
    # fn.body.end(name=heap_block.set_size, suffix=00000002)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.next(ptr:i32) => i32
.section .text
.F.heap_block.next:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.next, suffix=00000003)

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.load
    mov rax, [rip + .Memory.addr]
    mov r10d, dword ptr [rbp-16]
    add r10, rax
    mov eax, dword ptr [r10+4]
    mov dword ptr [rbp-16], eax

.L.brNext.heap_block.next.00000003:
    # fn.body.end(name=heap_block.next, suffix=00000003)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov eax,  dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax # ret.0

    # 将返回值变量复制到寄存器
    mov eax, [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.set_next(ptr:i32,next:i32)
.section .text
.F.heap_block.set_next:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.set_next, suffix=00000004)

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get next i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    mov rax, [rip + .Memory.addr]
    mov r10d, dword ptr [rbp-8]
    add r10, rax
    mov eax, dword ptr [rbp-16]
    mov dword ptr [r10+4], eax

.L.brNext.heap_block.set_next.00000004:
    # fn.body.end(name=heap_block.set_next, suffix=00000004)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.data(ptr:i32) => i32
.section .text
.F.heap_block.data:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.data, suffix=00000005)

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-16], eax

.L.brNext.heap_block.data.00000005:
    # fn.body.end(name=heap_block.data, suffix=00000005)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov eax,  dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax # ret.0

    # 将返回值变量复制到寄存器
    mov eax, [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.init(ptr:i32,size:i32,next:i32)
.section .text
.F.heap_block.init:
    push rbp
    mov  rbp, rsp
    sub  rsp, 0

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1
    mov [rbp+32], r8d # save arg.2

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=heap_block.init, suffix=00000006)

.L.brNext.heap_block.init.00000006:
    # fn.body.end(name=heap_block.init, suffix=00000006)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_reuse_varying(nbytes:i32) => i32
.section .text
.F.heap_reuse_varying:
    # local prevp: i32
    # local remaining: i32
    # local p: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-16], 0 # local prevp = 0
    mov dword ptr [rbp-24], 0 # local remaining = 0
    mov dword ptr [rbp-32], 0 # local p = 0

    # fn.body.begin(name=heap_reuse_varying, suffix=00000007)

    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-40], eax

    # local.set prevp i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-16], eax

    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-40], eax

    # local.set p i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-32], eax

    # loop.begin(name=continue, suffix=00000008)
.L.brNext.continue.00000008:
    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-40], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-56], eax

    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-56]
    mov dword ptr [rbp-48], eax

    # i32.ge_s
    mov   r10d, dword ptr [rbp-40]
    mov   r11d, dword ptr [rbp-48]
    cmp   r10d, r11d
    setge al
    movzx eax, al
    mov   dword ptr [rbp-40], eax

    # if.begin(name=, suffix=00000009)
    mov eax, [rbp-40]
    cmp eax, 0
    je  .L.brNext.00000009

    # if.body(name=, suffix=00000009)
    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # call heap_block.data(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    call .F.heap_block.data
    mov dword ptr [rbp-40], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.add
    mov eax, dword ptr [rbp-40]
    add eax, dword ptr [rbp-48]
    mov dword ptr [rbp-40], eax

    # local.set remaining i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-24], eax

    # local.get remaining i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-40], eax

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-48], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-48] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-48], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    call .F.heap_block.set_next

    # local.get remaining i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-40], eax

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-48], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-48] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-48], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-56], eax

    # i32.sub
    mov eax, dword ptr [rbp-48]
    sub eax, dword ptr [rbp-56]
    mov dword ptr [rbp-48], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-56], eax

    # i32.sub
    mov eax, dword ptr [rbp-48]
    sub eax, dword ptr [rbp-56]
    mov dword ptr [rbp-48], eax

    # call heap_block.set_size(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    call .F.heap_block.set_size

    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # local.get remaining i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-48], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    call .F.heap_block.set_next

    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # global.set __heap_l128_freep i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rip+.G.__heap_l128_freep], eax

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-56], eax

    # call heap_block.init(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    mov r8d, dword ptr [rbp-56] # arg 2
    call .F.heap_block.init

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # return
    jmp .L.brNext.heap_reuse_varying.00000007

    jmp .L.brNext.00000009

.L.brNext.00000009:
    # if.end(name=, suffix=00000009)

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-40], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.ge_s
    mov   r10d, dword ptr [rbp-40]
    mov   r11d, dword ptr [rbp-48]
    cmp   r10d, r11d
    setge al
    movzx eax, al
    mov   dword ptr [rbp-40], eax

    # if.begin(name=, suffix=0000000A)
    mov eax, [rbp-40]
    cmp eax, 0
    je  .L.brNext.0000000A

    # if.body(name=, suffix=0000000A)
    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # block.begin(name=, suffix=0000000B)

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-48], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-48] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-48], eax

.L.brNext.0000000B:
    # block.end(name=, suffix=0000000B)

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    call .F.heap_block.set_next

    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # global.set __heap_l128_freep i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rip+.G.__heap_l128_freep], eax

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-48], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    mov edx, dword ptr [rbp-48] # arg 1
    call .F.heap_block.set_next

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # return
    jmp .L.brNext.heap_reuse_varying.00000007

    jmp .L.brNext.0000000A

.L.brNext.0000000A:
    # if.end(name=, suffix=0000000A)

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-48], eax

    # i32.eq
    mov   r10d, dword ptr [rbp-40]
    mov   r11d, dword ptr [rbp-48]
    cmp   r10d, r11d
    sete  al      # al = (r10d==r11d)? 1: 0
    movzx eax, al # eax = al
    mov   dword ptr [rbp-40], eax

    # if.begin(name=, suffix=0000000C)
    mov eax, [rbp-40]
    cmp eax, 0
    je  .L.brNext.0000000C

    # if.body(name=, suffix=0000000C)
    # i32.const 0
    mov eax, 0
    mov [rbp-40], eax

    # return
    jmp .L.brNext.heap_reuse_varying.00000007

    jmp .L.brNext.0000000C

.L.brNext.0000000C:
    # if.end(name=, suffix=0000000C)

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # local.set prevp i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-40] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-40], eax

    # local.set p i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-32], eax

    # br continue
    jmp .L.brNext.continue.00000008

    # loop.end(name=continue, suffix=00000008)

    call .Runtime.panic # unreachable

.L.brNext.heap_reuse_varying.00000007:
    # fn.body.end(name=heap_reuse_varying, suffix=00000007)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov eax,  dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax # ret.0

    # 将返回值变量复制到寄存器
    mov eax, [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

