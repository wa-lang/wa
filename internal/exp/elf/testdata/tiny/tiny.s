// $ clang -nostdlib -static tiny.s -o tiny

.globl _start
.text

_start:
	// push + pop is 3 bytes rather than 4 for the usual mov $60, %ax
	pushq $60
	popq  %rax

	// mov $60, %ax

	xorl %edi, %edi
	syscall
