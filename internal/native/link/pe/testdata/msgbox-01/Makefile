# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

FLAGS:=-nostdlib
FLAGS+=-Wl,-e,_start
FLAGS+=-Wl,--subsystem,console
FLAGS+=-Wl,--file-alignment,16
FLAGS+=-Wl,--section-alignment,16
FLAGS+=-Wl,--build-id=none
FLAGS+=-Wl,--disable-dynamicbase
FLAGS+=-s
FLAGS+=-lkernel32
FLAGS+=-luser32

build:
	gcc msgbox.s $(FLAGS) -o msgbox.exe
	msgbox.exe

gen:
	go run gen-golden-exe.go -o a.out.exe
	a.out.exe

peinfo:
	go run peinfo.go msgbox.exe

objdump:
	objdump -p msgbox.exe

clean:
