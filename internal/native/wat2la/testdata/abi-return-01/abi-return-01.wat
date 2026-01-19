;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $abi_return_01
    (import "env" "get_multi_values" 
        (func $env.get_multi_values (param i64) (result i64 i64 i64))
    )
    (import "env" "print_i64" 
        (func $env.print_i64 (param i64))
    )

    (memory 1)

    (func $main (export "_start")
        (local $input i64)

        i64.const 100
        local.set $input

        ;; 调用多返回值函数
        local.get $input
        call $env.get_multi_values
        
        ;; 此时栈上有 3 个值 [v1, v2, v3]
        call $env.print_i64 ;; v3
        call $env.print_i64 ;; v2
        call $env.print_i64 ;; v1
    )
)