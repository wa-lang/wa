# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.section .data
.align 8
hello.msg: .ascii "hello, world!\n"
hello.len: .quad 14

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    mov eax, 1
    syscall
    ret


# void _Wa_Runtime_exit(int status)
.section .text
.globl .Wa.Runtime.exit
.Wa.Runtime.exit:
    mov eax, 60
    syscall


# 汇编程序入口函数
.section .text
.globl _start
_start:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stdout, msg, len)
    mov  rdi, 1 # arg.0: stdout
    lea  rsi, qword ptr [rip + hello.msg]
    mov  rdx, qword ptr [rip + hello.len]
    call .Wa.Runtime.write

    # runtime.exit(0)
    mov  rdi, 0
    call .Wa.Runtime.exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

