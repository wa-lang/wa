# 版权 @2024 凹语言 作者。保留所有权利。

default:
	wa -v
	go version
	wasmer -V

	wa build -optimize -target=wasi -output=fib_wa.wasm fib_wa.wa
	GOOS=wasip1 GOARCH=wasm go build -o fib_go.wasm

	du -sh fib_*.wasm

	time wasmer fib_wa.wasm
	time wasmer fib_go.wasm

clean:
