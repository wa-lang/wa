# 源文件: fib-01.wat, ABI: x64-Windows
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
# .extern .Wa.Import.env.print_i64

# 定义内存
.section .data
.align 8
.globl .Wa.Memory.addr
.Wa.Memory.addr: .quad 0
.globl .Wa.Memory.pages
.Wa.Memory.pages: .quad 1
.globl .Wa.Memory.maxPages
.Wa.Memory.maxPages: .quad 1

# 内存初始化函数
.section .text
.globl .Wa.Memory.initFunc
.Wa.Memory.initFunc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 分配内存
    mov  rcx, qword ptr [rip+.Wa.Memory.maxPages]
    shl  rcx, 16
    call .Wa.Runtime.malloc
    mov  qword ptr [rip+.Wa.Memory.addr], rax

    # 内存清零
    mov  rcx, qword ptr [rip+.Wa.Memory.addr]
    mov  rdx, 0
    mov  r8, qword ptr [rip+.Wa.Memory.maxPages]
    shl  r8, 16
    call .Wa.Runtime.memset

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# 定义表格
.section .data
.align 8
.globl .Wa.Table.addr
.Wa.Table.addr: .quad 0
.globl .Wa.Table.size
.Wa.Table.size: .quad 1
.globl .Wa.Table.maxSize
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
    mov  rcx, 0
    call .Wa.Runtime.exit

    # exit 后这里不会被执行, 但是依然保留
    mov rsp, rbp
    pop rbp
    ret

.extern .Wa.Memory.addr

# kernel32.dll

.extern ExitProcess
.extern GetStdHandle
.extern VirtualAlloc
.extern WriteFile


# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # 保存寄存器
    mov [rbp-8], r12
    mov [rbp-16], r13
    mov [rbp-24], r14

    # 保存参数
    mov r12, rcx # fd
    mov r13, rdx # buf
    mov r14, r8  # count

    # 获取标准输出句柄
    cmp rcx, 2
    je  .Wa.L.Runtime.write.stderr
    jmp .Wa.L.Runtime.write.stdout

.Wa.L.Runtime.write.stdout:
    mov ecx, -11 # STD_OUTPUT_HANDLE
    jmp .Wa.L.Runtime.write.gethandle

.Wa.L.Runtime.write.stderr:
    mov ecx, -12 # STD_ERROR_HANDLE
    jmp .Wa.L.Runtime.write.gethandle

.Wa.L.Runtime.write.gethandle:
    # rax = GetStdHandle(nStdHandle)
    sub  rsp, 32
    call GetStdHandle
    add  rsp, 32

    # BOOL WriteFile(
    #   [in]                HANDLE       hFile,
    #   [in]                LPCVOID      lpBuffer,
    #   [in]                DWORD        nNumberOfBytesToWrite,
    #   [out, optional]     LPDWORD      lpNumberOfBytesWritten,
    #   [in, out, optional] LPOVERLAPPED lpOverlapped
    # );

    sub  rsp, 48
    mov  rcx, rax               # arg.0: fd
    mov  rdx, r13               # arg.1: buf
    mov  r8, r14                # arg2: count
    mov  qword ptr [rsp+40], 0  # arg.3: (*lpNumberOfBytesWritten) = 0
    lea  r9, qword ptr [rsp+40] # arg.3: lpNumberOfBytesWritten
    mov  qword ptr [rsp+32], 0  # arg.4: lpOverlapped = NULL
    call WriteFile
    add  rsp, 48

    # return nWrite
    mov eax, dword ptr [rsp+40]
    
    # 恢复寄存器
    mov r12, [rbp-8]
    mov r13, [rbp-16]
    mov r14, [rbp-24]

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret


# void _Wa_Runtime_exit(int status)
.section .text
.globl .Wa.Runtime.exit
.Wa.Runtime.exit:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # void ExitProcess(
    #   [in] UINT uExitCode
    # );

    call ExitProcess

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret


# void* _Wa_Runtime_malloc(int size)
.section .text
.globl .Wa.Runtime.malloc
.Wa.Runtime.malloc:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # LPVOID VirtualAlloc(
    #   [in, optional] LPVOID lpAddress,
    #   [in]           SIZE_T dwSize,
    #   [in]           DWORD  flAllocationType,
    #   [in]           DWORD  flProtect
    # );

    mov  rdx, rcx      # dwSize
    xor  rcx, rcx      # lpAddress = NULL
    mov  r8,  0x3000   # MEM_COMMIT | MEM_RESERVE
    mov  r9,  0x04     # PAGE_READWRITE
    call VirtualAlloc # rax = allocated memory

    mov rsp, rbp
    pop rbp
    ret


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memcpy
.Wa.Runtime.memcpy:
    mov  rax, rcx
    test r8, r8
    jz   .Wa.L.memcpy.done

.Wa.L.memcpy.loop:
    mov r9b, byte ptr [rdx]
    mov byte ptr [rcx], r9b
    inc rcx
    inc rdx
    dec r8
    jnz .Wa.L.memcpy.loop

.Wa.L.memcpy.done:
    ret


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
.section .text
.globl .Wa.Runtime.memmove
.Wa.Runtime.memmove:
    mov  rax, rcx # 备份 dst 用于返回
    test r8, r8   # n == 0 ?
    jz   .Wa.L.memmove.done

    cmp rcx, rdx
    je  .Wa.L.memmove.done
    jb  .Wa.L.memmove.forward  # dst < src → 前向拷贝

    # =========================
    # 后向拷贝 (dst > src)
    # =========================

    push rdi
    push rsi

    mov rdi, rcx
    mov rsi, rdx

    add rdi, r8
    dec rdi
    add rsi, r8
    dec rsi

    mov rcx, r8 # 计数器

.Wa.L.memmove.back_loop:
    mov r9b, byte ptr [rsi]
    mov byte ptr [rdi], r9b
    dec rdi
    dec rsi
    dec rcx
    jnz .Wa.L.memmove.back_loop

    pop rsi
    pop rdi
    jmp .Wa.L.memmove.done

.Wa.L.memmove.forward:

    # =========================
    # 前向拷贝 (dst < src)
    # =========================

    push rdi
    push rsi

    mov rdi, rcx
    mov rsi, rdx
    mov rcx, r8

.Wa.L.memmove.fwd_loop:
    mov r9b, byte ptr [rsi]
    mov byte ptr [rdi], r9b
    inc rdi
    inc rsi
    dec rcx
    jnz .Wa.L.memmove.fwd_loop

    pop rsi
    pop rdi

.Wa.L.memmove.done:
    ret


# void* _Wa_Runtime_memset(void* s, int c, int n)
.section .text
.globl .Wa.Runtime.memset
.Wa.Runtime.memset:
    mov  rax, rcx # 返回 s
    test r8, r8
    jz   .Wa.L.memset.done

.Wa.L.memset.loop:
    mov byte ptr [rcx], dl # c 的低 8 位
    inc rcx
    dec r8
    jnz .Wa.L.memset.loop

.Wa.L.memset.done:
    ret


# void _Wa_Import_syscall_windows_print_str (uint32_t ptr, int32_t len)
.section .text
.globl .Wa.Import.syscall_windows.print_str
.Wa.Import.syscall_windows.print_str:
    mov r9, rdx # 保存 len

    # r8 = base + ptr
    mov r8, qword ptr [rip+.Wa.Memory.addr]
    add r8, rcx

    # 调用 write(fd=1, buf, len)
    mov  rcx, 1  # fd
    mov  rdx, r8 # buffer
    mov  r8,  r9 # len
    call .Wa.Runtime.write
    ret


# void _Wa_Import_syscall_windows_proc_exit(int32_t code)
.section .text
.globl .Wa.Import.syscall_windows.proc_exit
.Wa.Import.syscall_windows.proc_exit:
    jmp .Wa.Runtime.exit


