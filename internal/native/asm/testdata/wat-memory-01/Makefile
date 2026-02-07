# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

ENC_S:=../../assets/native-env-linux-la64.wa.s

dev:
	$(WA) build -arch=loong64 -target=linux -o a.wa.s memory-01.wat
	$(WA) ast -arch=loong64 a.wa.s > a.out2.s

build:
	$(WA) build -arch=loong64 -target=linux -o a.wa.s memory-01.wat
	$(WA) asm2elf -arch=loong64 a.wa.s

gcc:
	$(WA) build -arch=loong64 -target=linux -o a.wa.s memory-01.wat
	gcc -nostdlib -static a.wa.s
	./a.out

clean:
