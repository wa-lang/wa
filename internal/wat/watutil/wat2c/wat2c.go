// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"bytes"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/parser"
	"wa-lang.org/wa/internal/wat/token"
)

const DebugMode = false

type Options struct {
	Prefix  string            // 输出名字的前缀
	Exports map[string]string // 导出函数, 可能改名
}

func Wat2C(filename string, source []byte, opt Options) (code, header []byte, err error) {
	m, err := parser.ParseModule(filename, source)
	if err != nil {
		return nil, nil, err
	}

	worker := newWat2cWorker(m, opt)
	code, header, err = worker.BuildCode()
	return
}

type wat2cWorker struct {
	opt Options

	m *ast.Module

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType

	localNames      []string      // 参数和局部变量名
	localTypes      []token.Token // 参数和局部变量类型
	scopeLabels     []string      // 嵌套的label查询, if/block/loop
	scopeStackBases []int         // if/block/loop, 开始的栈位置

	useMathX  bool // 是否使用了 math_x 部分函数
	use_R_u32 bool // R_u32
	use_R_u16 bool // R_u16
	use_R_u8  bool // R_u8

	trace bool // 调试开关
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2cWorker(mWat *ast.Module, opt Options) *wat2cWorker {
	p := &wat2cWorker{m: mWat, trace: DebugMode}

	p.opt.Prefix = toCName(opt.Prefix)

	if p.opt.Exports == nil {
		p.opt.Exports = map[string]string{}
	}
	for k, v := range opt.Exports {
		p.opt.Exports[k] = v
	}

	// 更新 export 信息
	for _, fn := range p.m.Funcs {
		if exportName, ok := p.opt.Exports[fn.Name]; ok {
			if exportName != "" {
				fn.ExportName = exportName
			} else {
				fn.ExportName = fn.Name
			}
		}
	}

	// 如果 start 字段为空, 则尝试用 _start 导出函数替代
	if p.m.Start == "" {
		for _, fn := range p.m.Funcs {
			if fn.ExportName == "_start" {
				p.m.Start = fn.Name
			}
		}
	}

	return p
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

	// 删除头文件中多个连续的空白行
	headerCode := h.Bytes()
	for {
		if !bytes.Contains(headerCode, []byte("\n\n\n")) {
			break
		}
		headerCode = bytes.ReplaceAll(headerCode, []byte("\n\n\n"), []byte("\n\n"))
	}

	// 删除空白行
	bodyCode := bytes.ReplaceAll(c.Bytes(), []byte("\n\n\n"), []byte("\n\n"))

	return bodyCode, headerCode, nil
}
