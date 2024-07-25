(module $hello_wasi
	;; 导入参数带名字
	;; type iov struct { iov_base, iov_len int32 }
	;; func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)
	(import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param $a i32) (param $b i32) (param i32) (param i32) (result i32)))
)
