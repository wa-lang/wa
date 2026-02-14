# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# gcc -nostdlib -static -z noexecstack hello.s

# 查看退出状态码:
# echo $?

.section .text
.globl _start
_start:
    mov edi, 123 # bf 7b 00 00 00, 匹配 B8+rd id 规则. rdi 编号是 7, 所以 B8+7=BF
    mov eax, 60  # b8 3c 00 00 00, 匹配 B8+rd id 规则. rax 编号是 0, 所以 B8+0=B8
    syscall      # 0f 05, syscall 固定双字节操作码

