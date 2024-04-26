// 版权 @2024 凹语言 作者。保留所有权利。

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
