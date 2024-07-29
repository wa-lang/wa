(module $__walang__

;; Copyright 2023 The Wa Authors. All rights reserved.

(memory $memory 1024)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
(global $__heap_base i32 (i32.const 1048576))     ;; index=1
(global $__heap_max  i32       (i32.const 67108864)) ;; 64MB, 1024 page

;; Copyright 2023 The Wa Authors. All rights reserved.

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

(func $runtime.heapMax(result i32)
	global.get $__heap_max
)

) ;; module