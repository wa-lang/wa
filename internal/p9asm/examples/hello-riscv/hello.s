    .section .text
    .globl _start

# QEMU virt 机器的 UART0 基地址是 0x10000000
UART0 = 0x10000000

_start:
    la   a0, message      # a0 = 字符串地址

print_loop:
    lbu  a1, 0(a0)        # 取一个字节
    beqz a1, finished     # 如果是0则结束
    li   t0, UART0        # t0 = UART0 地址
    sb   a1, 0(t0)        # 写到UART寄存器
    addi a0, a0, 1        # 下一个字符
    j    print_loop

finished:
    # 停机: 无限循环
    wfi
    j finished

    .section .rodata
message:
    .asciz "Hello RISC-V Baremetal!\n"
