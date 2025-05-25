package main

//#cgo CFLAGS: -DNDEBUG
//#include "wa-app.h"
import "C"

func main() {
	C.wasm_main()
}
