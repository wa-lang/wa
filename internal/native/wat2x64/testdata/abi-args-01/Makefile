# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

golden:
	gcc env.c golden-win64.s -o a.out0.exe
	a.out0.exe

linux:
	gcc env.c golden-linux.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=x64 -target=windows abi-args-01.wat
	gcc env.c abi-args-01.wa.s -o a.out1.exe
	a.out1.exe

clean:
