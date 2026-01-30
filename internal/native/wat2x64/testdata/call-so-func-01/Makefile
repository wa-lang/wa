# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

build:
	gcc -nostdlib -static -o a.out main.s

run:
	./a.out
	echo $$?

lib:
	gcc -shared -fPIC -o libhello.so libhello.c

clean:
