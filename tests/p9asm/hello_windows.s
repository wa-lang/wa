// msg 长度
#define msg_len 13

// 数据段
DATA msg+0(SB)/13, $"Hello World!\n"
GLOBL msg(SB), (RODATA), $13

// 用于接收写入长度
DATA nw+0(SB)/8, $0
GLOBL nw(SB), (NOPTR), $8

// 声明需要用到的外部符号（Windows API）
GLOBL GetStdHandle(SB),DUPOK, $0
GLOBL WriteConsoleA(SB),DUPOK, $0

// main.main
TEXT main·main(SB),$0-0
    // 获取标准输出句柄（这里用 -11，即 STD_OUTPUT_HANDLE）
    MOVQ $-11, CX
    CALL GetStdHandle(SB)

    // 返回值在 AX 中（句柄）

    // 准备参数：句柄、字符串地址、长度、写入长度地址、保留
    MOVQ AX, CX                    // 第一个参数：句柄

    LEAQ msg(SB), DX               // 第二个参数：字符串地址
    MOVQ $msg_len, R8              // 第三个参数：长度
    LEAQ nw(SB), R9                // 第四个参数：写入长度地址

    MOVQ $0, R10                    // 第五个参数：保留

    CALL WriteConsoleA(SB)

    RET
