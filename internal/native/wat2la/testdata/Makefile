# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

wat2la:
	$(WA) wat2la a_step0.wat

la64-build:
	$(WA) asm2elf -arch=loong64 hello.s
	$(WA) objdump hello.elf.exe

la64-run:
	./hello.elf.exe

gcc:
	gcc --save-temps hello-gcc.c -o a.out.exe
	./a.out.exe

gcc2:
	gcc --save-temps hello2-gcc.c -o a.out.exe
	./a.out.exe

clean:
	-rm hello.elf.exe a.out*
