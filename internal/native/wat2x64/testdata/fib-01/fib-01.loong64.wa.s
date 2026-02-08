# 源文件: fib-01.wat, ABI: loong64
# 自动生成的代码, 不要手动修改!!!

# 运行时函数
# .extern .Wa.Runtime.write
# .extern .Wa.Runtime.exit
# .extern .Wa.Runtime.malloc
# .extern .Wa.Runtime.memcpy
# .extern .Wa.Runtime.memset
# .extern .Wa.Runtime.memmove

# 导入函数(外部库定义)

# 定义内存
.section .data
.align 3
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
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    addi.d $sp, $sp, -32

    # 分配内存
    pcalau12i $t0, %pc_hi20(.Wa.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a0, $t0, 16
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.malloc)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.malloc)
    jirl      $ra, $t0, 0
    pcalau12i $t1, %pc_hi20(.Wa.Memory.addr)
    addi.d    $t1, $t1, %pc_lo12(.Wa.Memory.addr)
    st.d      $a0, $t1, 0

    # 内存清零
    addi.d    $a1, $zero, 0 # a1 = 0
    pcalau12i $t0, %pc_hi20(.Wa.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a2, $t0, 16
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.memset)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.memset)
    jirl      $ra, $t0, 0

    # 函数返回
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

# 定义表格
.section .data
.align 3
.globl .Wa.Table.addr
.Wa.Table.addr: .quad 0
.globl .Wa.Table.size
.Wa.Table.size: .quad 1
.globl .Wa.Table.maxSize
.Wa.Table.maxSize: .quad 1

# 汇编程序入口函数
.section .text
.globl _start
_start:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    addi.d $sp, $sp, -32

    pcalau12i $t0, %pc_hi20(.Wa.Memory.initFunc)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Memory.initFunc)
    jirl      $ra, $t0, 0
    pcalau12i $t0, %pc_hi20(.Wa.F.main)
    addi.d    $t0, $t0, %pc_lo12(.Wa.F.main)
    jirl      $ra, $t0, 0

    # runtime.exit(0)
    addi.d    $a0, $zero, 0 # a0 = 0
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.exit)
    jirl      $ra, $t0, 0

    # exit 后这里不会被执行, 但是依然保留
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.extern .Wa.Memory.addr

# void _Wa_Import_syscall_write (int fd, uint32_t ptr, int32_t len)
.section .text
.globl .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    pcalau12i $t0, %pc_hi20(.Wa.Memory.addr)
    ld.d      $t0, $t0, %pc_lo12(.Wa.Memory.addr)

    # a0 = fd
    add.d $a1, $a1, $t0 # a1 = base + ptr
    # a2 = len
    addi.d  $a7, $zero, 64 # sys_write
    syscall 0
    jirl    $zero, $ra, 0


# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    addi.d  $a7, $zero, 64 # sys_write
    syscall 0
    jirl    $zero, $ra, 0


# void _Wa_Runtime_exit(int status)
.section .text
.globl .Wa.Runtime.exit
.Wa.Runtime.exit:
    addi.d  $a7, $zero, 93 # sys_exit
    syscall 0


# void* _Wa_Runtime_malloc(int size)
.section .text
.globl .Wa.Runtime.malloc
.Wa.Runtime.malloc:
    or      $a1, $a0, $zero # length = size
    addi.d  $a0, $zero, 0   # addr = NULL
    addi.d  $a2, $zero, 3   # prot = PROT_READ | PROT_WRITE
    addi.d  $a3, $zero, 34  # flags = MAP_PRIVATE | MAP_ANONYMOUS
    addi.d  $a4, $zero, -1  # fd = -1
    addi.d  $a5, $zero, 0   # offset = 0
    addi.d  $a7, $zero, 222 # sys_mmap (222)
    syscall 0
    jirl    $zero, $ra, 0


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
    or  $t0, $a0, $zero # 备份 dst
    beq $a2, $zero, .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    ld.b   $t1, $a1, 0 # 字节读取
    st.b   $t1, $a0, 0 # 字节写入
    addi.d $a0, $a0, 1
    addi.d $a1, $a1, 1
    addi.d $a2, $a2, -1
    bne    $a2, $zero, .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    or   $a0, $t0, $zero # 返回 dst
    jirl $zero, $ra, 0


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memmove
.Wa.Runtime.memmove:
    beq $a0, $a1, .Wa.L.memmove.done
    # 如果 a0 < a1 (无符号比较), 跳转到向前拷贝
    bltu $a0, $a1, .Wa.Runtime.memcpy

    # 后向拷贝 (dst > src)
    or    $t0, $a0, $zero
    add.d $a0, $a0, $a2
    add.d $a1, $a1, $a2
