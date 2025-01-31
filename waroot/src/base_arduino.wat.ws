(memory $memory 1)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))  ;; 1024

(global $__heap_lfixed_cap i32 (i32.const 0)) ;; 固定尺寸空闲链表最大长度, 满时回收; 也可用于关闭 fixed 策略
