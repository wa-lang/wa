# 版权 @2022 凹语言 作者。保留所有权利。

default:
	GOOS=js GOARCH=wasm go build -o wa.wasm
	mv wa.wasm ../../../docs

dev-chai:
	-rm wa.wasm
	GOOS=js GOARCH=wasm go build -o wa.wasm
	node run-go-wasm.js wa.wasm > a.out.wat

clean:
