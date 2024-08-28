# 版权 @2024 P5-hello 作者。保留所有权利。

dev:
	go run ../../../main.go build -target=wasm4 .

run:
	go run ../../../main.go run -target=wasm4 .

publish:
	go run ../../../main.go build -target=wasm4 .
	-rm -rf ../../../docs/w4-2048/
	mkdir -p ../../../docs/w4-2048
	cp ./output/* ../../../docs/w4-2048/

clean:
