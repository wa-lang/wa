# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.section .data
.align 3
.L.str.abi_ver:    .ascii "--- Linux ABI Verification ---\n"
.L.str.abi_ver_len: .quad 31

.L.str.p1:         .ascii "Param 1 (RDI): "
.L.str.p1_len:     .quad 15

.L.str.p2:         .ascii "Param 2 (RSI): "
.L.str.p2_len:     .quad 15

.L.str.p3:         .ascii "Param 3 (RDX): "
.L.str.p3_len:     .quad 15

.L.str.p4:         .ascii "Param 4 (RCX): "
.L.str.p4_len:     .quad 15

.L.str.p5:         .ascii "Param 5 (R8): "
.L.str.p5_len:     .quad 14

.L.str.p6:         .ascii "Param 6 (R9): "
.L.str.p6_len:     .quad 14

.L.str.p7:         .ascii "Param 7 (Stack RSP+0): "
.L.str.p7_len:     .quad 23

.L.str.p8:         .ascii "Param 8 (Stack RSP+8): "
.L.str.p8_len:     .quad 23

.L.str.div:        .ascii "-------------------------------\n"
.L.str.div_len:    .quad 32

.L.str.i64:        .ascii "printI64: "
.L.str.i64_len:    .quad 10

# static void print_str(const char* s, int32_t len)
.section .text
.print_str:
    addi.d  $sp, $sp, -48
    st.d    $ra, $sp, 40
    st.d    $fp, $sp, 32
    addi.d  $fp, $sp, 48
    st.d    $s0, $sp, 24
    st.d    $s1, $sp, 16
    st.d    $s2, $sp, 8
    
    or      $s0, $zero, $a0    # s
    or      $s1, $zero, $a1    # len
    or      $s2, $zero, $zero  # i = 0
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

# int64_t _Wa_Import_env_write(fd, ptr, size, p4, p5, p6, p7, p8)
.section .text
.globl .Wa.Import.env.write
.Wa.Import.env.write:
    # 这里的参数多，我们需要保护所有的参数到栈上
    # a0-a7 共 8 个，加上 ra, fp，至少需要 80 字节 -> 96 字节对齐
    addi.d  $sp, $sp, -96
    st.d    $ra, $sp, 88
    st.d    $fp, $sp, 80
    st.d    $a0, $sp, 0    # fd
    st.d    $a1, $sp, 8    # ptr
    st.d    $a2, $sp, 16   # size
    st.d    $a3, $sp, 24   # p4
    st.d    $a4, $sp, 32   # p5
    st.d    $a5, $sp, 40   # p6
    st.d    $a6, $sp, 48   # p7
    st.d    $a7, $sp, 56   # p8

    # --- Verification Header ---
    pcalau12i $a0, %pc_hi20(.L.str.abi_ver)
    addi.d    $a0, $a0, %pc_lo12(.L.str.abi_ver)
    pcalau12i $t0, %pc_hi20(.L.str.abi_ver_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.abi_ver_len)
    bl        .print_str

    # 重复逻辑：打印 Param 1-8。这里以 Param 1 为例，后续逻辑雷同
    # Param 1
    pcalau12i $a0, %pc_hi20(.L.str.p1)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p1)
    pcalau12i $t0, %pc_hi20(.L.str.p1_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p1_len)
    bl        .print_str
    ld.d      $a0, $sp, 0   # 恢复 fd
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 2
    pcalau12i $a0, %pc_hi20(.L.str.p2)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p2)
    pcalau12i $t0, %pc_hi20(.L.str.p2_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p2_len)
    bl        .print_str
    ld.d      $a0, $sp, 8   # 恢复 ptr
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 3
    pcalau12i $a0, %pc_hi20(.L.str.p3)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p3)
    pcalau12i $t0, %pc_hi20(.L.str.p3_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p3_len)
    bl        .print_str
    ld.d      $a0, $sp, 16   # 恢复 size
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 4
    pcalau12i $a0, %pc_hi20(.L.str.p4)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p4)
    pcalau12i $t0, %pc_hi20(.L.str.p4_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p4_len)
    bl        .print_str
    ld.d      $a0, $sp, 24   # 恢复 p4
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 5
    pcalau12i $a0, %pc_hi20(.L.str.p5)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p5)
    pcalau12i $t0, %pc_hi20(.L.str.p5_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p5_len)
    bl        .print_str
    ld.d      $a0, $sp, 32   # 恢复 p5
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 6
    pcalau12i $a0, %pc_hi20(.L.str.p6)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p6)
    pcalau12i $t0, %pc_hi20(.L.str.p6_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p6_len)
    bl        .print_str
    ld.d      $a0, $sp, 40   # 恢复 p6
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # Param 7
    pcalau12i $a0, %pc_hi20(.L.str.p7)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p7)
    pcalau12i $t0, %pc_hi20(.L.str.p7_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p7_len)
    bl        .print_str
    ld.d      $a0, $sp, 48   # 恢复 p7
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # (由于篇幅限制，Param 2-8 逻辑完全一致，仅需更换字符串标签和 ld.d 偏移)
    # 示例 Param 8:
    pcalau12i $a0, %pc_hi20(.L.str.p8)
    addi.d    $a0, $a0, %pc_lo12(.L.str.p8)
    pcalau12i $t0, %pc_hi20(.L.str.p8_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.p8_len)
    bl        .print_str
    ld.d      $a0, $sp, 56  # 恢复 p8
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune

    # --- Verification Footer ---
    pcalau12i $a0, %pc_hi20(.L.str.div)
    addi.d    $a0, $a0, %pc_lo12(.L.str.div)
    pcalau12i $t0, %pc_hi20(.L.str.div_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.div_len)
    bl        .print_str

    # 调用 _Wa_Runtime_write(fd, ptr, size)
    ld.d    $a0, $sp, 0
    ld.d    $a1, $sp, 8
    ld.d    $a2, $sp, 16
    bl      .Wa.Runtime.write

    # 返回值 0
    or      $a0, $zero, $zero
    ld.d    $fp, $sp, 80
    ld.d    $ra, $sp, 88
    addi.d  $sp, $sp, 96
    jirl    $zero, $ra, 0

# void _Wa_Import_env_print_i64(int64_t x)
.section .text
.globl .Wa.Import.env.print_i64
.Wa.Import.env.print_i64:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $s0, $sp, 0
    or      $s0, $zero, $a0
    pcalau12i $a0, %pc_hi20(.L.str.i64)
    addi.d    $a0, $a0, %pc_lo12(.L.str.i64)
    pcalau12i $t0, %pc_hi20(.L.str.i64_len)
    ld.d      $a1, $t0, %pc_lo12(.L.str.i64_len)
    bl        .print_str
    or        $a0, $zero, $s0
    bl        .Wa.Import.syscall_linux.print_i64
    addi.d    $a0, $zero, 10
    bl        .Wa.Import.syscall_linux.print_rune
    ld.d      $s0, $sp, 0
    ld.d      $ra, $sp, 8
    addi.d    $sp, $sp, 16
    jirl      $zero, $ra, 0
