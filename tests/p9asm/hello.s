#include "textflag.h"

TEXT main·main(SB), NOSPLIT, $0
    MOVQ $60, AX     // syscall number for exit
    MOVQ $42, DI     // exit code
    SYSCALL          // invoke syscall
    RET              // 理论上不会执行