.Wa.L.memmove.back_loop:
    beq    $a2, $zero, .Wa.L.memmove.ret
    addi.d $a0, $a0, -1
    addi.d $a1, $a1, -1
    ld.b   $t1, $a1, 0
    st.b   $t1, $a0, 0
    addi.d $a2, $a2, -1
    b      .Wa.L.memmove.back_loop
.Wa.L.memmove.ret:
    or $a0, $t0, $zero
.Wa.L.memmove.done:
    jirl $zero, $ra, 0


# void* _Wa_Runtime_memset(void* s, int c, int n)
.section .text
.globl .Wa.Runtime.memset
.Wa.Runtime.memset:
    or  $t0, $a0, $zero
    beq $a2, $zero, .Wa.L.memset.done
.Wa.L.memset.loop:
    st.b   $a1, $a0, 0
    addi.d $a0, $a0, 1
    addi.d $a2, $a2, -1
    bne    $a2, $zero, .Wa.L.memset.loop
.Wa.L.memset.done:
    or   $a0, $t0, $zero
    jirl $zero, $ra, 0


# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
.section .text
.globl .Wa.Import.syscall_linux.print_str
.Wa.Import.syscall_linux.print_str:
    pcalau12i $t0, %pc_hi20(.Wa.Memory.addr)
    ld.d      $t0, $t0, %pc_lo12(.Wa.Memory.addr)

    # a0 = ptr, a1 = len
    or      $a2, $a1, $zero # count = len
    add.d   $a1, $t0, $a0   # buf = base + ptr
    addi.d  $a0, $zero, 1   # fd = stdout (1)
    addi.d  $a7, $zero, 64  # sys_write
    syscall 0
    jirl    $zero, $ra, 0


# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
.section .text
.globl .Wa.Import.syscall_linux.proc_exit
.Wa.Import.syscall_linux.proc_exit:
    b .Wa.Runtime.exit


# void _Wa_Import_syscall_linux_print_rune(int32_t c)
.section .text
.globl .Wa.Import.syscall_linux.print_rune
.Wa.Import.syscall_linux.print_rune:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    addi.d $sp, $sp, -16

    st.b $a0, $sp, 0

    addi.d  $a0, $zero, 1  # arg.0: stdout
    addi.d  $a1, $sp, 0    # arg.1: buffer
    addi.d  $a2, $zero, 1  # arg.2: count
    addi.d  $a7, $zero, 64 # sys_write
    syscall 0

    addi.d $sp, $fp, 0
    ld.d   $fp, $sp, 0
    ld.d   $ra, $sp, 8
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0


# void _Wa_Import_syscall_linux_print_i64(int64_t val)
.section .text
.globl .Wa.Import.syscall_linux.print_i64
.Wa.Import.syscall_linux.print_i64:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    addi.d $sp, $sp, -32

    or     $t0, $zero, $a0 # t0 = 工作变量 (val)
    addi.d $t1, $fp, -1    # t1 为缓冲区指针 (从后往前填)
    addi.d $t2, $zero, 10  # 除数

    # 1. 处理负数
    bge   $t0, $zero, .Wa.L.syscall_linux.print_i64.convert
    sub.d $t0, $zero, $t0 # t0 = abs(t0)

.Wa.L.syscall_linux.print_i64.convert:
    div.d  $t3, $t0, $t2                              # t3 = 商
    mod.d  $t4, $t0, $t2                              # t4 = 余数
    addi.w $t4, $t4, 48                               # 加上 '0' 的 ASCII 码
    st.b   $t4, $t1, 0                                # 存入缓冲区
    addi.d $t1, $t1, -1                               # 指针前移
    or     $t0, $zero, $t3                            # 更新待处理的数字
    bnez   $t0, .Wa.L.syscall_linux.print_i64.convert # 如果商不为 0 则继续

    # 2. 补负号
    bge    $a0, $zero, .Wa.L.syscall_linux.print_i64.print
    addi.d $t4, $zero, 45 # '-'
    st.b   $t4, $t1, 0
    addi.d $t1, $t1, -1

.Wa.L.syscall_linux.print_i64.print:
    addi.d  $a0, $zero, 1  # arg.0: stdout
    addi.d  $a1, $t1, 1    # arg.1: buffer
    sub.d   $a2, $fp, $a1  # arg.2: count
    addi.d  $a7, $zero, 64 # sys_write
    syscall 0

    addi.d $sp, $fp, 0
    ld.d   $fp, $sp, 0
    ld.d   $ra, $sp, 8
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0


