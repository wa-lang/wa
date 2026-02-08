# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.section .data
.align 3
.L.str.printI64: .ascii "printI64: "
.L.str.printI64.len: .quad 10

# static void print_str(const char* s, int32_t len)
.section .text
.print_str:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    st.d   $s0, $sp, 0      # 保存 s
    st.d   $s1, $sp, 8      # 保存 len
    st.d   $s2, $sp, 16     # 存放 i
    
    or   $s0, $zero, $a0
    or   $s1, $zero, $a1
    or   $t0, $zero, $zero       # i = 0

.L.loop_begin:
    bge    $t0, $s1, .L.loop_done

    or     $s2, $zero, $t0

    ld.b   $a0, $s0, 0      # *s
    bl     .Wa.Import.syscall_linux.print_rune
    
    addi.d $s0, $s0, 1      # s++
    addi.d $s2, $s2, 1      # i++
    or     $t0, $zero, $s2
    b      .L.loop_begin

.L.loop_done:
    ld.d   $s0, $sp, 0
    ld.d   $s1, $sp, 8
    ld.d   $s2, $sp, 16

    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0


# void _Wa_Import_env_print_i64(int64_t x)
.section .text
.globl .Wa.Import.env.print_i64
.Wa.Import.env.print_i64:
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8
    st.d   $s0, $sp, 0 # save x
    or     $s0, $zero, $a0

    # print_str("printI64: ", len)
    pcalau12i $a0, %pc_hi20(.L.str.printI64)
    addi.d    $a0, $a0, %pc_lo12(.L.str.printI64)
    pcalau12i $t0, %pc_hi20(.L.str.printI64.len) # len
    ld.d      $a1, $t0, %pc_lo12(.L.str.printI64.len)
    bl        .print_str

    # _Wa_Import_syscall_linux_print_i64(x)
    or        $a0, $zero, $s0
    bl        .Wa.Import.syscall_linux.print_i64

    # _Wa_Import_syscall_linux_print_rune('\n')
    addi.d    $a0, $zero, 10 # '\n'
    bl        .Wa.Import.syscall_linux.print_rune

    ld.d   $s0, $sp, 0
    ld.d   $ra, $sp, 8
    addi.d $sp, $sp, 16
    jirl   $zero, $ra, 0

