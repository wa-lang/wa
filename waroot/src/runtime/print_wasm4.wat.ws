;; Copyright 2023 The Wa Authors. All rights reserved.

;; 打印字符串
(func $$runtime.waPuts (param $str i32) (param $len i32)
	;; {{$$runtime.waPuts/body/begin}}

	(local $sp i32)
	(local $p_iov i32)
	(local $p_nwritten i32)
	(local $stdout i32)

	;; 保存栈指针状态
	global.get $__stack_ptr
	local.set $sp

	;; 分配 iov 结构体
	i32.const 8
	call $runtime.stackAlloc
	local.set $p_iov

	;; 返回地址
	i32.const 4
	call $runtime.stackAlloc
	local.set $p_nwritten

	;; 设置字符串指针和长度
	local.get $p_iov
	local.get $str
	i32.store offset=0 align=1

	local.get $p_iov
	local.get $len
	i32.store offset=4 align=1

	;; 标准输出
	i32.const 1
	local.set $stdout

	;; 输出字符串
	local.get $stdout
	local.get $p_iov
	i32.const 1
	local.get $p_nwritten
	call $$runtime.fdWrite

	;; 重置栈指针
	local.get $sp
	global.set $__stack_ptr
	drop

	;; {{$$runtime.waPuts/body/end}}
)
