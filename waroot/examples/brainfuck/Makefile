# 版权 @2022 凹语言 作者。保留所有权利。

web:
	go run ../../../main.go run -web .

console:
	go run ../../../main.go run .

publish:
	go run ../../../main.go build .
	-rm -rf ../../../docs/brainfuck/
	mkdir -p ../../../docs/brainfuck
	cp ./output/* ../../../docs/brainfuck/

clean:
