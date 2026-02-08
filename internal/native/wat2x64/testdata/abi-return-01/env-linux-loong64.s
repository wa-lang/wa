# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.section .data
.align 3
.L.str.cside:    .ascii "C side: input="
.L.str.cside.len: .quad 14

.L.str.retprefix: .ascii ", returning ["
.L.str.retprefix.len: .quad 13

.L.str.printi64: .ascii "printI64: "
.L.str.printi64.len: .quad 10

# static void print_str(const char* s, int32_t len)
.section .text
.print_str:
    addi.d  $sp, $sp, -48
    st.d    $ra, $sp, 40
    st.d    $fp, $sp, 32
    st.d    $s0, $sp, 24      # s
    st.d    $s1, $sp, 16      # len
    st.d    $s2, $sp, 8       # i
    
    move    $s0, $a0
    move    $s1, $a1
    move    $s2, $zero
.L.ps_loop:
    bge     $s2, $s1, .L.ps_done
    ld.b    $a0, $s0, 0
    bl      .Wa.Import.syscall_linux.print_rune
    addi.d  $s0, $s0, 1
    addi.d  $s2, $s2, 1
    b       .L.ps_loop
.L.ps_done:
    ld.d    $s2, $sp, 8
    ld.d    $s1, $sp, 16
    ld.d    $s0, $sp, 24
    ld.d    $fp, $sp, 32
    ld.d    $ra, $sp, 40
    addi.d  $sp, $sp, 48
    jirl    $zero, $ra, 0

# EnvMultiRet _Wa_Import_env_get_multi_values(int64_t input_val)
.section .text
.globl .Wa.Import.env.get_multi_values
.Wa.Import.env.get_multi_values:
    # 栈分配：需要保存 ra, s0(hidden_ptr), s1(input_val)
    addi.d  $sp, $sp, -32
    st.d    $ra, $sp, 24
    st.d    $s0, $sp, 16      # 保存隐藏指针 (hidden_ptr)
    st.d    $s1, $sp, 8       # 保存输入值 (input_val)
    
    move    $s0, $a0          # s0 = r (hidden_ptr)
    move    $s1, $a1          # s1 = input_val

    # 逻辑计算：r.v1 = input_val + 1; r.v2 = + 2; r.v3 = + 3;
    addi.d  $t0, $s1, 1
    st.d    $t0, $s0, 0       # r.v1
    addi.d  $t1, $s1, 2
    st.d    $t1, $s0, 8       # r.v2
    addi.d  $t2, $s1, 3
    st.d    $t2, $s0, 16      # r.v3

    # 打印 "C side: input="
    pcalau12i $a0, %pc_hi20(.L.str.cside)
    addi.d    $a0, $a0, %pc_lo12(.L.str.cside)
    pcalau12i $t0, %pc_hi20(.L.str.cside.len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.cside.len)
    bl        .print_str

    # 打印 input_val
    move      $a0, $s1
    bl        .Wa.Import.syscall_linux.print_i64

    # 打印 ", returning ["
    pcalau12i $a0, %pc_hi20(.L.str.retprefix)
    addi.d    $a0, $a0, %pc_lo12(.L.str.retprefix)
    pcalau12i $t0, %pc_hi20(.L.str.retprefix.len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.retprefix.len)
    bl        .print_str

    # 打印 r.v1, r.v2, r.v3 (带逗号和空格)
    ld.d      $a0, $s0, 0
    bl        .Wa.Import.syscall_linux.print_i64
    li.d      $a0, 44         # ','
    bl        .Wa.Import.syscall_linux.print_rune
    li.d      $a0, 32         # ' '
    bl        .Wa.Import.syscall_linux.print_rune

    ld.d      $a0, $s0, 8
    bl        .Wa.Import.syscall_linux.print_i64
    li.d      $a0, 44         # ','
    bl        .Wa.Import.syscall_linux.print_rune
    li.d      $a0, 32         # ' '
    bl        .Wa.Import.syscall_linux.print_rune

    ld.d      $a0, $s0, 16
    bl        .Wa.Import.syscall_linux.print_i64
    
    # 打印 "]\n"
    li.d      $a0, 93         # ']'
    bl        .Wa.Import.syscall_linux.print_rune
    li.d      $a0, 10         # '\n'
    bl        .Wa.Import.syscall_linux.print_rune

    # 根据 ABI，返回值必须放回 $a0 (即隐藏指针地址)
    move    $a0, $s0
    ld.d    $s1, $sp, 8
    ld.d    $s0, $sp, 16
    ld.d    $ra, $sp, 24
    addi.d  $sp, $sp, 32
    jirl    $zero, $ra, 0

# void _Wa_Import_env_print_i64(int64_t x)
.section .text
.globl .Wa.Import.env.print_i64
.Wa.Import.env.print_i64:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $s0, $sp, 0
    move    $s0, $a0

    pcalau12i $a0, %pc_hi20(.L.str.printi64)
    addi.d    $a0, $a0, %pc_lo12(.L.str.printi64)
    pcalau12i $t0, %pc_hi20(.L.str.printi64.len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.printi64.len)
    bl        .print_str

    move      $a0, $s0
    bl        .Wa.Import.syscall_linux.print_i64
    li.d      $a0, 10
    bl        .Wa.Import.syscall_linux.print_rune

    ld.d      $s0, $sp, 0
    ld.d      $ra, $sp, 8
    addi.d    $sp, $sp, 16
    jirl      $zero, $ra, 0

