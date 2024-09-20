// 版权 @2024 凹语言 作者。保留所有权利。

// https://github.com/WebAssembly/wabt/blob/1.0.29/src/binary.h
// https://github.com/WebAssembly/extended-name-section/blob/main/proposals/extended-name-section/Overview.md

package watutil

import (
	"context"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/binary"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/parser"
	"wa-lang.org/wa/internal/wat/token"
)

func Wat2Wasm(filename string, source []byte) (wasmBytes []byte, err error) {
	defer func() {
		if err == nil {
			rt := wazero.NewRuntime(context.Background())
			m, errx := rt.CompileModule(context.Background(), wasmBytes)
			if errx == nil {
				m.Close(context.Background())
			} else {
				err = errx
			}
		}
	}()

	m, err := parser.ParseModule(filename, source)
	if err != nil {
		return nil, err
	}

	return newWat2wasmWorker(m).EncodeWasm(true)
}

type wat2wasmWorker struct {
	mWat *ast.Module

	mWasm *wasm.Module

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType

	labelScope []string // 嵌套的lebel查询, if/block/loop
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2wasmWorker(mWat *ast.Module) *wat2wasmWorker {
	return &wat2wasmWorker{mWat: mWat}
}

func (p *wat2wasmWorker) EncodeWasm(enableDebugNames bool) ([]byte, error) {
	names := &wasm.NameSection{
		ModuleName:    p.mWat.Name,
		FunctionNames: wasm.NameMap{},
		LocalNames:    wasm.IndirectNameMap{},
	}

	p.mWasm = &wasm.Module{NameSection: names}

	// ID: 1
	if err := p.buildTypeSection(); err != nil {
		return nil, err
	}

	// ID: 2
	if err := p.buildImportSection(); err != nil {
		return nil, err
	}

	// ID: 3
	if err := p.buildFunctionSection(); err != nil {
		return nil, err
	}

	// ID: 4
	if err := p.buildTableSection(); err != nil {
		return nil, err
	}

	// ID: 5
	if err := p.buildMemorySection(); err != nil {
		return nil, err
	}

	// ID: 6
	if err := p.buildGlobalSection(); err != nil {
		return nil, err
	}

	// ID: 7
	if err := p.buildExportSection(); err != nil {
		return nil, err
	}

	// ID: 8
	if err := p.buildStartSection(); err != nil {
		return nil, err
	}

	// ID: 9
	if err := p.buildElementSection(); err != nil {
		return nil, err
	}

	// ID: 10
	if err := p.buildCodeSection(); err != nil {
		return nil, err
	}

	// ID: 11
	if err := p.buildDataSection(); err != nil {
		return nil, err
	}

	// ID: 0
	if enableDebugNames {
		if err := p.buildNameSection(); err != nil {
			return nil, err
		}
	} else {
		p.mWasm.NameSection = nil
	}

	return binary.EncodeModule(p.mWasm), nil
}

// 1. 处理类型段
func (p *wat2wasmWorker) buildTypeSection() error {
	p.mWasm.TypeSection = []*wasm.FunctionType{}

	// Type段类型
	for _, x := range p.mWat.Types {
		if err := p.buildFuncType(x.Name, x.Type); err != nil {
			return err
		}
	}

	// 导入函数类型
	for _, spec := range p.mWat.Imports {
		if spec.ObjKind == token.FUNC {
			if err := p.buildFuncType(spec.FuncName, spec.FuncType); err != nil {
				return err
			}
		}
	}

	// 函数类型
	for _, fn := range p.mWat.Funcs {
		if err := p.buildFuncType(fn.Name, fn.Type); err != nil {
			return err
		}
		if err := p.buildFuncBodyTypes(fn.Name, fn.Body.Insts); err != nil {
			return err
		}
	}

	return nil
}

// 2. 处理导入
func (p *wat2wasmWorker) buildImportSection() error {
	p.mWasm.ImportSection = []*wasm.Import{}

	for _, x := range p.mWat.Imports {
		spec := &wasm.Import{
			Module: x.ObjModule,
			Name:   x.ObjName,
		}

		switch x.ObjKind {
		case token.FUNC:
			spec.Type = wasm.ExternTypeFunc
			spec.DescFunc = p.mustFindFuncTypeIndex(x.FuncType)
		case token.TABLE:
			spec.Type = wasm.ExternTypeTable
			spec.DescTable = &wasm.Table{}
			panic("import unsupport table")
		case token.MEMORY:
			spec.Type = wasm.ExternTypeMemory
			spec.DescMem = &wasm.Memory{
				Min: uint32(x.Memory.Pages),
				Max: uint32(x.Memory.MaxPages),
			}
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

	for _, fn := range p.mWat.Funcs {
		p.mWasm.FunctionSection = append(
			p.mWasm.FunctionSection, p.mustFindFuncTypeIndex(fn.Type),
		)
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
			idx := p.findFuncIndex(ident)
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
			p.buildInstruction(fnCode, fn, ins)
		}

		fnCode.Body = append(fnCode.Body, wasm.OpcodeEnd)

		p.mWasm.CodeSection = append(p.mWasm.CodeSection, fnCode)
	}

	return nil
}
func (p *wat2wasmWorker) buildNameSection() error {
	var funcNames wasm.NameMap
	var localNames wasm.IndirectNameMap
	var importFuncCount int

	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.FUNC {
			var localNameMap wasm.NameMap
			for j, local := range x.FuncType.Params {
				if local.Name != "" {
					localNameMap = append(localNameMap, &wasm.NameAssoc{
						Index: wasm.Index(j),
						Name:  local.Name,
					})
				}
			}

			funcNames = append(funcNames, &wasm.NameAssoc{
				Index: wasm.Index(importFuncCount),
				Name:  x.FuncName,
			})
			localNames = append(localNames, &wasm.NameMapAssoc{
				Index:   wasm.Index(importFuncCount),
				NameMap: localNameMap,
			})
			importFuncCount++
		}
	}
	for _, typ := range p.mWat.Types {
		var localNameMap wasm.NameMap
		for j, local := range typ.Type.Params {
			if local.Name != "" {
				localNameMap = append(localNameMap, &wasm.NameAssoc{
					Index: wasm.Index(j),
					Name:  local.Name,
				})
			}
		}
		localNames = append(localNames, &wasm.NameMapAssoc{
			Index:   wasm.Index(importFuncCount),
			NameMap: localNameMap,
		})
	}
	for i, fn := range p.mWat.Funcs {
		var localNameMap wasm.NameMap
		for j, local := range fn.Type.Params {
			if local.Name != "" {
				localNameMap = append(localNameMap, &wasm.NameAssoc{
					Index: wasm.Index(j),
					Name:  local.Name,
				})
			}
		}
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

	p.mWasm.NameSection.FunctionNames = funcNames
	p.mWasm.NameSection.LocalNames = localNames

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
			spec.Index = p.findFuncIndex(x.FuncIdx)
		case token.MEMORY:
			spec.Type = wasm.ExternTypeMemory
			spec.Index = p.findMemoryIndex(x.MemoryIdx)
		case token.TABLE:
			spec.Type = wasm.ExternTypeTable
			spec.Index = p.findTableIndex(x.TableIdx)
		case token.GLOBAL:
			spec.Type = wasm.ExternTypeGlobal
			spec.Index = p.findGlobalIndex(x.GlobalIdx)
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
