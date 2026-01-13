# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

golden:
	gcc golden-win64.s -o a.out0.exe
	a.out0.exe

dev:
	$(WA) build -arch=x64 -target=windows hello-02.wat

clean:
