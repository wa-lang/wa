// 版权 @2022 凹语言 作者。保留所有权利。

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"testing/fstest"

	"github.com/wa-lang/wa/api"
)

func init() {
	*api.FlagEnableTrace_loader = true
}

func main() {
	output, err := api.RunVFS(&api.PkgVFS{
		App: fstest.MapFS{
			"wa.mod.json": &fstest.MapFile{
				Data: []byte(waModJson),
			},
			"src/main.wa": &fstest.MapFile{
				Data: []byte(code),
			},
		},
	}, "myapp")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(string(output))
}

const waModJson = `
{
	"name": "_examples/hello",
	"pkgpath": "myapp"
}
`

const code = `
fn main() {
	println("你好, 凹语言!")
}
`
