;; 版权 @2022 凹语言 作者。保留所有权利。

(module $walang
	;; WASI
	(import "wasi_snapshot_preview1" "fd_write"
		(func $fd_write (param i32 i32 i32 i32) (result i32))
	)

	(memory $memory 1)

	(export "memory" (memory $memory))
	(export "_start" (func $_start))

	;; stack 状态
	(global $__stack_base i32 (i32.const 1024))      ;; index=0
	(global $__stack_ptr (mut i32) (i32.const 1024)) ;; index=1

	;; heap 状态
	(global $__heap_base i32 (i32.const 1024))      ;; index=2
	(global $__heap_ptr (mut i32) (i32.const 1024)) ;; index=3
	(global $__heap_top (mut i32) (i32.const 1024)) ;; index=4

	;; ASCII值
	(global $k_new_line i32 (i32.const 10)) ;; '\n'
	(global $k_space i32 (i32.const 32))    ;; ' '
	(global $k_0 i32 (i32.const 48))        ;; '0'
	(global $k_a i32 (i32.const 97))        ;; 'a'
	(global $k_A i32 (i32.const 65))        ;; 'A'

	(global $stdin i32 (i32.const 0))
	(global $stdout i32 (i32.const 1))
	(global $stderr i32 (i32.const 2))

	;; 获取栈顶
	(func $waStackPtr (result i32)
		(return (get_global $__stack_ptr))
	)

	;; 重置栈
	(func $waStackReset (param $sp i32)
		(set_global $__stack_ptr (get_local $sp))
	)

	;; 栈上分配空间
	(func $waStackAlloc (param $size i32) (result i32)
		;; $__stack_ptr -= $size
		(set_global $__stack_ptr (i32.sub (get_global $__stack_ptr) (get_local $size)))
		;; return $__stack_ptr
		(return (get_global $__stack_ptr))
	)

	;; 释放栈上的空间
	(func $waStackFree (param $size i32)
		;; $__stack_ptr += $size
		(set_global $__stack_ptr (i32.add (get_global $__stack_ptr) (get_local $size)))
	)

	;; 堆上分配内存(没有记录大小)
	(func $waHeapAlloc (param $size i32) (result i32)
		;; local $ptr = $__heap_ptr
		(local $ptr i32) (set_local $ptr (get_global $__heap_ptr))
		;; $__heap_ptr = $__heap_ptr + $size
		(set_global $__heap_ptr (i32.add (get_global $__heap_ptr) (get_local $size)))
		;; return $ptr
		(return (get_local $ptr))
	)

	;; 重置堆上的内存
	(func $waHeapReset (param $ptr i32)
		;; $__heap_ptr = $ptr
		(set_global $__heap_ptr (get_local $ptr))
	)

	;; 分配内存
	(func $waBlockAlloc (param $size i32) (result i32)
		(return (i32.const 0)) ;; todo
	)

	;; 内存复用(引用加一)
	(func $waBlockRetain(param $ptr i32) (result i32)
		(return (get_local $ptr)) ;; todo
	)

	;; 释放内存(引用减一)
	(func $waBlockFree (param $ptr i32)
		;; todo
	)

	;; 打印字符串
	(func $puts (param $str i32) (param $len i32)
		(local $sp i32)
		(local $p_iov i32)
		(local $p_nwritten i32)

		;; 保存栈指针状态
		(set_local $sp (call $waStackPtr))

		;; 分配 iov 结构体
		(set_local $p_iov (call $waStackAlloc (i32.const 8)))

		;; 返回地址
		(set_local $p_nwritten (call $waStackAlloc (i32.const 4)))

		;; 设置字符串指针和长度
		(i32.store offset=0 align=1 (get_local $p_iov) (get_local $str))
		(i32.store offset=4 align=1 (get_local $p_iov) (get_local $len))

		;; 输出字符串
		(call $fd_write
			(get_global $stdout)
			(get_local $p_iov) (i32.const 1)
			(get_local $p_nwritten)
		)

		;; 重置栈指针
		(call $waStackReset (get_local $sp))
		drop
	)

	;; 打印字符
	(func $putchar (param $ch i32)
		(local $sp i32)
		(local $p_ch i32)

		;; 保存栈指针状态
		(set_local $sp (call $waStackPtr))

		;; 分配字符
		(set_local $p_ch (call $waStackAlloc (i32.const 4)))
		(i32.store offset=0 align=1 (get_local $p_ch) (get_local $ch))

		;; 输出字符
		(call $puts (get_local $p_ch) (i32.const 1))

		;; 重置栈指针
		(call $waStackReset (get_local $sp))
	)

	;; 打印整数
	(func $putint (param $x i32)
		;; todo
	)

	(func $print_i32 (param $x i32)
		;; todo
	)

	;; 字符串常量
	(global $str.hello.ptr i32 (i32.const 10))
	(global $str.hello.len i32 (i32.const 12))
	(data (i32.const 10) "hello world\n")

	(func $_start
		(call $puts (get_global $str.hello.ptr) (get_global $str.hello.len))
		(call $putchar (get_global $k_A))
		(call $putchar (get_global $k_a))
		(call $putchar (get_global $k_A))
		(call $putchar (get_global $k_new_line))
	)
)
