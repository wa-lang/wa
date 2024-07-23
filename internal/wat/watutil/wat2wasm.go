// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/leb128"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/u64"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/binary"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/text"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/parser"
	"wa-lang.org/wa/internal/wat/token"
)

func Wat2Wasm(filename string, source []byte) ([]byte, error) {
	m, err := parser.ParseModule(filename, source)
	if err != nil {
		return nil, err
	}
	return newWat2wasmWorker(m).EncodeWasm()
}

type wat2wasmWorker struct {
	mWat *ast.Module

	mWasm *wasm.Module

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2wasmWorker(mWat *ast.Module) *wat2wasmWorker {
	return &wat2wasmWorker{mWat: mWat}
}

func (p *wat2wasmWorker) wat2wasm_bak(source []byte) ([]byte, error) {
	if m, err := text.DecodeModule(source, api.CoreFeaturesV2); err != nil {
		return nil, err
	} else {
		return binary.EncodeModule(m), nil
	}
}

func (p *wat2wasmWorker) EncodeWasm() ([]byte, error) {
	names := &wasm.NameSection{
		ModuleName: p.mWat.Name,
	}

	p.mWasm = &wasm.Module{NameSection: names}

	if err := p.buildTypeSection(); err != nil {
		return nil, err
	}
	if err := p.buildImportSection(); err != nil {
		return nil, err
	}
	if err := p.buildMemorySection(); err != nil {
		return nil, err
	}
	if err := p.buildTableSection(); err != nil {
		return nil, err
	}
	if err := p.buildGlobalSection(); err != nil {
		return nil, err
	}
	if err := p.buildFunctionSection(); err != nil {
		return nil, err
	}

	if err := p.buildDataSection(); err != nil {
		return nil, err
	}
	if err := p.buildElementSection(); err != nil {
		return nil, err
	}
	if err := p.buildCodeSection(); err != nil {
		return nil, err
	}
	if err := p.buildNameSection(); err != nil {
		return nil, err
	}
	if err := p.buildExportSection(); err != nil {
		return nil, err
	}
	if err := p.buildStartSection(); err != nil {
		return nil, err
	}

	if names.ModuleName == "" && names.FunctionNames == nil && names.LocalNames == nil {
		p.mWasm.NameSection = nil
	}

	return binary.EncodeModule(p.mWasm), nil
}

// 1. 处理类型段
func (p *wat2wasmWorker) buildTypeSection() error {
	p.mWasm.TypeSection = []*wasm.FunctionType{}

	// 导入函数类型
	for _, spec := range p.mWat.Imports {
		if spec.ObjKind == token.FUNC {
			typ := p.buildFuncType(spec.FuncType)
			p.mWasm.TypeSection = append(p.mWasm.TypeSection, typ)
		}
	}

	// Type段类型
	for _, x := range p.mWat.Types {
		typ := p.buildFuncType(x.Type)
		p.mWasm.TypeSection = append(p.mWasm.TypeSection, typ)
	}

	// 函数类型
	for _, x := range p.mWat.Funcs {
		typ := p.buildFuncType(x.Type)
		p.mWasm.TypeSection = append(p.mWasm.TypeSection, typ)
	}

	return nil
}

// 2. 处理导入
func (p *wat2wasmWorker) buildImportSection() error {
	p.mWasm.ImportSection = []*wasm.Import{}

	for i, x := range p.mWat.Imports {
		spec := &wasm.Import{
			Module: x.ObjModule,
			Name:   x.FuncName,
		}

		switch x.ObjKind {
		case token.FUNC:
			spec.Type = wasm.ExternTypeFunc
			spec.DescFunc = wasm.Index(i)
		case token.TABLE:
			spec.Type = wasm.ExternTypeTable
			spec.DescTable = &wasm.Table{}
			panic("import unsupport table")
		case token.MEMORY:
			spec.Type = wasm.ExternTypeMemory
			spec.DescMem = &wasm.Memory{}
			panic("import unsupport memory")
		case token.GLOBAL:
			spec.Type = wasm.ExternTypeGlobal
			spec.DescGlobal = &wasm.GlobalType{}
			panic("import unsupport global")
		default:
			panic("unreachable")
		}

		p.mWasm.ImportSection = append(p.mWasm.ImportSection, spec)
	}

	return nil
}

func (p *wat2wasmWorker) buildMemorySection() error {
	if p.mWat.Memory != nil {
		p.mWasm.MemorySection = &wasm.Memory{
			Min:          uint32(p.mWat.Memory.Pages),
			Max:          uint32(p.mWat.Memory.MaxPages),
			IsMaxEncoded: p.mWat.Memory.MaxPages > 0,
		}
	}

	return nil
}

func (p *wat2wasmWorker) buildTableSection() error {
	p.mWasm.TableSection = []*wasm.Table{}

	if p.mWat.Table != nil {
		if p.mWat.Table.Type != token.FUNCREF {
			panic("table only support funcref type")
		}

		tab := &wasm.Table{
			Min:  uint32(p.mWat.Table.Size),
			Type: wasm.RefTypeFuncref,
		}
		if n := uint32(p.mWat.Table.MaxSize); n != 0 {
			tab.Max = &n
		}

		p.mWasm.TableSection = append(p.mWasm.TableSection, tab)
	}

	return nil
}

func (p *wat2wasmWorker) buildGlobalSection() error {
	p.mWasm.GlobalSection = []*wasm.Global{}

	for _, g := range p.mWat.Globals {
		spec := &wasm.Global{
			Name: g.Name,
			Type: &wasm.GlobalType{
				ValType: p.buildValueType(g.Type),
				Mutable: g.Mutable,
			},
			Init: p.buildConstantExpression(g),
		}

		p.mWasm.GlobalSection = append(p.mWasm.GlobalSection, spec)
	}

	return nil
}

// 构建函数段
func (p *wat2wasmWorker) buildFunctionSection() error {
	p.mWasm.FunctionSection = []wasm.Index{}

	for i := range p.mWat.Funcs {
		p.mWasm.FunctionSection = append(p.mWasm.FunctionSection, wasm.Index(i))
	}

	return nil
}

func (p *wat2wasmWorker) buildDataSection() error {
	return nil // todo
}
func (p *wat2wasmWorker) buildElementSection() error {
	return nil // todo
}
func (p *wat2wasmWorker) buildCodeSection() error {
	return nil // todo
}
func (p *wat2wasmWorker) buildNameSection() error {
	return nil // todo
}
func (p *wat2wasmWorker) buildExportSection() error {
	return nil // todo
}
func (p *wat2wasmWorker) buildStartSection() error {
	return nil // todo
}

// 构建函数类型
func (p *wat2wasmWorker) buildFuncType(in *ast.FuncType) *wasm.FunctionType {
	t := &wasm.FunctionType{}
	for _, x := range in.Params {
		t.Params = append(t.Params, p.buildValueType(x.Type))
	}
	for _, x := range in.Results {
		t.Results = append(t.Results, p.buildValueType(x))
	}
	return t
}

// 构建值类型
func (p *wat2wasmWorker) buildValueType(x token.Token) wasm.ValueType {
	switch x {
	case token.I32:
		return wasm.ValueTypeI32
	case token.I64:
		return wasm.ValueTypeI64
	case token.F32:
		return wasm.ValueTypeF32
	case token.F64:
		return wasm.ValueTypeF64
	default:
		panic("unreachable")
	}
}

// 全局变量初始化指令
func (p *wat2wasmWorker) buildConstantExpression(g *ast.Global) *wasm.ConstantExpression {
	x := &wasm.ConstantExpression{}
	switch g.Type {
	case token.I32:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = leb128.EncodeInt32(g.I32Value)
	case token.I64:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = leb128.EncodeInt64(g.I64Value)
	case token.F32:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = u64.LeBytes(api.EncodeF32(g.F32Value))
	case token.F64:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = u64.LeBytes(api.EncodeF64(g.F64Value))
	default:
		panic("unreachable")
	}
	return x
}
