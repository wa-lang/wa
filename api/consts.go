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

	WaBackend_clang = config.WaBackend_clang // 输出 c
	WaBackend_wat   = config.WaBackend_wat   // 输出 wat
)

// 目标平台类型, 可管理后缀名
const (
	WaOS_Default = config.WaOS_Default // 默认

	WaOS_arduino = config.WaOS_arduino // Arduino 平台
	WaOS_chrome  = config.WaOS_chrome  // Chrome 浏览器
	WaOS_js      = config.WaOS_js      // 浏览器 js
	WaOS_wasi    = config.WaOS_wasi    // WASI 接口
)

// 体系结构类型
const (
	WaArch_Default = config.WaArch_Default // 默认

	WaArch_386     = config.WaArch_386     // 386 平台
	WaArch_amd64   = config.WaArch_amd64   // amd64 平台
	WaArch_arm64   = config.WaArch_arm64   // arm64 平台
	WaArch_riscv64 = config.WaArch_riscv64 // riscv64 平台
	WaArch_wasm    = config.WaArch_wasm    // wasm 平台
)
