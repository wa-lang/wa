    .section .text
    .global _start

_start:
    # 调用 write(1, msg, len)
    li   a7, 64          # sys_write 系统调用号 (Linux RISC-V)
    li   a0, 1           # fd = 1 (stdout)
    la   a1, msg         # buf = msg
    li   a2, 13          # len = 13
    ecall

    # 调用 exit(0)
    li   a7, 93          # sys_exit
    li   a0, 0
    ecall

    .section .rodata
msg:
    .asciz "Hello RISC-V\n"
