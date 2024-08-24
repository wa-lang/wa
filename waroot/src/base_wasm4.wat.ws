(import "env" "memory" (memory 1))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 14656)) ;; 6560+8096
(global $__heap_max  i32       (i32.const 65536)) ;; 64KB, 1 page
