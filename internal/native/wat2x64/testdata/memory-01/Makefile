# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

FLAGS:=-nostdlib
FLAGS+=-Wl,-e,_start
FLAGS+=-Wl,--subsystem,console
FLAGS+=-Wl,--file-alignment,16
FLAGS+=-Wl,--section-alignment,16
FLAGS+=-Wl,--build-id=none
FLAGS+=-Wl,--disable-dynamicbase
FLAGS+=-s
FLAGS+=-lkernel32

windows:
	gcc env-win64.s golden-win64.s $(FLAGS) -o a.out0.exe
	a.out0.exe

linux:
	gcc -nostdlib -static env-linux-x64.s golden-linux.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=x64 -target=windows memory-01.wat
	gcc memory-01.wa.s $(FLAGS) -o a.out1.exe
	a.out1.exe

dev-linux:
	$(WA) build -arch=x64 -target=linux -o memory-01.linux.wa.s memory-01.wat
	gcc -nostdlib -static memory-01.linux.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
