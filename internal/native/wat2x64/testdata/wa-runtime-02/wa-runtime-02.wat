;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

;; 在 wat 最简单的例子基础上, 增加凹语言 runtime 底层简单的叶子函数

(module $wa_runtime_01
    (import "syscall" "write" (func $syscall.write (param i64 i64 i64) (result i64)))

    (memory 1)(export "memory" (memory 0))
    (data (i32.const 8) "hello world\n")

    (func $main (export "_start")
        i64.const 1  ;; stdout
        i64.const 8  ;; ptr
        i64.const 12 ;; size
        call $syscall.write
        drop ;; drop return
    )

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))
(global $__heap_base i32 (i32.const 1048576))

(global $__heap_lfixed_cap i32 (i32.const 64)) ;; 固定尺寸空闲链表最大长度, 满时回收; 也可用于关闭 fixed 策略

;; ---------------------------------------------------------
;; package: runtime
;; ---------------------------------------------------------

(func $runtime.throw
	unreachable
)

(func $runtime.getStackPtr (result i32)
	global.get $__stack_ptr
)

(func $runtime.setStackPtr (param $sp i32)
	local.get $sp
	global.set $__stack_ptr
)

(func $runtime.stackAlloc (param $size i32) (result i32)
	;; $__stack_ptr -= $size
	global.get $__stack_ptr
	local.get  $size
	i32.sub 
	global.set $__stack_ptr 

	;; return $__stack_ptr
	global.get $__stack_ptr
	return
)

(func $runtime.stackFree (param $size i32)
	;; $__stack_ptr += $size
	global.get $__stack_ptr
	local.get $size
	i32.add
	global.set $__stack_ptr 
)

(func $runtime.heapBase(result i32)
	global.get $__heap_base
)

(func $runtime.HeapAlloc (export "runtime.HeapAlloc") (param $nbytes i32) (result i32) ;;result = ptr
	(local $ptr i32)

	local.get $nbytes
	i32.eqz
	if
		i32.const 0
		return
	end

	local.get $nbytes
	i32.const 7
	i32.add
	i32.const 8
	i32.div_u
	i32.const 8
	i32.mul
	local.set $nbytes

	local.get $nbytes
	call $runtime.malloc
	local.set $ptr

	loop $zero
		local.get $nbytes
		i32.const 8
		i32.sub
		local.tee $nbytes
		local.get $ptr
		i32.add

		i64.const 0
		i64.store

		local.get $nbytes
		if
			br $zero
		end
	end ;;loop $zero

	local.get $ptr
)

(func $runtime.HeapFree (export "runtime.HeapFree") (param $ptr i32)
	local.get $ptr
	call $runtime.free
)

(func $runtime.malloc (param $size i32) (result i32)
	i32.const 0 ;; todo
)

(func $runtime.free (param $ptr i32)
	;; todo
)
) ;; EOF
