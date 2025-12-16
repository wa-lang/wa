# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

default:
	gcc -O2 --save-temps main.c print.s -o a.out.exe
	./a.out.exe

dev:
	as hello_la64.s -o hello_la64.o
	ld -static hello_la64.o -o hello_la64.out.exe
	./hello_la64.out.exe

objdump:
	$(WA) objdump -prog a.out.exe

clean:
	-rm *.out*
