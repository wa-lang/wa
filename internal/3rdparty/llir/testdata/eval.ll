@format = constant [6 x i8] c"%08X\0A\00"

define i32 @add(i32 %a, i32 %b) {
0:
	%result = add i32 %a, %b
	ret i32 %result
}

define i32 @sub(i32 %0, i32 %1) {
2:
	%result = sub i32 %0, %1
	ret i32 %result
}

define i32 @f(i32 %a, i32 %b) {
0:
	%tmp1 = add i32 %a, %b
	%tmp2 = sub i32 %tmp1, 1
	%tmp3 = mul i32 %tmp2, 12345678
	%tmp4 = udiv i32 %tmp3, 2
	%tmp5 = sdiv i32 %tmp4, 3
	%tmp6 = urem i32 %tmp5, 14594
	%tmp7 = srem i32 %tmp6, 1000
	%tmp8 = shl i32 %tmp7, 1
	%tmp9 = lshr i32 %tmp8, 2
	%tmp10 = mul i32 %tmp9, -1
	%tmp11 = ashr i32 %tmp10, 2
	%tmp12 = and i32 %tmp11, 249  ; 0b11111001
	%tmp13 = or i32 %tmp12, 4     ; 0b00000100
	%result = xor i32 %tmp13, 255 ; 0b11111111
	call i32(i8*, ...) @printf(i8* getelementptr ([6 x i8], [6 x i8]* @format, i32 0, i32 0), i32 %result)
	ret i32 %result
}

define i32 @main() {
0:
	%tmp1 = call i32 @add(i32 -1, i32 3)
	%tmp2 = call i32 @sub(i32 13, i32 5)
	%result = call i32 @f(i32 %tmp1, i32 %tmp2)
	ret i32 %result
}

declare i32 @printf(i8*, ...)
