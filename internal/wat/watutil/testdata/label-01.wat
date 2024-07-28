(module $func.04
	(func $foo (result i32)
		i32.const 1
		if $label (result i32)
			i32.const 3
		else
			i32.const 4
		end
		return
	)
)
