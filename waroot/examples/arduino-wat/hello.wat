;; 版权 @2025 arduino-wat 作者。保留所有权利。

(module $wa-arduino-hello
	(import "arduino" "getPinLED" (func $getPinLED (result i32)))
	(import "arduino" "pinMode" (func $pinMode (param $pin i32) (param $mode i32)))
	(import "arduino" "digitalWrite" (func $digitalWrite (param $pin i32) (param $value i32)))
	(import "arduino" "delay" (func $delay (param $ms i32)))

	(memory (;0;) 1)
	(export "memory" (memory 0))

	(func $main (export "_main") (local $tmp i32)
		;; var LED i32 = getPinLED()
		;; &LED == 1024
		i32.const 1024
		call $getPinLED
		local.tee $tmp
		i32.store

		;; pinMode(LED, 1)
		local.get $tmp
		i32.const 1
		call $pinMode

		;; for {}
		loop $label0 ;; label = @1
			;; digitalWrite(LED, 1)
			i32.const 1024
			i32.load
			i32.const 1
			call $digitalWrite

			;; delay(100)
			i32.const 100
			call $delay

			;; digitalWrite(LED, 0)
			i32.const 1024
			i32.load
			i32.const 0
			call $digitalWrite

			;; delay(900)
			i32.const 900
			call $delay

			;; continue
			br $label0 (;@1;)
		end

		;; panic("unreachable")
		unreachable
	)
)
