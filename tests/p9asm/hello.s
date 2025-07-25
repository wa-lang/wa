// 函数不再包含flags字段, 默认都是固定栈

TEXT main·main(SB), $0
    MOVQ $60, AX     // syscall number for exit
    MOVQ $42, DI     // exit code
    SYSCALL          // invoke syscall
    RET              // 理论上不会执行
