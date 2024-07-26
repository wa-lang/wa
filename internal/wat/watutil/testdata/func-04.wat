(module $func.04
	(memory 1)

	(func $foo (result i32)
		i32.const 13
		return
	)

	(func $main
		block $label (result i32 i64)
			i32.const 4
			i64.const 15
		end

		drop
		drop

		call $foo
		drop
	)
)
