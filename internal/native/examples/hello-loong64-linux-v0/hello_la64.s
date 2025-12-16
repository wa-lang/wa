# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.data
msg:
    .asciz "Hello, LoongArch 64!\n"
len = . - msg

.equ STDOUT, 1
.equ SYS_write, 64
.equ SYS_exit, 93

.text
.align 2
.globl _start

_start:
    # SYS_write(STDOUT, msg, len)
    addi.d    $a0, $zero, STDOUT
    pcalau12i $a1, %pc_hi20(msg)
    addi.d    $a1, $a1, %pc_lo12(msg)
    addi.d    $a2, $zero, len
    addi.d    $a7, $zero, SYS_write
    syscall   0

    # SYS_exit(0)
    addi.d    $a0, $zero, 0
    addi.d    $a7, $zero, SYS_exit
    syscall   0

