;; Copyright 2023 The Wa Authors. All rights reserved.

(memory $memory 1024)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
;;(global $__heap_base i32 (i32.const 1048576))     ;; index=1

(global $__heap_lfixed_cap i32 (i32.const 64)) ;; 固定尺寸空闲链表最大长度, 满时回收; 也可用于关闭 fixed 策略
