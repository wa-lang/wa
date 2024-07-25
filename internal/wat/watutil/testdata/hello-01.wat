(module $hello_wasi
	;; type iov struct { iov_base, iov_len int32 }
	;; func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)
	(import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32) (param i32) (param i32) (param i32) (result i32)))
)
