# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

windows:
	$(WA) build -arch=x64 -target=windows wa-02.wa
	wa-02.exe

linux:
	$(WA) build -arch=x64 -target=linux wa-02.wa
	./wa-02.exe

loong64:
	$(WA) build -arch=loong64 -target=linux wa-02.wa
	./wa-02.exe

run:
	$(WA) run wa-02.wa

clean:
