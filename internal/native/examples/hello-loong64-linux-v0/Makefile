# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

build:
	as hello_la64.s -o hello_la64.o
	ld -static hello_la64.o -o hello_la64.out.exe
	./hello_la64.out.exe

dev:
	gcc --save-temps test.c add.s -o test_add.out.exe
	./test_add.out.exe

objdump:
	$(WA) objdump -prog test_add.out.exe

clean:
	-rm *.out.exe
