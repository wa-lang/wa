# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	$(WA) build -arch=x64 -target=windows wa-01.wa
	wa-01.exe

wz:
	$(WA) build -arch=x64 -target=windows wa-01.wz
	wa-01.exe

linux:
	$(WA) build -arch=x64 -target=linux wa-01.wa
	./wa-01.exe

linux-wz:
	$(WA) build -arch=x64 -target=linux wa-01.wz
	./wa-01.exe

loong64:
	$(WA) build -arch=loong64 -target=linux wa-01.wa
	./wa-01.exe

loong64-wz:
	$(WA) build -arch=loong64 -target=linux wa-01.wz
	./wa-01.exe

clean:
