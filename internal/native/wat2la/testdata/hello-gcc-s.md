# C代码对应的汇编

构建C函数调用，参数和返回值要够多，需要借助栈传递：

```c
struct ret_t {
    int v0, v1, v2, v3;
};

struct ret_t add(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9) {
    struct ret_t v;
    v.v0 = a0;
    return v;
}

int main() {
    struct ret_t x = add(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);
    return 0;
}
```

构建输出对应的汇编:

```bash
$ gcc --save-temps hello-gcc.c -o a.out.exe
```

汇编代码:

```s
    .file    "hello-gcc.c"

# ---------------------------------------------------------
# sizeof(int) == 4

    .text
    .align  2
    .globl  add
    .type   add, @function
add:
    # 栈帧大小 64 字节, sp 必须是 16 的倍数
    # 其中 ra 没有变化没有保存
    # 备份之前的 fp, 然后赋值为当前的 sp
    # 调整 sp, 分配栈帧空间

    addi.d  $sp, $sp, -64
    st.d    $fp, $sp, 56
    addi.d  $fp, $sp, 64

    # 将 a0-a7 备份到栈上
    # 对于非叶子函数, 这些寄存器可能会在调用其他函数时破坏

    or      $t0, $a0, $zero
    or      $t7, $a1, $zero
    or      $t6, $a2, $zero
    or      $t5, $a3, $zero
    or      $t4, $a4, $zero
    or      $t3, $a5, $zero
    or      $t2, $a6, $zero
    or      $t1, $a7, $zero

    st.w    $t0, $fp, -36
    or      $t0, $t7, $zero
    st.w    $t0, $fp, -40
    or      $t0, $t6, $zero
    st.w    $t0, $fp, -44
    or      $t0, $t5, $zero
    st.w    $t0, $fp, -48
    or      $t0, $t4, $zero
    st.w    $t0, $fp, -52
    or      $t0, $t3, $zero
    st.w    $t0, $fp, -56
    or      $t0, $t2, $zero
    st.w    $t0, $fp, -60
    or      $t0, $t1, $zero
    st.w    $t0, $fp, -64

    ld.w    $t0, $fp, -36
    st.w    $t0, $fp, -32
    ld.d    $t0, $fp, -32
    ld.d    $t1, $fp, -24
    or      $t2, $t0, $zero
    or      $t3, $t1, $zero
    or      $a0, $t2, $zero
    or      $a1, $t3, $zero

    # 恢复 fp 和 sp, 函数返回

    ld.d    $fp, $sp, 56
    addi.d  $sp, $sp, 64
    jr      $ra

# ---------------------------------------------------------

    .align  2
    .globl  main
    .type   main, @function
main:
    # 栈帧大小 48 字节, sp 必须是 16 的倍数
    # 备份之前的 ra, 调用函数后会被覆盖
    # 备份之前的 fp, 然后赋值为当前的 sp
    # 调整 sp, 分配栈帧空间

    addi.d  $sp, $sp, -48
    st.d    $ra, $sp, 40
    st.d    $fp, $sp, 32
    addi.d  $fp, $sp, 48

    # 超出寄存器参数范围, 通过栈传递
    addi.w  $t0, $zero, 9      # 0x9
    st.d    $t0, $sp, 8
    addi.w  $t0, $zero, 8      # 0x8
    stptr.d $t0, $sp, 0

    # 前 8 个参数用寄存器传递
    addi.w  $a7, $zero, 7      # 0x7
    addi.w  $a6, $zero, 6      # 0x6
    addi.w  $a5, $zero, 5      # 0x5
    addi.w  $a4, $zero, 4      # 0x4
    addi.w  $a3, $zero, 3      # 0x3
    addi.w  $a2, $zero, 2      # 0x2
    addi.w  $a1, $zero, 1      # 0x1
    or      $a0, $zero, $zero  # 0x0

    # 调用 add 函数
    # 此时 sp 指向第 9/10 个参数在栈上的位置

    bl      add

    # 接收函数的返回值
    # 超出的部分在栈上传递

    or      $t0, $a0, $zero
    or      $t1, $a1, $zero
    st.d    $t0, $fp, -32
    st.d    $t1, $fp, -24
    or      $t0, $zero, $zero
    or      $a0, $t0, $zero

    # 恢复 ra/fp/sp, 函数返回

    ld.d    $ra, $sp, 40
    ld.d    $fp, $sp, 32
    addi.d  $sp, $sp, 48
    jr      $ra
```
