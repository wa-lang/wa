# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

loong64:
	gcc env.c golden-loong64.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=loong64 -target=linux -o table-01.wa.s table-01.wat
	gcc env.c table-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
