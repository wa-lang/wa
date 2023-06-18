;; Copyright 2023 The Wa Authors. All rights reserved.

;; 打印字符串
(func $puts (param $str i32) (param $len i32)
	;; {{$puts/body/begin}}

	(local $sp i32)
	(local $p_iov i32)
	(local $p_nwritten i32)
	(local $stdout i32)

	;; 保存栈指针状态
	(local.set $sp (global.get $__stack_ptr))

	;; 分配 iov 结构体
	(local.set $p_iov (call $$waStackAlloc (i32.const 8)))

	;; 返回地址
	(local.set $p_nwritten (call $$waStackAlloc (i32.const 4)))

	;; 设置字符串指针和长度
	(i32.store offset=0 align=1 (local.get $p_iov) (local.get $str))
	(i32.store offset=4 align=1 (local.get $p_iov) (local.get $len))

	;; 标准输出
	(local.set $stdout (i32.const 1))

	;; 输出字符串
	(call $$runtime.fdWrite
		(local.get $stdout)
		(local.get $p_iov) (i32.const 1)
		(local.get $p_nwritten)
	)

	;; 重置栈指针
	(global.set $__stack_ptr (local.get $sp))
	drop

	;; {{$puts/body/end}}
)

;; 打印字符
(func $putchar (param $ch i32)
	;; {{$putchar/body/begin}}

	(local $sp i32)
	(local $p_ch i32)

	;; 保存栈指针状态
	(local.set $sp (global.get $__stack_ptr))

	;; 分配字符
	(local.set $p_ch (call $$waStackAlloc (i32.const 4)))
	(i32.store offset=0 align=1 (local.get $p_ch) (local.get $ch))

	;; 输出字符
	(call $puts (local.get $p_ch) (i32.const 1))

	;; 重置栈指针
	(global.set $__stack_ptr (local.get $sp))

	;; {{$putchar/body/begin}}
)
