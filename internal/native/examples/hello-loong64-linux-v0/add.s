# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.text
.align 2

.globl add_f
.type  add_f,@function

add_f:
    add.w $a0, $a0, $a1
    add.w $a0, $a0, $a2
    add.w $a0, $a0, $a3
    jr    $ra
    .size add_f, .-add_f
