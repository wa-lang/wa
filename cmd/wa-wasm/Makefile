# 版权 @2024 凹语言 作者。保留所有权利。

# Go 1.21+

default:
	GOOS=wasip1 GOARCH=wasm go build -o wa.wasm

	# wasmer wa.wasm

list:
	GOOS=wasip1 GOARCH=wasm go list -json

clean:
	-rm wa.wasm
