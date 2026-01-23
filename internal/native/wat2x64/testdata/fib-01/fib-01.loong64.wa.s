# 源文件: fib-01.wat, ABI: loong64
# 自动生成的代码, 不要手动修改!!!

# 运行时函数
.extern write
.extern exit
.extern malloc
.extern memcpy
.extern memset
.extern memmove
.set .Runtime.write, write
.set .Runtime.exit, exit
.set .Runtime.malloc, malloc
.set .Runtime.memcpy, memcpy
.set .Runtime.memset, memset
.set .Runtime.memmove, memmove

# 导入函数(外部库定义)
.extern wat2la_env_print_i64
.set .Import.env.print_i64, wat2la_env_print_i64

# 定义内存
.section .data
.align 3
.globl .Memory.addr
.globl .Memory.pages
.globl .Memory.maxPages
.Memory.addr: .quad 0
.Memory.pages: .quad 1
.Memory.maxPages: .quad 1

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
    # local N: i64
    # local i: i64

    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -32

    # 没有参数需要备份到栈

    # 没有返回值变量需要初始化为0

    # 将局部变量初始化为0
    st.d   $zero, $fp, -8 # local N = 0
    st.d   $zero, $fp, -16 # local i = 0

    # i64.const 10
    addi.d $t0, $zero, 10
    st.d   $t0, $fp, -24

   # local.set N
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -8

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -24

   # local.set i
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16

    # loop.begin(name=my_loop, suffix=00000000)
.L.brNext.my_loop.00000000:
    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # call fib(...)
    ld.d $a0, $fp, -24 # arg 0
    pcalau12i $t0, %pc_hi20(.F.fib)
    addi.d $t0, $t0, %pc_lo12(.F.fib)
    jirl $ra, $t0, 0
    st.d $a0, $fp, -24
    # call env.print_i64(...)
    ld.d $a0, $fp, -24 # arg 0
    pcalau12i $t0, %pc_hi20(.Import.env.print_i64)
    addi.d $t0, $t0, %pc_lo12(.Import.env.print_i64)
    jirl $ra, $t0, 0
    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -32

    # i64.add
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    add.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

   # local.set i
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16

    # local.get i
    ld.d $t0, $fp, -16
    st.d $t0, $fp, -24

    # local.get N
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -32

    # i64.lt_u
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sltu    $t1, $t0, $t1
    st.w    $t1, $fp, -24

    # br_if my_loop
    ld.w $t0, $fp, -24
    beqz $t0, .L.brFallthrough.my_loop.00000000
    b .L.brNext.my_loop.00000000
.L.brFallthrough.my_loop.00000000:

    # loop.end(name=my_loop, suffix=00000000)


    # 根据ABI处理返回值
.L.return.main:

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

# func fib(n:i64) => i64
.section .text
.globl .F.fib
.F.fib:
    addi.d  $sp, $sp, -16
    st.d    $ra, $sp, 8
    st.d    $fp, $sp, 0
    addi.d  $fp, $sp, 0
    addi.d  $sp, $sp, -48

    # 将寄存器参数备份到栈
    st.d $a0, $fp, -8 # save arg.0
    # 将返回值变量初始化为0
    st.d $zero, $fp, -16 # ret.0 = 0

    # 没有局部变量需要初始化为0

    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -24

    # i64.const 2
    addi.d $t0, $zero, 2
    st.d   $t0, $fp, -32

    # i64.le_u
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sltu    $t1, $t1, $t0
    xori    $t1, $t1, 1
    st.w    $t1, $fp, -24

    # if.begin(name=, suffix=00000001)
    ld.w $t0, $fp, -24
    beqz $t0, .L.else.00000001
    # if.body(name=, suffix=00000001)
    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -24

    b   .L.brNext.00000001

    # if.else(name=, suffix=00000001)
.L.else.00000001:
    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -24

    # i64.const 1
    addi.d $t0, $zero, 1
    st.d   $t0, $fp, -32

    # i64.sub
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    sub.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

    # call fib(...)
    ld.d $a0, $fp, -24 # arg 0
    pcalau12i $t0, %pc_hi20(.F.fib)
    addi.d $t0, $t0, %pc_lo12(.F.fib)
    jirl $ra, $t0, 0
    st.d $a0, $fp, -24
    # local.get n
    ld.d $t0, $fp, -8
    st.d $t0, $fp, -32

    # i64.const 2
    addi.d $t0, $zero, 2
    st.d   $t0, $fp, -40

    # i64.sub
    ld.d    $t0, $fp, -32
    ld.d    $t1, $fp, -40
    sub.d   $t0, $t0, $t1
    st.d    $t0, $fp, -32

    # call fib(...)
    ld.d $a0, $fp, -32 # arg 0
    pcalau12i $t0, %pc_hi20(.F.fib)
    addi.d $t0, $t0, %pc_lo12(.F.fib)
    jirl $ra, $t0, 0
    st.d $a0, $fp, -32
    # i64.add
    ld.d    $t0, $fp, -24
    ld.d    $t1, $fp, -32
    add.d   $t0, $t0, $t1
    st.d    $t0, $fp, -24

.L.brNext.00000001:
    # if.end(name=, suffix=00000001)

    # return
    # copy result from stack
    ld.d $t0, $fp, -24
    st.d $t0, $fp, -16 # ret.0

    # 根据ABI处理返回值
.L.return.fib:
    ld.d $a0, $fp, -16 # ret .F.ret.0

    # 函数返回
    addi.d  $sp, $fp, 0
    ld.d    $ra, $sp, 8
    ld.d    $fp, $sp, 0
    addi.d  $sp, $sp, 16
    jirl    $zero, $ra, 0

