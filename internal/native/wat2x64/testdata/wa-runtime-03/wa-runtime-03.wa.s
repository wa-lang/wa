# 源文件: wa-runtime-03.wat, ABI: x64-Windows
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

# 定义表格
.section .data
.align 8
.globl .Table.addr
.globl .Table.size
.globl .Table.maxSize
.Table.addr: .quad 0
.Table.size: .quad 1
.Table.maxSize: .quad 1

# 函数列表
# 保持连续并填充全部函数
.section .data
.align 8
.Table.funcIndexList:
.Table.funcIndexList.0: .quad .Import.syscall.write
.Table.funcIndexList.1: .quad .F.main
.Table.funcIndexList.2: .quad .F.runtime.throw
.Table.funcIndexList.3: .quad .F.runtime.getStackPtr
.Table.funcIndexList.4: .quad .F.runtime.setStackPtr
.Table.funcIndexList.5: .quad .F.runtime.stackAlloc
.Table.funcIndexList.6: .quad .F.runtime.stackFree
.Table.funcIndexList.7: .quad .F.runtime.heapBase
.Table.funcIndexList.8: .quad .F.runtime.HeapAlloc
.Table.funcIndexList.9: .quad .F.runtime.HeapFree
.Table.funcIndexList.10: .quad .F.runtime.Block.Init
.Table.funcIndexList.11: .quad .F.runtime.Block.SetFinalizer
.Table.funcIndexList.12: .quad .F.runtime.Block.HeapAlloc
.Table.funcIndexList.13: .quad .F.runtime.DupI32
.Table.funcIndexList.14: .quad .F.runtime.SwapI32
.Table.funcIndexList.15: .quad .F.runtime.Block.Retain
.Table.funcIndexList.16: .quad .F.runtime.Block.Release
.Table.funcIndexList.17: .quad .F.$wa.runtime.i32_ref_to_ptr
.Table.funcIndexList.18: .quad .F.$wa.runtime.i64_ref_to_ptr
.Table.funcIndexList.19: .quad .F.$wa.runtime.slice_to_ptr
.Table.funcIndexList.20: .quad .F.runtime.malloc
.Table.funcIndexList.21: .quad .F.runtime.free
.Table.funcIndexList.end: .quad 0

# 表格初始化函数
.section .text
.globl .Table.initFunc
.Table.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配表格
    mov  rcx, [rip + .Table.maxSize]
    shl  rcx, 3 # sizeof(i64) == 8
    call .Runtime.malloc
    mov  [rip + .Table.addr], rax

    # 表格填充 0xFF
    mov  rcx, [rip + .Table.addr]
    mov  rdx, 0xFF
    mov  r8, [rip + .Table.maxSize]
    shl  r8, 3 # sizeof(i64) == 8
    call .Runtime.memset

    # 初始化表格

    # 加载表格地址
    mov rax, [rip + .Table.addr]
    # elem[0]: table[0+0] = syscall.write
    mov qword ptr [rax+0], 0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义全局变量
.section .data
.align 8
.G.__stack_ptr: .long 1024
.G.__heap_base: .long 1048576
.G.__heap_lfixed_cap: .long 64

# 汇编程序入口函数
.section .text
.globl main
main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    call .Memory.initFunc
    call .Table.initFunc
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

# func runtime.getStackPtr => i32
.section .text
.F.runtime.getStackPtr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.runtime.getStackPtr:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.setStackPtr(sp:i32)
.section .text
.F.runtime.setStackPtr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get sp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__stack_ptr], eax

    # 根据ABI处理返回值
.L.return.runtime.setStackPtr:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.stackAlloc(size:i32) => i32
.section .text
.F.runtime.stackAlloc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-16], eax
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.sub
    mov eax, dword ptr [rbp-24]
    sub eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rip+.G.__stack_ptr], eax
    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-16], eax
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.stackAlloc

    # 根据ABI处理返回值
.L.return.runtime.stackAlloc:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.stackFree(size:i32)
.section .text
.F.runtime.stackFree:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-8], eax
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # global.set __stack_ptr i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__stack_ptr], eax

    # 根据ABI处理返回值
.L.return.runtime.stackFree:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.heapBase => i32
.section .text
.F.runtime.heapBase:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.runtime.heapBase:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.HeapAlloc(nbytes:i32) => i32
.section .text
.globl .F.runtime.HeapAlloc
.F.runtime.HeapAlloc:
    # local ptr: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-16], 0 # local ptr = 0

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.eqz
    mov  eax, dword ptr [rbp-24]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-24], eax
