# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.data
msg:
    .asciz "Hello, LoongArch 64!\n"
msg_len = . - msg

.text
.align 2

.extern printstring,@function

.globl _start
.type  _start,@function

# void _start();
_start:
    pcalau12i $a0, %pc_hi20(msg)
    addi.d    $a0, $a0, %pc_lo12(msg)
    addi.d    $a1, $zero, msg_len
    bl        printstring

    # exit(0)
    addi.d $a0, $zero, 0
    addi.d $a7, $zero, 93
    syscall 0
