# 龙芯 测试

安装工具链和qemu:

- https://www.loongnix.cn/zh/toolchain/GNU/
  - i686 Windows (MinGW-w64) 二进制
  - loongson-gnu-toolchain-8.3-i686-mingw-loongarch64-linux-gnu-rc1.6.zip
- https://qemu.weilnetz.de/w64/
  - qemu-w64-setup-20250814.exe

## 代码段

TODO

## 数据段

TODO

<!--
qemu-system-loongarch64 -machine virt -cpu loongarch64 -nographic -serial mon:stdio -kernel hello_loong64.ws.elf.exe
entry point: 0x80000000
uart 0x10000000
-->

