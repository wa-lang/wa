// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"strings"
)

var (
	DebugMode = false

	EnableTrace_api      bool
	EnableTrace_app      bool
	EnableTrace_compiler bool
	EnableTrace_loader   bool
)

func SetDebugMode() {
	DebugMode = true
}

// 打开全部 trace 调试开关
func EnableTraceAll() {
	setEnableTrace("*")
}

// 打开 trace 调试开关
func SetEnableTrace(parten string) {
	for _, s := range strings.Split(strings.TrimSpace(parten), ",") {
		setEnableTrace(s)
	}
}

func setEnableTrace(parten string) {
	switch strings.ToLower(parten) {
	case "app":
		EnableTrace_app = true
	case "compiler":
		EnableTrace_compiler = true
	case "loader":
		EnableTrace_loader = true
	case "*":
		EnableTrace_app = true
		EnableTrace_compiler = true
		EnableTrace_loader = true
	default:
		panic("unknown trace flag: " + parten)
	}
}
