# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

build:
	$(WA) asm2elf -arch=loong64 hello_loong64.ws

build-zh:
	$(WA) asm2elf -arch=loong64 hello_loong64_zh.ws

objdump:
	$(WA) objdump -prog hello_loong64.ws.elf.exe

run:
	./hello_loong64.ws.elf.exe

clean:
	-rm hello.o hello.elf.exe
