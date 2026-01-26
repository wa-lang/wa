# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build wa-01.wa
	$(WA) build -arch=x64 -target=windows wa-01.wat
	gcc env.c wa-01.wa.s -o a.out1.exe
	a.out1.exe

wz:
	$(WA) build wa-01.wz
	$(WA) build -arch=x64 -target=windows wa-01.wat
	gcc env.c wa-01.wa.s -o a.out1.exe
	a.out1.exe

linux:
	$(WA) build wa-01.wa
	$(WA) build -arch=x64 -target=linux wa-01.wat
	gcc env.c wa-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

linux-wz:
	$(WA) build wa-01.wz
	$(WA) build -arch=x64 -target=linux wa-01.wat
	gcc env.c wa-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

loong64:
	$(WA) build wa-01.wa
	$(WA) build -arch=loong64 -target=linux wa-01.wat
	gcc env.c wa-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

loong64-wz:
	$(WA) build wa-01.wz
	$(WA) build -arch=loong64 -target=linux wa-01.wat
	gcc env.c wa-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
