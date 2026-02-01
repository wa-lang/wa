# 凹汇编语言(The Wa Native Assembly Language) - GAS 风格

凹汇编语言用于描述本地的汇编语言程序, 目前仅支持龙芯架构. 以下是 LOGO:

```
+---+    +---+
| o |    | o |
|   +----+   |
|            |
|    NASM    |
|            |
+------------+
```

汇编代码的后缀名 `.wa.s` 对应 GAS 风格语法, `.wz.s` 对应中文语法, 英文名字缩写: NASM.

## 你好世界 - 龙芯版本

```s
.section .data
.align 3
.app.hello.str: .asciz "hello\n"
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

## 注释

`#` 开头的一行为注释

## 标识符

和C语言的标识符类似, 以字母下划线开头, 后续可以是字母和数字和下划线. 同时`.`和`$`也被视作字母.

## 常量

`A=123` 定义 `A` 常量, 值为 123. 常量一般在头部定义, 常量本身不占内存空间.

## 全局变量

全局变量需要绑定到数据段(`.data`和`.rodata`), 并且指令地址的对齐方式:

```s
.section .data
.align 3
.Table.funcIndexList.0: .quad .Import.syscall.write
```

表示在 `.data` 数据段定义了 `.Table.funcIndexList.0` 的全局变量, 对应8个字节内存空间, 初始值是 `.Import.syscall.write` 函数的地址. 对于龙芯, 变量的地址是 `2^3=8` 字节对齐(X64平台会有差异, 目前忽略).

全局变量有以下类型:

- `.byte`, 1个字节
- `.short`, 2个字节
- `.long`, 4个字节
- `.quad`, 8个字节
- `.ascii`, ASCII字符串
- `.asciz`, ASCII字符串, 在结尾添加 `\0`
- `.incbin`, 包含外部的数据文件

预留一定字节的空间:

```s
.section .data
.align 3
.some.table: .skip 100
```

以上是预留100字节的空间, 地址是`2^3=8`字节对齐.

包含外部的资源文件:

```s
.image.data: .incbin "lena.jpg"
```

## 定义函数

函数需要绑定到代码段(`.text`), 并且必须4字节对齐:

```s
# 内存初始化函数
.section .text
.align 3
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32
    ...
```

表示在 `.text` 代码段定义了 `.Wa.Memory.initFunc` 函数, 函数开始地址在龙芯平台为 `2^2=4` 字节对齐, 并且通过 `.globl` 导出了函数的名字. 最后以函数名开头的标号, 后面是函数的指令.

## 特殊情况

凹语言自带的汇编器不支持链接外部的符号, 对于 `.extern _Wa_Runtime_write` 语法凹语言汇编器可以解析但是不能链接, 这种情况可以用gcc汇编器.

对于未来可能会支持的X64汇编语言, 必须在开头加上 `.intel_syntax noprefix` 指令, 表示使用 intel 的指令语法风格, 同时去掉寄存器和立即数开头的`%`和`$`等前缀.

## 补充说明

凹汇编语言是 GAS 汇编语法的子集, 只选取了凹语言输出汇编和语言特性依赖的最小语法. 该子集确保可以被GCC的汇编器和凹语言自带的汇编器构建.

