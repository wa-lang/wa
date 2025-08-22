# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

riscv64-run:
	go run main.go
	qemu-system-riscv64 -machine virt -nographic -bios none -kernel a.out.exe

clean:
	-rm hello.o a.outexe
