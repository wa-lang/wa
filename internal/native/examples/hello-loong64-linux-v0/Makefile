# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

default:
	gcc -nostdlib -nostartfiles -static hello_la64.s -o hello_la64.out.exe
	./hello_la64.out.exe

dev:
	gcc -O2 --save-temps -nostdlib -nostartfiles -static start.s print.s -o a.out.exe
	./a.out.exe

objdump:
	$(WA) objdump -prog hello_la64.out.exe

info:
	go run hello.go -f=hello_la64.out.exe -l
	go run hello.go -f=hello_la64.out.exe -s=.text
	go run hello.go -f=hello_la64.out.exe -s=.data

info2:
	objdump -d hello_la64.out.exe

clean:
	-rm *.out*
