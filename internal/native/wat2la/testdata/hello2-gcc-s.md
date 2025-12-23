# C代码对应的汇编

构建C函数调用，参数和返回值要够多，需要借助栈传递：

```c
struct ret_t {
    int v0, v1, v2, v3;
};

int add2(int a0, int a1) {
    return a0+a1;
}
struct ret_t add10(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9) {
    struct ret_t v;
    v.v3 = a0+a9;
    return v;
}

int main() {
    int v0 = add2(1000, 200); 
    struct ret_t x = add10(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);
    return 0;
}
```

构建输出对应的汇编:

```bash
$ gcc --save-temps hello2-gcc.c -o a.out.exe
```

汇编代码:

```s
	.file	"hello2-gcc.c"
	.text
	.align	2
	.globl	add2
	.type	add2, @function
add2:
.LFB0 = .
	.cfi_startproc
	addi.d	$r3,$r3,-32
	.cfi_def_cfa_offset 32
	st.d	$r22,$r3,24
	.cfi_offset 22, -8
	addi.d	$r22,$r3,32
	.cfi_def_cfa 22, 0
	or	$r12,$r4,$r0
	or	$r13,$r5,$r0
	st.w	$r12,$r22,-20
	or	$r12,$r13,$r0
	st.w	$r12,$r22,-24
	ld.w	$r13,$r22,-20
	ld.w	$r12,$r22,-24
	add.w	$r12,$r13,$r12
	slli.w	$r12,$r12,0
	or	$r4,$r12,$r0
	ld.d	$r22,$r3,24
	.cfi_restore 22
	addi.d	$r3,$r3,32
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE0:
	.size	add2, .-add2
	.align	2
	.globl	add10
	.type	add10, @function
add10:
.LFB1 = .
	.cfi_startproc
	addi.d	$r3,$r3,-64
	.cfi_def_cfa_offset 64
	st.d	$r22,$r3,56
	.cfi_offset 22, -8
	addi.d	$r22,$r3,64
	.cfi_def_cfa 22, 0
	or	$r12,$r4,$r0
	or	$r19,$r5,$r0
	or	$r18,$r6,$r0
	or	$r17,$r7,$r0
	or	$r16,$r8,$r0
	or	$r15,$r9,$r0
	or	$r14,$r10,$r0
	or	$r13,$r11,$r0
	st.w	$r12,$r22,-36
	or	$r12,$r19,$r0
	st.w	$r12,$r22,-40
	or	$r12,$r18,$r0
	st.w	$r12,$r22,-44
	or	$r12,$r17,$r0
	st.w	$r12,$r22,-48
	or	$r12,$r16,$r0
	st.w	$r12,$r22,-52
	or	$r12,$r15,$r0
	st.w	$r12,$r22,-56
	or	$r12,$r14,$r0
	st.w	$r12,$r22,-60
	or	$r12,$r13,$r0
	st.w	$r12,$r22,-64
	ld.w	$r13,$r22,-36
	ld.w	$r12,$r22,8
	add.w	$r12,$r13,$r12
	slli.w	$r12,$r12,0
	st.w	$r12,$r22,-20
	ld.d	$r12,$r22,-32
	ld.d	$r13,$r22,-24
	or	$r14,$r12,$r0
	or	$r15,$r13,$r0
	or	$r4,$r14,$r0
	or	$r5,$r15,$r0
	ld.d	$r22,$r3,56
	.cfi_restore 22
	addi.d	$r3,$r3,64
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE1:
	.size	add10, .-add10
	.align	2
	.globl	main
	.type	main, @function
main:
.LFB2 = .
	.cfi_startproc
	addi.d	$r3,$r3,-64
	.cfi_def_cfa_offset 64
	st.d	$r1,$r3,56
	st.d	$r22,$r3,48
	.cfi_offset 1, -8
	.cfi_offset 22, -16
	addi.d	$r22,$r3,64
	.cfi_def_cfa 22, 0
	addi.w	$r5,$r0,200			# 0xc8
	addi.w	$r4,$r0,1000			# 0x3e8
	bl	add2
	or	$r12,$r4,$r0
	st.w	$r12,$r22,-36
	addi.w	$r12,$r0,9			# 0x9
	st.d	$r12,$r3,8
	addi.w	$r12,$r0,8			# 0x8
	stptr.d	$r12,$r3,0
	addi.w	$r11,$r0,7			# 0x7
	addi.w	$r10,$r0,6			# 0x6
	addi.w	$r9,$r0,5			# 0x5
	addi.w	$r8,$r0,4			# 0x4
	addi.w	$r7,$r0,3			# 0x3
	addi.w	$r6,$r0,2			# 0x2
	addi.w	$r5,$r0,1			# 0x1
	or	$r4,$r0,$r0
	bl	add10
	or	$r12,$r4,$r0
	or	$r13,$r5,$r0
	st.d	$r12,$r22,-32
	st.d	$r13,$r22,-24
	or	$r12,$r0,$r0
	or	$r4,$r12,$r0
	ld.d	$r1,$r3,56
	.cfi_restore 1
	ld.d	$r22,$r3,48
	.cfi_restore 22
	addi.d	$r3,$r3,64
	.cfi_def_cfa_register 3
	jr	$r1
	.cfi_endproc
.LFE2:
	.size	main, .-main


	.ident	"GCC: (GNU) 14.3.0 20250523 (AOSC OS, Core)"
	.section	.note.GNU-stack,"",@progbits
```