# void _Wa_Import_syscall_windows_print_rune(int32_t c)
.section .text
.globl .Wa.Import.syscall_windows.print_rune
.Wa.Import.syscall_windows.print_rune:
    sub rsp, 40

    mov byte ptr [rsp+32], cl  # 存 1 字节到临时空间

    mov rcx, 1                 # fd = stdout
    lea rdx, byte ptr [rsp+32] # buf
    mov r8, 1                  # len = 1
    call .Wa.Runtime.write

    add rsp, 40
    ret


# void _Wa_Import_syscall_windows_print_i64(int64_t val)
.section .text
.globl .Wa.Import.syscall_windows.print_i64
.Wa.Import.syscall_windows.print_i64:
    sub rsp, 72 # 32 shadow + 32 buffer + 8 对齐

    mov rax, rcx              # rax = val
    lea r9, byte ptr [rsp+63] # r9 指向缓冲区末尾
    mov r8, 10                # 除数 = 10

    # ----------------------------
    # 1. 处理负数
    # ----------------------------
    test rax, rax
    jns  .Wa.L.syscall_windows.print_i64.convert
    neg  rax

.Wa.L.syscall_windows.print_i64.convert:
    xor  rdx, rdx
    div  r8 # rax=商 rdx=余数
    add  dl, '0'
    mov  byte ptr [r9], dl
    dec  r9
    test rax, rax
    jnz  .Wa.L.syscall_windows.print_i64.convert

    # ----------------------------
    # 2. 补负号
    # ----------------------------
    test rcx, rcx
    jge  .Wa.L.syscall_windows.print_i64.setup_print
    mov  byte ptr [r9], '-'
    dec  r9

.Wa.L.syscall_windows.print_i64.setup_print:
    inc r9 # r9 = 字符串起始地址

    lea rdx, byte ptr [rsp+64] # 计算长度
    sub rdx, r9                # rdx = len

    # ----------------------------
    # 3. 调用 Runtime.write
    # ----------------------------
    mov  rcx, 1   # fd = stdout
    mov  r8,  rdx # len
    mov  rdx, r9  # buffer
    call .Wa.Runtime.write

    add rsp, 72
    ret


.section .text
.globl .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    mov rax, qword ptr [rip+.Wa.Memory.addr]
    add rdx, rax # rdx = base + offset

    # 参数已经匹配:
    # rcx = fd
    # rdx = buffer
    # r8  = len

    # 直接尾调用
    jmp .Wa.Runtime.write


.section .data
.align 8
.Wa.Runtime.panic.message: .ascii "panic"
.Wa.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Wa.Runtime.panic
.Wa.Runtime.panic:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    # runtime.write(stderr, panicMessage, size)
    mov  rcx, 2 # stderr
    lea  rdx, qword ptr [rip+.Wa.Runtime.panic.message]
    mov  r8, qword ptr [rip+.Wa.Runtime.panic.messageLen] # size
    call .Wa.Runtime.write

    # 退出程序
    mov  rcx, 1 # 退出码
    call .Wa.Runtime.exit

    # return
    mov rsp, rbp
    pop rbp
    ret

# func main
.section .text
.Wa.F.main:
    # local N: i64
    # local i: i64

    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    mov dword ptr [rbp-8], 0 # local N = 0
    mov dword ptr [rbp-16], 0 # local i = 0

    # fn.body.begin(name=main, suffix=00000000)

    # i64.const 10
    movabs rax, 10
    mov    qword ptr [rbp-24], rax

    # local.set N i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-8], rax

    # i64.const 1
    movabs rax, 1
    mov    qword ptr [rbp-24], rax

    # local.set i i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # loop.begin(name=my_loop, suffix=00000001)
.Wa.L.brNext.my_loop.00000001:
    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.F.fib
    mov qword ptr [rbp-24], rax

    # call env.print_i64(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.Import.env.print_i64

    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # i64.const 1
    movabs rax, 1
    mov    qword ptr [rbp-32], rax

    # i64.add
    mov rax, qword ptr [rbp-24]
    add rax, qword ptr [rbp-32]
    mov qword ptr [rbp-24], rax

    # local.set i i64
    mov rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # local.get i i64
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-24], rax

    # local.get N i64
    mov rax, qword ptr [rbp-8]
    mov qword ptr [rbp-32], rax

    # i64.lt_u
    mov   r10, qword ptr [rbp-24]
    mov   r11, qword ptr [rbp-32]
    cmp   r10, r11
    setb  al
    movzx eax, al
    mov   dword ptr [rbp-24], eax

    # br_if my_loop.00000001
    mov eax, dword ptr [rbp-24]
    cmp eax, 0
    je  .Wa.L.brFallthrough.my_loop.00000001
    jmp .Wa.L.brNext.my_loop.00000001
.Wa.L.brFallthrough.my_loop.00000001:

    # loop.end(name=my_loop, suffix=00000001)

.Wa.L.brNext.main.00000000:
    # fn.body.end(name=main, suffix=00000000)

    # 根据ABI处理返回值

    # 将返回值变量复制到寄存器

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

# func fib(n:i64) => i64
.section .text
.globl .Wa.F.fib
.Wa.F.fib:
    push rbp
    mov  rbp, rsp
    sub  rsp, 64

    # 将寄存器参数备份到栈
    mov qword ptr [rbp+16], rcx # save arg.0

    # 将返回值变量初始化为0
    mov dword ptr [rbp-8], 0 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # fn.body.begin(name=fib, suffix=00000002)

    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-16], rax

    # i64.const 2
    movabs rax, 2
    mov    qword ptr [rbp-24], rax

    # i64.le_u
    mov   r10, qword ptr [rbp-16]
    mov   r11, qword ptr [rbp-24]
    cmp   r10, r11
    setbe al
    movzx eax, al
    mov   dword ptr [rbp-16], eax

    # if.begin(name=, suffix=00000003)
    mov eax, dword ptr [rbp-16]
    cmp eax, 0
    je  .Wa.L.else.00000003

    # if.body(name=, suffix=00000003)
    # i64.const 1
    movabs rax, 1
    mov    qword ptr [rbp-16], rax

    jmp .Wa.L.brNext.00000003

.Wa.L.else.00000003:
    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-16], rax

    # i64.const 1
    movabs rax, 1
    mov    qword ptr [rbp-24], rax

    # i64.sub
    mov rax, qword ptr [rbp-16]
    sub rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-16] # arg 0
    call .Wa.F.fib
    mov qword ptr [rbp-16], rax

    # local.get n i64
    mov rax, qword ptr [rbp+16]
    mov qword ptr [rbp-24], rax

    # i64.const 2
    movabs rax, 2
    mov    qword ptr [rbp-32], rax

    # i64.sub
    mov rax, qword ptr [rbp-24]
    sub rax, qword ptr [rbp-32]
    mov qword ptr [rbp-24], rax

    # call fib(...)
    mov rcx, qword ptr [rbp-24] # arg 0
    call .Wa.F.fib
    mov qword ptr [rbp-24], rax

    # i64.add
    mov rax, qword ptr [rbp-16]
    add rax, qword ptr [rbp-24]
    mov qword ptr [rbp-16], rax

.Wa.L.brNext.00000003:
    # if.end(name=, suffix=00000003)

.Wa.L.brNext.fib.00000002:
    # fn.body.end(name=fib, suffix=00000002)

    # 根据ABI处理返回值

    # 将栈上数据复制到返回值变量
    mov rax, qword ptr [rbp-16]
    mov qword ptr [rbp-8], rax # ret.0

    # 将返回值变量复制到寄存器
    mov rax, qword ptr [rbp-8] # ret.0

    # 函数返回
    mov rsp, rbp
    pop rbp
    ret

