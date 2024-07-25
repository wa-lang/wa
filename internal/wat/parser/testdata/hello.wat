(module $hello_wasi
	;; type iov struct { iov_base, iov_len int32 }
	;; func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)
	(import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32) (param i32) (param i32) (param i32) (result i32)))

	(memory 1)(export "memory" (memory 0))

	;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
	(data (i32.const 8) "hello world\n")

	;; _start 类似 main 函数, 自动执行
	(func $main (export "_start")
		;; iov.iov_base - 字符串地址为 8
		i32.const 0
		i32.const 8
		i32.store

		;; iov.iov_len  - 字符串长度
		block $label (result i32 i32)
			i32.const 4
			i32.const 12
		end
		i32.store

		;; 输出到 stdout
		i32.const 1  ;; 1 对应 stdout
		i32.const 0  ;; *iovs - 前 8 个字节保留给 iov 数组
		i32.const 1  ;; len(iovs) - 只有1个字符串
		i32.const 20 ;; nwritten - 指针, 里面是要写到数据长度
		call $fd_write

		;; 忽略返回值
		drop
	)
)
