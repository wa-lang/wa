# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

ast:
	$(WA) ast -arch=loong64 app.wa.s

ast-zh:
	$(WA) ast -arch=loong64 app.wz.s

build:
	$(WA) asm2elf -arch=loong64 app.wa.s

build-zh:
	$(WA) asm2elf -arch=loong64 app.wz.s

golden:
	gcc -nostdlib -static app.wa.s -o a.out0.exe -z noexecstack
	./a.out0.exe

objdump:
	$(WA) objdump -prog a.out.elf

objdump-zh:
	$(WA) objdump -prog -zh a.out.elf

clean:
