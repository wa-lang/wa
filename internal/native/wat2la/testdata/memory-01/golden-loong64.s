# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 运行时函数
.extern .Wa.Runtime.write
.extern .Wa.Runtime.exit
.extern .Wa.Runtime.malloc
.extern .Wa.Runtime.memcpy
.extern .Wa.Runtime.memset

# 导入函数(外部库定义)
.extern .Wa.Import.syscall.write

# 定义内存
.section .data
.align 3
.globl .Wa.Memory.addr
.globl .Wa.Memory.pages
.globl .Wa.Memory.maxPages
.Wa.Memory.addr: .quad 0
.Wa.Memory.pages: .quad 1
.Wa.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 3
# memcpy(&Memory[8], data[0], size)
.Wa.Memory.dataOffset.0: .quad 8
.Wa.Memory.dataSize.0: .quad 12
.Wa.Memory.dataPtr.0: .asciz "hello world\n"

# 内存初始化函数
.section .text
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

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
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# 汇编程序入口函数
.section .text
.globl main
main:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    pcalau12i $t0, %pc_hi20(.Wa.Memory.initFunc)
    addi.d    $t0, $t0, %pc_lo12(.Memory.initFunc)
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
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

.section .data
.align 3
.Wa.Runtime.panic.message: .asciz "panic"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

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
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# func main
.section .text
.Wa.F.main:
    addi.d  $sp, $sp, -16 
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

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
    ld.d $a0, $fp, -8 # arg 0
    ld.d $a1, $fp, -16 # arg 1
    ld.d $a2, $fp, -24 # arg 2
    pcalau12i $t0, %pc_hi20(.Wa.Import.syscall.write)
    addi.d $t0, $t0, %pc_lo12(.Wa.Import.syscall.write)
    jirl $ra, $t0, 0
    st.d $a0, $fp, -8
    addi.w $zero, $zero, 0 # drop [fp-8]

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

