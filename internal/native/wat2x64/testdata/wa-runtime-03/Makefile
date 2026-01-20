# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=x64 -target=windows wa-runtime-03.wat
	gcc env.c wa-runtime-03.wa.s -o a.out1.exe
	a.out1.exe

linux:
	$(WA) build -arch=x64 -target=linux wa-runtime-03.wat
	gcc env.c wa-runtime-03.wa.s -o a.out1.exe
	./a.out1.exe

clean:
