# Wa Native 汇编语言

主要是内部用于表示 RISCV 等本地的汇编语言程序.

## 注释

- 支持 `#` 开头的单行注释

## 语句

- 语句以分号分割
- 每个换行符号可能会自动插入分号
- 同一个指令不会跨越多行

## 名字

- 关键字位小写字母, 比如 `func` 表示定义函数, `global` 表示定义全局变量
- 全局标识符以 `$` 开头, 比如 `$main` 表示 main 函数
- 局部标识符以 `%` 开头, 比如 `%a` 表示局部变量名
- 指令不区分大小写, 比如 `CALL` 和 `call` 等价
- 寄存器不区分大小写, 比如 `PC` 和 `pc` 等价
- `_.$`和字母数字都是合法的字符

## 关键字

- `i32`: `int32`
- `i64`: `int64`
- `f32`: `float32`
- `f64`: `float63`
- `const`: 定义全局常量
- `global`: 定义全局变量
- `local`: 定义局部变量
- `func`: 定义函数
- 寄存器: 平台相关定义
- 指令: 平台相关定义

## 字面值

- 字符常量: `'A'`/`'1'`, 支持 ASCII 单字节的字符
- 整数常量: `123`, 数字部分支持二进制/十进制/十六进制等
- 浮点数常量: `f32(1.5)`/`f64(123.456)`, 数字部分支持普通浮点数/科学计数法/十六进制等
- 字符串常量: 支持单行字符串, 转义规则参考凹语言
- 地址常量: 标识符表示的地址, `i64` 类型
- 默认整数类型: `123` 为 `i32`
- 默认浮点数类型: `12.3` 为 `f32`
- 不支持表达式: 不支持 `1+2` 写法

## 全局常量

- 常量只有整数和浮点数

```go
# QEMU virt 机器 UART0 和 exit device 的基地址
const $UART0       = 0x10000000
const $EXIT_DEVICE = 0x100000 
```

## 全局变量

```go
# 默认的i32
global $age: i32 = 5
global $ch: i32 = 'A'

# 默认的f32
global $size: f32 = 12.3

# 整数, 4字节, 12345678 
global $i32: i32 = 12345678

# 整数, 8字节, 0x12345678
global $i64: i64 = 0x12345678

# float32
global $f32: f32 = 12.5

# float64
global $f64: f64 = 12.34567

# utf8 编码的字符串, 结尾自动添加 `\0`, 地址和长度自动对齐到 8 字节
global $name = "wa native assembly language"
```

全局变量可以指定更大的长度:

```go
global $f32: 20 = f32(12.5)

# 指定长度的字符串, 不足部分补充0
global $str: 100 = "abc"
```

分段初始化(类似结构体):

```go
# 1024 字节结构体
# pos: value, 只能是用基础常量面值初始化
global $info: 1024 = {
    5: "abc",    # 从第5字节开始 `abc\0`
    9: i32(123), # 从第9字节开始
}
```

用常量初始化:

```go
global $ptr = $UART0 # 常量没有地址
```

用`$info`地址初始化:

```go
global $ptr = $info # $info 的地址
```

## 函数和寄存器

- `SB`: 数据段的开始位置, 用于全局变量定位
- `PC`: 当前指令位置, 用于 Label 等局部定位
- `FP`: 函数帧指针, 用于参数/返回值/局部变量定位
- `SP`: 函数栈指针, 用于临时变量, 调用参数等

定义函数:

```go
# 参数和返回值根据调用约定对齐
# 参数和返回值只支持基础的数据类型, 不支持结构体和数组值传入
func $add(%a:i32, %b:i32, %c:i32) => f64 {
    local %d: i32 # 局部变量必须先声明, i32 大小的空间

    # 指令
Loop:
}
```

## 地址表达式

- `MOVQ X1, $info   # X0 = $info`, 加载全局地址
- `MOVQ X1, $info+2 # X0 = $info`, 加载全局地址, 加偏移量
- `MOVQ %e, X1      # e[0] = X1`, 局部变量地址
- `MOVQ %e+10, X1   # e[9] = X1`, 局部变量地址, 加偏移量
- `MOVQ X2, Loop`, 局部标号

## 指令格式

- `ADD X1, X2, X3 # X1 = X2 + X3`

## 例子

QEMU 裸机输出字符串的例子:

```go
# QEMU virt 机器 UART0 和 exit device 的基地址
const $UART0       = 0x10000000
const $EXIT_DEVICE = 0x100000

# 字符串
global $message = "Hello RISC-V Baremetal!\n"

# 主函数
func $main() {
%start:
    # a0 = 字符串地址
    MOVQ a1, $message

%print_loop:
    lbu  a1, 0(a0)        # 取一个字节
    beq  a1, x0, finished # 如果是0则结束
   
    MOVQ t0, $UART0       # t0 = $UART0 地址
    sb   a1, 0(t0)        # 写到UART寄存器
    addi a0, a0, 1        # 下一个字符
    jal  x0, %print_loop

%finished:
    # 写退出码 0 到 EXIT_DEVICE，让 QEMU 退出
    MOVQ t0, $EXIT_DEVICE  # t0 = $EXIT_DEVICE
    MOVQ t1, 0x5555        # t1 = 0x5555
    sw t1, 0(t0)           # [t0] = t1

%forever:
    # 如果 QEMU 不支持 exit 设备，就进入并死循环
    jal x0, %forever
}
```
