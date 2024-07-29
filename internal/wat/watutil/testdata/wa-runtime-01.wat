(module $__walang__

;; Copyright 2023 The Wa Authors. All rights reserved.

(memory $memory 1024)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
;;(global $__heap_base i32 (i32.const 1048576))     ;; index=1
(global $__heap_max  i32       (i32.const 67108864)) ;; 64MB, 1024 page


)