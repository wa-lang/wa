;; 版权 @2022 凹语言 作者。保留所有权利。

(module $walang
	;; WASI
	(import "wasi_snapshot_preview1" "fd_write"
		(func $fd_write (param i32 i32 i32 i32) (result i32))
	)

	(memory $memory 1 1)

	(export "memory" (memory $memory))
	(export "_start" (func $_start))

	;; stack 状态
	(global $__stack_base i32 (i32.const 1000000))
	(global $SP (mut i32) (i32.const 1000000)) ;; __stack_pointer
	(global $FP (mut i32) (i32.const 1000000)) ;; __frame_pointer

	;; heap 状态
	(global $__heap_base i32 (i32.const 1000000)) ;;
	(global $__heap_ptr (mut i32) (i32.const 1000000)) ;; __heap_ptr
	(global $__heap_top (mut i32) (i32.const 1000000)) ;; __heap_top

	;; ASCII值
	(global $k_new_line i32 (i32.const 10)) ;; '\n'
	(global $k_space i32 (i32.const 32))    ;; ' '
	(global $k_0 i32 (i32.const 48))        ;; '0'
	(global $k_a i32 (i32.const 97))        ;; 'a'
	(global $k_A i32 (i32.const 65))        ;; 'A'

	;; 栈上分配空间
	(func $__stackAlloc (param $size i32) (result i32)
		;; $SP = $SP - $size
		(set_global $SP (i32.sub (get_global $SP) (get_local $size)))
		;; return $SP
		(return (get_global $SP))
	)

	;; 保存栈帧
	(func $__stackSave
		;; $SP = $SP - 4
		(set_global $SP (i32.sub (get_global $SP) (i32.const 4)))
		;; Mem[$SP] = $FP
		(i32.store (get_global $SP) (get_global $FP))
		;; $FP = $SP + 4
		(set_global $FP (i32.add (get_global $SP) (i32.const 4)))
	)

	;; 恢复栈帧
	(func $__stackRestore
		;; $SP = $FP
		(set_global $SP (get_global $FP))
		;; $FP = Mem[$SP-4]
		(set_global $FP (i32.load (i32.sub (get_global $SP) (i32.const 4))))
	)

	;; 堆上分配内存(没有记录大小)
	(func $__heapAlloc (param $size i32) (result i32)
		;; local $ptr = $__heap_ptr
		(local $ptr i32) (set_local $ptr (get_global $__heap_ptr))
		;; $__heap_ptr = $__heap_ptr + $size
		(set_global $__heap_ptr (i32.add (get_global $__heap_ptr) (get_local $size)))
		;; return $ptr
		(return (get_local $ptr))
	)

	;; 打印字符串
	(func $puts (param $str i32) (param $len i32)
		;; todo
	)

	;; 打印整数
	(func $putint (param $x i32)
		;; todo
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

		(call $print_char (get_global $k_A))
		(call $print_char (get_global $k_new_line))

		(call $say_hello)
	)
)
