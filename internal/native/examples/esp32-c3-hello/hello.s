// 需要针对 ESP32-C3 硬件微调。

.section .text
.globl _start

_start:
    // 1. 设置堆栈指针 (非常重要!)
    // 假设 RAM 从 0x3FC80000 开始，使用末尾作为堆栈
    li sp, 0x3FC8FFFF 
    
    // 2. 将 UART0 基地址加载到 t0 (UART_BASE)
    // 必须查阅手册确认!!! 假设为 0x60000000
    li t0, 0x60000000 
    
    // 3. 调用主程序逻辑
    call main

// 循环发送字符串
main:
    li t1, msg        // t1 = 字符串地址
    
loop:
    lb t2, 0(t1)      // t2 = *t1 (加载一个字节)
    beqz t2, done     // 如果是零字节 (字符串结束), 跳转到 done
    
    // 假设 UART0 的发送数据寄存器在基地址 0 偏移处
    // 实际需要检查 UART FIFO 状态，但最简示例暂时跳过
    sb t2, 0(t0)      // 将 t2 中的字符写入 UART0 (地址 t0 + 0)
    
    addi t1, t1, 1    // t1++ (指向下一个字符)
    j loop
    
done:
    // 程序结束，进入休眠/无限循环
    wfi               // Wait For Interrupt (低功耗模式)
    j done


.section .rodata
msg:
    .string "Hello from RISC-V Naked Assembly!\n"
