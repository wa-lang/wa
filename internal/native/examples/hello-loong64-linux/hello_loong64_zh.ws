# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 标准输出
常量 $STDOUT = 1

# 系统调用编号
常量 $SYS_write = 64
常量 $SYS_exit = 93

# 用于输出的字符串
全局 $信札 = "你好, 龙的传人(Loongson Linux)!\n"

# 主函数
函数 _启动:
    # SYS_write(STDOUT, $信札, %内存字节数($信札))
    加立.长          参甲格, 零格, 1
    计算术左移正12立 参乙格, %相对.高20($信札)
    加立.长          参乙格, 参乙格, %相对.低12($信札)
    加立.长          参丙格, 零格, %内存字节数($信札)
    加立.长          参辛格, 零格, $SYS_write
    系统调用 0

    # SYS_exit(0)
    加立.长  参甲格, 零格, 0
    加立.长  参辛格, 零格, $SYS_exit
    系统调用 0
完毕