.L.if.begin.00000000:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000000 # if eax != 0 { jmp end }
.L.if.body.00000000:
    # i32.const 0
    mov eax, 0
    mov [rbp-24], eax

    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.HeapAlloc
.L.next.00000000:
.L.if.end.00000000:
    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.const 7
    mov eax, 7
    mov [rbp-32], eax

    # i32.add
    mov eax, dword ptr [rbp-32]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-32], eax

    # i32.div_u
    mov edx, dword ptr [rbp-24]
    mov eax, dword ptr [rbp-32]
    xor edx, edx
    div dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
    mov edx, dword ptr [rbp+8]
    # i32.const 8
    mov eax, 8
    mov [rbp-32], eax

    # i32.mul
    mov eax, dword ptr [rbp-32]
    imul eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
    # local.set nbytes i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp+16], eax

    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call runtime.malloc(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.runtime.malloc
    mov dword ptr [rbp-24], eax
    # local.set ptr i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-16], eax

.L.loop.begin.zero00000001:
.L.next.zero00000001:
    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-32], eax

    # i32.sub
    mov eax, dword ptr [rbp-32]
    sub eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
    # local.tee nbytes i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp+16], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-32], eax

    # i32.add
    mov eax, dword ptr [rbp-32]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
    # i64.const 0
    movabs rax, 0
    mov    [rbp-32], rax

    # i64.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov rax, qword [rbp-32]
    mov qword [r10 +0], rax
    # local.get nbytes i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

.L.if.begin.00000002:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000002 # if eax != 0 { jmp end }
.L.if.body.00000002:
    jmp .L.next.zero00000001
.L.next.00000002:
.L.if.end.00000002:
.L.loop.end.zero00000001:
    # local.get ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-24], eax


    # 根据ABI处理返回值
.L.return.runtime.HeapAlloc:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.HeapFree(ptr:i32)
.section .text
.globl .F.runtime.HeapFree
.F.runtime.HeapFree:
    push rbp
    mov  rbp, rsp
    sub  rsp, 48

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # call runtime.free(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    call .F.runtime.free

    # 根据ABI处理返回值
.L.return.runtime.HeapFree:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.Block.Init(ptr:i32,item_count:i32,release_func:i32,item_size:i32) => i32
.section .text
.F.runtime.Block.Init:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1
    mov [rbp+32], r8d # save arg.2
    mov [rbp+40], r9d # save arg.3

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

.L.if.begin.00000003:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000003 # if eax != 0 { jmp end }
.L.if.body.00000003:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.const 1
    mov eax, 1
    mov [rbp-32], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov eax, dword [rbp-32]
    mov dword [r10 +0], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # local.get item_count i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-32], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov eax, dword [rbp-32]
    mov dword [r10 +4], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # local.get release_func i32
    mov eax, dword ptr [rbp+32]
    mov dword ptr [rbp-32], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov eax, dword [rbp-32]
    mov dword [r10 +8], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # local.get item_size i32
    mov eax, dword ptr [rbp+40]
    mov dword ptr [rbp-32], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov eax, dword [rbp-32]
    mov dword [r10 +12], eax
.L.next.00000003:
.L.if.end.00000003:

    # 根据ABI处理返回值
.L.return.runtime.Block.Init:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.Block.SetFinalizer(ptr:i32,release_func:i32)
.section .text
.F.runtime.Block.SetFinalizer:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

.L.if.begin.00000004:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000004 # if eax != 0 { jmp end }
.L.if.body.00000004:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get release_func i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-8]
    add r10, rax
    mov eax, dword [rbp-16]
    mov dword [r10 +8], eax
.L.next.00000004:
.L.if.end.00000004:

    # 根据ABI处理返回值
.L.return.runtime.Block.SetFinalizer:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.Block.HeapAlloc(item_count:i32,release_func:i32,item_size:i32) => i32,i32)
.section .text
.globl .F.runtime.Block.HeapAlloc
.F.runtime.Block.HeapAlloc:
    # local b: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 80

    # 将返回地址备份到栈
    mov [rbp+16], rcx # return address

    # 将寄存器参数备份到栈
    mov [rbp+16], edx # save arg.0
    mov [rbp+24], r8d # save arg.1
    mov [rbp+32], r9d # save arg.2

    # 将返回值变量初始化为0
    mov dword ptr [rbp+48], 0 # ret.0 = 0
    mov dword ptr [rbp+56], 0 # ret.1 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local b = 0

    # local.get item_count i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get item_size i32
    mov eax, dword ptr [rbp+32]
    mov dword ptr [rbp-24], eax

    # i32.mul
    mov eax, dword ptr [rbp-24]
    imul eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # i32.const 16
    mov eax, 16
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # call runtime.HeapAlloc(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.runtime.HeapAlloc
    mov dword ptr [rbp-16], eax
    # local.get item_count i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # local.get release_func i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-32], eax

    # local.get item_size i32
    mov eax, dword ptr [rbp+32]
    mov dword ptr [rbp-40], eax

    # call runtime.Block.Init(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    mov r8d, dword ptr [rbp-32] # arg 2
    mov r9d, dword ptr [rbp-40] # arg 3
    call .F.runtime.Block.Init
    mov dword ptr [rbp-16], eax
    # local.tee b i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax
    # local.get b i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # i32.const 16
    mov eax, 16
    mov [rbp-32], eax

    # i32.add
    mov eax, dword ptr [rbp-32]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax

    # 根据ABI处理返回值
