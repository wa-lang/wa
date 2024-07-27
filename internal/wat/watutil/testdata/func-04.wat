(module $func.04
	(func $inc (param $a i32) (result i32)
		i32.const 1
		local.get $a
		i32.add
		return
	)
)
