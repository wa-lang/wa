// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import "testing"

func Test_ZhAnames_重名(t *testing.T) {
	mZh := make(map[string]bool)
	for _, s := range _ZhAnames {
		if s != "" {
			if mZh[s] {
				t.Fatalf("%s 重名", s)
			}
			mZh[s] = true
		}
	}
}
