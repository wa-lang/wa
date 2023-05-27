# 版权 @2019 凹语言 作者。保留所有权利。

hello:
	go run main.go run hello.wa

prime:
	go run main.go run _examples/prime

build-wasm:
	GOARCH=wasm GOOS=js go build -o wa.out.wasm ./main_wasm.go

win-exe-icon:
	windres -o main_rc_windows.syso main.rc

arduino-run:
	go run main.go -target=arduino arduino.wa

arduino-build:
	go run main.go build -target=arduino arduino.wa
	wat2wasm a.out.wat -o a.out.wasm
	xxd -i a.out.wasm > app.wasm.h

ci-test-all:
	go test ./...

	@echo "== std test begin =="
	go run main.go test fmt
	@echo "== std ok =="

	go run main.go hello.wa

	make -C ./_examples ci-test-all

clean:
