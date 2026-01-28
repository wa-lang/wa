# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

windows:
	$(WA) build -arch=x64 -target=windows wz-01.wz
	wz-01.exe

linux:
	$(WA) build -arch=x64 -target=linux wz-01.wz
	./wz-01.exe

loong64:
	$(WA) build -arch=loong64 -target=linux wz-01.wz
	./wz-01.exe

run:
	$(WA) run wz-01.wz

clean:
