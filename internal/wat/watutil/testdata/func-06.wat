(module $func.06
	(memory 1)

	(func $foo (result i32)
		i32.const 13
		return
	)

	(func $main
		block $label (result i32 f64)
			i32.const 4
			f64.const 7

			block (result i32 i64 i64)
				i32.const 4
				i64.const 15
				i64.const 17
			end

			drop
			drop
			drop

		end

		drop
		drop


		block (result i64 i64)
			i64.const 14
			i64.const 15
		end
		drop
		drop

		call $foo
		drop
	)

	(func $foo2 (result i64 i32)
		i64.const 1
		i32.const 2
		return
	)

	;; 相同类型被合并
	(func $foo3 (result i64 i32)
		i64.const 1
		i32.const 2
		return
	)
)
