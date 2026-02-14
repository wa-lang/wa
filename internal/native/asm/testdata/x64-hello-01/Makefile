# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

build:
	gcc -nostdlib -static -z noexecstack hello.s
	./a.out
	echo $$?

dev:
	$(WA) asm2elf -arch=x64 -o a.out.elf hello.s
	$(WA) objdump -prog a.out.elf > 11.out.txt
	./a.out.elf

objdump:
	$(WA) objdump -prog a.out > 00.out.txt

clean:
