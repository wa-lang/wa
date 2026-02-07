# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

ENC_S:=../../assets/native-env-linux-la64.wa.s

build:
	$(WA) asm2elf -arch=loong64 hello-02.wa.s
	$(WA) objdump -prog hello-02.wa.elf.exe > 00.out.txt
	./hello-02.wa.elf.exe

ast:
	$(WA) ast -arch=loong64 hello-02.wa.s > hello-02.out2.s

gcc:
	gcc -nostdlib -static hello-02.wa.s
	./a.out

clean:
