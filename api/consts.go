// 版权 @2023 凹语言 作者。保留所有权利。

package api

import (
	"wa-lang.org/wa/internal/config"
)

// wa 代码:
// xxx.wa
// xxx_$OS.wa
// xxx_$ARCH.wa
// xxx_$OS_$ARCH.wa

// wz 中文代码:
// xxx.wz
// xxx_$OS.wz
// xxx_$ARCH.wz
// xxx_$OS_$ARCH.wz

// ws 汇编代码:
// xxx.$BACKEND.ws
// xxx_$OS.$BACKEND.ws
// xxx_$ARCH.$BACKEND.ws
// xxx_$OS_$ARCH.$BACKEND.ws

// 编译器后端类型
const (
	WaBackend_Default = config.WaBackend_Default // 默认

	WaBackend_wat = config.WaBackend_wat // 输出 wat
)

// 目标平台类型, 可管理后缀名
const (
	WaOS_Default = config.WaOS_Default // 默认

	WaOS_js      = config.WaOS_js      // 浏览器 js
	WaOS_wasi    = config.WaOS_wasi    // WASI 接口
	WaOS_wasm4   = config.WaOS_wasm4   // WASM4 接口
	WaOS_unknown = config.WaOS_unknown // Unknown
)

// 体系结构类型
const (
	WaArch_Default = config.WaArch_Default // 默认
	WaArch_wasm    = config.WaArch_wasm    // wasm 平台
)
