// 版权 @2022 凹语言 作者。保留所有权利。

package wabt

import _ "embed"

//go:embed internal/wabt-1.0.29-ubuntu/bin/wat2wasm
var wat2wasm string
