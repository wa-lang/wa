(module $table.01
	(table 3 funcref)
	(elem (i32.const 1) $$u8.$$block.$$onFree)
	(elem (i32.const 2) $$string.underlying.$$onFree)

	(func $$u8.$$block.$$onFree (param $ptr i32))
	(func $$string.underlying.$$onFree (param $$ptr i32))
)