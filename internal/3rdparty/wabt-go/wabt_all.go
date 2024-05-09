// 版权 @2022 凹语言 作者。保留所有权利。

package wabt

import _ "embed"

const Wat2WasmName = "wa.wat2wasm.exe"

//go:embed internal/wabt-1.0.29-macos/bin/wat2wasm
var Wat2wasm_macos string

//go:embed internal/wabt-1.0.29-ubuntu/bin/wat2wasm
var Wat2wasm_ubuntu string

//go:embed internal/wabt-1.0.29-windows/bin/wat2wasm.exe
var Wat2wasm_windows string
