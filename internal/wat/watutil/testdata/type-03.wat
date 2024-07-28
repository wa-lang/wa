;; 确保类型的顺序
(module $type.03
	(import "env" "foo" (func))
	(func $bar
		block (result i32 i64)
			i32.const 12
			i64.const 34
		end
		drop
		drop
	)
	(func $bar2 (param i32) (result i32 i32)
		block (result i32 i32)
			i32.const 56
			i32.const 78
		end
		return
	)
)
