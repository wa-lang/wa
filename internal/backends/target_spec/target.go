// 版权 @2019 凹语言 作者。保留所有权利。

package target_spec

// 目标机器
type Machine string

const (
	Machine_Wasm32_wa   Machine = "wasm32-wa"   // 凹语言定义的 WASM 规范
	Machine_Wasm32_wasi Machine = "wasm32-wasi" // WASI 定义的 WASM 规范
)
