# 版权 @2024 w4-rocket 作者。保留所有权利。

run:
	go run ../../../main.go run .

dev:
	go run ../../../main.go build .

publish:
	go run ../../../main.go build -target=wasm4 .
	-rm -rf ../../../docs/w4-rocket/
	mkdir -p ../../../docs/w4-rocket
	cp ./output/* ../../../docs/w4-rocket/

clean:
