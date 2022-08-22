// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_wasm

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
