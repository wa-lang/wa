// 版权 @2021 凹语言 作者。保留所有权利。

package builtin

import _ "embed"

//go:embed _builtin.ll
var llBuiltin string

//go:embed _builtin_wasm.ll
var llBuiltin_wasm string

func GetBuiltinLL(goos, goarch string) string {
	switch goos {
	case "wasm":
		return llBuiltin_wasm
	case "darwin":
	case "linux":
	case "windows":
	}
	return llBuiltin
}

const Header = `
declare i32 @ugo_builtin_println(i32)
declare i32 @ugo_builtin_exit(i32)
`

const MainMain = `
define i32 @main() {
	call i32() @ugo_main_init()
	call i32() @ugo_main_main()
	ret i32 0
}
`
