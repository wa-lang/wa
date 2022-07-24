declare i32 @ugo_builtin_exit(i32)

define i32 @ugo_main_main() {
	%t0 = add i32 0, 40    ; t0 = 40
	%t1 = add i32 0, 2     ; t1 = 2
	%t2 = add i32 %t0, %t1 ; t2 = t1 + t1
	call i32(i32) @ugo_builtin_exit(i32 %t2)
	ret i32 0
}

define i32 @main() {
	call i32() @ugo_main_main()
	ret i32 0
}
