# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

windows:
	gcc env-win64.c golden-win64.s -o a.out0.exe
	a.out0.exe

linux:
	gcc -nostdlib -static env-linux-x64.s golden-linux.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=x64 -target=windows memory-01.wat
	gcc env-win64.c memory-01.wa.s -o a.out1.exe
	a.out1.exe

dev-linux:
	$(WA) build -arch=x64 -target=linux -o memory-01.linux.wa.s memory-01.wat
	gcc -nostdlib -static env-linux-x64.s memory-01.linux.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
