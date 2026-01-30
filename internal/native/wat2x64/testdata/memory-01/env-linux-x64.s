# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.extern .Wa.Memory.addr

.section .text
.global _Wa_Import_syscall_write
_Wa_Import_syscall_write:
    mov  rdi, rdi # arg.0: fd
    mov  rax, [rip + .Wa.Memory.addr]
    add  rsi, rax # arg.1: buffer
    mov  eax, 1 # arg.2: len
    syscall
    ret

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.global _Wa_Runtime_write
_Wa_Runtime_write:
    mov eax, 1
    syscall
    ret

# void _Wa_Runtime_exit(int status)
.section .text
.global _Wa_Runtime_exit
_Wa_Runtime_exit:
    mov eax, 60
    syscall

# void* _Wa_Runtime_malloc(int size)
.section .text
.global _Wa_Runtime_malloc
_Wa_Runtime_malloc:
    mov r8, -1   # fd = -1
    mov r9, 0    # offset = 0
    mov r10, 34  # flags = MAP_PRIVATE | MAP_ANONYMOUS (0x02 | 0x20)
    mov rdx, 3   # prot = PROT_READ | PROT_WRITE (0x01 | 0x02)
    mov rsi, rdi # length = size
    xor rdi, rdi # addr = NULL
    mov eax, 9   # sys_mmap
    syscall
    ret

# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.global _Wa_Runtime_memcpy
_Wa_Runtime_memcpy:
    mov  rax, rdi
    test rdx, rdx
    jz   .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    mov r8b, [rsi]
    mov [rdi], r8b
    inc rdi
    inc rsi
    dec rdx
    jnz .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    ret

# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.global _Wa_Runtime_memmove
_Wa_Runtime_memmove:
    mov rax, rdi              # 备份 dst 用于返回
    cmp rdi, rsi              # 比较 dst 和 src
    je  .Wa.L.memmove.done    # 如果相等, 直接结束
    jb  .Wa.L.memmove.forward # 如果 dst < src，进行前向拷贝 (同 memcpy)

    # --- 后向拷贝逻辑 (dst > src) ---
    mov rcx, rdx              # 将 n 放入计数器
    add rdi, rdx              # 将 dst 指针移到末尾 (dst + n)
    add rsi, rdx              # 将 src 指针移到末尾 (src + n)
    dec rdi                   # 指向最后一个字节
    dec rsi
    std                       # 设置方向标志位 (Direction Flag = 1)
                              # 这会让接下来的 movsb 指令每拷贝一个字节后, 指针自动递减
    rep movsb                 # 硬件加速后向拷贝
    cld                       # 清除方向标志位, 恢复为默认的自动递增模式
    jmp .Wa.L.memmove.done

.Wa.L.memmove.forward:
    # --- 前向拷贝逻辑 (dst < src) ---
    mov rcx, rdx
    rep movsb           # cld 模式下 (默认), 指针自动递增

.Wa.L.memmove.done:
    ret

# void* _Wa_Runtime_memset(void* s, int c, int n)
.global _Wa_Runtime_memset
_Wa_Runtime_memset:
    mov rax, rdi        # 返回 s
    test rdx, rdx
    jz .Wa.L.memset.done
.Wa.L.memset.loop:
    mov [rdi], sil      # sil 是 rsi 的低 8 位 (字符 c)
    inc rdi
    dec rdx
    jnz .Wa.L.memset.loop
.Wa.L.memset.done:
    ret

# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
.global _Wa_Import_syscall_linux_print_str
_Wa_Import_syscall_linux_print_str:
    # rax = base + ptr
    mov rax, [rip + .Wa.Memory.addr]
    add rax, rdi

    mov rdx, rsi # arg.2: len
    mov rsi, rax # arg.1: base + ptr
    mov rdi, 1   # arg.0: stdout

    mov eax, 1   # sys_write
    syscall
    ret

# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
.global _Wa_Import_syscall_linux_proc_exit
_Wa_Import_syscall_linux_proc_exit:
    jmp _Wa_Runtime_exit

# void _Wa_Import_syscall_linux_print_rune(int32_t c)
.section .text
.global _Wa_Import_syscall_linux_print_rune
_Wa_Import_syscall_linux_print_rune:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16
    
    mov  [rsp], dil
    
    mov  rax, 1   # sys_write
    mov  rdi, 1   # stdout
    mov  rsi, rsp # buf
    mov  rdx, 1   # count
    syscall
    
    add  rsp, 16
    pop  rbp
    ret

# void _Wa_Import_syscall_linux_print_i64(int64_t val)
.section .text
.global _Wa_Import_syscall_linux_print_i64
_Wa_Import_syscall_linux_print_i64:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32            # 分配缓冲区

    mov  rax, rdi           # rax = val
    lea  rcx, [rbp - 1]     # rcx 从缓冲区末尾向前移动
    mov  r8, 10             # 除数

    # 1. 处理负数特殊情况
    test rax, rax
    jns  .Wa.L.syscall.linux.print_i64.convert
    neg  rax                # 如果是负数, 取反

.Wa.L.syscall.linux.print_i64.convert:
    xor  rdx, rdx           # 每次除法前必须清空 rdx
    div  r8                 # rdx:rax / 10 -> rax=商, rdx=余数
    add  dl, '0'            # 余数转 ASCII
    mov  [rcx], dl
    dec  rcx
    test rax, rax           # 商是否为 0
    jnz  .Wa.L.syscall.linux.print_i64.convert

    # 2. 补负号
    cmp  rdi, 0
    jge  .Wa.L.syscall.linux.print_i64.setup_print
    mov  byte ptr [rcx], '-'
    dec  rcx

.Wa.L.syscall.linux.print_i64.setup_print:
    # 3. 计算长度并打印
    # rcx 现在指向字符串第一个字符的前一个位置
    inc  rcx                # 指向第一个字符
    mov  rdx, rbp
    sub  rdx, rcx           # rdx = rbp - rcx (当前字符串长度)

    mov  rax, 1             # sys_write
    mov  rsi, rcx           # buffer
    mov  rdi, 1             # fd = stdout
    syscall

    add  rsp, 32
    pop  rbp
    ret

