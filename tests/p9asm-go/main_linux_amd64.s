// +build ignore

// 目标：调用 Linux 系统调用 exit(42)
// syscall: rax=60 (exit), rdi=返回码

// TEXT ·main(SB), NOSPLIT, $0

// func _start()
TEXT main·_start(SB), 4, $0
    MOVQ $60, AX     // syscall number for exit
    MOVQ $42, DI     // exit code
    SYSCALL          // invoke syscall
    RET              // （理论上不会执行）

// 添加这个符号！用于告诉 Wa linker：这是 package main
// RODATA=8
GLOBL  type··importpath·main(SB), 8, $0
