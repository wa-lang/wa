// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

func strEqualFold(a string, more ...string) bool {
	if len(more) == 0 {
		return true
	}
	for _, s := range more {
		if !singleStrEqualFold(a, s) {
			return false
		}
	}
	return true
}

func singleStrEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca == cb {
			continue
		}
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return false
}

func int2str(n int) string {
	if n == 0 {
		return "0"
	}

	buf := [20]byte{}
	i := len(buf)
	negative := n < 0

	// 处理负数
	u := uint(n)
	if negative {
		u = uint(-n)
	}

	// 从后往前填充
	for u > 0 {
		i--
		buf[i] = byte(u%10) + '0'
		u /= 10
	}

	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}

func strLastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func strHasSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
