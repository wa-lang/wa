(memory $memory 1)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))  ;; 1024
(global $__heap_max  i32       (i32.const 65536)) ;; 64KB, 1 page
