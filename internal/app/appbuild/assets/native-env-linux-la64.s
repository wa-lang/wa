# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.extern .Wa.Memory.addr

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.global _Wa_Runtime_write
_Wa_Runtime_write:
    addi.d $a7, $zero, 64 # sys_write
    syscall 0
    jirl   $zero, $ra, 0

# void _Wa_Runtime_exit(int status)
.section .text
.global _Wa_Runtime_exit
_Wa_Runtime_exit:
    addi.d $a7, $zero, 93 # sys_exit
    syscall 0

# void* _Wa_Runtime_malloc(int size)
.section .text
.global _Wa_Runtime_malloc
_Wa_Runtime_malloc:
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
.global _Wa_Runtime_memcpy
_Wa_Runtime_memcpy:
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
.global _Wa_Runtime_memmove
_Wa_Runtime_memmove:
    beq    $a0, $a1, .Wa.L.memmove.done
    # 如果 a0 < a1 (无符号比较), 跳转到向前拷贝
    bltu   $a0, $a1, _Wa_Runtime_memcpy 

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
.global _Wa_Runtime_memset
_Wa_Runtime_memset:
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
.global _Wa_Import_syscall_linux_print_str
_Wa_Import_syscall_linux_print_str:
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
.global _Wa_Import_syscall_linux_proc_exit
_Wa_Import_syscall_linux_proc_exit:
    b _Wa_Runtime_exit
