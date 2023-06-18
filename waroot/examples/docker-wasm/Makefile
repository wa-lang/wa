# 版权 @2019 凹语言 作者。保留所有权利。

# go install wa-lang.org/wa@dev-wasi
# WA:=wa

WA:=go run ../../../main.go

run:
	$(WA) run hello.wa

wasmer-run:
	make clean
	$(WA) build hello.wa
	wasmer a.out.wasm

docker-image-list:
	docker image list

docker-build:
	make clean
	$(WA) build hello.wa
	docker buildx build --platform wasi/wasm32 -t wa-lang/hello-world .

docker-run:
	docker run --rm \
		--name=wasm-hello \
		--runtime=io.containerd.wasmedge.v1 \
		--platform=wasi/wasm32 \
		docker.io/wa-lang/hello-world:latest

clean:
	-rm *.wasm a.out*
