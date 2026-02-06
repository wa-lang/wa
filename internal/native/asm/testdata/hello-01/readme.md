# 龙芯64 可执行文件

该测试用官方的gcc生成, 可以正常执行. 该测试用于自研汇编器的输出elf文件的对比参考.

## 示例程序: GAS 语法

该程序gcc和凹语义汇编器均可编译.

```s
.section .data
.align 3
.app.hello.str: .ascii "hello\n\000"
.app.hello.len: .quad 6

.section .text
.globl _start
_start:
    # main
    pcalau12i $t0, %pc_hi20(main)
    addi.d $t0, $t0, %pc_lo12(main)
    jirl $ra, $t0, 0

    # exit
    addi.d $a7, $zero, 93 # sys_exit
    syscall 0

.section .text
.globl main
main:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0

    # write(stdout, str, len)
    addi.d    $a0, $zero, 1 # arg.0 stdout
    pcalau12i $a1, %pc_hi20(.app.hello.str) # arg.1: ptr
    addi.d    $a1, $a1, %pc_lo12(.app.hello.str)
    pcalau12i $a2, %pc_hi20(.app.hello.len) # arg.2: len
    addi.d    $a2, $a2, %pc_lo12(.app.hello.len)
    ld.d      $a2, $a2, 0
    addi.d    $a7, $zero, 64 # sys_write
    syscall   0

    # return 0
    addi.d $a0, $zero, 0

    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0
```

## 示例程序: 中文汇编语法

中文版本的语义和GAS语法的程序完全一致, 输出结果也应该完全一致.

```
全局 信札: 字串 = "hello\n"

函数 _启动:
    # 主控
    计齐加高12立 $暂甲格, %相对.高20(主控)
    加立.长 $暂甲格, $暂甲格, %相对.低12(主控)
    链接跳转 $回格, $暂甲格, 0

    # exit
    加立.长 $参辛格, $零格, 93 # sys_exit
    系统调用 0
完毕

函数 主控:
    加立.长  $栈格, $栈格, -16
    存储.长  $回格, $栈格, 8
    存储.长  $帧格, $栈格, 0
    加立.长  $帧格, $栈格, 0

    # write(stdout, str, len)
    加立.长      $参甲格, $零格, 1 # arg.0 stdout
    计齐加高12立 $参乙格, %相对.高20(信札) # arg.1: ptr
    加立.长      $参乙格, $参乙格, %相对.低12(信札)
    加立.长      $参丙格, $零格, %内存字节数(信札) # arg.2: len
    装载.长      $参丙格, $参丙格, 0
    加立.长      $参辛格, $零格, 64 # sys_write
    系统调用     0

    # return 0
    加立.长 $参甲格, $零格, 0

    加立.长  $栈格, $帧格, 0
    装载.长  $回格, $栈格, 8
    装载.长  $帧格, $栈格, 0
    加立.长  $栈格, $栈格, 16
    链接跳转 $零格, $回格, 0
完毕
```

## 输出的ELF文件

```
$ wa objdump -prog a.out.elf
Class  : ELFCLASS64
Version: 1
OS/ABI : UNIX System V ABI
Machine: EM_LOONGARCH
Entry  : 120000144

[.text.]
120000144: 1A00000C # pcalau12i $t0, 0
120000148: 02C5618C # addi.d    $t0, $t0, 344
12000014C: 4C000181 # jirl      $ra, $t0, 0
120000150: 02C1740B # addi.d    $a7, $zero, 93
120000154: 002B0000 # syscall   0
120000158: 02FFC063 # addi.d    $sp, $sp, -16
12000015C: 29C02061 # st.d      $ra, $sp, 8
120000160: 29C00076 # st.d      $fp, $sp, 0
120000164: 02C00076 # addi.d    $fp, $sp, 0
120000168: 02C00404 # addi.d    $a0, $zero, 1
12000016C: 1A000205 # pcalau12i $a1, 16
120000170: 02C680A5 # addi.d    $a1, $a1, 416
120000174: 1A000206 # pcalau12i $a2, 16
120000178: 02C69CC6 # addi.d    $a2, $a2, 423
12000017C: 28C000C6 # ld.d      $a2, $a2, 0
120000180: 02C1000B # addi.d    $a7, $zero, 64
120000184: 002B0000 # syscall   0
120000188: 02C00004 # addi.d    $a0, $zero, 0
12000018C: 02C002C3 # addi.d    $sp, $fp, 0
120000190: 28C02061 # ld.d      $ra, $sp, 8
120000194: 28C00076 # ld.d      $fp, $sp, 0
120000198: 02C04063 # addi.d    $sp, $sp, 16
12000019C: 4C000020 # jirl      $zero, $ra, 0

[.data.]   00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
1200101A0: 68 65 6C 6C 6F 0A 00 06 00 00 00 00 00 00 00    hello..........
```

中文指令格式反汇编:

```
$ wa objdump -prog -zh a.out.elf
Class  : ELFCLASS64
Version: 1
OS/ABI : UNIX System V ABI
Machine: EM_LOONGARCH
Entry  : 120000144

[.text.]
120000144: 1A00000C # 计齐加高12立 $暂甲格, 0
120000148: 02C5618C # 加立.长      $暂甲格, $暂甲格, 344
12000014C: 4C000181 # 链接跳转     $回格, $暂甲格, 0
120000150: 02C1740B # 加立.长      $参辛格, $零格, 93
120000154: 002B0000 # 系统调用     0
120000158: 02FFC063 # 加立.长      $栈格, $栈格, -16
12000015C: 29C02061 # 存储.长      $回格, $栈格, 8
120000160: 29C00076 # 存储.长      $帧格, $栈格, 0
120000164: 02C00076 # 加立.长      $帧格, $栈格, 0
120000168: 02C00404 # 加立.长      $参甲格, $零格, 1
12000016C: 1A000205 # 计齐加高12立 $参乙格, 16
120000170: 02C680A5 # 加立.长      $参乙格, $参乙格, 416
120000174: 1A000206 # 计齐加高12立 $参丙格, 16
120000178: 02C69CC6 # 加立.长      $参丙格, $参丙格, 423
12000017C: 28C000C6 # 装载.长      $参丙格, $参丙格, 0
120000180: 02C1000B # 加立.长      $参辛格, $零格, 64
120000184: 002B0000 # 系统调用     0
120000188: 02C00004 # 加立.长      $参甲格, $零格, 0
12000018C: 02C002C3 # 加立.长      $栈格, $帧格, 0
120000190: 28C02061 # 装载.长      $回格, $栈格, 8
120000194: 28C00076 # 装载.长      $帧格, $栈格, 0
120000198: 02C04063 # 加立.长      $栈格, $栈格, 16
12000019C: 4C000020 # 链接跳转     $零格, $回格, 0

[.data.]   00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
1200101A0: 68 65 6C 6C 6F 0A 00 06 00 00 00 00 00 00 00    hello..........
```