.L.return.runtime.Block.HeapAlloc:
    mov rax, [rbp+48] # ret address

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.DupI32(a:i32) => i32,i32)
.section .text
.F.runtime.DupI32:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将返回地址备份到栈
    mov [rbp+16], rcx # return address

    # 将寄存器参数备份到栈
    mov [rbp+16], edx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp+48], 0 # ret.0 = 0
    mov dword ptr [rbp+56], 0 # ret.1 = 0

    # 没有局部变量需要初始化为0

    # local.get a i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get a i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax


    # 根据ABI处理返回值
.L.return.runtime.DupI32:
    mov rax, [rbp+48] # ret address

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.SwapI32(a:i32,b:i32) => i32,i32)
.section .text
.F.runtime.SwapI32:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将返回地址备份到栈
    mov [rbp+16], rcx # return address

    # 将寄存器参数备份到栈
    mov [rbp+16], edx # save arg.0
    mov [rbp+24], r8d # save arg.1

    # 将返回值变量初始化为0
    mov dword ptr [rbp+48], 0 # ret.0 = 0
    mov dword ptr [rbp+56], 0 # ret.1 = 0

    # 没有局部变量需要初始化为0

    # local.get b i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-8], eax

    # local.get a i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax


    # 根据ABI处理返回值
.L.return.runtime.SwapI32:
    mov rax, [rbp+48] # ret address

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.Block.Retain(ptr:i32) => i32
.section .text
.globl .F.runtime.Block.Retain
.F.runtime.Block.Retain:
    push rbp
    mov  rbp, rsp
    sub  rsp, 48

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

.L.if.begin.00000005:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000005 # if eax != 0 { jmp end }
.L.if.body.00000005:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-32]
    add r10, rax
    mov eax, dword [r10 +0]
    mov dword [rbp-32], eax
    # i32.const 1
    mov eax, 1
    mov [rbp-40], eax

    # i32.add
    mov eax, dword ptr [rbp-40]
    add eax, dword ptr [rbp-32]
    mov dword ptr [rbp-32], eax
    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-24]
    add r10, rax
    mov eax, dword [rbp-32]
    mov dword [r10 +0], eax
.L.next.00000005:
.L.if.end.00000005:

    # 根据ABI处理返回值
.L.return.runtime.Block.Retain:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.Block.Release(ptr:i32)
.section .text
.globl .F.runtime.Block.Release
.F.runtime.Block.Release:
    # local ref_count: i32
    # local item_count: i32
    # local free_func: i32
    # local item_size: i32
    # local data_ptr: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local ref_count = 0
    mov dword ptr [rbp-16], 0 # local item_count = 0
    mov dword ptr [rbp-24], 0 # local free_func = 0
    mov dword ptr [rbp-32], 0 # local item_size = 0
    mov dword ptr [rbp-40], 0 # local data_ptr = 0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-56], eax

    # i32.eq
    mov r10d, dword ptr [rbp-56]
    mov r11d, dword ptr [rbp-48]
    cmp r10d, r11d
    sete al
    movzx eax, al
    mov dword ptr [rbp-48], eax
.L.if.begin.00000006:
    mov eax, [rbp-48]
    test eax, eax
    je .L.if.end.00000006 # if eax != 0 { jmp end }
.L.if.body.00000006:
    jmp .L.return.runtime.Block.Release
.L.next.00000006:
.L.if.end.00000006:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-48]
    add r10, rax
    mov eax, dword [r10 +0]
    mov dword [rbp-48], eax
    # i32.const 1
    mov eax, 1
    mov [rbp-56], eax

    # i32.sub
    mov eax, dword ptr [rbp-56]
    sub eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # local.set ref_count i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-8], eax

    # local.get ref_count i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-48], eax

.L.if.begin.00000007:
    mov eax, [rbp-48]
    test eax, eax
    jne .L.if.body.00000007 # if eax != 0 { jmp body }
