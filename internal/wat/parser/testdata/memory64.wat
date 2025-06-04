;; https://spidermonkey.dev/blog/2025/01/15/is-memory64-actually-worth-using.html

(module
  ;; Declare an i64 memory with a size of 1 page (64KiB, or 65536 bytes)
  (memory i64 1)

  ;; Declare, and export, our store function
  (func (export "storeAt16") (param i32)
    i64.const 16  ;; push address 16 to the stack
    local.get 0   ;; get the i32 param and push it to the stack
    i32.store     ;; store the value to the address
  )

  ;; Declare, and export, our load function
  (func (export "loadFrom16") (result i32)
    i64.const 16  ;; push address 16 to the stack
    i32.load      ;; load from the address
  )
)
