(module $label.01
	(func $foo (result i64)
		i32.const 1
		if $label (result i64 i32)
			i64.const 3
			i32.const 4
		else
			i64.const 5
			i32.const 6
		end
		drop
		return
	)
)
