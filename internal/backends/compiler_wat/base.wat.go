// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_wat

import (
	_ "embed"
)

//go:embed base.wat
var __base_wat_data string

// wasm 内存 1 页大小
const _WASM_PAGE_SIZE = 65536

// WASM 约定栈和内存管理
// 相关全局变量地址必须和 base_wasm.go 保持一致
const (
	// ;; heap 和 stack 状态(__heap_base 只读)
	// ;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
	// (global $__stack_ptr (mut i32) (i32.const 1024)) ;; index=0
	// (global $__heap_base i32 (i32.const 2048))       ;; index=1
	__stack_ptr_index = 0
	__heap_base_index = 1
)

// 内置函数名字
const (
	// 栈函数
	_waStackPtr   = "waStackPtr"
	_waStackAlloc = "waStackAlloc"
	_waStackFree  = "waStackFree"

	// 堆管理函数
	_waHeapPtr = "waHeapPtr"
	_waAlloc   = "waAlloc"
	_waRetain  = "waRetain"
	_waFree    = "waFree"

	// 输出函数
	_waPrintChar   = "putchar"
	_waPrintString = "puts"
	_waPrintInt32  = "print_i32"

	// 开始函数
	_waStart = "_start"
)

const modBaseWat = `
;; ----------------------------------------------------
;; import 必须最先定义
;; ----------------------------------------------------

;; WASI 最小子集
;; 用于输出字符串
(import "wasi_snapshot_preview1" "fd_write"
	(func $fd_write (param i32 i32 i32 i32) (result i32))
)

;; ----------------------------------------------------
;; 内存和入口
;; ----------------------------------------------------

(memory $memory 1)

(export "memory" (memory $memory))
(export "_start" (func $_start))
;; (export "main.main" (func $main.main))

;; ----------------------------------------------------
;; WASM 约定栈和内存管理
;; 相关全局变量地址必须和 base_wasm.go 保持一致
;; ----------------------------------------------------

;; heap 和 stack 状态(__heap_base 只读)
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
(global $__stack_ptr (mut i32) (i32.const 1024)) ;; index=0
(global $__heap_base i32 (i32.const 2048))       ;; index=1

;; ----------------------------------------------------
;; Stack 辅助函数
;; ----------------------------------------------------

;; 获取栈顶地址
(func $waStackPtr (result i32)
	(global.get $__stack_ptr)
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

;; ----------------------------------------------------
;; Heap 辅助函数
;; ----------------------------------------------------

;; 获取堆地址
(func $waHeapPtr (result i32)
	(global.get $__heap_base)
)

;; 堆上分配内存(没有记录大小)
(func $waAlloc (param $size i32) (result i32)
	;; {{$waAlloc/body/begin}}
	unreachable
	;; {{$waAlloc/body/end}}
)

;; 内存复用(引用加一)
(func $waRetain(param $ptr i32) (result i32)
	;; {{$waRetain/body/begin}}
	unreachable
	;; {{$waRetain/body/end}}
)

;; 释放内存(引用减一)
(func $waFree (param $ptr i32)
	;; {{$waFree/body/begin}}
	unreachable
	;; {{$waFree/body/end}}
)

;; ----------------------------------------------------
;; 输出函数
;; ----------------------------------------------------

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
	(global.set $__stack_ptr (local.get $sp))
	drop

	;; {{$puts/body/end}}
)

;; 打印字符
(func $__print_char (param $ch i32)
	;; {{$putchar/body/begin}}

	(local $sp i32)
	(local $p_ch i32)

	;; 保存栈指针状态
	(local.set $sp (global.get $__stack_ptr))

	;; 分配字符
	(local.set $p_ch (call $waStackAlloc (i32.const 4)))
	(i32.store offset=0 align=1 (local.get $p_ch) (local.get $ch))

	;; 输出字符
	(call $puts (local.get $p_ch) (i32.const 1))

	;; 重置栈指针
	(global.set $__stack_ptr (local.get $sp))

	;; {{$putchar/body/begin}}
)

;; 打印整数
(func $__print_i32 (param $x i32)
	;; {{$print_i32/body/begin}}

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

	;; {{$print_i32/body/end}}
)

(func $putchar (param $ch i32)
	local.get $ch
	call $__print_char
)
(func $print_i32 (param $ch i32)
	local.get $ch
	call $__print_i32
)

;; ----------------------------------------------------
;; _start 函数
;; ----------------------------------------------------

;; _start 函数
(func $_start
	;; {{$_start/body/begin}}
	;; (call $main.init)
	(call $main)
	;; {{$_start/body/end}}
)

`