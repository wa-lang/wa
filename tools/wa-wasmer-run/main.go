// 版权 @2022 凹语言 作者。保留所有权利。

package main

import (
	"encoding/binary"
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

	var memory *wasmer.Memory

	// type iov struct { iov_base, iov_len int32 }
	// func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)

	wasiFdWrite := wasmer.NewFunction(store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			var argv []int
			for _, arg := range args {
				argv = append(argv, int(arg.I32()))
			}

			iov := argv[1]
			iov_base := binary.LittleEndian.Uint32(memory.Data()[iov:][:4])
			iov_len := binary.LittleEndian.Uint32(memory.Data()[iov+4:][:4])

			msg := string(memory.Data()[iov_base:][:iov_len])
			fmt.Print(msg)

			return []wasmer.Value{wasmer.NewI32(len(msg))}, nil
		},
	)
	waBuiltinWrite := wasmer.NewFunction(store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			if len(args) > 0 {
				fmt.Println(args[0].I32())
			}
			return nil, nil
		},
	)

	importObject := wasmer.NewImportObject()
	importObject.Register("wasi_snapshot_preview1", map[string]wasmer.IntoExtern{
		"fd_write": wasiFdWrite,
	})
	importObject.Register("env", map[string]wasmer.IntoExtern{
		"__wa_builtin_print_int32": waBuiltinWrite,
	})

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return err
	}

	memory, err = instance.Exports.GetMemory("memory")
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
