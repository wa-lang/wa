# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

loong64-run:
	loongarch64-linux-gnu-as hello.S -o hello.o
	loongarch64-linux-gnu-ld -nostdlib -nostartfiles -Ttext=0x80000000 hello.o -o hello.elf.exe
	qemu-system-loongarch64 -machine virt -nographic -bios none -kernel hello.elf.exe

clean:
	-rm hello.o hello.elf.exe
