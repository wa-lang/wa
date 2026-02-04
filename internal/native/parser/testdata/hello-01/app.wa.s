# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.section .data
.align 3
.app.hello.str: .ascii "hello\n\000"
.app.hello.len: .quad 6

.section .text
.globl _start
_start:
    # main
    pcalau12i $t0, %pc_hi20(main)
    addi.d $t0, $t0, %pc_lo12(main)
    jirl $ra, $t0, 0

    # exit
    addi.d $a7, $zero, 93 # sys_exit
    syscall 0


.section .text
.globl main
main:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0

    # write(stdout, str, len)
    addi.d    $a0, $zero, 1 # arg.0 stdout
    pcalau12i $a1, %pc_hi20(.app.hello.str) # arg.1: ptr
    addi.d    $a1, $a1, %pc_lo12(.app.hello.str)
    pcalau12i $a2, %pc_hi20(.app.hello.len) # arg.2: len
    addi.d    $a2, $a2, %pc_lo12(.app.hello.len)
    ld.d      $a2, $a2, 0
    addi.d    $a7, $zero, 64 # sys_write
    syscall   0

    # set return value
    addi.d $a0, $zero, 0

    # return
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

