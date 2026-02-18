.extern .Wa.Memory.addr

# kernel32.dll

.extern ExitProcess
.extern GetStdHandle
.extern VirtualAlloc
.extern WriteFile

# int _Wa_Import_syscall_write(int fd, int ptr, int size) {
.section .text
.globl .Wa.Import.syscall.write
.Wa.Import.syscall.write:
    push rbp
    mov  r10, [rip+.Wa.Memory.addr]
    add  rdx, r10 # buf
    call .Wa.Runtime.write
    pop  rbp
    ret


# int _Wa_Runtime_write(int fd, void *buf, int count)
.section .text
.globl .Wa.Runtime.write
.Wa.Runtime.write:
    push rbp
    mov  rbp, rsp
    sub  rsp, 32

    // 保存寄存器
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
    
    // 恢复寄存器
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
    sub rsp, 40
    call ExitProcess


# void* _Wa_Runtime_malloc(int size)
.section .text
.globl .Wa.Runtime.malloc
.Wa.Runtime.malloc:
    sub rsp, 40

    mov rdx, rcx      # dwSize
    xor rcx, rcx      # lpAddress = NULL
    mov r8,  0x3000   # MEM_COMMIT | MEM_RESERVE
    mov r9,  0x04     # PAGE_READWRITE
    call VirtualAlloc # rax = allocated memory

    add rsp, 40
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

