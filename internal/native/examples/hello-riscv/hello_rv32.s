# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# ==============================================================================
# 目标：ESP32-C3 (RISC-V 32-bit) 裸机 UART0 打印
# 仅使用基础指令集 (RV32I)
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
    # 1. 加载字符串地址 (message) 到 a0
    # li a0, message  <-- 替换为 lui/addi 组合
    # --------------------------------------------------------
    
    # a0 = message 地址的高 20 位
    lui     a0, %hi(message) 
    # a0 = message 地址的完整 32 位
    addi    a0, a0, %lo(message) 

print_loop:
    # a1 = 当前字符
    lbu     a1, 0(a0)           # 取一个字节 (Load Byte Unsigned)
    
    # 结束检查
    # 如果 a1 是 0 ('\0')，跳转到 finished
    beq     a1, x0, finished    

    # --------------------------------------------------------
    # 2. 加载 UART0 TX FIFO 地址到 t0
    # li t0, UART0_TX_FIFO  <-- 替换为 lui/addi 组合
    # --------------------------------------------------------
    
    # t0 = UART0_TX_FIFO 地址的高 20 位
    lui     t0, %hi(UART0_TX_FIFO) 
    # t0 = UART0_TX_FIFO 地址的完整 32 位 (0x60000000)
    addi    t0, t0, %lo(UART0_TX_FIFO) 

    # 发送字符
    # sb (Store Byte) 将 a1 中的字符写入 UART0 TX FIFO 寄存器
    sb      a1, 0(t0)           
    
    # 移动到下一个字符
    addi    a0, a0, 1           
    
    # 循环
    jal     x0, print_loop

finished:
    # 程序结束，进入死循环
forever:
    jal x0, forever
