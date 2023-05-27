default:
	go run ../../main.go -target=arduino app.wa

build:
	go run ../../main.go build -target=arduino app.wa
	wat2wasm a.out.wat -o app.wasm
	xxd -i app.wasm > app.wasm.h

clean:
	-rm *.h *.wasm *.out*