.section .data
.align 3
.Wa.Runtime.panic.message: .ascii "panic"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    addi.d $sp, $sp, -32

    # runtime.write(stderr, panicMessage, size)
    addi.d    $a0, $zero, 2
    pcalau12i $a1, %pc_hi20(.Wa.Runtime.panic.message)
    addi.d    $a1, $a1, %pc_lo12(.Wa.Runtime.panic.message)
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.panic.messageLen)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.panic.messageLen)
    ld.d      $a2, $t0, 0
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.write)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.write)
    jirl      $ra, $t0, 0

    # 退出程序
    addi.d    $a0, $zero, 1 # 退出码
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.exit)
    jirl      $ra, $t0, 0

    # return
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

# func main
.section .text
.Wa.F.main:
    # local N: i64
    # local i: i64

    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    # $sp = $sp - 32
    addi.d $sp, $sp, -32

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    st.d   $zero, $fp, -8 # local N = 0
    st.d   $zero, $fp, -16 # local i = 0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 10
    addi.d $t0, $zero, 10
    st.d   $t0, $fp, -24

   # local.set N
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -8

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -24

   # local.set i
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16

    # loop.begin(name=my_loop, suffix=00000001)
.Wa.L.brNext.my_loop.00000001:
    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # call fib(...)
    ld.d $a0, $fp, -24
    pcalau12i $t0, %pc_hi20(.Wa.F.fib)
    addi.d    $t0, $t0, %pc_lo12(.Wa.F.fib)
    jirl      $ra, $t0, 0
    st.d $a0, $fp, -24

    # call env.print_i64(...)
    ld.d $a0, $fp, -24
    pcalau12i $t0, %pc_hi20(.Wa.Import.env.print_i64)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Import.env.print_i64)
    jirl      $ra, $t0, 0

    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -32

    # i64.add
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    add.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

   # local.set i
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16

    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # local.get N
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -32

    # i64.lt_u
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sltu    $t1, $t0, $t1
    st.w    $t1, $fp, -24

    # br_if my_loop.00000001
    ld.w $t0, $fp, -24
    beqz $t0, .Wa.L.brFallthrough.my_loop.00000001
    b .Wa.L.brNext.my_loop.00000001
.Wa.L.brFallthrough.my_loop.00000001:

    # loop.end(name=my_loop, suffix=00000001)

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

# func fib(n:i64) => i64
.section .text
.globl .Wa.F.fib
.Wa.F.fib:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    # $sp = $sp - 48
    addi.d $sp, $sp, -48

    # 将寄存器参数备份到栈
    st.d $a0, $fp, -8 # save arg.0

    # 将返回值变量初始化为0
    st.d $zero, $fp, -16 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=fib, suffix=00000002)

    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -24

    # i64.const 2
    addi.d $t0, $zero, 2
    st.d   $t0, $fp, -32

    # i64.le_u
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sltu    $t1, $t1, $t0
    xori    $t1, $t1, 1
    st.w    $t1, $fp, -24

    # if.begin(name=, suffix=00000003)
    ld.w $t0, $fp, -24
    beqz $t0, .Wa.L.else.00000003
    # if.body(name=, suffix=00000003)
    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -24

    b .Wa.L.brNext.00000003

.Wa.L.else.00000003:
    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -24

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -32

    # i64.sub
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sub.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

    # call fib(...)
    ld.d $a0, $fp, -24
    pcalau12i $t0, %pc_hi20(.Wa.F.fib)
    addi.d    $t0, $t0, %pc_lo12(.Wa.F.fib)
    jirl      $ra, $t0, 0
    st.d $a0, $fp, -24

    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -32

    # i64.const 2
    addi.d $t0, $zero, 2
    st.d   $t0, $fp, -40

    # i64.sub
    ld.d    $t0, $fp, -32
    ld.d    $t1, $fp, -40
    sub.d   $t0, $t0, $t1
    st.d    $t0, $fp, -32

    # call fib(...)
    ld.d $a0, $fp, -32
    pcalau12i $t0, %pc_hi20(.Wa.F.fib)
    addi.d    $t0, $t0, %pc_lo12(.Wa.F.fib)
    jirl      $ra, $t0, 0
    st.d $a0, $fp, -32

    # i64.add
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    add.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

.Wa.L.brNext.00000003:
    # if.end(name=, suffix=00000003)

.Wa.L.brNext.fib.00000002:
    # fn.body.end(name=fib, suffix=00000002)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16 # ret.0

    # 将返回值变量复制到寄存器
    ld.d $a0, $fp, -16 # ret .Wa.F.ret.0

    # 函数返回
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

