# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 注: 尽量避免使用伪指令和相关特性

    .section .rodata
message:
    .asciz "Hello LoongArch Baremetal!\n"

    .section .text
    .globl _start

# 假设的 LoongArch 裸机机器设备地址
# LoongArch 裸机平台设备地址可能不同，这里仅为示例
UART0_BASE      = 0x10000000 # 假设 UART 基地址
EXIT_DEVICE_BASE = 0x10000400 # 假设退出设备基地址

_start:
    # $a0 = 字符串地址
    # LoongArch PC相对寻址: pcalau12i + addi.d
    # $a0 高20位 = 当前PC + 偏移 (PCALAU12I 类似于 AUIPC)
    pcalau12i $a0, %pc_hi20(message)
    # $a0 低12位 = $a0 + 偏移 (ADDI.D 是 ADDI 的 64位版本)
    addi.d $a0, $a0, %pc_lo12(message) 

print_loop:
    # 取一个字节 (LBU: Load Byte Unsigned)
    lbu $a1, $a0, 0
    # 如果是0则结束 (BEQ: Branch on Equal)
    beq $a1, $zero, finished 

    # $t0 = UART0 地址 (使用绝对地址寻址)
    # LoongArch 加载32位/64位立即数/地址通常使用 `lu12i.w` + `ori` 或 `ld.d` 等
    # 为保持和RISC-V的 `lui` + `addi` 结构相似，我们分两步加载地址。
    # LoongArch 的地址加载指令通常为 `lu12i.w` (Load Upper 12-bit Immediate to Word/32位)
    # 或更简化的 `li.d` (伪指令). 这里使用 `pcalau12i`/`addi.d` 组合加载绝对地址:

    # 加载 UART0_BASE 地址到 $t0 (类似PC相对加载，但这里加载的是绝对地址的常量)
    # LoongArch 标准指令集没有像 RISC-V 那样直接的绝对地址加载组合指令
    # 避免伪指令 `li.d`，我们使用以下方法加载64位常量:

    # $t0 = UART0_BASE (0x10000000)
    lu12i.w $t0, %hi(UART0_BASE)  # $t0 = UART0_BASE 的高 12 位，低 20 位清零
    ori $t0, $t0, %lo(UART0_BASE) # $t0 = $t0 | UART0_BASE 的低 12 位
    
    # 写到UART寄存器 (SB: Store Byte)
    st.b $a1, $t0, 0 
    
    # 下一个字符
    addi.d $a0, $a0, 1 
    # 跳转到循环 (JAL: Jump and Link. JAL $zero, print_loop)
    j print_loop

finished:
    # 写退出码 0 到 EXIT_DEVICE，让 QEMU 退出
    
    # $t0 = EXIT_DEVICE_BASE 地址 (0x10000400)
    lu12i.w $t0, %hi(EXIT_DEVICE_BASE)
    ori $t0, $t0, %lo(EXIT_DEVICE_BASE)

    # $t1 = 0x5555
    # LoongArch lu12i.w + ori 加载32位常量
    lu12i.w $t1, 0x0
    ori $t1, $t1, 0x5555 # $t1 = 0x5555 (ORI 立即数范围是 [0, 4095])

    # 退出码写入 EXIT_DEVICE
    st.w $t1, $t0, 0 # SW -> ST.W (Store Word)

    # 如果 QEMU 不支持 exit 设备，就进入并死循环
forever:
    b forever # JAL x0, forever -> B forever (Branch)
