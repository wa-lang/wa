build:
	riscv64-unknown-elf-as hello.S -o hello.o
	riscv64-unknown-elf-ld -T link.ld hello.o -o hello.elf.exe

run:
	qemu-system-riscv64 -machine virt -nographic -bios none -kernel hello.elf.exe

linux-build:
	riscv64-unknown-elf-as hello.S -o hello.o
	riscv64-unknown-elf-ld -T link.ld hello.o -o hello.elf.exe

linux-run:
	qemu-riscv64 hello.elf

clean:
