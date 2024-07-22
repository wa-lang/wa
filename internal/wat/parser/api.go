// 版权 @2024 凹语言 作者。保留所有权利。

// wat 格式的子集
// - 函数指令不支持折叠
// - 单指令之间不支持注释

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
)

func ParseModule(path string, src []byte) (f *ast.Module, err error) {
	p := newParser(path, src)
	f, err = p.ParseModule()
	return
}
