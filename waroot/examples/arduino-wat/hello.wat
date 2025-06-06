;; 版权 @2025 arduino-wat 作者。保留所有权利。

(module $wa-arduino-hello
	(import "arduino" "HIGH" (global $HIGH i32))
	(import "arduino" "LOW" (global $LOW i32))

	(import "arduino" "INPUT" (global $INPUT i32))
	(import "arduino" "OUTPUT" (global $OUTPUT i32))

	(import "arduino" "LED_BUILTIN" (global $LED_BUILTIN i32))

	(import "arduino" "delay" (func $delay (param $ms i32)))
	(import "arduino" "pinMode" (func $pinMode (param $pin i32) (param $mode i32)))
	(import "arduino" "digitalWrite" (func $digitalWrite (param $pin i32) (param $value i32)))
	(import "arduino" "println" (func $println (param $ptr i32) (param $len i32)))

	(memory $memory 1)

	(global $hello_ptr i32 (i32.const 8))
	(global $hello_len i32 (i32.const 14))
	(data (i32.const 8) "hello arduino!\00")

	(func $_start (start)
		;; pinMode(LED_BUILTIN, OUTPUT)
		global.get $LED_BUILTIN
		global.get $OUTPUT
		call $pinMode

		;; println("hello arduino!")
		global.get $hello_ptr
		global.get $hello_len
		call $println
	)

	(func $loop (export "loop")
		;; digitalWrite(LED_BUILTIN, HIGH)
		global.get $LED_BUILTIN
		global.get $HIGH
		call $digitalWrite

		;; delay(100)
		i32.const 100
		call $delay

		;; digitalWrite(LED_BUILTIN, LOW)
		global.get $LED_BUILTIN
		global.get $LOW
		call $digitalWrite

		;; delay(900)
		i32.const 900
		call $delay	
	)
)
