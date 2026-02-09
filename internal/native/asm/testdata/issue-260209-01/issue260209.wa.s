# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# void _Wa_Import_env_print_char(int64_t x)
.section .text
.globl .Wa.Import.env.print_char
.Wa.Import.env.print_char:
    b        .Wa.Import.syscall_linux.print_rune

# 汇编程序入口函数
.section .text
.globl _start
_start:
    pcalau12i $s0, %pc_hi20(.Wa.Import.env.print_char)
    addi.d    $s0, $s0, %pc_lo12(.Wa.Import.env.print_char)
    
    addi.d    $a0, $zero, 49 # '1'
    jirl      $ra, $s0, 0

    addi.d    $a0, $zero, 50 # '2'
    jirl      $ra, $s0, 0

    addi.d    $a0, $zero, 51 # '3'
    jirl      $ra, $s0, 0

    # runtime.exit(0)
    addi.d    $a0, $zero, 0 # a0 = 0
    addi.d  $a7, $zero, 93 # sys_exit
    syscall 0

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
