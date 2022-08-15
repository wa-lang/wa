# 版权 @2019 凹语言 作者。保留所有权利。

hello:
	go run main.go run _examples/hello

hi:
	go run main.go run _examples/hi

prime:
	go run main.go run _examples/prime

wasm:
	go run main.go ssa _examples/hi
	go run main.go wasm _examples/hi
	wasm2wat a.out.wasm -o a.out.wast

clean:
