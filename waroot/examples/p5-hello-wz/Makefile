# 版权 @2024 P5-hello 作者。保留所有权利。

dev:
	go run ../../../main.go build .

run:
	go run ../../../main.go run .

publish:
	go run ../../../main.go build  .
	-rm -rf ../../../docs/p5-hello-wz/
	mkdir -p ../../../docs/p5-hello-wz
	cp ./output/* ../../../docs/p5-hello-wz/

clean:
