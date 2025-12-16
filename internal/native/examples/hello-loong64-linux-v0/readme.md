# 龙芯汇编示例

- SP 寄存器必须 16 字节对齐

## 打印字符串

```
.text
.align 2

.globl printstring
.type  printstring,@function

# void printstring(const char* s, int len);
printstring:
    # $sp = $sp - 16, sp 需要 16 字节对齐
    # $ra = $sp + 8
    addi.d $sp, $sp, -16
    st.d   $ra, $sp, 8

    # mem[$sp+7] = $a0
    st.b $a0, $sp, 7

    # write(1, $a0, $a1)
    or $a2, $a1, $zero   # $a2 = $a1
    or $a1, $a0, $zero   # $a1 = $a0
    addi.d $a0, $zero, 1 # $a0 = 1
    addi.d $a7, $zero, 64
    syscall 0

    # return
    ld.d $ra, $sp, 8
    addi.d $sp, $sp, 16
    jr $ra
```

