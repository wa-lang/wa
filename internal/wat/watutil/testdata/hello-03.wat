(module $hello_wasi.03
	;; type iov struct { iov_base, iov_len int32 }
	;; func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)
	(import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32) (param i32) (param i32) (param i32) (result i32)))

	(memory 1)(export "memory" (memory 0))

	;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
	(data (i32.const 8) "hello world\n")

	;; _start 类似 main 函数, 自动执行
	(func $main (export "_start"))
)
