// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	EnableTrace_wat2x64  bool
	EnableTrace_wat2la   bool
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
	case "wat2x64":
		EnableTrace_wat2x64 = true
	case "wat2la":
		EnableTrace_wat2la = true
	case "*":
		EnableTrace_app = true
		EnableTrace_compiler = true
		EnableTrace_loader = true
	default:
		panic("unknown trace flag: " + parten)
	}
}
