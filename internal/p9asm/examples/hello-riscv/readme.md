# RISCV 测试

安装工具链和qemu:

- https://github.com/sifive/freedom-tools/releases
  - riscv64-unknown-elf-toolchain-10.2.0-2020.12.8-x86_64-w64-mingw32.zip
- https://qemu.weilnetz.de/w64/
  - qemu-w64-setup-20250814.exe


## 裸机模式(不需要Linux系统)

构建 hello.s 例子:

```s
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
```

准备Link配置脚本link.ld文件:

```
ENTRY(_start)

SECTIONS
{
  . = 0x80000000;    /* QEMU virt 默认从0x80000000开始执行 */

  .text : {
    *(.text*)
  }

  .rodata : {
    *(.rodata*)
  }

  .data : {
    *(.data*)
  }

  .bss : {
    *(.bss*)
    *(COMMON)
  }
}
```

编译出elf格式的可执行程序:

```
$ riscv64-unknown-elf-as hello.S -o hello.o
$ riscv64-unknown-elf-ld -T link.ld hello.o -o hello.elf.exe
```

命令行执行:

```
$ qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello.elf.exe
```

执行后可以输出文本, 但是无法用Ctrl-C退出程序, 需要手动在任务管理器中退出qemu进程。

## Linux用户模式(有待验证)

构建 hello.s 例子:

```s
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
```

编译程序:

```
$ riscv64-unknown-elf-as hello.S -o hello.o
$ riscv64-unknown-elf-ld hello.o -o hello.elf
```

生成的是riscv指令的linux可执行程序, 需要通过qemu等工具执行:

```
$ qemu-riscv64 hello.elf
```
