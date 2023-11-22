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
		watCode: ``,
		errMsg:  `error: unexpected token "EOF", expected a module field or a module.`,
	},
	{
		watCode: `()`,
		errMsg:  `error: unexpected token ")", expected a module field or a module.`,
	},

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
		watCode: t_hello_wasi,
		errMsg:  "",
	},
	{
		watCode: tBuildWat(t_hello_wa),
		errMsg:  "",
	},
}

const t_hello_wasi = `
(module $hello_wasi
    ;; type iov struct { iov_base, iov_len int32 }
    ;; func fd_write(id *iov, iovs_len int32, nwritten *int32) (written int32)
    (import "wasi_unstable" "fd_write" (func $fd_write (param i32 i32 i32 i32) (result i32)))

    (memory 1)(export "memory" (memory 0))

    ;; 前 8 个字节保留给 iov 数组, 字符串从地址 8 开始
    (data (i32.const 8) "hello world\n")

    ;; _start 类似 main 函数, 自动执行
    (func $main (export "_start")
        (i32.store (i32.const 0) (i32.const 8))  ;; iov.iov_base - 字符串地址为 8
        (i32.store (i32.const 4) (i32.const 12)) ;; iov.iov_len  - 字符串长度

        (call $fd_write
            (i32.const 1)  ;; 1 对应 stdout
            (i32.const 0)  ;; *iovs - 前 8 个字节保留给 iov 数组
            (i32.const 1)  ;; len(iovs) - 只有1个字符串
            (i32.const 20) ;; nwritten - 指针, 里面是要写到数据长度
        )
        drop ;; 忽略返回值
    )
)
`

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
