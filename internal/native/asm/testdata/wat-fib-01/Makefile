# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=loong64 -target=linux -o fib-01.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux-loong64.s fib-01.wa.s -o a.out1.exe
	./a.out1.exe

build:
	$(WA) build -arch=loong64 -target=linux fib-01.wat
	cat env-linux-loong64.s fib-01.wa.s > fib-01-all.wa.s
	$(WA) asm2elf -arch=loong64 fib-01-all.wa.s
	./fib-01-all.wa.elf.exe

build-gcc:
	$(WA) build -arch=loong64 -target=linux fib-01.wat
	cat env-linux-loong64.s fib-01.wa.s > fib-01-all.wa.s
	gcc -nostdlib -static -z noexecstack -o fib-01-all.wa.elf.exe fib-01-all.wa.s
	./fib-01-all.wa.elf.exe

gcc:
	$(WA) build -arch=loong64 -target=linux -o fib-01.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.c fib-01.wa.s -o a.out1.exe
	./a.out1.exe

clean:
