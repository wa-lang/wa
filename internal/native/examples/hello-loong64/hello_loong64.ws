# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# QEMU virt 机器 UART0 和 exit device 的基地址
const $UART0 = 0x10000000
const $EXIT_DEVICE = 0x100000

# 用于输出的字符串
global $message = "Hello LOONG64 Baremetal!\n\x00"

# 主函数
func _start {
%begin:
    # a0 = 字符串地址
    lu12i.w a0, %hi($message)     # 高20位
    ori     a0, a0, %lo($message) # 低12位(不能使用addi.w)

%print_loop:
    ld.bu  a1, 0(a0)         # 取一个字节
    beq  a1, zero, %finished # 如果是0则结束

    # t0 = UART0 地址
    lu12i.w t0, %hi($UART0)           # UART0 高20位
    ori     t0, t0, %lo($UART0)       # UART0 低12位

    st.b   a1, 0(t0)        # 写到UART寄存器
    addi.w a0, a0, 1        # 下一个字符
    b      %print_loop

%finished:
    b      %finished
}
