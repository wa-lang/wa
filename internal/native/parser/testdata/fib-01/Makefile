# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

linux:
	$(WA) build -arch=x64 -target=linux -o a.out.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.c a.out.wa.s -o a.out1.exe
	./a.out1.exe

ast:
	$(WA) build -arch=x64 -target=linux -o a.out.wa.s fib-01.wat
	$(WA) ast -arch=x64 a.out.wa.s > a.out1.s
	gcc -nostdlib -static env-linux.c a.out1.s -z noexecstack -o a.out1.exe
	./a.out1.exe

loong64:
	$(WA) build -arch=loong64 -target=linux -o a.out.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.c a.out.wa.s -o a.out1.exe
	./a.out1.exe

clean:
