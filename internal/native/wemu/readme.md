# Wa RISCV64 模拟器

这是一个非常简陋的实现, 主要用于开发中测试一些代码片段, 因此内置了一些基础的调试命令.

模拟器启动时需要一个特殊的可执行的 elf 文件, 比如有以下汇编程序:

```s
    .section .rodata
message:
    .asciz "Hello RISC-V Baremetal!\n"

    .section .text
    .globl _start

# QEMU virt 机器 UART0 和 exit device 的基地址
UART0      = 0x10000000
EXIT_DEVICE = 0x100000

_start:
    # a0 = 字符串地址
    auipc   a0, %pcrel_hi(message)     # 高20位 = 当前PC + 偏移
    addi    a0, a0, %pcrel_lo(_start)  # 低12位

print_loop:
    lbu  a1, 0(a0)        # 取一个字节
    beq  a1, x0, finished # 如果是0则结束

    # t0 = UART0 地址
    lui     t0, %hi(UART0)           # UART0 高20位
    addi    t0, t0, %lo(UART0)       # UART0 低12位

    sb   a1, 0(t0)        # 写到UART寄存器
    addi a0, a0, 1        # 下一个字符
    jal  x0, print_loop

finished:
    # 写退出码 0 到 EXIT_DEVICE，让 QEMU 退出
    lui     t0, %hi(EXIT_DEVICE)     # exit device 地址
    addi    t0, t0, %lo(EXIT_DEVICE)

    # t1 = 0x5555
    # addi rd, rs1, imm 的 imm 范围是 [-2048, +2047](12 位有符号立即数)
    lui   t1, 0x5             # 高 20 位 (0x5 << 12 = 0x5000)
    addi  t1, t1, 0x555       # 结果 = 0x5000 + 0x555 = 0x5555

    sw   t1, 0(t0)

    # 如果 QEMU 不支持 exit 设备，就进入并死循环
forever:
    jal x0, forever
```

然后通过以下命令构建elf文件:

```
$ riscv64-unknown-elf-as hello.S -o hello.o
$ riscv64-unknown-elf-ld -Ttext=0x80000000 hello.o -o hello.elf.exe
```

需要注意的是以上构建命令设置了代码段的开始地址. 该地址是 Qemu 模拟器的默认地址, 因此构建成功后可以选择Qemu等模拟器测试下:

```
$ qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello.elf.exe
Hello RISC-V Baremetal!
```

如果正常就可以尝试通过 Wa RISCV64 模拟器运行:

```
$ wa wemu hello.elf.exe
Hello RISC-V Baremetal!
```

只要 elf 文件能够正常加载, 就可以进入调试模式执行:

```
$ wa wemu -d hello.elf.exe
Debug (enter h for help)...

Enter command: h
Commands are:
  h)elp           show help command list
  g)o             run instructions until power off
  s)tep  <n>      run n (default 1) instructions
  j)ump  <b>      jump to the b (default is current location)
  r)egs           print the contents of the registers
  i)Mem  <b <n>>  print n iMem locations starting at b
  d)Mem  <b <n>>  print n dMem locations starting at b
  a)lter <b <v>>  change the memory value at b
  t)race          toggle instruction trace
  p)rint          toggle print of total instructions executed
  c)lear          reset VM
  q)uit           exit

Enter command:
```

比如可以查看寄存器和当前指令信息:

```
Enter command: i
mem[80000000]: AUIPC A0, 0x0
Enter command: r
PC = 0x80000000
REG_X[0] = 0x00000000
REG_X[1] = 0x00000000
REG_X[2] = 0x00000000
REG_X[3] = 0x81000000
...
Enter command:
```

通过`go`命令执行, `q`命令退出:

```
Enter command: go
Hello RISC-V Baremetal!
Enter command: q
$
```

