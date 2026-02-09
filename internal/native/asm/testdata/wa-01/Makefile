# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=loong64 -target=linux wa-01.wa
	$(WA) asm2elf -arch=loong64 wa-01.wa.s
	./wa-01.wa.elf.exe

gcc:
	$(WA) build -arch=loong64 -target=linux wa-01.wa
	$(WA) asm2elf -arch=loong64 wa-01.wa.s
	$(WA) ast -arch=loong64 wa-01.wa.s > wa-01.out2.s
	gcc -nostdlib -static -z noexecstack -o wa-01.out2.exe wa-01.out2.s
	./wa-01.out2.exe

loong64:
	$(WA) build -arch=loong64 -target=linux wa-01.wa
	./wa-01.exe

loong64-wz:
	$(WA) build -arch=loong64 -target=linux wa-01.wz
	./wa-01.exe

clean:
