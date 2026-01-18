;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $abi_args_01
    ;; 导入一个 6 参数的函数
    ;; 验证点：
    ;; 1-4 参数通过寄存器传递
    ;; 5-6 参数必须通过栈传递（在影子空间之上）
    (import "env" "write" 
        (func $env.write 
            (param i64 i64 i64 i64 i64 i64 i64 i64)
            (result i64)
        )
    )
    (import "env" "print_i64" (func $env.print_i64 (param i64)))

    (memory 1)
    (data (i32.const 8) "hello world\n")

    (func $main (export "_start")
        ;;                           Windos ABI            Linux ABI
        i64.const 1         ;; p1 -> RCX (stdout)          RDI
        i64.const 8         ;; p2 -> RDX (buffer offset)   RSI
        i64.const 12        ;; p3 -> R8  (size)            RDX
        i64.const 100       ;; p4 -> R9  (extra info)      RCX
        i64.const 200       ;; p5 -> Stack [RSP + 32]      R8
        i64.const 300       ;; p6 -> Stack [RSP + 40]      R9
        i64.const 400       ;; p7 -> Stack [RSP + 48]      Stack [RSP + 0]
        i64.const 500       ;; p8 -> Stack [RSP + 56]      Stack [RSP + 8]

        call $env.write
        drop
    )
)