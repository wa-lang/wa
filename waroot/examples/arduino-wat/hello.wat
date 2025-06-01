;; 版权 @2025 arduino-wat 作者。保留所有权利。

(module $wa-arduino-hello
	(import "arduino" "getPinLED" (func $getPinLED (result i32)))
	(import "arduino" "pinMode" (func $pinMode (param $pin i32) (param $mode i32)))
	(import "arduino" "digitalWrite" (func $digitalWrite (param $pin i32) (param $value i32)))
	(import "arduino" "delay" (func $delay (param $ms i32)))

	(memory (;0;) 0)
	(export "memory" (memory 0))

	(global $LED (mut i32) (i32.const 0))


	(func $_start (start)
		(local $tmp i32)

		;; global LED i32 = getPinLED()
		call $getPinLED
		global.set $LED

		;; pinMode(LED, 1)
		global.get $LED
		i32.const 1
		call $pinMode
	)

	(func $loop (export "loop")
		;; digitalWrite(LED, 1)
		global.get $LED
		i32.const 1
		call $digitalWrite

		;; delay(100)
		i32.const 100
		call $delay

		;; digitalWrite(LED, 0)
		global.get $LED
		i32.const 0
		call $digitalWrite

		;; delay(900)
		i32.const 900
		call $delay	
	)
)
