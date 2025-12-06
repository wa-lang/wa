# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

dev:
	$(WA) -debug asm2elf -arch=loong64 hello_loong64.ws

debug-run:
	$(WA) wemu -debug hello_loong64.ws.elf.exe

objdump:
	$(WA) objdump -prog hello_loong64.ws.elf.exe

qemu-run:
	qemu-system-loongarch64 -machine virt -nographic -serial mon:stdio -kernel hello_loong64.ws.elf.exe

asm2elf:
	$(WA) asm2elf -arch=loong64 hello_loong64.ws

run:
	$(WA) wemu hello_loong64.ws.elf.exe

clean:
	-rm hello.o hello.elf.exe
