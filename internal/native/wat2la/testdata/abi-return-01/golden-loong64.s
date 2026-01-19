# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

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
.extern wat2la_env_get_multi_values
.extern wat2la_env_print_i64
.set .Import.env.get_multi_values, wat2la_env_get_multi_values
.set .Import.env.print_i64, wat2la_env_print_i64

# 定义内存
.section .data
.align 3
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
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # 分配内存
    pcalau12i $t0, %pc_hi20(.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a0, $t0, 16
    pcalau12i $t0, %pc_hi20(.Runtime.malloc)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.malloc)
    jirl      $ra, $t0, 0
    pcalau12i $t1, %pc_hi20(.Memory.addr)
    addi.d    $t1, $t1, %pc_lo12(.Memory.addr)
    st.d      $a0, $t1, 0

    # 内存清零
    addi.d    $a1, $zero, 0 # a1 = 0
    pcalau12i $t0, %pc_hi20(.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a2, $t0, 16
    pcalau12i $t0, %pc_hi20(.Runtime.memset)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.memset)
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

    pcalau12i $t0, %pc_hi20(.Memory.initFunc)
    addi.d    $t0, $t0, %pc_lo12(.Memory.initFunc)
    jirl      $ra, $t0, 0
    pcalau12i $t0, %pc_hi20(.F.main)
    addi.d    $t0, $t0, %pc_lo12(.F.main)
    jirl      $ra, $t0, 0

    # runtime.exit(0)
    addi.d    $a0, $zero, 0 # a0 = 0
    pcalau12i $t0, %pc_hi20(.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.exit)
    jirl      $ra, $t0, 0

    # exit 后这里不会被执行, 但是依然保留
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

.section .data
.align 3
.Runtime.panic.message: .asciz "panic"
.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Runtime.panic
.Runtime.panic:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # runtime.write(stderr, panicMessage, size)
    addi.d    $a0, $zero, 2
    pcalau12i $a1, %pc_hi20(.Runtime.panic.message)
    addi.d    $a1, $a1, %pc_lo12(.Runtime.panic.message)
    pcalau12i $t0, %pc_hi20(.Runtime.panic.messageLen)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.panic.messageLen)
    ld.d      $a2, $t0, 0
    pcalau12i $t0, %pc_hi20(.Runtime.write)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.write)
    jirl      $ra, $t0, 0

    # 退出程序
    addi.d    $a0, $zero, 1 # 退出码
    pcalau12i $t0, %pc_hi20(.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.exit)
    jirl      $ra, $t0, 0

    # return
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# func main
.section .text
.F.main:
    # local input: i64

    addi.d  $sp, $sp, -16 
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    st.d   $zero, $fp, -8 # local input = 0

    # i64.const 100
    addi.d $t0, $zero, 100
    st.d   $t0, $fp, -16

    # local.set input
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -8

    # local.get input
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -16

    # call env.get_multi_values(...)
    addi.d $a0, $sp, 0 # return address
    ld.d $a1, $fp, -16 # arg 0
    pcalau12i $t0, %pc_hi20(.Import.env.get_multi_values)
    addi.d $t0, $t0, %pc_lo12(.Import.env.get_multi_values)
    jirl $ra, $t0, 0
    ld.d $t0, $a0, 0
    st.d $t0, $fp, -16
    ld.d $t0, $a0, 8
    st.d $t0, $fp, -24
    ld.d $t0, $a0, 16
    st.d $t0, $fp, -32
    # call env.print_i64(...)
    ld.d $a0, $fp, -32 # arg 0
    pcalau12i $t0, %pc_hi20(.Import.env.print_i64)
    addi.d $t0, $t0, %pc_lo12(.Import.env.print_i64)
    jirl $ra, $t0, 0
    # call env.print_i64(...)
    ld.d $a0, $fp, -24 # arg 0
    pcalau12i $t0, %pc_hi20(.Import.env.print_i64)
    addi.d $t0, $t0, %pc_lo12(.Import.env.print_i64)
    jirl $ra, $t0, 0
    # call env.print_i64(...)
    ld.d $a0, $fp, -16 # arg 0
    pcalau12i $t0, %pc_hi20(.Import.env.print_i64)
    addi.d $t0, $t0, %pc_lo12(.Import.env.print_i64)
    jirl $ra, $t0, 0

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

