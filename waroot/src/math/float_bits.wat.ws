;; Copyright 2023 The Wa Authors. All rights reserved.

(func $$waFloat64bits (param $x f64) (result i64)
	local.get $x
	i64.reinterpret_f64
)

(func $$waFloat64frombits (param $x i64) (result f64)
	local.get $x
	f64.reinterpret_i64
)

(func $$waFloat32bits (param $x f32) (result i32)
	local.get $x
	i32.reinterpret_f32
)

(func $$waFloat32frombits (param $x i32) (result f32)
	local.get $x
	f32.reinterpret_i32
)