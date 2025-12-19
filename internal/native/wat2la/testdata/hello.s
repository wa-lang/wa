# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 标准输出
const $STDOUT = 1

# 系统调用编号
const $SYS_write = 64
const $SYS_exit = 93

# 用于输出的字符串
global $message = "Hello, LoongArch 64!\n\x00"

# 主函数
func _start {
%begin:
    # SYS_write(STDOUT, $message, %sizeof($message))
    addi.d a0, zero, $STDOUT
    pcalau12i a1, %pc_hi20($message)   # 高20位
    addi.d    a1, a1, %pc_lo12($message) # 低12位
    addi.d    a2, zero, %sizeof($message)
    addi.d    a7, zero, $SYS_write
    syscall   0

    # SYS_exit(0)
    addi.d  a0, zero, 0
    addi.d  a7, zero, $SYS_exit
    syscall 0
}
