build:
	-@rm ./web/*.wasm
	#go run ../../../main.go fmt ./src/main.wa
	#go run ../../../main.go fmt ./src/canvas/canvas.wa

	go run ../../../main.go build -o ./web/snake.wasm .

run:
	make build
	cd web && python3 -m http.server

publish:
	make build
	-rm -rf ../../../docs/snake/
	mkdir -p ../../../docs/snake
	cp ./web/* ../../../docs/snake/
	cp ./output/snake.wasm ../../../docs/snake/snake.wasm

clean:
