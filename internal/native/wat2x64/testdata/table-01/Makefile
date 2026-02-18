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

golden:
	gcc env-win64.s golden-win64.s $(FLAGS) -o a.out0.exe
	a.out0.exe

linux:
	gcc -nostdlib -static env-linux-x64.s golden-linux.s -o a.out0.exe -z noexecstack
	./a.out0.exe

dev:
	$(WA) build -arch=x64 -target=windows table-01.wat
	gcc table-01.wa.s $(FLAGS) -o a.out1.exe
	a.out1.exe

dev-linux:
	$(WA) build -arch=x64 -target=linux -o table-01.linux.wa.s table-01.wat
	gcc -nostdlib -static env-linux-x64.s table-01.linux.wa.s -z noexecstack -o a.out1.exe
	./a.out1.exe

clean:
