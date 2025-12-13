# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 系统调用编号
const $SYS_write = 64
const $SYS_exit = 93

# 用于输出的字符串
global $message = "Hello Linux-loong64!\n\x00"

# 主函数
func _start {
%begin:
    # a0 = stdout(1)
    ori a0, zero, 1

    # a1 = $message
    pcaddu12i a1, %pcrel_hi($message)   # 高20位
    addi.w    a1, a1, %pcrel_lo(%begin) # 低12位

    # a2 = strlen($message)
    ori a2, zero, %strlen($message)

    # 执行系统调用
    # a7 = SYS_write(64)
    ori a7, zero, $SYS_write
    syscall 0

    # 执行系统调用
    # a0 = Success(0)
    # a7 = SYS_exit(93)
    ori a0, zero, 0
    ori a7, zero, $SYS_exit
    syscall 0
}
