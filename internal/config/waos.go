// 版权 @2022 凹语言 作者。保留所有权利。

package config

import (
	wasrc "wa-lang.org/wa/waroot/src"
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
	WaBackend_Default = wasrc.WaBackend_wat // 默认

	WaBackend_wat = wasrc.WaBackend_wat // 输出 wat
)

// 目标平台类型, 可管理后缀名
const (
	WaOS_Default = wasrc.WaOS_js // 默认

	WaOS_js      = wasrc.WaOS_js      // 浏览器 JS
	WaOS_wasi    = wasrc.WaOS_wasi    // WASI 接口
	WaOS_wasm4   = wasrc.WaOS_wasm4   // WASM4 接口
	WaOS_arduino = wasrc.WaOS_arduino // Arduino
	WaOS_unknown = wasrc.WaOS_unknown // Unknown
)

// 体系结构类型
const (
	WaArch_Default = wasrc.WaArch_wasm // 默认
	WaArch_wasm    = wasrc.WaArch_wasm // wasm 平台
)

// 后端列表
var WaBackend_List = wasrc.WaBackend_List

// OS 列表
var WaOS_List = wasrc.WaOS_List

// CPU 列表
var WaArch_List = wasrc.WaArch_List

// 检查 OS 值是否 OK
func CheckWaOS(os string) bool {
	return wasrc.CheckWaOS(os)
}
