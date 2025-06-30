// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appbase

import "testing"

func TestHasExt(t *testing.T) {
	var tests = []struct {
		path string
		exts []string
		ok   bool
	}{
		{"", nil, false},
		{"a", nil, true},

		{"a.wa", []string{".wa"}, true},
		{"a.wa", []string{".WA"}, true},
		{"a.Wa", []string{".wa"}, true},

		{"a.wa", []string{".wa", ".wz"}, true},

		{".wa", []string{".wa"}, false},
		{".wa", []string{}, true},
	}
	for i, tt := range tests {
		expect := tt.ok
		got := HasExt(tt.path, tt.exts...)
		if expect != got {
			t.Fatalf("%d: expect = %v, got = %v; // %v", i, expect, got, tt)
		}
	}
}
