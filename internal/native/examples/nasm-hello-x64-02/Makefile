# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

build:
	nasm -f elf64 -g -F dwarf hello.wa.s -l hello.lst
	gcc -o a.out hello.wa.o -no-pie -Wl,-z,noexecstack
	./a.out

clean:
	-rm a.out *.o
