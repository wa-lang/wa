# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

run:
	$(WA) asm2elf -arch=loong64 hello_loong64.ws
	$(WA) wemu hello_loong64.ws.elf.exe

run-zh:
	$(WA) asm2elf -arch=loong64 hello_loong64_zh.ws
	$(WA) wemu hello_loong64_zh.ws.elf.exe

debug-run:
	$(WA) wemu -debug hello_loong64.ws.elf.exe

objdump:
	$(WA) objdump -prog hello_loong64.ws.elf.exe

qemu-run:
	qemu-system-loongarch64 -machine virt -nographic -serial mon:stdio -kernel hello_loong64.ws.elf.exe

asm2elf:
	$(WA) asm2elf -arch=loong64 hello_loong64.ws

clean:
	-rm hello.o hello.elf.exe
