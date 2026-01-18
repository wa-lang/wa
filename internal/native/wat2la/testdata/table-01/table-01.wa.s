# 源文件: table-01.wat, ABI: loong64
# 自动生成的代码, 不要手动修改!!!

# 运行时函数
.extern write
.extern exit
.extern malloc
.extern memcpy
.extern memset
.set .Runtime.write, write
.set .Runtime.exit, exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset

# 导入函数(外部库定义)
.extern wat2la_syscall_write
.set .Import.syscall.write, wat2la_syscall_write

# 定义内存
.section .data
.align 3
.globl .Memory.addr
.globl .Memory.pages
.globl .Memory.maxPages
.Memory.addr: .quad 0
.Memory.pages: .quad 1
.Memory.maxPages: .quad 1

# 内存数据
.section .data
.align 3
# memcpy(&Memory[8], data[0], size)
.Memory.dataOffset.0: .quad 8
.Memory.dataSize.0: .quad 12
.Memory.dataPtr.0: .asciz "hello world\n"

# 内存初始化函数
.section .text
.globl .Memory.initFunc
.Memory.initFunc:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # 分配内存
    pcalau12i $t0, %pc_hi20(.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a0, $t0, 16
    pcalau12i $t0, %pc_hi20(.Runtime.malloc)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.malloc)
    jirl      $ra, $t0, 0
    pcalau12i $t1, %pc_hi20(.Memory.addr)
    addi.d    $t1, $t1, %pc_lo12(.Memory.addr)
    st.d      $a0, $t1, 0

    # 内存清零
    addi.d    $a1, $zero, 0 # a1 = 0
    pcalau12i $t0, %pc_hi20(.Memory.maxPages)
    addi.d    $t0, $t0, %pc_lo12(.Memory.maxPages)
    ld.d      $t0, $t0, 0
    slli.d    $a2, $t0, 16
    pcalau12i $t0, %pc_hi20(.Runtime.memset)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.memset)
    jirl      $ra, $t0, 0

    # 初始化内存

    # memcpy(&Memory[8], data[0], size)
    pcalau12i $t1, %pc_hi20(.Memory.addr)
    addi.d    $t1, $t1, %pc_lo12(.Memory.addr)
    ld.d      $t1, $t1, 0
    pcalau12i $t0, %pc_hi20(.Memory.dataOffset.0)
    addi.d    $t0, $t0, %pc_lo12(.Memory.dataOffset.0)
    ld.d      $t0, $t0, 0
    add.d     $a0, $t1, $t0
    pcalau12i $a1, %pc_hi20(.Memory.dataPtr.0)
    addi.d    $a1, $a1, %pc_lo12(.Memory.dataPtr.0)
    pcalau12i $t0, %pc_hi20(.Memory.dataSize.0)
    addi.d    $t0, $t0, %pc_lo12(.Memory.dataSize.0)
    ld.d      $a2, $t0, 0
    pcalau12i $t0, %pc_hi20(.Runtime.memcpy)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.memcpy)
    jirl      $ra, $t0, 0

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# 定义表格
.section .data
.align 3
.globl .Table.addr
.globl .Table.size
.globl .Table.maxSize
.Table.addr: .quad 0
.Table.size: .quad 1
.Table.maxSize: .quad 1

# 函数列表
# 保持连续并填充全部函数
.section .data
.align 3
.Table.funcIndexList:
.Table.funcIndexList.0: .quad .Import.syscall.write
.Table.funcIndexList.1: .quad .F.main
.Table.funcIndexList.end: .quad 0

# 表格初始化函数
.section .text
.globl .Table.initFunc
.Table.initFunc:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # 分配表格
    # mov  rcx, [rip + .Table.maxSize]
    # shl  rcx, 3 # sizeof(i64) == 8
    # call .Runtime.malloc
    # mov  [rip + .Table.addr], rax

    # 表格填充 0xFF
    # mov  rcx, [rip + .Table.addr]
    # mov  rdx, 0xFF
    # mov  r8, [rip + .Table.maxSize]
    # shl  r8, 3 # sizeof(i64) == 8
    # call .Runtime.memset

    # 初始化表格

    # 加载表格地址
    # mov rax, [rip + .Table.addr]
    # elem[0]: table[0+0] = syscall.write
    # mov qword ptr [rax+0], 0

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# 汇编程序入口函数
.section .text
.globl main
main:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    pcalau12i $t0, %pc_hi20(.Memory.initFunc)
    addi.d    $t0, $t0, %pc_lo12(.Memory.initFunc)
    jirl      $ra, $t0, 0
    pcalau12i $t0, %pc_hi20(.Table.initFunc)
    addi.d    $t0, $t0, %pc_lo12(.Table.initFunc)
    jirl      $ra, $t0, 0
    pcalau12i $t0, %pc_hi20(.F.main)
    addi.d    $t0, $t0, %pc_lo12(.F.main)
    jirl      $ra, $t0, 0

    # runtime.exit(0)
    addi.d    $a0, $zero, 0 # a0 = 0
    pcalau12i $t0, %pc_hi20(.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.exit)
    jirl      $ra, $t0, 0

    # exit 后这里不会被执行, 但是依然保留
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

.section .data
.align 3
.Runtime.panic.message: .asciz "panic"
.Runtime.panic.messageLen: .quad 5

.section .text
.globl .Runtime.panic
.Runtime.panic:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # runtime.write(stderr, panicMessage, size)
    addi.d    $a0, $zero, 2
    pcalau12i $a1, %pc_hi20(.Runtime.panic.message)
    addi.d    $a1, $a1, %pc_lo12(.Runtime.panic.message)
    pcalau12i $t0, %pc_hi20(.Runtime.panic.messageLen)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.panic.messageLen)
    ld.d      $a2, $t0, 0
    pcalau12i $t0, %pc_hi20(.Runtime.write)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.write)
    jirl      $ra, $t0, 0

    # 退出程序
    addi.d    $a0, $zero, 1 # 退出码
    pcalau12i $t0, %pc_hi20(.Runtime.exit)
    addi.d    $t0, $t0, %pc_lo12(.Runtime.exit)
    jirl      $ra, $t0, 0

    # return
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# func main
.section .text
.F.main:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 没有局部变量需要初始化为0

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -8

    # i64.const 8
    addi.d $t0, $zero, 8
    st.d   $t0, $fp, -16

    # i64.const 12
    addi.d $t0, $zero, 12
    st.d   $t0, $fp, -24

    # i32.const 0
    addi.d $t0, $zero, 0
    st.d   $t0, $fp, -32

    ld.d t2, fp, -24
    ld.d R4, fp, 0    ld.d R5, fp, 8    ld.d R6, fp, 16    bl t2
    st.d R4, fp, 0
    addi.w $zero, $zero, 0 # drop [fp-8]

    # 根据ABI处理返回值
.L.return:

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

