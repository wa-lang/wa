# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

loong64:
	$(WA) build -arch=loong64 -target=linux wa-01.wa
	$(WA) ast -arch=loong64 wa-01.wa.s > wa-01.out.s

x64:
	$(WA) build -arch=x64 -target=linux wa-01.wa
	$(WA) ast -arch=x64 wa-01.wa.s > wa-01.out.s
	gcc -nostdlib -static wa-01.out.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
