# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

build:
	gcc -nostdlib -static -z noexecstack hello.s
	./a.out
	echo $$?

clean:
