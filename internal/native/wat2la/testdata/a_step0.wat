;; Copyright (C) 2025 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $hello_la64
    ;; func syscall.write(fd int64, data *byte, size int64) => int64
    (import "syscall" "write" (func $syscall.write (param i64 i64 i64) (result i64)))
    (import "syscall" "exit" (func $syscall.exit (param i64)))

    (memory 1)(export "memory" (memory 0))

    ;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
    (data (i32.const 8) "hello world\n")

    ;; _start 类似 main 函数, 自动执行
    (func $main (export "_start")
        ;; todo
    )
)
