# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../main.go

build-windows:
	$(WA) build -arch=x64 -target=windows hello.wa

build-linux-x64:
	$(WA) build -arch=x64 -target=linux hello.wa

build-loong64:
	$(WA) build -arch=loong64 -target=linux hello.wa

run:
	$(WA) run hello.wa

clean:
