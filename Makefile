# 版权 @2019 凹语言 作者。保留所有权利。

.PHONY: wa hello prime build-wasm ci-test-all clean

GOBIN = ./build/bin/wa

wa:
	go build -o $(GOBIN)
	@echo "Done building."
	@echo "Run \"$(GOBIN)\" to launch wa."

hello:
	go install
	cd waroot && go run ../main.go run hello.wa

prime:
	cd waroot && go run ../main.go run examples/prime

build-wasm:
	GOARCH=wasm GOOS=js go build -o wa.out.wasm ./main_wasm.go

ci-test-all:
	go install
	go test ./...

	@echo "== std test begin =="
	go run main.go test std
	@echo "== std ok =="

	go run main.go run ./waroot/hello.wa
	cd waroot && go run ../main.go run hello.wa

	make -C ./waroot/examples ci-test-all
	wa -v

build-dist:
	go run ./builder

clean:
	-rm a.out*
