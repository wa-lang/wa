# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../main.go

run:
	$(WA) rv2elf hello_riscv64.ws

	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello_riscv64.ws.elf.exe

run-zh:
	$(WA) rv2elf hello_riscv64_zh.ws

	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello_riscv64_zh.ws.elf.exe

fmt:
	$(WA) fmt -riscv hello_riscv64_zh.ws > a.out.ws

debug:
	$(WA) -debug rv2elf hello_riscv64.ws

clean:
	-rm hello.elf.exe
