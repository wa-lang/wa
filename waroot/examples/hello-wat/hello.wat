;; 版权 @2025 arduino-wat 作者。保留所有权利。

;; wa run hello.wat

(module $wa-hello-wat
    (import "syscall_js" "print_i32" (func $print_i32 (param $i i32)))
    (import "syscall_js" "print_rune" (func $print_rune (param $ch i32)))

    (func $_main (export "_main")
        i32.const 42
        call $print_i32
    )
)
