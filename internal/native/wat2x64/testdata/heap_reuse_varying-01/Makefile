# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=x64 -target=windows heap_reuse_varying-01.wat
	gcc env.c heap_reuse_varying-01.wa.s -o a.out1.exe
	a.out1.exe

debug:
	$(WA) -trace=wat2x64 build -arch=x64 -target=windows heap_reuse_varying-01.wat

linux:
	$(WA) build -arch=x64 -target=linux -o heap_reuse_varying-01.linux.wa.s heap_reuse_varying-01.wat
	gcc env.c heap_reuse_varying-01.linux.wa.s -o a.out1.exe
	./a.out1.exe

loong64:
	$(WA) build -arch=loong64 -target=linux -o heap_reuse_varying-01.loong64.wa.s heap_reuse_varying-01.wat
	gcc env-loong64.c heap_reuse_varying-01.loong64.wa.s -o a.out1.exe
	./a.out1.exe

clean:
