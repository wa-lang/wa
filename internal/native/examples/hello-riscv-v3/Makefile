# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

run:
	$(WA) asm2elf -arch=riscv32 hello_riscv64.ws
	$(WA) wemu hello_riscv64.ws.elf.exe
	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello_riscv64.ws.elf.exe

run32:
	$(WA) asm2elf -arch=riscv32 hello_riscv32.ws
	$(WA) wemu hello_riscv32.ws.elf.exe
	qemu-system-riscv32 -machine virt -nographic -bios none -kernel hello_riscv32.ws.elf.exe

run-zh:
	$(WA) asm2elf hello_riscv64_zh.ws
	$(WA) wemu hello_riscv64_zh.ws.elf.exe

run32-zh:
	$(WA) asm2elf -arch=riscv32 hello_riscv32_zh.ws
	$(WA) wemu hello_riscv32_zh.ws.elf.exe

fmt:
	$(WA) fmt -riscv hello_riscv64_zh.ws > a.out.ws

debug:
	$(WA) -debug asm2elf hello_riscv64.ws

clean:
	-rm hello.elf.exe
