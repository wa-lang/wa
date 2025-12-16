# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.text
.align 2

.globl printchar
.type  printchar,@function

.globl printstring
.type  printstring,@function

# void printchar(int c);
printchar:
    # $sp = $sp - 16, sp 需要 16 字节对齐
    # $ra = $sp + 8
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8

    # mem[$sp+7] = $a0
    st.b $a0, $sp, 7

    # write(1, $sp+7, size)
    addi.d $a0, $zero, 1
    addi.d $a1, $sp, 7
    addi.d $a2, $zero, 1
    addi.d $a7, $zero, 64
    syscall 0

    # return
    ld.d $ra, $sp, 8
    addi.d $sp, $sp, 16
    jr $ra

# void printstring(const char* s, int len);
printstring:
    # $sp = $sp - 16, sp 需要 16 字节对齐
    # $ra = $sp + 8
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8

    # mem[$sp+7] = $a0
    st.b $a0, $sp, 7

    # write(1, $a0, $a1)
    or $a2, $a1, $zero   # $a2 = $a1
    or $a1, $a0, $zero   # $a1 = $a0
    addi.d $a0, $zero, 1 # $a0 = 1
    addi.d $a7, $zero, 64
    syscall 0

    # return
    ld.d $ra, $sp, 8
    addi.d $sp, $sp, 16
    jr $ra
