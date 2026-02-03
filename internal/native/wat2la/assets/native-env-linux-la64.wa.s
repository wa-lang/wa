# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.extern .Wa.Memory.addr

# void _Wa_Import_syscall_write (int fd, uint32_t ptr, int32_t len)
.section .text
.global .Wa.Import.syscall_linux.print_str
.Wa.Import.syscall_linux.print_str:
    pcalau12i $t0, %pc_hi20(.Wa.Memory.addr)
    ld.d      $t0, $t0, %pc_lo12(.Wa.Memory.addr)

    # a0 = fd
    add.d  $a1, $a1, $t0 # a1 = base + ptr
    # a2 = len
    addi.d $a7, $zero, 64 # sys_write
    syscall 0
    jirl   $zero, $ra, 0

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.global .Wa.Runtime.write
.Wa.Runtime.write:
    addi.d $a7, $zero, 64 # sys_write
    syscall 0
    jirl   $zero, $ra, 0

# void _Wa_Runtime_exit(int status)
.section .text
.global .Wa.Runtime.exit
.Wa.Runtime.exit:
    addi.d $a7, $zero, 93 # sys_exit
    syscall 0

# void* _Wa_Runtime_malloc(int size)
.section .text
.global .Wa.Runtime.malloc
.Wa.Runtime.malloc:
    or     $a1, $a0, $zero   # length = size
    addi.d $a0, $zero, 0     # addr = NULL
    addi.d $a2, $zero, 3     # prot = PROT_READ | PROT_WRITE
    addi.d $a3, $zero, 34    # flags = MAP_PRIVATE | MAP_ANONYMOUS
    addi.d $a4, $zero, -1    # fd = -1
    addi.d $a5, $zero, 0     # offset = 0
    addi.d $a7, $zero, 222   # sys_mmap (222)
    syscall 0
    jirl   $zero, $ra, 0

# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.global .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
    or     $t0, $a0, $zero   # 备份 dst
    beq    $a2, $zero, .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    ld.b   $t1, $a1, 0       # 字节读取
    st.b   $t1, $a0, 0       # 字节写入
    addi.d $a0, $a0, 1
    addi.d $a1, $a1, 1
    addi.d $a2, $a2, -1
    bne    $a2, $zero, .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    or     $a0, $t0, $zero   # 返回 dst
    jirl   $zero, $ra, 0

# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.global .Wa.Runtime.memmove
.Wa.Runtime.memmove:
    beq    $a0, $a1, .Wa.L.memmove.done
    # 如果 a0 < a1 (无符号比较), 跳转到向前拷贝
    bltu   $a0, $a1, .Wa.Runtime.memcpy 

    # 后向拷贝 (dst > src)
    or     $t0, $a0, $zero
    add.d  $a0, $a0, $a2
    add.d  $a1, $a1, $a2
.Wa.L.memmove.back_loop:
    beq    $a2, $zero, .Wa.L.memmove.ret
    addi.d $a0, $a0, -1
    addi.d $a1, $a1, -1
    ld.b   $t1, $a1, 0
    st.b   $t1, $a0, 0
    addi.d $a2, $a2, -1
    b      .Wa.L.memmove.back_loop
.Wa.L.memmove.ret:
    or     $a0, $t0, $zero
.Wa.L.memmove.done:
    jirl   $zero, $ra, 0

# void* _Wa_Runtime_memset(void* s, int c, int n)
.global .Wa.Runtime.memset
.Wa.Runtime.memset:
    or     $t0, $a0, $zero
    beq    $a2, $zero, .Wa.L.memset.done
.Wa.L.memset.loop:
    st.b   $a1, $a0, 0
    addi.d $a0, $a0, 1
    addi.d $a2, $a2, -1
    bne    $a2, $zero, .Wa.L.memset.loop
.Wa.L.memset.done:
    or     $a0, $t0, $zero
    jirl   $zero, $ra, 0

# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
.global .Wa.Import.syscall_linux.print_str
.Wa.Import.syscall_linux.print_str:
    pcalau12i $t0, %pc_hi20(.Wa.Memory.addr)
    ld.d      $t0, $t0, %pc_lo12(.Wa.Memory.addr)

    # a0 = ptr, a1 = len
    or     $a2, $a1, $zero    # count = len
    add.d  $a1, $t0, $a0      # buf = base + ptr
    addi.d $a0, $zero, 1      # fd = stdout (1)
    addi.d $a7, $zero, 64     # sys_write
    syscall 0
    jirl   $zero, $ra, 0

# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
.global .Wa.Import.syscall_linux.proc_exit
.Wa.Import.syscall_linux.proc_exit:
    b .Wa.Runtime.exit

# void _Wa_Import_syscall_linux_print_rune(int32_t c)
.section .text
.global .Wa.Import.syscall_linux.print_rune
.Wa.Import.syscall_linux.print_rune:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -16

    st.b    $a0, $sp, 0

    addi.d  $a0, $zero, 1  # arg.0: stdout
    addi.d  $a1, $sp, 0    # arg.1: buffer
    addi.d  $a2, $zero, 1  # arg.2: count
    addi.d  $a7, $zero, 64 # sys_write
    syscall 0

    addi.d  $sp, $fp, 0
    ld.d    $fp, $sp, 0
    ld.d    $ra, $sp, 8
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# void _Wa_Import_syscall_linux_print_i64(int64_t val)
.section .text
.global .Wa.Import.syscall_linux.print_i64
.Wa.Import.syscall_linux.print_i64:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32
    
    or      $t0, $zero, $a0  # t0 = 工作变量 (val)
    addi.d  $t1, $fp, -1     # t1 为缓冲区指针 (从后往前填)
    addi.d  $t2, $zero, 10   # 除数

    # 1. 处理负数
    bge     $t0, $zero, .Wa.L.syscall_linux.print_i64.convert
    sub.d   $t0, $zero, $t0  # t0 = abs(t0)

.Wa.L.syscall_linux.print_i64.convert:
    div.d   $t3, $t0, $t2    # t3 = 商
    mod.d   $t4, $t0, $t2    # t4 = 余数
    addi.w  $t4, $t4, 48     # 加上 '0' 的 ASCII 码
    st.b    $t4, $t1, 0      # 存入缓冲区
    addi.d  $t1, $t1, -1     # 指针前移
    or      $t0, $zero, $t3  # 更新待处理的数字
    bnez    $t0, .Wa.L.syscall_linux.print_i64.convert  # 如果商不为 0 则继续

    # 2. 补负号
    bge     $a0, $zero, .Wa.L.syscall_linux.print_i64.print
    addi.d  $t4, $zero, 45   # '-'
    st.b    $t4, $t1, 0
    addi.d  $t1, $t1, -1

.Wa.L.syscall_linux.print_i64.print:
    addi.d  $a0, $zero, 1    # arg.0: stdout
    addi.d  $a1, $t1, 1      # arg.1: buffer
    sub.d   $a2, $fp, $a1    # arg.2: count
    addi.d  $a7, $zero, 64   # sys_write
    syscall 0

    addi.d  $sp, $fp, 0
    ld.d    $fp, $sp, 0
    ld.d    $ra, $sp, 8
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

