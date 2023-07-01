# 版权 @2019 凹语言 作者。保留所有权利。

hello:
	go install
	cd waroot && go run ../main.go run hello.wa

prime:
	cd waroot && go run ../main.go run examples/prime

build-wasm:
	GOARCH=wasm GOOS=js go build -o wa.out.wasm ./main_wasm.go

arduino-run:
	go run main.go -target=arduino arduino.wa

arduino-build:
	go run main.go build -target=arduino arduino.wa
	wat2wasm a.out.wat -o a.out.wasm
	xxd -i a.out.wasm > app.wasm.h

ci-test-all:
	go install
	go test ./...

	@echo "== std test begin =="
	go run main.go test std
	@echo "== std ok =="

	go run main.go ./waroot/hello.wa
	cd waroot && go run ../main.go hello.wa

	make -C ./waroot/examples ci-test-all

clean:
	-rm a.out*
