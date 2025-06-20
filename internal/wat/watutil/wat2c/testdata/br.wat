(module
  ;; br 从多层 block 跳出, 需要处理栈的问题
  ;; return 111
  (func (export "foo") (result i32)
    (block $A (result i32) (result i32)
      (block $B
        (block $C (result i32)
            i32.const 111
            i32.const 222
            br $A
        )
        drop ;; 丢弃 $C 的返回值
      )
      ;; $B 没有返回值, 以下代码虽然不会执行, 但是必须保留, 否则静态分析报错
      i32.const 333
      i32.const 444
    )
    ;; $A 有两个返回值, 必须丢弃一个
    drop
  )
)
