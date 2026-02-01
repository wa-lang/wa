# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

.intel_syntax noprefix

.section .data
.balign 16
msg: .ascii "hello windows kernel.dll\n"
.set msg_len, . - msg
s_GetStdHandle: .asciz "GetStdHandle"
s_WriteFile:    .asciz "WriteFile"
s_ExitProcess:  .asciz "ExitProcess"

.section .text
.globl _start

_start:
    push rbp
    mov rbp, rsp
    and rsp, -16
    sub rsp, 64

    # 1. 获取 kernel32.dll 基地址
    mov rax, gs:[0x60]            
    mov rax, [rax + 0x18]         
    mov rax, [rax + 0x10]         # InLoadOrderModuleList
    mov rax, [rax]                # exe
    mov rax, [rax]                # ntdll
    mov rbx, [rax + 0x30]         # RBX = kernel32 基地址

    # 2. 查找 API
    lea rdx, [rip + s_GetStdHandle]
    call find_api
    mov r12, rax         

    lea rdx, [rip + s_WriteFile]
    call find_api
    mov r13, rax         

    lea rdx, [rip + s_ExitProcess]
    call find_api
    mov r14, rax         

    # 验证是否查找到 (避免 call 0)
    test r12, r12
    jz _crash
    test r13, r13
    jz _crash

    # 3. 业务逻辑
    mov rcx, -11
    call r12
    
    mov rcx, rax                  
    lea rdx, [rip + msg]          
    mov r8, msg_len               
    lea r9, [rsp + 48]            
    mov qword ptr [rsp + 32], 0   
    call r13

    xor rcx, rcx
    call r14

_crash:
    # 简单的断点，方便调试
    .byte 0xCC 
    ret

# ------------------------------------------------------------
# 修正后的 find_api (PE32+)
# ------------------------------------------------------------
find_api:
    push rbx
    push rsi
    push rdi
    push r15
    push r12

    # R15 = PE Header
    mov r15d, [rbx + 0x3c]
    add r15, rbx

    # 关键修正：64位导出表偏移是 0x88
    mov r12d, [r15 + 0x88]  
    test r12, r12
    jz .fail
    add r12, rbx            # R12 = 导出目录地址 (取代之前的 rax)

    mov ecx, [r12 + 0x18]   # Number of Names
    mov r10d, [r12 + 0x20]  # AddressOfNames RVA
    add r10, rbx

.loop:
    jecxz .fail
    dec ecx

    mov esi, [r10 + rcx*4]
    add rsi, rbx
    mov rdi, rdx            # 目标字符串

    # 寄存器保护：比较字符串时不破坏 r12(导出目录) 和 rbx(基址)
.cmp_str:
    mov al, [rsi]
    mov r11b, [rdi]         # 使用 r11b 避免破坏 rbx
    cmp al, r11b
    jne .loop
    test al, al
    jz .found
    inc rsi
    inc rdi
    jmp .cmp_str

.found:
    mov r11d, [r12 + 0x24]  # Ordinals RVA
    add r11, rbx
    movzx ecx, word ptr [r11 + rcx*2]

    mov r11d, [r12 + 0x1c]  # Functions RVA
    add r11, rbx
    mov eax, [r11 + rcx*4]
    add rax, rbx
    jmp .done

.fail:
    xor rax, rax

.done:
    pop r12
    pop r15
    pop rdi
    pop rsi
    pop rbx
    ret

