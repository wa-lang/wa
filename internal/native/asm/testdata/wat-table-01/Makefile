# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

build:
	$(WA) build -arch=loong64 -target=linux -o table-01.wa.s table-01.wat
	$(WA) asm2elf -arch=loong64 table-01.wa.s
	./table-01.wa.elf.exe

gcc:
	$(WA) build -arch=loong64 -target=linux -o table-01.wa.s table-01.wat
	gcc -nostdlib -static table-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
