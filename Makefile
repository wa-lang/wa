# 版权 @2019 凹语言 作者。保留所有权利。

hello:
	CGO_ENABLED=0 go run main.go run _examples/hello

hi:
	CGO_ENABLED=0 go run main.go run _examples/hi

prime:
	CGO_ENABLED=0 go run main.go run _examples/prime

wasm:
	CGO_ENABLED=0 go run main.go ssa _examples/hi
	CGO_ENABLED=0 go run main.go wasm _examples/hi

	wasm2wat a.out.wasm -o a.out.wast
	cd ./tools/wa-wasmer-run && go run main.go -file=../../a.out.wasm

clean:
