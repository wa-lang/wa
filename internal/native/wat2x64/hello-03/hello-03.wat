;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $verify_win64_abi
    ;; 导入一个 6 参数的函数
    ;; 验证点：
    ;; 1-4 参数通过寄存器传递
    ;; 5-6 参数必须通过栈传递（在影子空间之上）
    (import "env" "env_write" 
        (func $env_write 
            (param i64 i64 i64 i64 i64 i64) 
            (result i64)
        )
    )

    (memory 1)
    (data (i32.const 0) "ABI Test")

    (func (export "_start")
        ;; 准备 6 个参数
        i64.const 1         ;; p1 -> RCX (stdout)
        i64.const 0         ;; p2 -> RDX (buffer offset)
        i64.const 8         ;; p3 -> R8  (size)
        i64.const 100       ;; p4 -> R9  (extra info)
        i64.const 200       ;; p5 -> Stack [RSP + 32]
        i64.const 300       ;; p6 -> Stack [RSP + 40]

        call $env_write
        drop
    )
)