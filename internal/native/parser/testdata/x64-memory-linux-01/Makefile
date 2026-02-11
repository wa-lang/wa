# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

linux:
	gcc -nostdlib -static env-linux-x64.s golden-linux.s -o a.out0.exe -z noexecstack
	./a.out0.exe

ast:
	$(WA) ast -arch=x64 golden-linux.s > a.golden.out.s
	gcc -nostdlib -static env-linux-x64.s a.golden.out.s -z noexecstack -o a.out1.exe
	./a.out1.exe

	$(WA) build -arch=x64 -target=linux -o a.out.wa.s memory-01.wat
	$(WA) ast -arch=x64 a.out.wa.s > a.out1.s
	gcc -nostdlib -static a.out1.s -z noexecstack -o a.out1.exe
	./a.out1.exe

dev-linux:
	$(WA) build -arch=x64 -target=linux -o a.out.wa.s memory-01.wat
	gcc -nostdlib -static a.out.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
