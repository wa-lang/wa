package scanner_test

import (
	"fmt"
	"testing"

	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wat/scanner"
	"wa-lang.org/wa/internal/wat/token"
)

// 构建 wat 目标
func tBuildWat(t *testing.T, filename string, src interface{}) string {
	cfg := config.DefaultConfig()
	prog, err := loader.LoadProgramFile(cfg, filename, src)
	if err != nil || prog == nil {
		t.Fatal(err)
	}

	watOut, err := compiler_wat.New().Compile(prog)
	if err != nil {
		t.Fatal(err)
	}

	return watOut
}

func TestHello(t *testing.T) {
	wat := tBuildWat(t, "hello.wa", `func main { println(123) }`)

	var src = []byte(wat)
	var file = token.NewFile("hello.wa", len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.ILLEGAL {
			t.Fatalf("failed: %v: %s %q", file.Position(pos), tok, lit)
		}
	}
}

func ExampleScanner_hello() {
	var src = []byte(tHello)
	var file = token.NewFile("", len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.ILLEGAL {
			fmt.Printf("failed: %v: %s %q\n", file.Position(pos), tok, lit)
			return
		}
	}
	fmt.Println("ok")

	// output:
	// ok
}

const tHello = `(module $hello_wasi
    ;; type iov struct { iov_base, iov_len int32 }
    ;; func fd_write(fd int32, id *iov, iovs_len int32, nwritten *int32) (errno int32)
    (import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32 i32 i32 i32) (result i32)))

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
)`
