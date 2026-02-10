# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-或-later

# TODO: 改为RISCV中文汇编

声明 .Wa.Memory.addr

# void _Wa_Import_syscall_write (int fd, uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall.write:
    auipc   t0, %相对高位(.Wa.Memory.addr)
    ld      t0, %相对低位(.Wa.Import.syscall.write)(t0)

    # a0 = fd
    add a1, a1, t0 # a1 = base + ptr
    # a2 = len
    addi a7, zero, 64 # sys_write
    ecall

    jalr zero, 0(ra)
完毕


# int _Wa_Runtime_write(int fd, void *buf, int count)
函数 .Wa.Runtime.write:
    addi a7, zero, 64 # sys_write
    ecall
    jalr zero, 0(ra)
完毕


# void _Wa_Runtime_exit(int status)
函数 .Wa.Runtime.exit:
    addi a7, zero, 93 # sys_exit
    ecall
完毕


# void* _Wa_Runtime_malloc(int size)
函数 .Wa.Runtime.malloc:
    addi  a1, a0, 0     # length = size
    addi  a0, zero, 0   # addr = NULL
    addi  a2, zero, 3   # prot = PROT_READ | PROT_WRITE
    addi  a3, zero, 34  # flags = MAP_PRIVATE | MAP_ANONYMOUS
    addi  a4, zero, -1  # fd = -1
    addi  a5, zero, 0   # offset = 0
    addi  a7, zero, 222 # sys_mmap (222)
    ecall
    jalr  zero, 0(ra)
完毕


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
函数 .Wa.Runtime.memcpy:
    addi t0, a0, 0 # 备份 dst
    beq  a2, zero, .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    lb   t1, 0(a1) # 字节读取
    sb   t1, 0(a0) # 字节写入
    addi a0, a0, 1
    addi a1, a1, 1
    addi a2, a2, -1
    bne  a2, zero, .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    addi a0, t0, 0 # 返回 dst
    jalr zero, 0(ra)
完毕


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
函数 .Wa.Runtime.memmove:
    beq a0, a1, .Wa.L.memmove.done
    # 如果 a0 < a1 (无符号比较), 跳转到向前拷贝
    bltu a0, a1, .Wa.Runtime.memcpy

    # 后向拷贝 (dst > src)
    addi t0, a0, 0
    add  a0, a0, a2
    add  a1, a1, a2
.Wa.L.memmove.back_loop:
    beq  a2, zero, .Wa.L.memmove.ret
    addi a0, a0, -1
    addi a1, a1, -1
    lb   t1, 0(a1)
    sb   t1, 0(a0)
    addi a2, a2, -1
    jal  zero, .Wa.L.memmove.back_loop
.Wa.L.memmove.ret:
    addi a0, t0, 0
.Wa.L.memmove.done:
    jalr zero, 0(ra)
完毕


# void* _Wa_Runtime_memset(void* s, int c, int n)
函数 .Wa.Runtime.memset:
    addi t0, a0, 0
    beq  a2, zero, .Wa.L.memset.done
.Wa.L.memset.loop:
    sb   a1, 0(a0)
    addi a0, a0, 1
    addi a2, a2, -1
    bne  a2, zero, .Wa.L.memset.loop
.Wa.L.memset.done:
    addi a0, t0, 0
    jalr zero, 0(ra)
完毕


# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall_linux.print_str:
    auipc t0, %相对高位(.Wa.Memory.addr)
    ld    t0, %相对低位(.Wa.Import.syscall_linux.print_str)(t0)

    # a0 = ptr, a1 = len
    addi a2, a1, 0    # count = len
    add  a1, t0, a0   # buf = base + ptr
    addi a0, zero, 1  # fd = stdout (1)
    addi a7, zero, 64 # sys_write
    ecall
    jalr zero, 0(ra)
完毕


# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
函数 .Wa.Import.syscall_linux.proc_exit:
    jal zero, .Wa.Runtime.exit
完毕


# void _Wa_Import_syscall_linux_print_rune(int32_t c)
函数 .Wa.Import.syscall_linux.print_rune:
    addi sp, sp, -16
    sd   ra, 8(sp)
    sd   s0, 0(sp)
    addi s0, sp, 0
    addi sp, sp, -16

    sb a0, 0(sp)

    addi a0, zero, 1  # arg.0: stdout
    addi a1, sp, 0    # arg.1: buffer
    addi a2, zero, 1  # arg.2: count
    addi a7, zero, 64 # sys_write
    ecall

    addi sp, s0, 0
    ld   s0, 0(sp)
    ld   ra, 8(sp)
    ld   ra, 8(sp)
    addi sp, sp, 16
完毕


# void _Wa_Import_syscall_linux_print_i64(int64_t val)
函数 .Wa.Import.syscall_linux.print_i64:
    addi sp, sp, -16
    sd   ra, 8(sp)
    sd   s0, 0(sp)
    addi s0, sp, 0
    addi sp, sp, -32

    addi t0, a0, 0    # t0 = 工作变量 (val)
    addi t1, s0, -1   # t1 为缓冲区指针 (从后往前填)
    addi t2, zero, 10 # 除数

    # 1. 处理负数
    bge t0, zero, .Wa.L.syscall_linux.print_i64.convert
    sub t0, zero, t0 # t0 = abs(t0)

.Wa.L.syscall_linux.print_i64.convert:
    div  t3, t0, t2 # t3 = 商
    rem  t4, t0, t2 # t4 = 余数
    addi t4, t4, 48 # 加上 '0' 的 ASCII 码
    sb   t4, 0(t1)  # 存入缓冲区
    addi t1, t1, -1 # 指针前移
    addi t0, t3, 0  # 更新待处理的数字
    bne  t0, zero, .Wa.L.syscall_linux.print_i64.convert

    # 2. 补负号
    bge  a0, zero, .Wa.L.syscall_linux.print_i64.print
    addi t4, zero, 45 # '-'
    sb   t4, 0(t1)
    addi t1, t1, -1

.Wa.L.syscall_linux.print_i64.print:
    addi a0, zero, 1  # arg.0: stdout
    addi a1, t1, 1    # arg.1: buffer
    sub  a2, s0, a1   # arg.2: count
    addi a7, zero, 64 # sys_write
    ecall

    addi sp, s0, 0
    ld   s0, 0(sp)
    ld   ra, 8(sp)
    addi sp, sp, 16
    jalr zero, 0(ra)
完毕

