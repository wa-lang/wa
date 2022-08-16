// 版权 @2012 凹语言 作者。保留所有权利。

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wasmerio/wasmer-go/wasmer"
)

var flagFile = flag.String("file", "a.out.wasm", "wasm file")

func main() {
	flag.Parse()
	wasmBytes, err := os.ReadFile(*flagFile)
	if err != nil {
		fmt.Printf("read %q failed: %v", *flagFile, err)
		os.Exit(1)
	}
	if err = Run(wasmBytes); err != nil {
		fmt.Printf("run %q failed: %v", *flagFile, err)
		os.Exit(1)
	}
}

func Run(wasmBytes []byte) error {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		return err
	}

	waBuiltinWrite := wasmer.NewFunction(store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			if len(args) > 0 {
				fmt.Println(args[0].I32())
			}
			return nil, nil
		},
	)

	importObject := wasmer.NewImportObject()
	importObject.Register("env", map[string]wasmer.IntoExtern{
		"__wa_builtin_print_int32": waBuiltinWrite,
	})

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return err
	}

	_start_fn, err := instance.Exports.GetFunction("_start")
	if err != nil {
		return err
	}

	if _, err = _start_fn(); err != nil {
		return err
	}

	return nil
}
