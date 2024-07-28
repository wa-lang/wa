// 版权 @2024 凹语言 作者。保留所有权利。

package watutil_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"wa-lang.org/wa/internal/wat/watutil"
)

var tTestWat2Wasm_files = []string{
	"testdata/empty.wat",
	"testdata/empty_with_module_name.wat",
	"testdata/memory-01.wat",
	"testdata/memory-02.wat",
	"testdata/memory-03.wat",
	"testdata/memory-04.wat",
	"testdata/func-01.wat",
	"testdata/func-02.wat",
	"testdata/func-03.wat",
	"testdata/hello-01.wat",
	"testdata/hello-02.wat",
	"testdata/hello-03.wat",

	// todo:
	// "testdata/func-04.wat",
	// "testdata/type-01.wat",
}

var tTestWat2WasmWithOptions_noname_files = []string{
	"testdata/empty.wat",
	"testdata/empty_with_module_name.wat",
	"testdata/memory-01.wat",
	"testdata/memory-02.wat",
	"testdata/memory-03.wat",
	"testdata/memory-04.wat",
	"testdata/func-01.wat",
	"testdata/func-02.wat",
	"testdata/func-03.wat",
	"testdata/func-04.wat",
	"testdata/hello-01.wat",
	"testdata/hello-02.wat",
	"testdata/hello-03.wat",
	"testdata/type-01.wat",
	"testdata/type-02.wat",

	// debug:
	"testdata/label-01.wat",

	//"testdata/func-05.wat",

	// "testdata/type-03.wat",
}

func TestWat2Wasm(t *testing.T) {
	for i, name := range tTestWat2Wasm_files {
		wat, expect, _ := tLoadWatWasm(t, name)
		got, err := watutil.Wat2Wasm(name, wat)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		tCmpBytes(t, name, expect, got)
	}
}

func TestWat2WasmWithOptions_noname(t *testing.T) {
	for i, name := range tTestWat2WasmWithOptions_noname_files {
		wat, _, expect := tLoadWatWasm(t, name)
		got, err := watutil.Wat2WasmWithOptions(name, wat, watutil.Options{DisableDebugNames: true})
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		tCmpBytes(t, name, expect, got)
	}
}

func tCmpBytes(t *testing.T, name string, expect, got []byte) {
	var i int
	for ; i < len(expect) && i < len(got); i++ {
		if expect[i] != got[i] {
			break
		}
	}
	if i == len(expect) && i == len(got) {
		return
	}

	os.WriteFile("testdata/a.out.wasm", got, 0666)

	t.Fatalf("%s:\nexpect[%04X]: %s\n   got[%04X]: %s",
		name,
		i, tHexString(expect, i, 16),
		i, tHexString(got, i, 16),
	)
}

func tHexString(b []byte, off, size int) string {
	var sb strings.Builder
	for i := 0; i < size && off+i < len(b); i++ {
		sb.WriteString(fmt.Sprintf("%02X ", b[off+i]))
	}
	return sb.String()
}

func tLoadWatWasm(t *testing.T, name string) (watBytes, wasmBytes, wasmNonameBytes []byte) {
	var err error
	watBytes, err = os.ReadFile(name)
	if err != nil {
		t.Fatalf("os.ReadFile %s failed: %v", name, err)
	}
	wasmBytes, err = os.ReadFile(name + ".wasm")
	if err != nil {
		t.Fatalf("os.ReadFile %s failed: %v", name, err)
	}
	wasmNonameBytes, err = os.ReadFile(name + ".noname.wasm")
	if err != nil {
		t.Fatalf("os.ReadFile %s failed: %v", name, err)
	}
	return
}
