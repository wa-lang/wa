# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) asm2elf -arch=loong64 fib-01-all.wa.s
	./fib-01-all.wa.elf.exe

gcc:
	gcc -nostdlib -static -z noexecstack -o fib-01-all.wa.elf.exe fib-01-all.wa.s
	./fib-01-all.wa.elf.exe

clean:
