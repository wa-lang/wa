// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"bytes"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/parser"
	"wa-lang.org/wa/internal/wat/token"
)

func Wat2C(filename string, source []byte) (code, header []byte, err error) {
	m, err := parser.ParseModule(filename, source)
	if err != nil {
		return nil, nil, err
	}

	worker := newWat2cWorker(m)
	code, header, err = worker.BuildCode()
	return
}

type wat2cWorker struct {
	m *ast.Module

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType

	localNames []string      // 参数和局部变量名
	localTypes []token.Token // 参数和局部变量类型
	labelScope []string      // 嵌套的lebel查询, if/block/loop
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2cWorker(mWat *ast.Module) *wat2cWorker {
	return &wat2cWorker{m: mWat}
}

func (p *wat2cWorker) BuildCode() (code, header []byte, err error) {
	var h bytes.Buffer
	var c bytes.Buffer

	if err := p.buildCode(&c); err != nil {
		return nil, nil, err
	}
	if err := p.buildHeader(&h); err != nil {
		return nil, nil, err
	}

	return c.Bytes(), h.Bytes(), nil
}
