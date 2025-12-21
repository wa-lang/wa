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

    .text
    .align  2
    .globl  add
    .type   add, @function
add:
    addi.d  $sp, $sp, -64
    st.d    $fp, $sp, 56
    addi.d  $fp, $sp, 64

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
    ld.d    $fp, $sp, 56
    addi.d  $sp, $sp, 64
    jr      $ra

# ---------------------------------------------------------

    .align  2
    .globl  main
    .type   main, @function
main:
    addi.d  $sp, $sp, -48
    st.d    $ra, $sp, 40
    st.d    $fp, $sp, 32
    addi.d  $fp, $sp, 48

    addi.w  $t0, $zero, 9      # 0x9
    st.d    $t0, $sp, 8
    addi.w  $t0, $zero, 8      # 0x8
    stptr.d $t0, $sp, 0
    addi.w  $a7, $zero, 7      # 0x7
    addi.w  $a6, $zero, 6      # 0x6
    addi.w  $a5, $zero, 5      # 0x5
    addi.w  $a4, $zero, 4      # 0x4
    addi.w  $a3, $zero, 3      # 0x3
    addi.w  $a2, $zero, 2      # 0x2
    addi.w  $a1, $zero, 1      # 0x1
    or      $a0, $zero, $zero  # 0x0
    bl      add
    or      $t0, $a0, $zero
    or      $t1, $a1, $zero
    st.d    $t0, $fp, -32
    st.d    $t1, $fp, -24
    or      $t0, $zero, $zero
    or      $a0, $t0, $zero

    ld.d    $ra, $sp, 40
    ld.d    $fp, $sp, 32
    addi.d  $sp, $sp, 48
    jr      $ra
```
