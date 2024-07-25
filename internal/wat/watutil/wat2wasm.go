// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"strings"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/binary"
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

func (p *wat2wasmWorker) EncodeWasm() ([]byte, error) {
	names := &wasm.NameSection{
		ModuleName:    strings.TrimPrefix(p.mWat.Name, "$"),
		FunctionNames: wasm.NameMap{},
		LocalNames:    wasm.IndirectNameMap{},
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
	p.mWasm.DataSection = []*wasm.DataSegment{}

	for _, x := range p.mWat.Data {
		p.mWasm.DataSection = append(p.mWasm.DataSection, &wasm.DataSegment{
			OffsetExpression: &wasm.ConstantExpression{
				Opcode: wasm.OpcodeI32Const,
				Data:   p.encodeInt32(int32(x.Offset)),
			},
			Init: x.Value,
		})
	}

	return nil
}

func (p *wat2wasmWorker) buildElementSection() error {
	p.mWasm.ElementSection = []*wasm.ElementSegment{}

	for _, x := range p.mWat.Elem {
		initList := []*wasm.Index{}
		for _, ident := range x.Values {
			idx := p.findFuncIdx(ident)
			initList = append(initList, &idx)
		}

		p.mWasm.ElementSection = append(p.mWasm.ElementSection, &wasm.ElementSegment{
			Type:       wasm.RefTypeFuncref,
			TableIndex: 0,
			OffsetExpr: &wasm.ConstantExpression{
				Opcode: wasm.OpcodeI32Const,
				Data:   p.encodeInt32(int32(x.Offset)),
			},
			Init: initList,
		})
	}

	return nil
}
func (p *wat2wasmWorker) buildCodeSection() error {
	p.mWasm.CodeSection = []*wasm.Code{}

	for _, fn := range p.mWat.Funcs {
		fnCode := &wasm.Code{}
		for _, local := range fn.Body.Locals {
			fnCode.LocalTypes = append(fnCode.LocalTypes, p.buildValueType(local.Type))
		}
		for _, ins := range fn.Body.Insts {
			fnCode.Body = p.buildInstruction(fnCode.Body, ins)
		}
	}

	return nil
}
func (p *wat2wasmWorker) buildNameSection() error {
	p.mWasm.NameSection.FunctionNames = nil
	p.mWasm.NameSection.LocalNames = nil

	var funcNames wasm.NameMap
	var localNames wasm.IndirectNameMap
	var importFuncCount int

	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.FUNC {
			funcNames = append(funcNames, &wasm.NameAssoc{
				Index: wasm.Index(importFuncCount),
				Name:  x.FuncName,
			})
			importFuncCount++
		}
	}
	for i, fn := range p.mWat.Funcs {
		var localNameMap wasm.NameMap
		for j, local := range fn.Body.Locals {
			localNameMap = append(localNameMap, &wasm.NameAssoc{
				Index: wasm.Index(j),
				Name:  local.Name,
			})
		}

		funcNames = append(funcNames, &wasm.NameAssoc{
			Index: wasm.Index(importFuncCount + i),
			Name:  fn.Name,
		})

		localNames = append(localNames, &wasm.NameMapAssoc{
			Index:   wasm.Index(importFuncCount + i),
			NameMap: localNameMap,
		})
	}

	return nil
}

func (p *wat2wasmWorker) buildExportSection() error {
	p.mWasm.ExportSection = []*wasm.Export{}

	if len(p.mWat.Exports) == 0 {
		p.mWasm.ExportSection = nil
		return nil
	}

	for _, x := range p.mWat.Exports {
		spec := &wasm.Export{Name: x.Name}
		switch x.Kind {
		case token.FUNC:
			spec.Type = wasm.ExternTypeFunc
			spec.Index = p.findFuncIdx(x.FuncIdx)
		case token.MEMORY:
			spec.Type = wasm.ExternTypeMemory
			spec.Index = p.findMemoryIdx(x.FuncIdx)
		case token.TABLE:
			spec.Type = wasm.ExternTypeTable
			spec.Index = p.findTableIdx(x.FuncIdx)
		case token.GLOBAL:
			spec.Type = wasm.ExternTypeGlobal
			spec.Index = p.findGlobalIdx(x.FuncIdx)
		default:
			panic("unreachable")
		}

		p.mWasm.ExportSection = append(p.mWasm.ExportSection, spec)
	}

	return nil
}

func (p *wat2wasmWorker) buildStartSection() error {
	p.mWasm.StartSection = nil

	if p.mWat.Start == "" {
		return nil
	}

	var startIdx wasm.Index
	var startFound = false

	for _, spec := range p.mWat.Imports {
		if spec.ObjKind == token.FUNC {
			if spec.FuncName == p.mWat.Start {
				startFound = true
				break
			}
			startIdx++
		}
	}
	if !startFound {
		for _, fn := range p.mWat.Funcs {
			if fn.Name == p.mWat.Start {
				startFound = true
				break
			}
		}
	}

	p.mWasm.StartSection = &startIdx
	return nil
}
