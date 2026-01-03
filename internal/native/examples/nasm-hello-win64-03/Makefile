# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

build:
	nasm -f win64 -g -F cv8 hello.wa.s -l hello.lst
	gcc -o a.out.exe hello.wa.obj
	./a.out.exe

clean:
	-rm a.out *.o
