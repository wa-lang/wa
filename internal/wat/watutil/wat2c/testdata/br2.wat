(module
  ;; return 555
  (func (export "foo") (result i32)
    (block $A (result i32) (result i32)
      i32.const 111
      i32.const 222
      (block $B
        i32.const 333
        (block $C (result i32)
            i32.const 444 ;; 被丢弃
            i32.const 555
            i32.const 666
            br $A ;; 会清空之前的栈, 保留需要返回的值
        )
        drop ;; 丢弃 block $C 的返回
        drop ;; 丢弃 333
      )
    )
    drop ;; 丢弃 666
  )
)
