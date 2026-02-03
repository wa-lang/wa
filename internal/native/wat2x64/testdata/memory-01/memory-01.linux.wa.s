# 源文件: memory-01.wat, ABI: X64-Unix
# 自动生成的代码, 不要手动修改!!!

.intel_syntax noprefix

# 运行时函数
# .extern .Wa.Runtime.write
# .extern .Wa.Runtime.exit
# .extern .Wa.Runtime.malloc
# .extern .Wa.Runtime.memcpy
# .extern .Wa.Runtime.memset
# .extern .Wa.Runtime.memmove

# 导入函数(外部库定义)
# .extern .Wa.Import.syscall.write

# 定义内存
.section .data
.align 8
.globl .Wa.Memory.addr
.globl .Wa.Memory.pages
.globl .Wa.Memory.maxPages
.Wa.Memory.addr: .quad 0
.Wa.Memory.pages: .quad 1
.Wa.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 8
# memcpy(&Memory[8], data[0], size)
.Wa.Memory.dataOffset.0: .quad 8
.Wa.Memory.dataSize.0: .quad 12
.Wa.Memory.dataPtr.0: .asciz "hello world\012"

# 内存初始化函数
.section .text
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rdi, [rip + .Wa.Memory.maxPages]
    shl  rdi, 16
    call .Wa.Runtime.malloc
    mov  [rip + .Wa.Memory.addr], rax

    # 内存清零
    mov  rdi, [rip + .Wa.Memory.addr]
    mov  rsi, 0
    mov  rdx, [rip + .Wa.Memory.maxPages]
    shl  rdx, 16
    call .Wa.Runtime.memset

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    mov  rax, [rip + .Wa.Memory.addr]
    mov  rdi, [rip + .Wa.Memory.dataOffset.0]
    add  rdi, rax
    lea  rsi, [rip + .Wa.Memory.dataPtr.0]
    mov  rdx, [rip + .Wa.Memory.dataSize.0]
    call .Wa.Runtime.memcpy

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义表格
.section .data
.align 8
.globl .Wa.Table.addr
.globl .Wa.Table.size
.globl .Wa.Table.maxSize
.Wa.Table.addr: .quad 0
.Wa.Table.size: .quad 1
.Wa.Table.maxSize: .quad 1

# 汇编程序入口函数
.section .text
.globl _start
_start:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    call .Wa.Memory.initFunc
    call .Wa.F.main

    # runtime.exit(0)
    mov  rdi, 0
    call .Wa.Runtime.exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.global .Wa.Runtime.write
.Wa.Runtime.write:
    mov eax, 1
    syscall
    ret

# void _Wa_Runtime_exit(int status)
.section .text
.global .Wa.Runtime.exit
.Wa.Runtime.exit:
    mov eax, 60
    syscall

# void* _Wa_Runtime_malloc(int size)
.section .text
.global .Wa.Runtime.malloc
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
.global .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
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
.global .Wa.Runtime.memmove
.Wa.Runtime.memmove:
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
.global .Wa.Runtime.memset
.Wa.Runtime.memset:
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
.global .Wa.Import.syscall_linux.print_str
.Wa.Import.syscall_linux.print_str:
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
.global .Wa.Import.syscall_linux.proc_exit
.Wa.Import.syscall_linux.proc_exit:
    jmp .Wa.Runtime.exit

# void _Wa_Import_syscall_linux_print_rune(int32_t c)
.section .text
.global .Wa.Import.syscall_linux.print_rune
.Wa.Import.syscall_linux.print_rune:
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
.global .Wa.Import.syscall_linux.print_i64
.Wa.Import.syscall_linux.print_i64:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32            # 分配缓冲区

    mov  rax, rdi           # rax = val
    lea  rcx, [rbp - 1]     # rcx 从缓冲区末尾向前移动
    mov  r8, 10             # 除数

    # 1. 处理负数特殊情况
    test rax, rax
    jns  .Wa.L.syscall_linux.print_i64.convert
    neg  rax                # 如果是负数, 取反

.Wa.L.syscall_linux.print_i64.convert:
    xor  rdx, rdx           # 每次除法前必须清空 rdx
    div  r8                 # rdx:rax / 10 -> rax=商, rdx=余数
    add  dl, '0'            # 余数转 ASCII
    mov  [rcx], dl
    dec  rcx
    test rax, rax           # 商是否为 0
    jnz  .Wa.L.syscall_linux.print_i64.convert

    # 2. 补负号
    cmp  rdi, 0
    jge  .Wa.L.syscall_linux.print_i64.setup_print
    mov  byte ptr [rcx], '-'
    dec  rcx

.Wa.L.syscall_linux.print_i64.setup_print:
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

.section .text
.global .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    mov  rdi, rdi # arg.0: fd
    mov  rax, [rip + .Wa.Memory.addr]
    add  rsi, rax # arg.1: buffer
    mov  eax, 1 # arg.2: len
    syscall
    ret

.section .data
.align 8
.Wa.Runtime.panic.message: .asciz "panic"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stderr, panicMessage, size)
    mov  rdi, 2 # stderr
    lea  rsi, [rip + .Wa.Runtime.panic.message]
    mov  rdx, [rip + .Wa.Runtime.panic.messageLen] # size
    call .Wa.Runtime.write

    # 退出程序
    mov  rdi, 1 # 退出码
    call .Wa.Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.Wa.F.main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 1
    movabs rax, 1
    mov    [rbp-8], rax

    # i64.const 8
    movabs rax, 8
    mov    [rbp-16], rax

    # i64.const 12
    movabs rax, 12
    mov    [rbp-24], rax

    # call syscall.write(...)
    mov rdi, qword ptr [rbp-8] # arg 0
    mov rsi, qword ptr [rbp-16] # arg 1
    mov rdx, qword ptr [rbp-24] # arg 2
    call .Wa.Import.syscall.write
    mov qword ptr [rbp-8], rax

    nop # drop [rbp-8]

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

