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
    pcaddu12i a0, %pcrel_hi($message)   # 高20位
    addi.w    a0, a0, %pcrel_lo(%begin) # 低12位

%print_loop:
    ld.bu  a1, a0, 0         # 取一个字节
    beq    a1, zero, %finished # 如果是0则结束

    # t0 = UART0 值
    lu12i.w t0, %hi($UART0)           # UART0 高20位
    addi.w  t0, t0, %lo($UART0)       # UART0 低12位

    st.b   a1, t0, 0        # 写到UART寄存器
    addi.w a0, a0, 1        # 下一个字符
    b      %print_loop

%finished:
    # 写退出码 0 到 EXIT_DEVICE, 退出
    lu12i.w t0, %hi($EXIT_DEVICE)     # exit device 地址
    addi.w  t0, t0, %lo($EXIT_DEVICE)

    # t1 = 0x5555
    lu12i.w t1, 0x5             # 高 20 位 (0x5 << 12 = 0x5000)
    addi.w  t1, t1, 0x555       # 结果 = 0x5000 + 0x555 = 0x5555

    st.w    t1, t0, 0

    # 如果不支持 exit 设备，就进入并死循环
%forever:
    b      %forever
}
