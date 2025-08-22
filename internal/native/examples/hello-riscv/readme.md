# RISCV 测试

安装工具链和qemu:

- https://github.com/sifive/freedom-tools/releases
  - riscv64-unknown-elf-toolchain-10.2.0-2020.12.8-x86_64-w64-mingw32.zip
- https://qemu.weilnetz.de/w64/
  - qemu-w64-setup-20250814.exe

## Makeilfe

```makefile
riscv64-run:
	riscv64-unknown-elf-as hello.S -o hello.o
	riscv64-unknown-elf-ld -Ttext=0x80000000 hello.o -o hello.elf.exe
	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello.elf.exe

riscv32-run:
	riscv32-unknown-elf-as hello.S -o hello.o
	riscv32-unknown-elf-ld -Ttext=0x80000000 hello.o -o hello.elf.exe
	qemu-system-riscv32 -machine virt -nographic -bios none -kernel hello.elf.exe

clean:
	-rm hello.o hello.elf.exe
```

## 代码段

```
Name     : .text
Addr     : 0x80000000
Addralign: 4
Offset   : 0x00001000
Size     : 60
FileSize : 60
--------------------------------------------------------
         00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
00001000 17 05 00 00 13 05 C5 03 83 45 05 00 63 8C 05 00
00001010 B7 02 00 10 93 82 02 00 23 80 B2 00 13 05 15 00
00001020 6F F0 9F FE B7 02 10 00 93 82 02 00 37 53 00 00
00001030 13 03 53 55 23 A0 62 00 6F 00 00 00
```

## 数据段

```
Name     : .data
Addr     : 0x80001055
Addralign: 1
Offset   : 0x00001055
Size     : 11
FileSize : 11
--------------------------------------------------------
         00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
00001055 00 00 00 00 00 00 00 00 00 00 00
```

```
Name     : .rodata
Addr     : 0x8000003c
Addralign: 1
Offset   : 0x0000103c
Size     : 25
FileSize : 25
--------------------------------------------------------
         00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
0000103c 48 65 6C 6C 6F 20 52 49 53 43 2D 56 20 42 61 72
0000104c 65 6D 65 74 61 6C 21 0A 00
```

```
Name     : .sdata
Addr     : 0x80001060
Addralign: 1
Offset   : 0x00001060
Size     : 0
FileSize : 0
```
