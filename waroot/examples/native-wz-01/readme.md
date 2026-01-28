# 凹语言构建本地程序

说明: 这是一个处于非常早期的实验性的特性, 只面向龙芯和X64芯片的操作系统. 在测试前请先确保本地安装了对应的 GCC, 用于目标汇编代码的连接.

# Windows/X64

```
C:\...\> chcp 65001
C:\...\> wa build -arch=x64 -target=windows hello.wz
C:\...\> ./hello.exe
你好, 2026 年!
```

其中 `chcp 65001` 是将命令行切换到 UTF8 编码, 否则输出的中文会出现乱码.

# Windows/Ubuntu

参考 $ROOT/docs/wsl.md 配置好 WSL 环境, 安装并进入 linux 环境. 切换到对应的目录:

```
$ wa build -arch=x64 -target=linux hello.wz
$ ./hello.exe
你好, 2026 年!
```

# 龙芯64/Linux

切换到对应的目录:

```
$ wa build -arch=loong64 -target=linux hello.wz
$ ./hello.exe
你好, 2026 年!
```