.L.if.body.00000007:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # local.get ref_count i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-56], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-48]
    add r10, rax
    mov eax, dword [rbp-56]
    mov dword [r10 +0], eax
.L.if.else.00000007:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-48]
    add r10, rax
    mov eax, dword [r10 +8]
    mov dword [rbp-48], eax
    # local.set free_func i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-24], eax

    # local.get free_func i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-48], eax

.L.if.begin.00000008:
    mov eax, [rbp-48]
    test eax, eax
    je .L.if.end.00000008 # if eax != 0 { jmp end }
.L.if.body.00000008:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-48]
    add r10, rax
    mov eax, dword [r10 +4]
    mov dword [rbp-48], eax
    # local.set item_count i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-16], eax

    # local.get item_count i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-48], eax

.L.if.begin.00000009:
    mov eax, [rbp-48]
    test eax, eax
    je .L.if.end.00000009 # if eax != 0 { jmp end }
.L.if.body.00000009:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-48]
    add r10, rax
    mov eax, dword [r10 +12]
    mov dword [rbp-48], eax
    # local.set item_size i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-32], eax

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.const 16
    mov eax, 16
    mov [rbp-56], eax

    # i32.add
    mov eax, dword ptr [rbp-56]
    add eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # local.set data_ptr i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-40], eax

.L.loop.begin.free_next0000000A:
.L.next.free_next0000000A:
    # local.get data_ptr i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-48], eax

    # local.get free_func i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-56], eax

    # 加载函数的地址

    # r10 = table[?]
    mov  rax, [rip+.Table.addr]
    mov  r10, [rbp-56]
    mov  r10, [rax+r10*8]

    # r11 = .Table.funcIndexList[r10]
    lea  rax, [rip+.Table.funcIndexList]
    mov  r11, [rax+r10*8]

    # call_indirect r11(...)
    # type (i32)
    mov ecx, dword ptr [rbp-48] # arg 0
    call r11
    # local.get item_count i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-48], eax

    # i32.const 1
    mov eax, 1
    mov [rbp-56], eax

    # i32.sub
    mov eax, dword ptr [rbp-56]
    sub eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # local.set item_count i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-16], eax

    # local.get item_count i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-48], eax

.L.if.begin.0000000B:
    mov eax, [rbp-48]
    test eax, eax
    je .L.if.end.0000000B # if eax != 0 { jmp end }
.L.if.body.0000000B:
    # local.get data_ptr i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-48], eax

    # local.get item_size i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-56], eax

    # i32.add
    mov eax, dword ptr [rbp-56]
    add eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # local.set data_ptr i32
    mov eax, dword ptr [rbp-48]
    mov dword ptr [rbp-40], eax

    jmp .L.next.free_next0000000A
.L.next.0000000B:
.L.if.end.0000000B:
.L.loop.end.free_next0000000A:
.L.next.00000009:
.L.if.end.00000009:
.L.next.00000008:
.L.if.end.00000008:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # call runtime.HeapFree(...)
    mov ecx, dword ptr [rbp-48] # arg 0
    call .F.runtime.HeapFree
.L.next.00000007:
.L.if.end.00000007:

    # 根据ABI处理返回值
.L.return.runtime.Block.Release:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func $wa.runtime.i32_ref_to_ptr(b:i32,d:i32) => i32
.section .text
.F.$wa.runtime.i32_ref_to_ptr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get d i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax


    # 根据ABI处理返回值
.L.return.$wa.runtime.i32_ref_to_ptr:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func $wa.runtime.i64_ref_to_ptr(b:i32,d:i32) => i32
.section .text
.F.$wa.runtime.i64_ref_to_ptr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get d i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax


    # 根据ABI处理返回值
.L.return.$wa.runtime.i64_ref_to_ptr:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func $wa.runtime.slice_to_ptr(b:i32,d:i32,l:i32,c:i32) => i32
.section .text
.F.$wa.runtime.slice_to_ptr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1
    mov [rbp+32], r8d # save arg.2
    mov [rbp+40], r9d # save arg.3

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get d i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax


    # 根据ABI处理返回值
.L.return.$wa.runtime.slice_to_ptr:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.malloc(size:i32) => i32
.section .text
.F.runtime.malloc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax


    # 根据ABI处理返回值
.L.return.runtime.malloc:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.free(ptr:i32)
.section .text
.F.runtime.free:
    push rbp
    mov  rbp, rsp
    sub  rsp, 0

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    jmp .L.return.runtime.free

    # 根据ABI处理返回值
.L.return.runtime.free:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

