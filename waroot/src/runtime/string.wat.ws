(func $$wa.runtime.string_to_ptr (param $b i32) (param $d i32) (param $l i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$wa.runtime.string_to_iter (param $b i32) (param $d i32) (param $l i32) (result i32 i32 i32)
	local.get $d
	local.get $l
	i32.const 0
)
