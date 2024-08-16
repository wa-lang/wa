# 版权 @2019 凹语言 作者。保留所有权利。

.PHONY: wa hello prime build-wasm ci-test-all clean

GOBIN = ./build/bin/wa
DOCKER_VOLUME=-v $(shell pwd):/root

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

build-docker:
	go run ./builder
	docker build -t wa-lang/wa .

docker-run:
	docker run --platform linux/amd64 --rm -it ${DOCKER_VOLUME} wa-lang/wa

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

wasm-js:
	GOOS=js GOARCH=wasm go build -o wa.wasm
	mv wa.wasm ./docs
	cd ./docs && zip wa.wasm.zip wa.wasm

clean:
	-rm a.out*
