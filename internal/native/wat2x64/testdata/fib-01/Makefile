# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=x64 -target=windows fib-01.wat
	gcc env.c fib-01.wa.s -o a.out1.exe
	a.out1.exe

linux:
	$(WA) build -arch=x64 -target=linux fib-01.wat
	gcc env.c fib-01.wa.s -o a.out1.exe
	./a.out1.exe

loong64:
	$(WA) build -arch=loong64 -target=linux fib-01.wat
	gcc env.c fib-01.wa.s -o a.out1.exe
	./a.out1.exe

clean:
