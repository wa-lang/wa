(module
    ;; func write(fd: i32, ptr: i32, size i32) => i32
    (import "syscall" "write" (func $syscall.write (param i32 i32 i32) (result i32)))

    (memory 1)(export "memory" (memory 0))

    ;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
    (data (i32.const 8) "hello world\n")

    ;; _start 类似 main 函数, 自动执行
    (func $main (export "_start")
        i32.const 8  ;; 字符串地址为 8
        i32.const 12 ;; 字符串长度
        call $syscall.write
        drop ;; 忽略返回值
    )
)
