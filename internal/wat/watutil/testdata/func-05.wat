(module $func.05
	(memory 1)

	(func $main
		;; iov.iov_base - 字符串地址为 8
		i32.const 0
		i32.const 8
		i32.store

		;; iov.iov_len  - 字符串长度
		block $label (result i32 i32)
			i32.const 4
			i32.const 12
		end
		i32.store

		;; 输出到 stdout
		i32.const 1  ;; 1 对应 stdout
		drop

		call $main   ;; 临时测试
	)
)
