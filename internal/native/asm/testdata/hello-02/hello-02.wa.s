# 源文件: memory-01.wat, ABI: loong64
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

# 内存数据
.section .data
.align 3
# memcpy(&Memory[8], data[0], size)
.Wa.Memory.dataOffset.0: .quad 8
.Wa.Memory.dataSize.0: .quad 12
.Wa.Memory.dataPtr.0: .ascii "hello world\012"

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

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    pcalau12i $t1, %pc_hi20(.Wa.Memory.addr)
    addi.d    $t1, $t1, %pc_lo12(.Wa.Memory.addr)
    ld.d      $t1, $t1, 0
    pcalau12i $t0, %pc_hi20(.Wa.Memory.dataOffset.0)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Memory.dataOffset.0)
    ld.d      $t0, $t0, 0
    add.d     $a0, $t1, $t0
    pcalau12i $a1, %pc_hi20(.Wa.Memory.dataPtr.0)
    addi.d    $a1, $a1, %pc_lo12(.Wa.Memory.dataPtr.0)
    pcalau12i $t0, %pc_hi20(.Wa.Memory.dataSize.0)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Memory.dataSize.0)
    ld.d      $a2, $t0, 0
    pcalau12i $t0, %pc_hi20(.Wa.Runtime.memcpy)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Runtime.memcpy)
    jirl      $ra, $t0, 0

    # 函数返回
    addi.d $sp, $fp, 0
    ld.d   $ra, $sp, 8
    ld.d   $fp, $sp, 0
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0


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


# func main
.section .text
.Wa.F.main:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $fp, $sp, 0
    addi.d $fp, $sp, 0
    # $sp = $sp - 32
    addi.d $sp, $sp, -32

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -8

    # i64.const 8
    addi.d $t0, $zero, 8
    st.d   $t0, $fp, -16

    # i64.const 12
    addi.d $t0, $zero, 12
    st.d   $t0, $fp, -24

    # call syscall.write(...)
    ld.d $a0, $fp, -8
    ld.d $a1, $fp, -16
    ld.d $a2, $fp, -24
    pcalau12i $t0, %pc_hi20(.Wa.Import.syscall.write)
    addi.d    $t0, $t0, %pc_lo12(.Wa.Import.syscall.write)
    jirl      $ra, $t0, 0
    st.d $a0, $fp, -8

    addi.w $zero, $zero, 0 # drop [fp-8]

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

