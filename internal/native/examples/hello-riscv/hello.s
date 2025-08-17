# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 注: 尽量避免使用伪指令和相关特性

    .section .rodata
message:
    .asciz "Hello RISC-V Baremetal!\n"

    .section .text
    .globl _start

# QEMU virt 机器 UART0 和 exit device 的基地址
UART0      = 0x10000000
EXIT_DEVICE = 0x100000

_start:
    # a0 = 字符串地址
    auipc   a0, %pcrel_hi(message)     # 高20位 = 当前PC + 偏移
    addi    a0, a0, %pcrel_lo(_start)  # 低12位

print_loop:
    lbu  a1, 0(a0)        # 取一个字节
    beq  a1, x0, finished # 如果是0则结束

    # t0 = UART0 地址
    lui     t0, %hi(UART0)           # UART0 高20位
    addi    t0, t0, %lo(UART0)       # UART0 低12位

    sb   a1, 0(t0)        # 写到UART寄存器
    addi a0, a0, 1        # 下一个字符
    jal  x0, print_loop

finished:
    # 写退出码 0 到 EXIT_DEVICE，让 QEMU 退出
    lui     t0, %hi(EXIT_DEVICE)     # exit device 地址
    addi    t0, t0, %lo(EXIT_DEVICE)

    # t1 = 0x5555
    # addi rd, rs1, imm 的 imm 范围是 [-2048, +2047](12 位有符号立即数)
    lui   t1, 0x5             # 高 20 位 (0x5 << 12 = 0x5000)
    addi  t1, t1, 0x555       # 结果 = 0x5000 + 0x555 = 0x5555

    sw   t1, 0(t0)

    # 如果 QEMU 不支持 exit 设备，就进入并死循环
forever:
    jal x0, forever
