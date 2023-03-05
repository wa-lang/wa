// 版权 @2022 凹语言 作者。保留所有权利。

package config

// 目标平台类型, 可管理后缀名
const (
	WaOS_Default = WaOS_Wasi // 默认

	WaOS_Arduino = "arduino" // Arduino 平台
	WaOS_Chrome  = "chrome"  // Chrome 浏览器
	WaOS_Wasi    = "wasi"    // WASI 接口
)

// 体系结构类型
const (
	WaArch_Wasm = "wasm" // wasm 平台
)
