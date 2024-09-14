;; Copyright 2024 The Wa Authors. All rights reserved.

(func $$syscall/wasi.__linkname__string_data_ptr (param $b i32) (param $d i32) (param $l i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$syscall/wasi.__linkname__slice_data_ptr (param $b i32) (param $d i32) (param $l i32) (param $c i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$syscall/wasi.__linkname__make_slice (param $blk i32) (param $ptr i32) (param $len i32) (param $cap i32) (result i32 i32 i32 i32)
	local.get $blk
	local.get $ptr
	local.get $len
	local.get $cap
	return
)
