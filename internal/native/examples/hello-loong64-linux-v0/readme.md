# 龙芯汇编示例

- SP 寄存器必须 16 字节对齐

## 打印字符串

构建 `hello_la64.s`文件:

```s
.data
msg:
    .asciz "Hello, LoongArch 64!\n"
len = . - msg

.text
.align 2
.globl _start

_start:
    # SYS_write(STDOUT, msg, len)
    addi.d    $a0, $zero, STDOUT
    pcalau12i $a1, %pc_hi20(msg)
    addi.d    $a1, $a1, %pc_lo12(msg)
    addi.d    $a2, $zero, len
    addi.d    $a7, $zero, 64
    syscall   0

    # SYS_exit(0)
    addi.d    $a0, $zero, 0
    addi.d    $a7, $zero, 93
    syscall   0
```

编译并执行：

```bash
$ gcc -nostdlib -nostartfiles -static hello_la64.s -o hello_la64.out.exe
$ ./hello_la64.out.exe
Hello, LoongArch 64!
```

## 代码段

```
Name     : .text
Addr     : 0x12000010c
Addralign: 4
Offset   : 0x0000010c
Size     : 36
FileSize : 36
--------------------------------------------------------
         00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F 
0000010c 04 04 C0 02 05 02 00 1A A5 C0 C4 02 06 58 C0 02 
0000011c 0B 00 C1 02 00 00 2B 00 04 00 C0 02 0B 74 C1 02 
0000012c 00 00 2B 00
```

指令解码:

```bash
$ objdump -d hello_la64.out.exe

hello_la64.out.exe：     文件格式 elf64-loongarch

Disassembly of section .text:

000000012000010c <_start>:
   12000010c:   02c00404        li.d            $a0, 1
   120000110:   1a000205        pcalau12i       $a1, 16
   120000114:   02c4c0a5        addi.d          $a1, $a1, 304
   120000118:   02c05806        li.d            $a2, 22
   12000011c:   02c1000b        li.d            $a7, 64
   120000120:   002b0000        syscall         0x0
   120000124:   02c00004        li.d            $a0, 0
   120000128:   02c1740b        li.d            $a7, 93
   12000012c:   002b0000        syscall         0x0
```

`wa objdump` 查看信息:

```bash
$ wa objdump -prog hello_la64.out.exe
Class  : ELFCLASS64
Version: 1
OS/ABI : UNIX System V ABI
Machine: EM_LOONGARCH
Entry  : 12000010c

[.text.]
12000010C: 02C00404 # ADDI.D A0, ZERO, 1
120000110: 1A000205 # PCALAU12I A1, 16
120000114: 02C4C0A5 # ADDI.D A1, A1, 304
120000118: 02C05806 # ADDI.D A2, ZERO, 22
12000011C: 02C1000B # ADDI.D A7, ZERO, 64
120000120: 002B0000 # SYSCALL 0
120000124: 02C00004 # ADDI.D A0, ZERO, 0
120000128: 02C1740B # ADDI.D A7, ZERO, 93
12000012C: 002B0000 # SYSCALL 0

[.data.]   00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
120010130: 48 65 6C 6C 6F 2C 20 4C 6F 6F 6E 67 41 72 63 68 Hello, LoongArch
120010140: 20 36 34 21 0A 00                                64!..
```

## 数据段

```
Name     : .data
Addr     : 0x120010130
Addralign: 1
Offset   : 0x00000130
Size     : 22
FileSize : 22
--------------------------------------------------------
         00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F 
00000130 48 65 6C 6C 6F 2C 20 4C 6F 6F 6E 67 41 72 63 68 
00000140 20 36 34 21 0A 00
```
