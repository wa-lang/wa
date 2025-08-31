# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

riscv64-run:
	riscv64-unknown-elf-as hello.S -o hello.o
	riscv64-unknown-elf-ld -Ttext=0x80000000 hello.o -o hello.elf.exe
	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello.elf.exe

clean:
	-rm hello.o hello.elf.exe
