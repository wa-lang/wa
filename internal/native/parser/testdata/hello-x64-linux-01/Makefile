# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

dev:
	gcc -nostdlib -static -z noexecstack app.wa.s
	./a.out

ast:
	$(WA) ast -arch=x64 app.wa.s > a.out1.s
	gcc -nostdlib -static -z noexecstack a.out1.s
	./a.out

clean:
