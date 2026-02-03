# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../../main.go

run:
	gcc env.s app.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) asm2elf -arch=loong64 app.s

clean:
