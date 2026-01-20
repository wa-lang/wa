# 源文件: wa-runtime-04.wat, ABI: x64-Windows
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
.Table.funcIndexList.20: .quad .F.heap_assert_valid_ptr
.Table.funcIndexList.21: .quad .F.heap_is_fixed_list_enabled
.Table.funcIndexList.22: .quad .F.heap_assert_fixed_list_enabled
.Table.funcIndexList.23: .quad .F.heap_is_fixed_size
.Table.funcIndexList.24: .quad .F.heap_alignment8
.Table.funcIndexList.25: .quad .F.heap_assert_align8
.Table.funcIndexList.26: .quad .F.heap_block.init
.Table.funcIndexList.27: .quad .F.heap_block.set_size
.Table.funcIndexList.28: .quad .F.heap_block.set_next
.Table.funcIndexList.29: .quad .F.heap_block.size
.Table.funcIndexList.30: .quad .F.heap_block.next
.Table.funcIndexList.31: .quad .F.heap_block.data
.Table.funcIndexList.32: .quad .F.heap_block.end
.Table.funcIndexList.33: .quad .F.heap_free_list.ptr_and_fixed_size
.Table.funcIndexList.34: .quad .F.wa_malloc_init_once
.Table.funcIndexList.35: .quad .F.runtime.malloc
.Table.funcIndexList.36: .quad .F.wa_malloc_reuse_fixed
.Table.funcIndexList.37: .quad .F.heap_reuse_varying
.Table.funcIndexList.38: .quad .F.heap_new_allocation
.Table.funcIndexList.39: .quad .F.runtime.free
.Table.funcIndexList.40: .quad .F.wa_lfixed_free_block
.Table.funcIndexList.41: .quad .F.wa_lfixed_free_all
.Table.funcIndexList.42: .quad .F.wa_l128_free
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
.G.__heap_ptr: .long 0
.G.__heap_top: .long 0
.G.__heap_l128_freep: .long 0
.G.__heap_init_flag: .long 0

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

# func heap_assert_valid_ptr(ptr:i32)
.section .text
.F.heap_assert_valid_ptr:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000000C:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.0000000C # if eax != 0 { jmp body }
.L.if.body.0000000C:
.L.if.else.0000000C:
    call .Runtime.panic # unreachable
.L.next.0000000C:
.L.if.end.0000000C:
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 4
    mov eax, 4
    mov [rbp-16], eax

    # i32.rem_s
    mov dword ptr [rbp+8], edx
    mov eax, dword ptr [rbp-16]
    cdq
    idiv dword ptr [rbp-8]
    mov dword ptr [rbp-8], edx
    mov edx, dword ptr [rbp+8]
    # i32.eqz
    mov  eax, dword ptr [rbp-8]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000000D:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.0000000D # if eax != 0 { jmp body }
.L.if.body.0000000D:
.L.if.else.0000000D:
    call .Runtime.panic # unreachable
.L.next.0000000D:
.L.if.end.0000000D:

    # 根据ABI处理返回值
.L.return.heap_assert_valid_ptr:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_is_fixed_list_enabled => i32
.section .text
.F.heap_is_fixed_list_enabled:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # global.get __heap_lfixed_cap i32
    mov eax, dword ptr [rip+.G.__heap_lfixed_cap]
    mov dword ptr [rbp-16], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-16]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000000E:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.0000000E # if eax != 0 { jmp body }
.L.if.body.0000000E:
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

.L.if.else.0000000E:
    # i32.const 1
    mov eax, 1
    mov [rbp-16], eax

.L.next.0000000E:
.L.if.end.0000000E:

    # 根据ABI处理返回值
.L.return.heap_is_fixed_list_enabled:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_assert_fixed_list_enabled
.section .text
.F.heap_assert_fixed_list_enabled:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # global.get __heap_lfixed_cap i32
    mov eax, dword ptr [rip+.G.__heap_lfixed_cap]
    mov dword ptr [rbp-8], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-8]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000000F:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.0000000F # if eax != 0 { jmp end }
.L.if.body.0000000F:
    call .Runtime.panic # unreachable
.L.next.0000000F:
.L.if.end.0000000F:

    # 根据ABI处理返回值
.L.return.heap_assert_fixed_list_enabled:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_is_fixed_size(size:i32) => i32
.section .text
.F.heap_is_fixed_size:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # call heap_is_fixed_list_enabled(...)
    call .F.heap_is_fixed_list_enabled
    mov dword ptr [rbp-16], eax
.L.if.begin.00000010:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.00000010 # if eax != 0 { jmp body }
.L.if.body.00000010:
.L.if.else.00000010:
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_is_fixed_size
.L.next.00000010:
.L.if.end.00000010:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.const 80
    mov eax, 80
    mov [rbp-24], eax

    # i32.le_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setle al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.00000011:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.00000011 # if eax != 0 { jmp body }
.L.if.body.00000011:
    # i32.const 1
    mov eax, 1
    mov [rbp-16], eax

.L.if.else.00000011:
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

.L.next.00000011:
.L.if.end.00000011:

    # 根据ABI处理返回值
.L.return.heap_is_fixed_size:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_alignment8(size:i32) => i32
.section .text
.F.heap_alignment8:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.const 7
    mov eax, 7
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.div_s
    mov dword ptr [rbp+8], edx
    mov eax, dword ptr [rbp-24]
    cdq
    idiv dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    mov edx, dword ptr [rbp+8]
    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.mul
    mov eax, dword ptr [rbp-24]
    imul eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.heap_alignment8:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_assert_align8(size:i32)
.section .text
.F.heap_assert_align8:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-16], eax

    # i32.rem_s
    mov dword ptr [rbp+8], edx
    mov eax, dword ptr [rbp-16]
    cdq
    idiv dword ptr [rbp-8]
    mov dword ptr [rbp-8], edx
    mov edx, dword ptr [rbp+8]
    # i32.eqz
    mov  eax, dword ptr [rbp-8]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000012:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.00000012 # if eax != 0 { jmp body }
.L.if.body.00000012:
.L.if.else.00000012:
    call .Runtime.panic # unreachable
.L.next.00000012:
.L.if.end.00000012:

    # 根据ABI处理返回值
.L.return.heap_assert_align8:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.init(ptr:i32,size:i32,next:i32)
.section .text
.F.heap_block.init:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1
    mov [rbp+32], r8d # save arg.2

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get size i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-8]
    add r10, rax
    mov eax, dword [rbp-16]
    mov dword [r10 +0], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get next i32
    mov eax, dword ptr [rbp+32]
    mov dword ptr [rbp-16], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-8]
    add r10, rax
    mov eax, dword [rbp-16]
    mov dword [r10 +4], eax

    # 根据ABI处理返回值
.L.return.heap_block.init:

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

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get size i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-8]
    add r10, rax
    mov eax, dword [rbp-16]
    mov dword [r10 +0], eax

    # 根据ABI处理返回值
.L.return.heap_block.set_size:

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

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get next i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # i32.store
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-8]
    add r10, rax
    mov eax, dword [rbp-16]
    mov dword [r10 +4], eax

    # 根据ABI处理返回值
.L.return.heap_block.set_next:

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

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-16]
    add r10, rax
    mov eax, dword [r10 +0]
    mov dword [rbp-16], eax

    # 根据ABI处理返回值
.L.return.heap_block.size:
    mov eax, [rbp-8] # ret .F.ret.0

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

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-16]
    add r10, rax
    mov eax, dword [r10 +4]
    mov dword [rbp-16], eax

    # 根据ABI处理返回值
.L.return.heap_block.next:
    mov eax, [rbp-8] # ret .F.ret.0

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

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.heap_block.data:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_block.end(ptr:i32) => i32
.section .text
.F.heap_block.end:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.load
    lea rax, [rip + .Memory.addr]
    mov r10, dword [rbp-16]
    add r10, rax
    mov eax, dword [r10 +0]
    mov dword [rbp-16], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax

    # 根据ABI处理返回值
.L.return.heap_block.end:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_free_list.ptr_and_fixed_size(size:i32) => i32,i32)
.section .text
.F.heap_free_list.ptr_and_fixed_size:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将返回地址备份到栈
    mov [rbp+16], rcx # return address

    # 将寄存器参数备份到栈
    mov [rbp+16], edx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp+48], 0 # ret.0 = 0
    mov dword ptr [rbp+56], 0 # ret.1 = 0

    # 没有局部变量需要初始化为0

    # call heap_is_fixed_list_enabled(...)
    call .F.heap_is_fixed_list_enabled
    mov dword ptr [rbp-8], eax
.L.if.begin.00000013:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.00000013 # if eax != 0 { jmp body }
.L.if.body.00000013:
.L.if.else.00000013:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 32
    mov eax, 32
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # call heap_alignment8(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.heap_alignment8
    mov dword ptr [rbp-16], eax
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size
.L.next.00000013:
.L.if.end.00000013:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 80
    mov eax, 80
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000014:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000014 # if eax != 0 { jmp end }
.L.if.body.00000014:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 32
    mov eax, 32
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
.L.block.begin.00000015:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # i32.const 128
    mov eax, 128
    mov [rbp-24], eax

    # i32.le_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setle al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.00000016:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.00000016 # if eax != 0 { jmp body }
.L.if.body.00000016:
    # i32.const 128
    mov eax, 128
    mov [rbp-16], eax

.L.if.else.00000016:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # call heap_alignment8(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.heap_alignment8
    mov dword ptr [rbp-16], eax
.L.next.00000016:
.L.if.end.00000016:
.L.next.00000015:
.L.block.end.00000015:
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size
.L.next.00000014:
.L.if.end.00000014:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 48
    mov eax, 48
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000017:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000017 # if eax != 0 { jmp end }
.L.if.body.00000017:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 24
    mov eax, 24
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # i32.const 80
    mov eax, 80
    mov [rbp-16], eax

    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size
.L.next.00000017:
.L.if.end.00000017:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 32
    mov eax, 32
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000018:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000018 # if eax != 0 { jmp end }
.L.if.body.00000018:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 16
    mov eax, 16
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # i32.const 48
    mov eax, 48
    mov [rbp-16], eax

    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size
.L.next.00000018:
.L.if.end.00000018:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # i32.const 24
    mov eax, 24
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000019:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000019 # if eax != 0 { jmp end }
.L.if.body.00000019:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # i32.const 32
    mov eax, 32
    mov [rbp-16], eax

    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size
.L.next.00000019:
.L.if.end.00000019:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 24
    mov eax, 24
    mov [rbp-16], eax

    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp+56], eax
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp+48], eax
    mov    rax, [rbp+16] # return address
    jmp    .L.return.heap_free_list.ptr_and_fixed_size

    # 根据ABI处理返回值
.L.return.heap_free_list.ptr_and_fixed_size:
    mov rax, [rbp+48] # ret address

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func wa_malloc_init_once
.section .text
.F.wa_malloc_init_once:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # global.get __heap_init_flag i32
    mov eax, dword ptr [rip+.G.__heap_init_flag]
    mov dword ptr [rbp-8], eax
.L.if.begin.0000001A:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.0000001A # if eax != 0 { jmp end }
.L.if.body.0000001A:
    jmp .L.return.wa_malloc_init_once
.L.next.0000001A:
.L.if.end.0000001A:
    # i32.const 1
    mov eax, 1
    mov [rbp-8], eax

    # global.set __heap_init_flag i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__heap_init_flag], eax
    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-8], eax
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000001B:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.0000001B # if eax != 0 { jmp body }
.L.if.body.0000001B:
.L.if.else.0000001B:
    call .Runtime.panic # unreachable
.L.next.0000001B:
.L.if.end.0000001B:
    # global.get __stack_ptr i32
    mov eax, dword ptr [rip+.G.__stack_ptr]
    mov dword ptr [rbp-8], eax
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-16], eax
    # i32.lt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setl al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000001C:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.0000001C # if eax != 0 { jmp body }
.L.if.body.0000001C:
.L.if.else.0000001C:
    call .Runtime.panic # unreachable
.L.next.0000001C:
.L.if.end.0000001C:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # call heap_assert_align8(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    call .F.heap_assert_align8
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 48
    mov eax, 48
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # global.set __heap_ptr i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__heap_ptr], eax
    # memory.size
    lea rax, [rip+.Memory.pages]
    mov [rbp-8], rax
    # i32.const 65536
    mov eax, 65536
    mov [rbp-16], eax

    # i32.mul
    mov eax, dword ptr [rbp-16]
    imul eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # global.set __heap_top i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__heap_top], eax
    # global.get __heap_top i32
    mov eax, dword ptr [rip+.G.__heap_top]
    mov dword ptr [rbp-8], eax
    # global.get __heap_ptr i32
    mov eax, dword ptr [rip+.G.__heap_ptr]
    mov dword ptr [rbp-16], eax
    # i32.gt_s
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.0000001D:
    mov eax, [rbp-8]
    test eax, eax
    jne .L.if.body.0000001D # if eax != 0 { jmp body }
.L.if.body.0000001D:
.L.if.else.0000001D:
    call .Runtime.panic # unreachable
.L.next.0000001D:
.L.if.end.0000001D:
    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

    # i32.const 48
    mov eax, 48
    mov [rbp-24], eax

    # memory.fill    # .Runtime.memset(&.Memory.addr[R0.i32], R1.i32, R2.i32);

    # global.get __heap_base i32
    mov eax, dword ptr [rip+.G.__heap_base]
    mov dword ptr [rbp-8], eax
    # i32.const 32
    mov eax, 32
    mov [rbp-16], eax

    # i32.add
    mov eax, dword ptr [rbp-16]
    add eax, dword ptr [rbp-8]
    mov dword ptr [rbp-8], eax
    # global.set __heap_l128_freep i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rip+.G.__heap_l128_freep], eax
    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-8], eax
    # i32.const 0
    mov eax, 0
    mov [rbp-16], eax

    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-24], eax
    # call heap_block.init(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    mov edx, dword ptr [rbp-16] # arg 1
    mov r8d, dword ptr [rbp-24] # arg 2
    call .F.heap_block.init

    # 根据ABI处理返回值
.L.return.wa_malloc_init_once:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.malloc(size:i32) => i32
.section .text
.F.runtime.malloc:
    # local free_list: i32
    # local b: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-16], 0 # local free_list = 0
    mov dword ptr [rbp-24], 0 # local b = 0

    # call wa_malloc_init_once(...)
    call .F.wa_malloc_init_once
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_alignment8(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_alignment8
    mov dword ptr [rbp-32], eax
    # local.set size i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp+16], eax

    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_free_list.ptr_and_fixed_size(...)
    lea rcx, [rsp+32] # return address
    mov edx, dword ptr [rbp-32] # arg 0
    call .F.heap_free_list.ptr_and_fixed_size
    mov r10d, dword ptr [rax+0]
    mov dword ptr [rbp-32], r10d
    mov r10d, dword ptr [rax+8]
    mov dword ptr [rbp-40], r10d
    # local.set size i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp+16], eax

    # local.set free_list i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-16], eax

    # call heap_is_fixed_list_enabled(...)
    call .F.heap_is_fixed_list_enabled
    mov dword ptr [rbp-32], eax
.L.if.begin.0000001E:
    mov eax, [rbp-32]
    test eax, eax
    je .L.if.end.0000001E # if eax != 0 { jmp end }
.L.if.body.0000001E:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_is_fixed_size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_is_fixed_size
    mov dword ptr [rbp-32], eax
.L.if.begin.0000001F:
    mov eax, [rbp-32]
    test eax, eax
    je .L.if.end.0000001F # if eax != 0 { jmp end }
.L.if.body.0000001F:
    # local.get free_list i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-32], eax

    # call wa_malloc_reuse_fixed(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.wa_malloc_reuse_fixed
    mov dword ptr [rbp-32], eax
    # local.tee b i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-24], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-32]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-32], eax
.L.if.begin.00000020:
    mov eax, [rbp-32]
    test eax, eax
    jne .L.if.body.00000020 # if eax != 0 { jmp body }
.L.if.body.00000020:
.L.if.else.00000020:
    # local.get b i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-32], eax

    # call heap_block.data(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.data
    mov dword ptr [rbp-32], eax
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.malloc
.L.next.00000020:
.L.if.end.00000020:
.L.next.0000001F:
.L.if.end.0000001F:
.L.next.0000001E:
.L.if.end.0000001E:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_reuse_varying(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_reuse_varying
    mov dword ptr [rbp-32], eax
    # local.tee b i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-24], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-32]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-32], eax
.L.if.begin.00000021:
    mov eax, [rbp-32]
    test eax, eax
    jne .L.if.body.00000021 # if eax != 0 { jmp body }
.L.if.body.00000021:
.L.if.else.00000021:
    # local.get b i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-32], eax

    # call heap_block.data(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.data
    mov dword ptr [rbp-32], eax
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.malloc
.L.next.00000021:
.L.if.end.00000021:
    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_new_allocation(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_new_allocation
    mov dword ptr [rbp-32], eax
    # local.tee b i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-24], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-32]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-32], eax
.L.if.begin.00000022:
    mov eax, [rbp-32]
    test eax, eax
    jne .L.if.body.00000022 # if eax != 0 { jmp body }
.L.if.body.00000022:
.L.if.else.00000022:
    # local.get b i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-32], eax

    # call heap_block.data(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.data
    mov dword ptr [rbp-32], eax
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.malloc
.L.next.00000022:
.L.if.end.00000022:
    # i32.const 0
    mov eax, 0
    mov [rbp-32], eax

    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-8], eax
    jmp .L.return.runtime.malloc

    # 根据ABI处理返回值
.L.return.runtime.malloc:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func wa_malloc_reuse_fixed(free_list:i32) => i32
.section .text
.F.wa_malloc_reuse_fixed:
    # local p: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 80

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-16], 0 # local p = 0

    # call heap_assert_fixed_list_enabled(...)
    call .F.heap_assert_fixed_list_enabled
    # local.get free_list i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-24], eax
    # i32.eqz
    mov  eax, dword ptr [rbp-24]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-24], eax
.L.if.begin.00000023:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000023 # if eax != 0 { jmp end }
.L.if.body.00000023:
    # i32.const 0
    mov eax, 0
    mov [rbp-24], eax

    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-8], eax
    jmp .L.return.wa_malloc_reuse_fixed
.L.next.00000023:
.L.if.end.00000023:
    # local.get free_list i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

.L.block.begin.00000024:
    # local.get free_list i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-32], eax
    # i32.const 1
    mov eax, 1
    mov [rbp-40], eax

    # i32.sub
    mov eax, dword ptr [rbp-40]
    sub eax, dword ptr [rbp-32]
    mov dword ptr [rbp-32], eax
.L.next.00000024:
.L.block.end.00000024:
    # call heap_block.set_size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    mov edx, dword ptr [rbp-32] # arg 1
    call .F.heap_block.set_size
    # local.get free_list i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # local.set p i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-16], eax

    # local.get free_list i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

.L.block.begin.00000025:
    # local.get p i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-32], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-32], eax
.L.next.00000025:
.L.block.end.00000025:
    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    mov edx, dword ptr [rbp-32] # arg 1
    call .F.heap_block.set_next
    # local.get p i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-24], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-32], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    mov edx, dword ptr [rbp-32] # arg 1
    call .F.heap_block.set_next
    # local.get p i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-24], eax

    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-8], eax
    jmp .L.return.wa_malloc_reuse_fixed

    # 根据ABI处理返回值
.L.return.wa_malloc_reuse_fixed:
    mov eax, [rbp-8] # ret .F.ret.0

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

.L.loop.begin.continue00000026:
.L.next.continue00000026:
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
    mov eax, dword ptr [rbp-56]
    add eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # i32.ge_s
    mov r10d, dword ptr [rbp-48]
    mov r11d, dword ptr [rbp-40]
    cmp r10d, r11d
    setge al
    movzx eax, al
    mov dword ptr [rbp-40], eax
.L.if.begin.00000027:
    mov eax, [rbp-40]
    test eax, eax
    je .L.if.end.00000027 # if eax != 0 { jmp end }
.L.if.body.00000027:
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
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
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
    mov eax, dword ptr [rbp-56]
    sub eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-56], eax

    # i32.sub
    mov eax, dword ptr [rbp-56]
    sub eax, dword ptr [rbp-48]
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

    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_reuse_varying
.L.next.00000027:
.L.if.end.00000027:
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
    mov r10d, dword ptr [rbp-48]
    mov r11d, dword ptr [rbp-40]
    cmp r10d, r11d
    setge al
    movzx eax, al
    mov dword ptr [rbp-40], eax
.L.if.begin.00000028:
    mov eax, [rbp-40]
    test eax, eax
    je .L.if.end.00000028 # if eax != 0 { jmp end }
.L.if.body.00000028:
    # local.get prevp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

.L.block.begin.00000029:
    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-48], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-48] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-48], eax
.L.next.00000029:
.L.block.end.00000029:
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

    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_reuse_varying
.L.next.00000028:
.L.if.end.00000028:
    # local.get p i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-48], eax
    # i32.eq
    mov r10d, dword ptr [rbp-48]
    mov r11d, dword ptr [rbp-40]
    cmp r10d, r11d
    sete al
    movzx eax, al
    mov dword ptr [rbp-40], eax
.L.if.begin.0000002A:
    mov eax, [rbp-40]
    test eax, eax
    je .L.if.end.0000002A # if eax != 0 { jmp end }
.L.if.body.0000002A:
    # i32.const 0
    mov eax, 0
    mov [rbp-40], eax

    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_reuse_varying
.L.next.0000002A:
.L.if.end.0000002A:
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

    jmp .L.next.continue00000026
.L.loop.end.continue00000026:
    call .Runtime.panic # unreachable

    # 根据ABI处理返回值
.L.return.heap_reuse_varying:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func heap_new_allocation(size:i32) => i32
.section .text
.F.heap_new_allocation:
    # local ptr: i32
    # local block_size: i32
    # local pages: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 将局部变量初始化为0
    mov dword ptr [rbp-16], 0 # local ptr = 0
    mov dword ptr [rbp-24], 0 # local block_size = 0
    mov dword ptr [rbp-32], 0 # local pages = 0

    # global.get __heap_ptr i32
    mov eax, dword ptr [rip+.G.__heap_ptr]
    mov dword ptr [rbp-40], eax
    # local.set ptr i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-16], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-40], eax

    # local.get size i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-48], eax

    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    # local.set block_size i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-24], eax

    # global.get __heap_ptr i32
    mov eax, dword ptr [rip+.G.__heap_ptr]
    mov dword ptr [rbp-40], eax
    # local.get block_size i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-48], eax

    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    # global.get __heap_top i32
    mov eax, dword ptr [rip+.G.__heap_top]
    mov dword ptr [rbp-48], eax
    # i32.ge_s
    mov r10d, dword ptr [rbp-48]
    mov r11d, dword ptr [rbp-40]
    cmp r10d, r11d
    setge al
    movzx eax, al
    mov dword ptr [rbp-40], eax
.L.if.begin.0000002B:
    mov eax, [rbp-40]
    test eax, eax
    je .L.if.end.0000002B # if eax != 0 { jmp end }
.L.if.body.0000002B:
    # local.get block_size i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-40], eax

    # i32.const 65535
    mov eax, 65535
    mov [rbp-48], eax

    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    # i32.const 65536
    mov eax, 65536
    mov [rbp-48], eax

    # i32.div_s
    mov dword ptr [rbp+8], edx
    mov eax, dword ptr [rbp-48]
    cdq
    idiv dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    mov edx, dword ptr [rbp+8]
    # local.set pages i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-32], eax

    # local.get pages i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-40], eax

    # memory.grow
    lea r10, [rip+.Memory.pages]
    lea r11, [rip+.Memory.maxPages]
    mov rax, [rbp-40]
    add rax, r10
    cmp rax, r11
    jg  .L.if.else.0000002C
.L.if.body.0000002C:
    mov [rbp-40], rax
    jmp .L.if.end.0000002C
.L.if.else.0000002C:
    mov rax, -1
    mov [rbp-40], rax
.L.if.end.0000002C:
    # i32.const 0
    mov eax, 0
    mov [rbp-48], eax

    # i32.lt_s
    mov r10d, dword ptr [rbp-48]
    mov r11d, dword ptr [rbp-40]
    cmp r10d, r11d
    setl al
    movzx eax, al
    mov dword ptr [rbp-40], eax
.L.if.begin.0000002D:
    mov eax, [rbp-40]
    test eax, eax
    je .L.if.end.0000002D # if eax != 0 { jmp end }
.L.if.body.0000002D:
    # i32.const 0
    mov eax, 0
    mov [rbp-40], eax

    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_new_allocation
.L.next.0000002D:
.L.if.end.0000002D:
    # global.get __heap_top i32
    mov eax, dword ptr [rip+.G.__heap_top]
    mov dword ptr [rbp-40], eax
.L.block.begin.0000002E:
    # local.get pages i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-48], eax

    # i32.const 65536
    mov eax, 65536
    mov [rbp-56], eax

    # i32.mul
    mov eax, dword ptr [rbp-56]
    imul eax, dword ptr [rbp-48]
    mov dword ptr [rbp-48], eax
.L.next.0000002E:
.L.block.end.0000002E:
    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    # global.set __heap_top i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rip+.G.__heap_top], eax
.L.next.0000002B:
.L.if.end.0000002B:
    # global.get __heap_ptr i32
    mov eax, dword ptr [rip+.G.__heap_ptr]
    mov dword ptr [rbp-40], eax
    # local.get block_size i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-48], eax

    # i32.add
    mov eax, dword ptr [rbp-48]
    add eax, dword ptr [rbp-40]
    mov dword ptr [rbp-40], eax
    # global.set __heap_ptr i32
    mov eax, dword ptr [rbp-40]
    mov dword ptr [rip+.G.__heap_ptr], eax
    # local.get ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # local.get size i32
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
    # local.get ptr i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    mov eax, dword ptr [rbp-40]
    mov dword ptr [rbp-8], eax
    jmp .L.return.heap_new_allocation

    # 根据ABI处理返回值
.L.return.heap_new_allocation:
    mov eax, [rbp-8] # ret .F.ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func runtime.free(ptr:i32)
.section .text
.F.runtime.free:
    # local size: i32
    # local block: i32
    # local freep: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 96

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local size = 0
    mov dword ptr [rbp-16], 0 # local block = 0
    mov dword ptr [rbp-24], 0 # local freep = 0

    # call wa_malloc_init_once(...)
    call .F.wa_malloc_init_once
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_assert_valid_ptr(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_assert_valid_ptr
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_assert_align8(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_assert_align8
    # local.get ptr i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # i32.const 8
    mov eax, 8
    mov [rbp-40], eax

    # i32.sub
    mov eax, dword ptr [rbp-40]
    sub eax, dword ptr [rbp-32]
    mov dword ptr [rbp-32], eax
    # local.set block i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-16], eax

    # local.get block i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-32], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-32], eax
    # local.set size i32
    mov eax, dword ptr [rbp-32]
    mov dword ptr [rbp-8], eax

    # call heap_is_fixed_list_enabled(...)
    call .F.heap_is_fixed_list_enabled
    mov dword ptr [rbp-32], eax
.L.if.begin.0000002F:
    mov eax, [rbp-32]
    test eax, eax
    je .L.if.end.0000002F # if eax != 0 { jmp end }
.L.if.body.0000002F:
    # local.get size i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-32], eax

    # call heap_is_fixed_size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_is_fixed_size
    mov dword ptr [rbp-32], eax
.L.if.begin.00000030:
    mov eax, [rbp-32]
    test eax, eax
    je .L.if.end.00000030 # if eax != 0 { jmp end }
.L.if.body.00000030:
    # local.get size i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-32], eax

    # call heap_free_list.ptr_and_fixed_size(...)
    lea rcx, [rsp+32] # return address
    mov edx, dword ptr [rbp-32] # arg 0
    call .F.heap_free_list.ptr_and_fixed_size
    mov r10d, dword ptr [rax+0]
    mov dword ptr [rbp-32], r10d
    mov r10d, dword ptr [rax+8]
    mov dword ptr [rbp-40], r10d
    nop # drop [rbp-40]
    # local.get block i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-40], eax

    # call wa_lfixed_free_block(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    mov edx, dword ptr [rbp-40] # arg 1
    call .F.wa_lfixed_free_block
    jmp .L.return.runtime.free
.L.next.00000030:
.L.if.end.00000030:
.L.next.0000002F:
.L.if.end.0000002F:
    # local.get block i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-32], eax

    # call wa_l128_free(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.wa_l128_free
    jmp .L.return.runtime.free

    # 根据ABI处理返回值
.L.return.runtime.free:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func wa_lfixed_free_block(freep:i32,block:i32)
.section .text
.F.wa_lfixed_free_block:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0
    mov [rbp+24], edx # save arg.1

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # call heap_assert_fixed_list_enabled(...)
    call .F.heap_assert_fixed_list_enabled
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-8], eax
    # global.get __heap_lfixed_cap i32
    mov eax, dword ptr [rip+.G.__heap_lfixed_cap]
    mov dword ptr [rbp-16], eax
    # i32.eq
    mov r10d, dword ptr [rbp-16]
    mov r11d, dword ptr [rbp-8]
    cmp r10d, r11d
    sete al
    movzx eax, al
    mov dword ptr [rbp-8], eax
.L.if.begin.00000031:
    mov eax, [rbp-8]
    test eax, eax
    je .L.if.end.00000031 # if eax != 0 { jmp end }
.L.if.body.00000031:
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # call wa_lfixed_free_all(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    call .F.wa_lfixed_free_all
.L.next.00000031:
.L.if.end.00000031:
    # local.get block i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-8], eax

.L.block.begin.00000032:
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-16], eax
.L.next.00000032:
.L.block.end.00000032:
    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    mov edx, dword ptr [rbp-16] # arg 1
    call .F.heap_block.set_next
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

    # local.get block i32
    mov eax, dword ptr [rbp+24]
    mov dword ptr [rbp-16], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    mov edx, dword ptr [rbp-16] # arg 1
    call .F.heap_block.set_next
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-8], eax

.L.block.begin.00000033:
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-16], eax
    # i32.const 1
    mov eax, 1
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
.L.next.00000033:
.L.block.end.00000033:
    # call heap_block.set_size(...)
    mov ecx, dword ptr [rbp-8] # arg 0
    mov edx, dword ptr [rbp-16] # arg 1
    call .F.heap_block.set_size

    # 根据ABI处理返回值
.L.return.wa_lfixed_free_block:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func wa_lfixed_free_all(freep:i32)
.section .text
.F.wa_lfixed_free_all:
    # local p: i32
    # local temp: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 80

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local p = 0
    mov dword ptr [rbp-16], 0 # local temp = 0

    # call heap_assert_fixed_list_enabled(...)
    call .F.heap_assert_fixed_list_enabled
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # local.set p i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-8], eax

.L.block.begin.break00000034:
.L.loop.begin.continue00000035:
.L.next.continue00000035:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # i32.eqz
    mov  eax, dword ptr [rbp-24]
    test eax, eax
    setz al
    movzx eax, al
    mov dword ptr [rbp-24], eax
.L.if.begin.00000036:
    mov eax, [rbp-24]
    test eax, eax
    je .L.if.end.00000036 # if eax != 0 { jmp end }
.L.if.body.00000036:
    jmp .L.next.break00000034
.L.next.00000036:
.L.if.end.00000036:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # local.set temp i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call wa_l128_free(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.wa_l128_free
    # local.get temp i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-24], eax

    # local.set p i32
    mov eax, dword ptr [rbp-24]
    mov dword ptr [rbp-8], eax

    jmp .L.next.continue00000035
.L.loop.end.continue00000035:
.L.next.break00000034:
.L.block.end.break00000034:
    # local.get freep i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-32], eax

    # i32.const 0
    mov eax, 0
    mov [rbp-40], eax

    # call heap_block.init(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    mov edx, dword ptr [rbp-32] # arg 1
    mov r8d, dword ptr [rbp-40] # arg 2
    call .F.heap_block.init

    # 根据ABI处理返回值
.L.return.wa_lfixed_free_all:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func wa_l128_free(bp:i32)
.section .text
.F.wa_l128_free:
    # local p: i32

    push rbp
    mov  rbp, rsp
    sub  rsp, 80

    # 将寄存器参数备份到栈
    mov [rbp+16], ecx # save arg.0

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local p = 0

    # global.get __heap_l128_freep i32
    mov eax, dword ptr [rip+.G.__heap_l128_freep]
    mov dword ptr [rbp-16], eax
    # local.set p i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax

.L.block.begin.break00000037:
.L.loop.begin.continue00000038:
.L.next.continue00000038:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.00000039:
    mov eax, [rbp-16]
    test eax, eax
    je .L.if.end.00000039 # if eax != 0 { jmp end }
.L.if.body.00000039:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # i32.lt_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setl al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000003A:
    mov eax, [rbp-16]
    test eax, eax
    je .L.if.end.0000003A # if eax != 0 { jmp end }
.L.if.body.0000003A:
    jmp .L.next.break00000037
.L.next.0000003A:
.L.if.end.0000003A:
.L.next.00000039:
.L.if.end.00000039:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # i32.ge_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setge al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000003B:
    mov eax, [rbp-16]
    test eax, eax
    je .L.if.end.0000003B # if eax != 0 { jmp end }
.L.if.body.0000003B:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # i32.gt_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setg al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000003C:
    mov eax, [rbp-16]
    test eax, eax
    je .L.if.end.0000003C # if eax != 0 { jmp end }
.L.if.body.0000003C:
    jmp .L.next.break00000037
.L.next.0000003C:
.L.if.end.0000003C:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # i32.lt_s
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    setl al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000003D:
    mov eax, [rbp-16]
    test eax, eax
    je .L.if.end.0000003D # if eax != 0 { jmp end }
.L.if.body.0000003D:
    jmp .L.next.break00000037
.L.next.0000003D:
.L.if.end.0000003D:
.L.next.0000003B:
.L.if.end.0000003B:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-16], eax
    # local.set p i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rbp-8], eax

    jmp .L.next.continue00000038
.L.loop.end.continue00000038:
.L.next.break00000037:
.L.block.end.break00000037:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-24], eax
    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # i32.eq
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    sete al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.0000003E:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.0000003E # if eax != 0 { jmp body }
.L.if.body.0000003E:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

.L.block.begin.0000003F:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-24], eax
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-32], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-32], eax
    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-32], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-40], eax

    # i32.add
    mov eax, dword ptr [rbp-40]
    add eax, dword ptr [rbp-32]
    mov dword ptr [rbp-32], eax
    # i32.add
    mov eax, dword ptr [rbp-32]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
.L.next.0000003F:
.L.block.end.0000003F:
    # call heap_block.set_size(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_size
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

.L.block.begin.00000040:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
.L.next.00000040:
.L.block.end.00000040:
    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_next
.L.if.else.0000003E:
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_next
.L.next.0000003E:
.L.if.end.0000003E:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-24], eax
    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-24], eax

    # i32.add
    mov eax, dword ptr [rbp-24]
    add eax, dword ptr [rbp-16]
    mov dword ptr [rbp-16], eax
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # i32.eq
    mov r10d, dword ptr [rbp-24]
    mov r11d, dword ptr [rbp-16]
    cmp r10d, r11d
    sete al
    movzx eax, al
    mov dword ptr [rbp-16], eax
.L.if.begin.00000041:
    mov eax, [rbp-16]
    test eax, eax
    jne .L.if.body.00000041 # if eax != 0 { jmp body }
.L.if.body.00000041:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

.L.block.begin.00000042:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-24], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-24], eax
    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-32], eax

    # call heap_block.size(...)
    mov ecx, dword ptr [rbp-32] # arg 0
    call .F.heap_block.size
    mov dword ptr [rbp-32], eax
    # i32.const 8
    mov eax, 8
    mov [rbp-40], eax

    # i32.add
    mov eax, dword ptr [rbp-40]
    add eax, dword ptr [rbp-32]
    mov dword ptr [rbp-32], eax
    # i32.add
    mov eax, dword ptr [rbp-32]
    add eax, dword ptr [rbp-24]
    mov dword ptr [rbp-24], eax
.L.next.00000042:
.L.block.end.00000042:
    # call heap_block.set_size(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_size
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.next(...)
    mov ecx, dword ptr [rbp-24] # arg 0
    call .F.heap_block.next
    mov dword ptr [rbp-24], eax
    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_next
.L.if.else.00000041:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # local.get bp i32
    mov eax, dword ptr [rbp+16]
    mov dword ptr [rbp-24], eax

    # call heap_block.set_next(...)
    mov ecx, dword ptr [rbp-16] # arg 0
    mov edx, dword ptr [rbp-24] # arg 1
    call .F.heap_block.set_next
.L.next.00000041:
.L.if.end.00000041:
    # local.get p i32
    mov eax, dword ptr [rbp-8]
    mov dword ptr [rbp-16], eax

    # global.set __heap_l128_freep i32
    mov eax, dword ptr [rbp-16]
    mov dword ptr [rip+.G.__heap_l128_freep], eax

    # 根据ABI处理返回值
.L.return.wa_l128_free:

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

