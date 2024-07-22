// 版权 @2024 凹语言 作者。保留所有权利。

// wat 格式的子集
// - 函数指令不支持折叠
// - 单指令之间不支持注释
// - 不支持多行注释

package parser

import (
	"errors"
	"os"

	"wa-lang.org/wa/internal/wat/ast"
)

var DebugMode = false

func ParseModule(path string, src interface{}) (f *ast.Module, err error) {
	data, err := readSource(path, src)
	if err != nil {
		return nil, err
	}

	p := newParser(path, data)
	p.trace = DebugMode

	f, err = p.ParseModule()
	return
}

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		}
		return nil, errors.New("invalid source")
	}

	return os.ReadFile(filename)
}
