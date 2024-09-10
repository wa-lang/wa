# 版权 @2024 P5-hello 作者。保留所有权利。

run:
	go run ../../../main.go run

build:
	go run ../../../main.go build

publish:
	go run ../../../main.go build -target=wasm4 .
	-rm -rf ../../../docs/w4-life/
	mkdir -p ../../../docs/w4-life
	cp ./output/* ../../../docs/w4-life/

clean:
