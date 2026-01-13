;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $memory_01
    (import "syscall" "_write" (func $syscall.write (param i64 i64 i64) (result i64)))

    (memory 1)(export "memory" (memory 0))
    (data (i32.const 8) "hello world\n")

    (func $main (export "_start")
        i64.const 1  ;; stdout
        i64.const 8  ;; ptr
        i64.const 12 ;; size
        call $syscall.write
        drop ;; drop return
    )
)
