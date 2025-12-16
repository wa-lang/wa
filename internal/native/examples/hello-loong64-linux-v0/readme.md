# 龙芯汇编示例

add.s

```s
.text
.align 2

.globl add_f
.type  add_f,@function

add_f:
    add.w $a0, $a0, $a1
    add.w $a0, $a0, $a2
    add.w $a0, $a0, $a3
    jr    $ra
    .size add_f, .-add_f
```

test.c

```c
extern int add_f(int a, int b, int c, int d);
extern int write(int fd, void* p, int size);

void printch(char c) {
    write(1, &c, 1);
}

void printint(int n) {
    if(n >= 10) {
        printint(n/10);
    }
    printch((n%10)+'0');
}

int main() {
    int ret = add_f(1, 2, 3, 4);
    printint(ret);
    printch('\n');
    return 0;
}
```

编译执行:

```bash
$ gcc --save-temps test.c add.s -o test_add.out.exe
$ ./test_add.out.exe
10
```

## C代码对应的汇编

test_add.out-test.s

```s
	.file	"test.c"
	.text
	.align	2
	.globl	printch
	.type	printch, @function
printch:
.LFB0 = .
	.cfi_startproc
	addi.d	$r3,$r3,-32
	.cfi_def_cfa_offset 32
	st.d	$r1,$r3,24
	st.d	$r22,$r3,16
	.cfi_offset 1, -8
	.cfi_offset 22, -16
	addi.d	$r22,$r3,32
	.cfi_def_cfa 22, 0
	or	$r12,$r4,$r0
	st.b	$r12,$r22,-17
	addi.d	$r12,$r22,-17
	addi.w	$r6,$r0,1			# 0x1
	or	$r5,$r12,$r0
	addi.w	$r4,$r0,1			# 0x1
	bl	%plt(write)
	nop
	ld.d	$r1,$r3,24
	.cfi_restore 1
	ld.d	$r22,$r3,16
	.cfi_restore 22
	addi.d	$r3,$r3,32
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE0:
	.size	printch, .-printch
	.align	2
	.globl	printint
	.type	printint, @function
printint:
.LFB1 = .
	.cfi_startproc
	addi.d	$r3,$r3,-32
	.cfi_def_cfa_offset 32
	st.d	$r1,$r3,24
	st.d	$r22,$r3,16
	.cfi_offset 1, -8
	.cfi_offset 22, -16
	addi.d	$r22,$r3,32
	.cfi_def_cfa 22, 0
	or	$r12,$r4,$r0
	st.w	$r12,$r22,-20
	ld.w	$r12,$r22,-20
	slli.w	$r13,$r12,0
	addi.w	$r12,$r0,9			# 0x9
	ble	$r13,$r12,.L3
	ld.w	$r13,$r22,-20
	addi.w	$r12,$r0,10			# 0xa
	slli.w	$r13,$r13,0
	slli.w	$r12,$r12,0
	div.w	$r13,$r13,$r12
	bne	$r12,$r0,1f
	break	7
1:
	or	$r12,$r13,$r0
	slli.w	$r12,$r12,0
	or	$r4,$r12,$r0
	bl	printint
.L3:
	ld.w	$r13,$r22,-20
	addi.w	$r12,$r0,10			# 0xa
	slli.w	$r13,$r13,0
	slli.w	$r12,$r12,0
	mod.w	$r13,$r13,$r12
	bne	$r12,$r0,1f
	break	7
1:
	or	$r12,$r13,$r0
	slli.w	$r12,$r12,0
	bstrpick.w	$r12,$r12,7,0
	addi.w	$r12,$r12,48
	bstrpick.w	$r12,$r12,7,0
	ext.w.b	$r12,$r12
	or	$r4,$r12,$r0
	bl	printch
	nop
	ld.d	$r1,$r3,24
	.cfi_restore 1
	ld.d	$r22,$r3,16
	.cfi_restore 22
	addi.d	$r3,$r3,32
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE1:
	.size	printint, .-printint
	.align	2
	.globl	main
	.type	main, @function
main:
.LFB2 = .
	.cfi_startproc
	addi.d	$r3,$r3,-32
	.cfi_def_cfa_offset 32
	st.d	$r1,$r3,24
	st.d	$r22,$r3,16
	.cfi_offset 1, -8
	.cfi_offset 22, -16
	addi.d	$r22,$r3,32
	.cfi_def_cfa 22, 0
	addi.w	$r7,$r0,4			# 0x4
	addi.w	$r6,$r0,3			# 0x3
	addi.w	$r5,$r0,2			# 0x2
	addi.w	$r4,$r0,1			# 0x1
	bl	%plt(add_f)
	or	$r12,$r4,$r0
	st.w	$r12,$r22,-20
	ldptr.w	$r12,$r22,-20
	or	$r4,$r12,$r0
	bl	printint
	addi.w	$r4,$r0,10			# 0xa
	bl	printch
	or	$r12,$r0,$r0
	or	$r4,$r12,$r0
	ld.d	$r1,$r3,24
	.cfi_restore 1
	ld.d	$r22,$r3,16
	.cfi_restore 22
	addi.d	$r3,$r3,32
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE2:
	.size	main, .-main


	.ident	"GCC: (GNU) 14.3.0 20250523 (AOSC OS, Core)"
	.section	.note.GNU-stack,"",@progbits
```
