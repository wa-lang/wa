;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $hello_x64
    ;; func syscall.write(fd int64, data *byte, size int64) => int64
    (import "syscall" "write" (func $syscall.write (param i64 i64 i64) (result i64)))

    (type $syscall.write.type (func (param i64 i64 i64) (result i64)))

    (memory 1)(export "memory" (memory 0))
    (data (i32.const 8) "hello world\n")

    (table 1 funcref)
    (elem (i32.const 0) $syscall.write)

    (func $main (export "_start")
        i64.const 1  ;; stdout
        i64.const 8  ;; ptr
        i64.const 12 ;; size

        i32.const 0  ;; table[0]: $syscall.write
        call_indirect (type $syscall.write.type)
        drop ;; drop return
    )
)
