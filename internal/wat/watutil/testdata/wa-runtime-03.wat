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

(global $$knr_basep (mut i32) (i32.const 0))
(global $$knr_freep (mut i32) (i32.const 0))

(func $runtime.knr_getBlockHeader (param $addr i32) (result i32 i32)
  local.get $addr
  i32.load offset=0 align=4
  local.get $addr
  i32.load offset=4 align=4
) ;;runtime.knr_getBlockHeader

(func $runtime.knr_setBlockHeader (param $addr i32) (param $data.0 i32) (param $data.1 i32)
  local.get $addr
  local.get $data.0
  i32.store offset=0 align=4
  local.get $addr
  local.get $data.1
  i32.store offset=4 align=4
) ;;runtime.knr_setBlockHeader



) ;; module