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

	;; 获取栈顶
	(func $waStackPtr (result i32)
		(return (global.get $__stack_ptr))
	)

	;; 重置栈
	(func $waStackReset (param $sp i32)
		(global.set $__stack_ptr (local.get $sp))
	)

	;; 栈上分配空间
	(func $waStackAlloc (param $size i32) (result i32)
		;; $__stack_ptr -= $size
		(global.set $__stack_ptr (i32.sub (global.get $__stack_ptr) (local.get  $size)))
		;; return $__stack_ptr
		(return (global.get $__stack_ptr))
	)

	;; 释放栈上的空间
	(func $waStackFree (param $size i32)
		;; $__stack_ptr += $size
		(global.set $__stack_ptr (i32.add (global.get $__stack_ptr) (local.get $size)))
	)

	;; 堆上分配内存(没有记录大小)
	(func $waHeapAlloc (param $size i32) (result i32)
		;; local $ptr = $__heap_ptr
		(local $ptr i32) (local.get $ptr (global.get $__heap_ptr))
		;; $__heap_ptr = $__heap_ptr + $size
		(global.set $__heap_ptr (i32.add (global.get $__heap_ptr) (local.get $size)))
		;; return $ptr
		(return (local.get $ptr))
	)

	;; 重置堆上的内存
	(func $waHeapReset (param $ptr i32)
		;; $__heap_ptr = $ptr
		(global.set $__heap_ptr (local.get $ptr))
	)

	;; 分配内存
	(func $waBlockAlloc (param $size i32) (result i32)
		(return (i32.const 0)) ;; todo
	)

	;; 内存复用(引用加一)
	(func $waBlockRetain(param $ptr i32) (result i32)
		(return (local.get $ptr)) ;; todo
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
		(local $stdout i32)

		;; 保存栈指针状态
		(local.set $sp (call $waStackPtr))

		;; 分配 iov 结构体
		(local.set $p_iov (call $waStackAlloc (i32.const 8)))

		;; 返回地址
		(local.set $p_nwritten (call $waStackAlloc (i32.const 4)))

		;; 设置字符串指针和长度
		(i32.store offset=0 align=1 (local.get $p_iov) (local.get $str))
		(i32.store offset=4 align=1 (local.get $p_iov) (local.get $len))

		;; 标准输出
		(local.set $stdout (i32.const 1))

		;; 输出字符串
		(call $fd_write
			(local.get $stdout)
			(local.get $p_iov) (i32.const 1)
			(local.get $p_nwritten)
		)

		;; 重置栈指针
		(call $waStackReset (local.get $sp))
		drop
	)

	;; 打印字符
	(func $putchar (param $ch i32)
		(local $sp i32)
		(local $p_ch i32)

		;; 保存栈指针状态
		(local.set $sp (call $waStackPtr))

		;; 分配字符
		(local.set $p_ch (call $waStackAlloc (i32.const 4)))
		(i32.store offset=0 align=1 (local.get $p_ch) (local.get $ch))

		;; 输出字符
		(call $puts (local.get $p_ch) (i32.const 1))

		;; 重置栈指针
		(call $waStackReset (local.get $sp))
	)

	;; 打印整数
	(func $print_i32 (param $x i32)
		(local $div i32)
		(local $rem i32)

		;; if $x == 0 { print '0'; return }
		(i32.eq (local.get $x) (i32.const 0))
		if
			(call $putchar (i32.const 48)) ;; '0'
			(call $putchar (i32.const 10)) ;; '\n'
			(return)
		end

		;; if $x < 0 { $x = 0-$x; print '-'; }
		(i32.lt_s (local.get $x) (i32.const 0))
		if 
			(local.set $x (i32.sub (i32.const 0) (local.get $x)))
			(call $putchar (i32.const 45)) ;; '-'
		end

		;; print_i32($x / 10)
		;; puchar($x%10 + '0')
		(local.set $div (i32.div_s (local.get $x) (i32.const 10)))
		(local.set $rem (i32.rem_s (local.get $x) (i32.const 10)))
		(call $print_i32 (local.get $div))
		(call $putchar (i32.add (local.get $rem) (i32.const 48))) ;; '0'
	)

	;; 字符串常量
	(global $str.hello.ptr i32 (i32.const 10))
	(global $str.hello.len i32 (i32.const 12))
	(data $str.hello (i32.const 10) "hello world\n")

	(func $_start
		(call $puts (global.get $str.hello.ptr) (global.get $str.hello.len))
		(call $putchar (i32.const 65)) ;; 'A'
		(call $putchar (i32.const 97)) ;; 'a'
		(call $putchar (i32.const 65)) ;; 'A'
		(call $putchar (i32.const 10)) ;; '\n'

		(call $print_i32 (i32.const 123))
		(call $putchar (i32.const 10)) ;; '\n'
	)
)
