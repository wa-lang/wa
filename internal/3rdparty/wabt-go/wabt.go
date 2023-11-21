// 版权 @2022 凹语言 作者。保留所有权利。

// wabt 可执行程序包装
package wabt

import _ "embed"

// Wabt 版本号
const Version = "1.0.29"

// 读取 wat2wasm 命令
func LoadWat2Wasm() []byte {
	return []byte(wat2wasm)
}
