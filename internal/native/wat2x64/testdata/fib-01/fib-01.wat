;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $fib_01
    (import "env" "print_i64" (func $env.print_i64 (param i64)))

    (memory 1)

    (func $main (export "_start")
        (local $N i64)
        (local $i i64)

        i64.const 5
        ;;call $fib
        call $env.print_i64

        ;;i64.const 10
        ;;local.set $N
        ;;
        ;;i64.const 1
        ;;local.set $i
        ;;
        ;;loop $my_loop
        ;;    local.get $i
        ;;    call $fib
        ;;    call $env.print_i64
        ;;
        ;;    local.get $i
        ;;    i64.const 1
        ;;    i64.add
        ;;    local.set $i
        ;;
        ;;    local.get $i
        ;;    local.get $N
        ;;    i64.lt_u
        ;;    br_if $my_loop
        ;;end
    )

    (func $fib (export "fib") (param $n i64) (result i64)
        local.get $n

        i64.const 2
        i64.le_u
        if (result i64)
            i64.const 1
        else
            local.get $n
            i64.const 1
            i64.sub
            call $fib ;; fib(n-1)

            local.get $n
            i64.const 2
            i64.sub
            call $fib ;; fib(n-2)

            i64.add
        end
    )
)
