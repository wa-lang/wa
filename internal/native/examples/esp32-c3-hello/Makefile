# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

run:
	$(WA) rv2elf -arch=riscv32 hello_riscv32_zh.ws
	$(WA) wemu -uart=esp32c3 hello_riscv32_zh.ws.elf.exe

fmt:
	$(WA) fmt -riscv hello_riscv32_zh.ws > a.out.ws

debug:
	$(WA) -debug rv2elf hello_riscv32_zh.ws

clean:
	-rm hello_riscv32_zh.ws.elf.exe
