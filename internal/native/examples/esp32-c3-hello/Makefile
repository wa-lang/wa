# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

run:
	$(WA) rv2elf -arch=riscv32 hello_riscv32_zh.ws
	$(WA) wemu -uart=esp32c3 hello_riscv32_zh.ws.elf.exe

build:
	$(WA) esp32build -arch=riscv32 hello_riscv32_zh.ws

write-flash:
	esptool.py --chip esp32c3  --port COM3  --baud 460800 write_flash -z 0x0000 hello_riscv32_zh.ws.esp32.bin

fmt:
	$(WA) fmt -riscv hello_riscv32_zh.ws > a.out.ws

debug:
	$(WA) -debug rv2elf hello_riscv32_zh.ws

clean:
	-rm hello_riscv32_zh.ws.elf.exe
