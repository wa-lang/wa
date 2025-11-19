# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# ==============================================================================
# 目标：ESP32-C3 (RISC-V 32-bit) 裸机 UART0 打印
# 假设：UART0 已由标准引导器初始化（例如波特率 115200）。
# ==============================================================================

.section .rodata
message:
    .asciz "Hello RISC-V Baremetal on ESP32-C3!\n"

.section .text
.globl _start

# ESP32-C3 UART0 TX FIFO 寄存器地址
# UART0 基地址 0x60000000，TX FIFO 偏移 0x0000
UART0_TX_FIFO = 0x60000000 

_start:
    # --------------------------------------------------------
    # 1. 初始化 (通常需要设置栈指针等，这里为了简化省略)
    # --------------------------------------------------------
    
    # a0 = 字符串地址
    # li 是一个伪指令，它会使用 lui 和 addi 来加载 32 位常数地址
    li      a0, message      # a0 = 字符串 "message" 的加载地址

print_loop:
    # a1 = 当前字符
    lbu     a1, 0(a0)           # 取一个字节 (Load Byte Unsigned)
    
    # 结束检查
    # 如果 a1 是 0 (字符串结束符 '\0')，跳转到 finished
    beq     a1, x0, finished    

    # t0 = UART0 TX FIFO 地址
    li      t0, UART0_TX_FIFO   # t0 = 0x60000000 (UART0 TX 寄存器)

    # 发送字符
    sb      a1, 0(t0)           # 将 a1 中的字符写入 UART0 TX FIFO (Store Byte)
    
    # 移动到下一个字符
    addi    a0, a0, 1           
    
    # 循环
    jal     x0, print_loop

finished:
    # ESP32 裸机程序完成后，通常进入死循环或触发软件复位。
    # 这里我们进入死循环。
forever:
    jal x0, forever
