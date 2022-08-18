;; 版权 @2022 凹语言 作者。保留所有权利。

(module $walang
	(type $__wa_builtin_print_int32_fn_t (func (param $x i32)))

	(import "wasi_snapshot_preview1" "fd_write"
		(func $fd_write (param i32 i32 i32 i32) (result i32))
	)

	(memory $memory 1 1)
	(global $__heap_base i32 (i32.const 1000000))
	(global $__stack_pointer (mut i32) (i32.const 1000000))

	(global $kChar_NewLine i32 (i32.const 10)) ;; '\n'
	(global $kChar_Space i32 (i32.const 32))   ;; ' '
	(global $kChar_Zero i32 (i32.const 48))    ;; '0'
	(global $kChar_a i32 (i32.const 97))  ;; 'a'
	(global $kChar_A i32 (i32.const 65))  ;; 'A'

	(export "memory" (memory $memory))
	(export "_start" (func $_start))

	(func $__stackAlloc (param $size i32) (result i32)
		(local $sp i32)
		get_global $__stack_pointer ;; 保存 SP 指针
		get_local $size
		set_global $__stack_pointer
		return
	)
	(func $__stackRestore)
	(func $__stackSave)

	(func $__heapAlloc (param $size i32) (result i32)
		i32.const 0
		return
	)

	;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
	(data (i32.const 4) "0000")
	(data (i32.const 12) "hello world\n")

	;; _start 类似 main 函数, 自动执行
	(func $say_hello
		(i32.store (i32.const 0) (i32.const 12)) ;; iov.iov_base - 字符串地址为 8
		(i32.store (i32.const 4) (i32.const 12)) ;; iov.iov_len  - 字符串长度

		(call $fd_write
			(i32.const 1)  ;; 1 对应 stdout
			(i32.const 0)  ;; *iovs - 前 8 个字节保留给 iov 数组
			(i32.const 1)  ;; len(iovs) - 只有1个字符串
			(i32.const 20) ;; nwritten - 指针, 里面是要写到数据长度
		)
		drop ;; 忽略返回值
	)

	(func $print_char (param $x i32)
		(i32.store (i32.const 8) (get_local $x))

		(i32.store (i32.const 0) (i32.const 8))  ;; iov.iov_base - 字符串地址为 8
		(i32.store (i32.const 4) (i32.const 1)) ;; iov.iov_len  - 字符串长度

		(call $fd_write
			(i32.const 1)  ;; 1 对应 stdout
			(i32.const 0)  ;; *iovs - 前 8 个字节保留给 iov 数组
			(i32.const 1)  ;; len(iovs) - 只有1个字符串
			(i32.const 20) ;; nwritten - 指针, 里面是要写到数据长度
		)
		drop ;; 忽略返回值
	)

	(func $print_i32 (param $x i32)
		;; todo
	)

	(func $_start
		(call $say_hello)

		(call $print_char (get_global $kChar_A))
		(call $print_char (get_global $kChar_NewLine))

		(call $say_hello)
	)
)
