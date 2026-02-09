# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) asm2elf -arch=loong64 issue260209.wa.s
	$(WA) objdump -prog issue260209.wa.elf.exe > 11.out.txt
	./issue260209.wa.elf.exe

gcc:
	gcc -nostdlib -static -z noexecstack -o issue260209.wa.elf.exe issue260209.wa.s
	$(WA) objdump -prog issue260209.wa.elf.exe > 00.out.txt
	./issue260209.wa.elf.exe

clean:
