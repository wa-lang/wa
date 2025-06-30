// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package lsp

import (
	"encoding/json"
)

func jsonMarshal(v interface{}) string {
	d, _ := json.Marshal(v)
	return string(d)
}

func jsonMarshalIndent(v interface{}) string {
	d, _ := json.MarshalIndent(v, "", "    ")
	return string(d)
}
