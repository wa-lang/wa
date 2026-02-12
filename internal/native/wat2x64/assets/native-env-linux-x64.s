# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    mov eax, 1
    syscall
    ret


# void _Wa_Runtime_exit(int status)
.section .text
.globl .Wa.Runtime.exit
.Wa.Runtime.exit:
    mov eax, 60
    syscall


# void* _Wa_Runtime_malloc(int size)
.section .text
.globl .Wa.Runtime.malloc
.Wa.Runtime.malloc:
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
.globl .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
    mov  rax, rdi
    test rdx, rdx
    jz   .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    mov r8b, byte ptr [rsi]
    mov byte ptr [rdi], r8b
    inc rdi
    inc rsi
    dec rdx
    jnz .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    ret


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memmove
.Wa.Runtime.memmove:
    mov  rax, rdi           # 备份 dst 用于返回 (Reg2Reg)
    test rdx, rdx           # 检查 n 是否为 0
    jz   .Wa.L.memmove.done # 如果 n=0, 直接结束

    cmp rdi, rsi              # 比较 dst 和 src (Reg2Reg)
    je  .Wa.L.memmove.done    # 如果相等, 直接结束
    jb  .Wa.L.memmove.forward # 如果 dst < src, 前向拷贝

    # --- 后向拷贝逻辑 (dst > src) ---

    # 从末尾开始: 指针移动到 (ptr + n - 1)

    add rdi, rdx # (Reg2Reg) rdi = rdi + rdx
    dec rdi      # (Reg)     rdi = rdi - 1
    add rsi, rdx # (Reg2Reg) rsi = rsi + rdx
    dec rsi      # (Reg)      rsi = rsi - 1
    mov rcx, rdx # 计数器 = n (Reg2Reg)

.Wa.L.memmove.back_loop:
    mov r8b, byte ptr [rsi]     # (Mem2Reg) 取 1 字节
    mov byte ptr [rdi], r8b     # (Reg2Mem) 存 1 字节
    dec rdi                     # (Reg) 指针前移
    dec rsi                     # (Reg) 指针前移
    dec rcx                     # (Reg) 计数减 1
    jnz .Wa.L.memmove.back_loop # (Imm) 如果 rcx != 0 继续
    jmp .Wa.L.memmove.done

.Wa.L.memmove.forward:
    # --- 前向拷贝逻辑 (dst < src) ---
    mov rcx, rdx # 计数器 = n

.Wa.L.memmove.fwd_loop:
    mov r8b, byte ptr [rsi]    # (Mem2Reg) 取 1 字节
    mov byte ptr [rdi], r8b    # (Reg2Mem) 存 1 字节
    inc rdi                    # (Reg) 指针后移
    inc rsi                    # (Reg) 指针后移
    dec rcx                    # (Reg) 计数减 1
    jnz .Wa.L.memmove.fwd_loop # (Imm)
    
.Wa.L.memmove.done:
    ret


# void* _Wa_Runtime_memset(void* s, int c, int n)
.section .text
.globl .Wa.Runtime.memset
.Wa.Runtime.memset:
    mov  rax, rdi # 返回 s
    test rdx, rdx
    jz   .Wa.L.memset.done
.Wa.L.memset.loop:
    mov byte ptr [rdi], sil # sil 是 rsi 的低 8 位 (字符 c)
    inc rdi
    dec rdx
    jnz .Wa.L.memset.loop
.Wa.L.memset.done:
    ret


# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
.section .text
.globl .Wa.Import.syscall_linux.print_str
.Wa.Import.syscall_linux.print_str:
    # rax = base + ptr
    mov rax, qword ptr [rip+.Wa.Memory.addr]
    add rax, rdi

    mov rdx, rsi # arg.2: len
    mov rsi, rax # arg.1: base + ptr
    mov rdi, 1   # arg.0: stdout

    mov eax, 1 # sys_write
    syscall
    ret


# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
.section .text
.globl .Wa.Import.syscall_linux.proc_exit
.Wa.Import.syscall_linux.proc_exit:
    jmp .Wa.Runtime.exit


# void _Wa_Import_syscall_linux_print_rune(int32_t c)
.section .text
.globl .Wa.Import.syscall_linux.print_rune
.Wa.Import.syscall_linux.print_rune:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16

    mov byte ptr [rsp], dil

    mov rax, 1   # sys_write
    mov rdi, 1   # stdout
    mov rsi, rsp # buf
    mov rdx, 1   # count
    syscall

    add rsp, 16
    pop rbp
    ret


# void _Wa_Import_syscall_linux_print_i64(int64_t val)
.section .text
.globl .Wa.Import.syscall_linux.print_i64
.Wa.Import.syscall_linux.print_i64:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32 # 分配缓冲区

    mov rax, rdi               # rax = val
    lea rcx, qword ptr [rbp-1] # rcx 从缓冲区末尾向前移动
    mov r8, 10                 # 除数

    # 1. 处理负数特殊情况
    test rax, rax
    jns  .Wa.L.syscall_linux.print_i64.convert
    neg  rax # 如果是负数, 取反

.Wa.L.syscall_linux.print_i64.convert:
    xor  rdx, rdx # 每次除法前必须清空 rdx
    div  r8       # rdx:rax / 10 -> rax=商, rdx=余数
    add  dl, '0'  # 余数转 ASCII
    mov  byte ptr [rcx], dl
    dec  rcx
    test rax, rax # 商是否为 0
    jnz  .Wa.L.syscall_linux.print_i64.convert

    # 2. 补负号
    cmp rdi, 0
    jge .Wa.L.syscall_linux.print_i64.setup_print
    mov byte ptr [rcx], '-'
    dec rcx

.Wa.L.syscall_linux.print_i64.setup_print:
    # 3. 计算长度并打印
    # rcx 现在指向字符串第一个字符的前一个位置
    inc rcx # 指向第一个字符
    mov rdx, rbp
    sub rdx, rcx # rdx = rbp - rcx (当前字符串长度)

    mov rax, 1   # sys_write
    mov rsi, rcx # buffer
    mov rdi, 1   # fd = stdout
    syscall

    add rsp, 32
    pop rbp
    ret


.section .text
.globl .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    mov rdi, rdi # arg.0: fd
    mov rax, qword ptr [rip+.Wa.Memory.addr]
    add rsi, rax # arg.1: buffer
    mov eax, 1   # arg.2: len
    syscall
    ret

