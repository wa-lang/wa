# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

# kernel32.dll

.extern ExitProcess

.section .text
.global _start
_start:
    sub rsp, 40  # Windows x64 要求 shadow space
    mov ecx, 123 # 参数 = 123
    call ExitProcess
