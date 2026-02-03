# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

ENC_S:=../../assets/native-env-linux-la64.wa.s

loong64:
	gcc $(ENC_S) golden-loong64.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=loong64 -target=linux -o memory-01.wa.s memory-01.wat
	gcc env.c memory-01.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
