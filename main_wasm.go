// 版权 @2022 凹语言 作者。保留所有权利。

//go:build wasm
// +build wasm

// 凹语言™ 命令行程序。
package main

import (
	"encoding/json"

	"github.com/wa-lang/wa/api"
)

func main() {
	println("hello wasm")
}

type WASMResut struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

//export walang_api_BuildCode
func walang_api_BuildCode(code string) string {
	var result WASMResut
	wat, err := api.BuildFile("hello.wa", code)
	if err != nil {
		result.Error = err.Error()
	}
	result.Result = string(wat)

	b, _ := json.Marshal(result)
	return string(b)
}
