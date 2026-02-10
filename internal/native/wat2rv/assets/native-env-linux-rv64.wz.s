# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-或-later

声明 .Wa.Memory.addr

# void _Wa_Import_syscall_write (int fd, uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall.write:
    计加常 暂甲格, %相对高位(.Wa.Memory.addr)
    读长整 暂甲格, %相对低位(.Wa.Import.syscall.write)(暂甲格)

    # 参甲格 = fd
    某加某 参乙格, 参乙格, 暂甲格 # 参乙格 = base + ptr
    # 参丙格 = len
    某加常 参辛格, 零格, 64 # sys_write
    调用某

    转某址 零格, 0(回格)
完毕


# int _Wa_Runtime_write(int fd, void *buf, int count)
函数 .Wa.Runtime.write:
    某加常 参辛格, 零格, 64 # sys_write
    调用某
    转某址 零格, 0(回格)
完毕


# void _Wa_Runtime_exit(int status)
函数 .Wa.Runtime.exit:
    某加常 参辛格, 零格, 93 # sys_exit
    调用某
完毕


# void* _Wa_Runtime_malloc(int size)
函数 .Wa.Runtime.malloc:
    某加常 参乙格, 参甲格, 0     # length = size
    某加常 参甲格, 零格, 0   # addr = NULL
    某加常 参丙格, 零格, 3   # prot = PROT_READ | PROT_WRITE
    某加常 参丁格, 零格, 34  # flags = MAP_PRIVATE | MAP_ANONYMOUS
    某加常 参戊格, 零格, -1  # fd = -1
    某加常 参己格, 零格, 0   # offset = 0
    某加常 参辛格, 零格, 222 # sys_mmap (222)
    调用某
    转某址 零格, 0(回格)
完毕


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
函数 .Wa.Runtime.memcpy:
    某加常 暂甲格, 参甲格, 0 # 备份 dst
    若等转 参丙格, 零格, .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    读微整 暂乙格, 0(参乙格) # 字节读取
    写微整 暂乙格, 0(参甲格) # 字节写入
    某加常 参甲格, 参甲格, 1
    某加常 参乙格, 参乙格, 1
    某加常 参丙格, 参丙格, -1
    若异转 参丙格, 零格, .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    某加常 参甲格, 暂甲格, 0 # 返回 dst
    转某址 零格, 0(回格)
完毕


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
函数 .Wa.Runtime.memmove:
    若等转 参甲格, 参乙格, .Wa.L.memmove.done
    # 如果 参甲格 < 参乙格 (无符号比较), 跳转到向前拷贝
    若低转 参甲格, 参乙格, .Wa.Runtime.memcpy

    # 后向拷贝 (dst > src)
    某加常 暂甲格, 参甲格, 0
    某加某 参甲格, 参甲格, 参丙格
    某加某 参乙格, 参乙格, 参丙格
.Wa.L.memmove.back_loop:
    若等转 参丙格, 零格, .Wa.L.memmove.ret
    某加常 参甲格, 参甲格, -1
    某加常 参乙格, 参乙格, -1
    读微整 暂乙格, 0(参乙格)
    写微整 暂乙格, 0(参甲格)
    某加常 参丙格, 参丙格, -1
    转常址 零格, .Wa.L.memmove.back_loop
.Wa.L.memmove.ret:
    某加常 参甲格, 暂甲格, 0
.Wa.L.memmove.done:
    转某址 零格, 0(回格)
完毕


# void* _Wa_Runtime_memset(void* s, int c, int n)
函数 .Wa.Runtime.memset:
    某加常 暂甲格, 参甲格, 0
    若等转 参丙格, 零格, .Wa.L.memset.done
.Wa.L.memset.loop:
    写微整 参乙格, 0(参甲格)
    某加常 参甲格, 参甲格, 1
    某加常 参丙格, 参丙格, -1
    若异转 参丙格, 零格, .Wa.L.memset.loop
.Wa.L.memset.done:
    某加常 参甲格, 暂甲格, 0
    转某址 零格, 0(回格)
完毕


# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall_linux.print_str:
    计加常 暂甲格, %相对高位(.Wa.Memory.addr)
    读长整    暂甲格, %相对低位(.Wa.Import.syscall_linux.print_str)(暂甲格)

    # 参甲格 = ptr, 参乙格 = len
    某加常 参丙格, 参乙格, 0    # count = len
    某加某 参乙格, 暂甲格, 参甲格   # buf = base + ptr
    某加常 参甲格, 零格, 1  # fd = stdout (1)
    某加常 参辛格, 零格, 64 # sys_write
    调用某
    转某址 零格, 0(回格)
完毕


# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
函数 .Wa.Import.syscall_linux.proc_exit:
    转常址 零格, .Wa.Runtime.exit
完毕


# void _Wa_Import_syscall_linux_print_rune(int32_t c)
函数 .Wa.Import.syscall_linux.print_rune:
    某加常 栈格, 栈格, -16
    写长整 回格, 8(栈格)
    写长整 守甲格, 0(栈格)
    某加常 守甲格, 栈格, 0
    某加常 栈格, 栈格, -16

    写微整 参甲格, 0(栈格)

    某加常 参甲格, 零格, 1  # arg.0: stdout
    某加常 参乙格, 栈格, 0    # arg.1: buffer
    某加常 参丙格, 零格, 1  # arg.2: count
    某加常 参辛格, 零格, 64 # sys_write
    调用某

    某加常 栈格, 守甲格, 0
    读长整 守甲格, 0(栈格)
    读长整 回格, 8(栈格)
    读长整 回格, 8(栈格)
    某加常 栈格, 栈格, 16
完毕


# void _Wa_Import_syscall_linux_print_i64(int64_t val)
函数 .Wa.Import.syscall_linux.print_i64:
    某加常 栈格, 栈格, -16
    写长整 回格, 8(栈格)
    写长整 守甲格, 0(栈格)
    某加常 守甲格, 栈格, 0
    某加常 栈格, 栈格, -32

    某加常 暂甲格, 参甲格, 0    # 暂甲格 = 工作变量 (val)
    某加常 暂乙格, 守甲格, -1   # 暂乙格 为缓冲区指针 (从后往前填)
    某加常 暂丙格, 零格, 10 # 除数

    # 1. 处理负数
    若大转 暂甲格, 零格, .Wa.L.syscall_linux.print_i64.convert
    某减某 暂甲格, 零格, 暂甲格 # 暂甲格 = abs(暂甲格)

.Wa.L.syscall_linux.print_i64.convert:
    某除某 暂丁格, 暂甲格, 暂丙格 # 暂丁格 = 商
    某模某 暂戊格, 暂甲格, 暂丙格 # 暂戊格 = 余数
    某加常 暂戊格, 暂戊格, 48 # 加上 '0' 的 ASCII 码
    写微整 暂戊格, 0(暂乙格)  # 存入缓冲区
    某加常 暂乙格, 暂乙格, -1 # 指针前移
    某加常 暂甲格, 暂丁格, 0  # 更新待处理的数字
    若异转 暂甲格, 零格, .Wa.L.syscall_linux.print_i64.convert

    # 2. 补负号
    若大转 参甲格, 零格, .Wa.L.syscall_linux.print_i64.print
    某加常 暂戊格, 零格, 45 # '-'
    写微整 暂戊格, 0(暂乙格)
    某加常 暂乙格, 暂乙格, -1

.Wa.L.syscall_linux.print_i64.print:
    某加常 参甲格, 零格, 1  # arg.0: stdout
    某加常 参乙格, 暂乙格, 1    # arg.1: buffer
    某减某 参丙格, 守甲格, 参乙格   # arg.2: count
    某加常 参辛格, 零格, 64 # sys_write
    调用某

    某加常 栈格, 守甲格, 0
    读长整 守甲格, 0(栈格)
    读长整 回格, 8(栈格)
    某加常 栈格, 栈格, 16
    转某址 零格, 0(回格)
完毕

