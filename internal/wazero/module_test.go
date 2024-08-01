// 版权 @2023 凹语言 作者。保留所有权利。

package wazero_test

import (
	"strings"
	"testing"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/wat/watutil"
	"wa-lang.org/wa/internal/wazero"
)

func TestModule(t *testing.T) {
	const wasmName = "main.wa"
	const waCode = `
func main {
	println("hello wazero")
}

func Add(a:i32, b:i32) => i32 {
	return a+b
}
`
	wasmBytes := tBuildWasm(t, waCode)

	m, err := wazero.BuildModule(wasmName, wasmBytes)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	// main 执行
	stdout, _, err := m.RunMain("__main__.main")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(stdout)); got != "hello wazero" {
		t.Fatalf("expect %q, got %q", "hello wazero", got)
	}

	// add 函数执行
	result, _, _, err := m.RunFunc("__main__.Add", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0] != 3 {
		t.Fatalf("expect %q, got %q", []uint64{3}, result)
	}
}

func tBuildWasm(t *testing.T, waCode string) []byte {
	_, watBytes, err := api.BuildFile(api.DefaultConfig(), "main.wa", waCode)
	if err != nil {
		t.Fatal(err)
	}
	wasmBytes, err := watutil.Wat2Wasm("a.out.wat", watBytes)
	if err != nil {
		t.Fatal(err)
	}
	return wasmBytes
}
