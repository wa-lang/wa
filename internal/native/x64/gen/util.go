// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import "strings"

// 判断是否为2个十六进制字符
// 可能有 +foo 格式的后缀
func isHex(s string) bool {
	if i := strings.Index(s, "+"); i >= 0 {
		s = s[:i]
	}
	if len(s) != 2 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if '0' <= c && c <= '9' || 'A' <= c && c <= 'F' {
			continue
		}
		return false
	}
	return true
}

// 判断是否为用 /n 格式表示的斜杠数字, 一般是对应 [0,7] 之间的数字
func isSlashNum(s string) bool {
	return len(s) == 2 && s[0] == '/' && '0' <= s[1] && s[1] <= '7'
}
