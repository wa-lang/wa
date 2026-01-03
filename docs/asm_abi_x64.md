# 凹汇编语言函数调用约定 - X64 平台

## Microsoft X64 ABI

- 栈指针 rsp 必须保持 16 字节对齐
- 整数参数: rcx, rdx, r8, r9, stack(右到左入栈)
- 浮点数参数: xmm0-xmm3, stack(右到左入栈)
- 返回值: rax(整数/地址), xmm0(浮点数)
- 调用者保存: rax, rcx, rdx, r8-r11, xmm0-xmm5
- 被调用者保存: rbx, rbp, rdi, rsi, rsp, r12-r15
- 调用者需要在栈上准备 32 字节的影子空间, 被调用函数可以使用
- 系统调用编号不稳定

## System V AMD64 ABI

- 栈指针 rsp 必须是对齐到 16 字节
- 整数参数: rdi, rsi, rdx, rcx, r8, r9, stack(右到左入栈)
- 浮点数参数: xmm0-xmm7, stack(右到左入栈)
- 返回值: rax(整数/地址), xmm0(浮点数)
- 调用者保存: rax, rdi, rsi, rdx, rcx, r8-r11
- 被调用者保存: RBX, RBP, R12-R15
- rsp 之下有 128 字节的区域被称为 RedZone, 可临时使用

