;; Copyright 2023 The Wa Authors. All rights reserved.

(func $$math.waSqrtF64 (param $x f64) (result f64)
	local.get $x
	f64.sqrt
)

(func $$math.waSqrtF32 (param $x f32) (result f32)
	local.get $x
	f32.sqrt
)
