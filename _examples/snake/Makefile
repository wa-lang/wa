build:
	-@rm ./web/*.wasm
	go run ../../main.go build -o ./web/snake.wasm .

publish:
	make build
	-rm -rf ../../docs/snake/
	mkdir -p ../../docs/snake
	cp -r ./web/ ../../docs/snake/
	#mv ../../docs/snake/snake.html ../../docs/snake/index.html

clean:
