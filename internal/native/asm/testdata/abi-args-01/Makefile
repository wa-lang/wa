# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=loong64 -target=linux abi-args-01.wat
	cat env-linux-loong64.s abi-args-01.wa.s > abi-args-01-all.wa.s
	$(WA) asm2elf -arch=loong64 abi-args-01-all.wa.s
	$(WA) objdump -prog abi-args-01-all.wa.elf.exe > 11.out.txt
	./abi-args-01-all.wa.elf.exe

gcc:
	$(WA) build -arch=loong64 -target=linux abi-args-01.wat
	gcc -nostdlib -static -z noexecstack env-linux.c abi-args-01.wa.s -o a.out1.exe
	$(WA) objdump -prog a.out1.exe > 00.out.txt
	./a.out1.exe

gcc2:
	$(WA) build -arch=loong64 -target=linux abi-args-01.wat
	gcc -nostdlib -static -z noexecstack env-linux-loong64.s abi-args-01.wa.s -o a.out1.exe
	$(WA) objdump -prog a.out1.exe > 00.out.txt
	./a.out1.exe

clean:
