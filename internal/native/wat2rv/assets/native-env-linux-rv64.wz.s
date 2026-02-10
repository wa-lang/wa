# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-或-later

# TODO: 改为RISCV中文汇编

声明 .Wa.Memory.addr

# void _Wa_Import_syscall_write (int fd, uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall.write:
    计齐加高12立 $暂甲格, %相对.高20(.Wa.Memory.addr)
    装载.长      $暂甲格, $暂甲格, %相对.低12(.Wa.Memory.addr)

    # 参甲格 = fd
    加.长 $参乙格, $参乙格, $暂甲格 # 参乙格 = base + ptr
    # 参丙格 = len
    加立.长  $参辛格, $零格, 64 # sys_write
    系统调用 0
    链接跳转 $零格, $回格, 0
完毕


# int _Wa_Runtime_write(int fd, void *buf, int count)
函数 .Wa.Runtime.write:
    加立.长  $参辛格, $零格, 64 # sys_write
    系统调用 0
    链接跳转 $零格, $回格, 0
完毕


# void _Wa_Runtime_exit(int status)
函数 .Wa.Runtime.exit:
    加立.长  $参辛格, $零格, 93 # sys_exit
    系统调用 0
完毕


# void* _Wa_Runtime_malloc(int size)
函数 .Wa.Runtime.malloc:
    或       $参乙格, $参甲格, $零格 # length = size
    加立.长  $参甲格, $零格, 0       # addr = NULL
    加立.长  $参丙格, $零格, 3       # prot = PROT_READ | PROT_WRITE
    加立.长  $参丁格, $零格, 34      # flags = MAP_PRIVATE | MAP_ANONYMOUS
    加立.长  $参戊格, $零格, -1      # fd = -1
    加立.长  $参己格, $零格, 0       # offset = 0
    加立.长  $参辛格, $零格, 222     # sys_mmap (222)
    系统调用 0
    链接跳转 $零格, $回格, 0
完毕


# void* _Wa_Runtime_memcpy(void* dst, const void* src, int n)
函数 .Wa.Runtime.memcpy:
    或        $暂甲格, $参甲格, $零格 # 备份 dst
    跳转.相等 $参丙格, $零格, .Wa.L.memcpy.done
.Wa.L.memcpy.loop:
    装载.微     $暂乙格, $参乙格, 0 # 字节读取
    存储.微     $暂乙格, $参甲格, 0 # 字节写入
    加立.长     $参甲格, $参甲格, 1
    加立.长     $参乙格, $参乙格, 1
    加立.长     $参丙格, $参丙格, -1
    跳转.不相等 $参丙格, $零格, .Wa.L.memcpy.loop
.Wa.L.memcpy.done:
    或       $参甲格, $暂甲格, $零格 # 返回 dst
    链接跳转 $零格, $回格, 0
完毕


# void* _Wa_Runtime_memmove(void* dst, const void* src, int n)
函数 .Wa.Runtime.memmove:
    跳转.相等 $参甲格, $参乙格, .Wa.L.memmove.done
    # 如果 参甲格 < 参乙格 (无符号比较), 跳转到向前拷贝
    跳转.可返.小于.正 $参甲格, $参乙格, .Wa.Runtime.memcpy

    # 后向拷贝 (dst > src)
    或    $暂甲格, $参甲格, $零格
    加.长 $参甲格, $参甲格, $参丙格
    加.长 $参乙格, $参乙格, $参丙格
.Wa.L.memmove.back_loop:
    跳转.相等 $参丙格, $零格, .Wa.L.memmove.ret
    加立.长   $参甲格, $参甲格, -1
    加立.长   $参乙格, $参乙格, -1
    装载.微   $暂乙格, $参乙格, 0
    存储.微   $暂乙格, $参甲格, 0
    加立.长   $参丙格, $参丙格, -1
    跳转      .Wa.L.memmove.back_loop
.Wa.L.memmove.ret:
    或 $参甲格, $暂甲格, $零格
.Wa.L.memmove.done:
    链接跳转 $零格, $回格, 0
完毕


# void* _Wa_Runtime_memset(void* s, int c, int n)
函数 .Wa.Runtime.memset:
    或        $暂甲格, $参甲格, $零格
    跳转.相等 $参丙格, $零格, .Wa.L.memset.done
.Wa.L.memset.loop:
    存储.微     $参乙格, $参甲格, 0
    加立.长     $参甲格, $参甲格, 1
    加立.长     $参丙格, $参丙格, -1
    跳转.不相等 $参丙格, $零格, .Wa.L.memset.loop
