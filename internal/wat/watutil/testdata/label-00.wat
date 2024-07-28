(module $label.00
	(func $foo (result i32)
		block
		end
		i32.const 1
		if $label (result i32)
			i32.const  1
		else
			i32.const  2
		end

		i32.const 42
		return
	)
)
