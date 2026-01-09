// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !amd64 && !loong64

package main

func wat2xxI32Add(a, b int32) int32 {
	return a + b
}