.Wa.L.memset.done:
    或       $参甲格, $暂甲格, $零格
    链接跳转 $零格, $回格, 0
完毕


# void _Wa_Import_syscall_linux_print_str (uint32_t ptr, int32_t len)
函数 .Wa.Import.syscall_linux.print_str:
    计齐加高12立 $暂甲格, %相对.高20(.Wa.Memory.addr)
    装载.长      $暂甲格, $暂甲格, %相对.低12(.Wa.Memory.addr)

    # 参甲格 = ptr, 参乙格 = len
    或       $参丙格, $参乙格, $零格   # count = len
    加.长    $参乙格, $暂甲格, $参甲格 # buf = base + ptr
    加立.长  $参甲格, $零格, 1         # fd = stdout (1)
    加立.长  $参辛格, $零格, 64        # sys_write
    系统调用 0
    链接跳转 $零格, $回格, 0
完毕


# void _Wa_Import_syscall_linux_proc_exit(int32_t code)
函数 .Wa.Import.syscall_linux.proc_exit:
    跳转 .Wa.Runtime.exit
完毕


# void _Wa_Import_syscall_linux_print_rune(int32_t c)
函数 .Wa.Import.syscall_linux.print_rune:
    加立.长 $栈格, $栈格, -16
    存储.长 $回格, $栈格, 8
    存储.长 $帧格, $栈格, 0
    加立.长 $帧格, $栈格, 0
    加立.长 $栈格, $栈格, -16

    存储.微 $参甲格, $栈格, 0

    加立.长  $参甲格, $零格, 1  # arg.0: stdout
    加立.长  $参乙格, $栈格, 0  # arg.1: buffer
    加立.长  $参丙格, $零格, 1  # arg.2: count
    加立.长  $参辛格, $零格, 64 # sys_write
    系统调用 0

    加立.长  $栈格, $帧格, 0
    装载.长  $帧格, $栈格, 0
    装载.长  $回格, $栈格, 8
    加立.长  $栈格, $栈格, 16
    链接跳转 $零格, $回格, 0
完毕


# void _Wa_Import_syscall_linux_print_i64(int64_t val)
函数 .Wa.Import.syscall_linux.print_i64:
    加立.长 $栈格, $栈格, -16
    存储.长 $回格, $栈格, 8
    存储.长 $帧格, $栈格, 0
    加立.长 $帧格, $栈格, 0
    加立.长 $栈格, $栈格, -32

    或      $暂甲格, $零格, $参甲格 # 暂甲格 = 工作变量 (val)
    加立.长 $暂乙格, $帧格, -1      # 暂乙格 为缓冲区指针 (从后往前填)
    加立.长 $暂丙格, $零格, 10      # 除数

    # 1. 处理负数
    跳转.大于等于 $暂甲格, $零格, .Wa.L.syscall_linux.print_i64.convert
    减.长         $暂甲格, $零格, $暂甲格 # 暂甲格 = abs(暂甲格)

.Wa.L.syscall_linux.print_i64.convert:
    除.长       $暂丁格, $暂甲格, $暂丙格                      # 暂丁格 = 商
    模.长       $暂戊格, $暂甲格, $暂丙格                      # 暂戊格 = 余数
    加立.字     $暂戊格, $暂戊格, 48                           # 加上 '0' 的 ASCII 码
    存储.微     $暂戊格, $暂乙格, 0                            # 存入缓冲区
    加立.长     $暂乙格, $暂乙格, -1                           # 指针前移
    或          $暂甲格, $零格, $暂丁格                        # 更新待处理的数字
    跳转.不等零 $暂甲格, .Wa.L.syscall_linux.print_i64.convert # 如果商不为 0 则继续

    # 2. 补负号
    跳转.大于等于 $参甲格, $零格, .Wa.L.syscall_linux.print_i64.print
    加立.长       $暂戊格, $零格, 45 # '-'
    存储.微       $暂戊格, $暂乙格, 0
    加立.长       $暂乙格, $暂乙格, -1

.Wa.L.syscall_linux.print_i64.print:
    加立.长  $参甲格, $零格, 1       # arg.0: stdout
    加立.长  $参乙格, $暂乙格, 1     # arg.1: buffer
    减.长    $参丙格, $帧格, $参乙格 # arg.2: count
    加立.长  $参辛格, $零格, 64      # sys_write
    系统调用 0

    加立.长  $栈格, $帧格, 0
    装载.长  $帧格, $栈格, 0
    装载.长  $回格, $栈格, 8
    加立.长  $栈格, $栈格, 16
    链接跳转 $零格, $回格, 0
完毕

