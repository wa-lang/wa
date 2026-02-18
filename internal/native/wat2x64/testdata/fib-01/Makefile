# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

FLAGS:=-nostdlib
FLAGS+=-Wl,-e,_start
FLAGS+=-Wl,--subsystem,console
FLAGS+=-Wl,--file-alignment,16
FLAGS+=-Wl,--section-alignment,16
FLAGS+=-Wl,--build-id=none
FLAGS+=-Wl,--disable-dynamicbase
FLAGS+=-s
FLAGS+=-lkernel32

windows:
	$(WA) build -arch=x64 -target=windows fib-01.wat
	gcc env-win64.s fib-01.wa.s $(FLAGS) -o a.out1.exe
	a.out1.exe

linux:
	$(WA) build -arch=x64 -target=linux -o fib-01.linux.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.s fib-01.linux.wa.s -o a.out1.exe
	./a.out1.exe

loong64:
	$(WA) build -arch=loong64 -target=linux -o fib-01.loong64.wa.s fib-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.c fib-01.loong64.wa.s -o a.out1.exe
	./a.out1.exe

clean:
