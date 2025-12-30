# 最小的 elf 示例文件

代码调整, 改变返回值(elf文件待测试):

```
// $ clang -nostdlib -static tiny.s -o tiny

.globl _start
.text

_start:
	// push + pop is 3 bytes rather than 4 for the usual mov $60, %ax
	// 60 is sys_exit syscall
	pushq $60
	popq  %rax

	// mov $60, %ax

	// xorl %edi, %edi // RDI = 0
	movl $42, %edi

	syscall // sys_exit(RDI)
```

对应的二进制汇编指令:

```
// | asm             | .text          |
// | --------------- | -------------- |
// | pushq \$60      | 6a 3c          |
// | popq %rax       | 58             |
// | movl \$42, %edi | bf 2a 00 00 00 | movl $42, %rdi ???
// | syscall         | 0f 05          |
// =>
// 6a 3c 58 bf 2a 00 00 00 0f 05
// =>
0x3c6a, 0xbf58, 0x002a, 0x0000, 0x050f,
```

## Windows 10 环境测试

- 下载 `ubuntu-24.04.2-wsl-amd64.wsl`: https://ubuntu.com/desktop/wsl
- 安装 `wsl --import Ubuntu-24.04 C:\WSL\Ubuntu24 ubuntu-24.04.2-wsl-amd64.wsl`
- 启动 `wsl -d Ubuntu-24.04`
- 执行 `./tiny.elf.bin`
- 验证返回值 `echo $?`

ubuntn环境查看elf文件:

```
# xxd tiny.elf.bin
# readelf -h tiny.elf.bin
# readelf -l tiny.elf.bin
# objdump -d tiny.elf.bin
```

```
# xxd tiny.elf.bin
00000000: 7f45 4c46 0201 0100 0000 0000 0000 0000  .ELF............
00000010: 0200 3e00 0100 0000 7800 4000 0000 0000  ..>.....x.@.....
00000020: 4000 0000 0000 0000 0000 0000 0000 0000  @...............
00000030: 0000 0000 4000 3800 0100 0000 0000 0000  ....@.8.........
00000040: 0100 0000 0700 0000 7800 0000 0000 0000  ........x.......
00000050: 7800 4000 0000 0000 0000 0000 0000 0000  x.@.............
00000060: 0700 0000 0000 0000 0a00 0000 0000 0000  ................
00000070: 0010 0000 0000 0000 6a3c 58bf 2a00 0000  ........j<X.*...
00000080: 0f05                                     ..

# readelf -h tiny.elf.bin
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              EXEC (Executable file)
  Machine:                           Advanced Micro Devices X86-64
  Version:                           0x1
  Entry point address:               0x400078
  Start of program headers:          64 (bytes into file)
  Start of section headers:          0 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           56 (bytes)
  Number of program headers:         1
  Size of section headers:           0 (bytes)
  Number of section headers:         0
  Section header string table index: 0

# readelf -l tiny.elf.bin

Elf file type is EXEC (Executable file)
Entry point 0x400078
There is 1 program header, starting at offset 64

Program Headers:
  Type           Offset             VirtAddr           PhysAddr
                 FileSiz            MemSiz              Flags  Align
  LOAD           0x0000000000000078 0x0000000000400078 0x0000000000000000
                 0x0000000000000007 0x000000000000000a  RWE    0x1000

# objdump -d tiny.elf.bin

tiny.elf.bin:     file format elf64-x86-64
```
