// 版权 @2023 凹语言 作者。保留所有权利。

package wabt_test

import (
	"strings"
	"testing"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/wabt"
)

func TestWat2Wasm(t *testing.T) {
	for i, tt := range watTests {
		_, err := wabt.Wat2Wasm([]byte(tt.watCode))
		if !tMatchErrMsg(err, tt.errMsg) {
			t.Fatalf("%d: check failed: %v", i, err)
		}
	}
}

func BenchmarkWat2Wasm(b *testing.B) {
	wat := tBuildWat(t_hello_wa)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := wabt.Wat2Wasm([]byte(wat)); err != nil {
			b.Fatal(err)
		}
	}
}

func tMatchErrMsg(err error, errMsg string) bool {
	if errMsg == "" {
		return err == nil
	}
	return strings.Contains(err.Error(), errMsg)
}

func tBuildWat(waCode string) string {
	_, watBytes, err := api.BuildFile(api.DefaultConfig(), "main.wa", waCode)
	if err != nil {
		return err.Error()
	}
	return string(watBytes)
}

var watTests = []struct {
	watCode string
	errMsg  string
}{
	{
		watCode: `(module)`,
		errMsg:  "",
	},
	{
		watCode: `
		(module
			(func)
			(memory 1)
		)`,
		errMsg: "",
	},

	{
		watCode: tBuildWat(t_hello_wa),
		errMsg:  "",
	},
}

const t_hello_wa = `
// 版权 @2019 凹语言 作者。保留所有权利。

import "fmt"
import "runtime"

global year: i32 = 2023

func main {
	println("你好，凹语言！", runtime.WAOS)
	println(add(40, 2), year)

	fmt.Println("1+1 =", 1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
`
