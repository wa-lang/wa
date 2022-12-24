build:
	-@rm ./web/*.wasm
	go run ../../main.go fmt ./src/main.wa
	go run ../../main.go fmt ./src/canvas/canvas.wa

	go run ../../main.go build -o ./web/snake.wasm .

publish:
	make build
	-rm -rf ../../docs/snake/
	mkdir -p ../../docs/snake
	cp ./web/* ../../docs/snake/
	mv ../../docs/snake/snake.html ../../docs/snake/index.html

clean:
