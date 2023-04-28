# 版权 @2023 凹语言 作者。保留所有权利。

dev:
	go run . -dbg=a1.out.dbg -wat=a1.out.wat

	-rm wa.wasm
	GOOS=js GOARCH=wasm go build -o wa.wasm
	node run-go-wasm.js wa.wasm -dbg=a2.out.dbg -wat=a2.out.wat

clean:
	-rm *.out.* *.wasm
