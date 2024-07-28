// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2wasmWorker) mustFindFuncTypeIndex(fnType *ast.FuncType) wasm.Index {
	typ := &wasm.FunctionType{}
	for _, x := range fnType.Params {
		typ.Params = append(typ.Params, p.buildValueType(x.Type))
	}
	for _, x := range fnType.Results {
		typ.Results = append(typ.Results, p.buildValueType(x))
	}
	for i, x := range p.mWasm.TypeSection {
		if typ.EqualsSignature(x.Params, x.Results) {
			return wasm.Index(i)
		}
	}
	panic("unreachable")
}

func (p *wat2wasmWorker) findBlockTypeIndex(results []token.Token) int32 {
	typ := &wasm.FunctionType{}
	for _, x := range results {
		typ.Results = append(typ.Results, p.buildValueType(x))
	}
	for i, x := range p.mWasm.TypeSection {
		if typ.EqualsSignature(x.Params, x.Results) {
			return int32(i)
		}
	}
	return 0
}

// 注册类型
// 相同的类型会被合并
func (p *wat2wasmWorker) registerFuncType(typ *wasm.FunctionType) {
	for _, x := range p.mWasm.TypeSection {
		if typ.EqualsSignature(x.Params, x.Results) {
			return
		}
	}
	p.mWasm.TypeSection = append(p.mWasm.TypeSection, typ)
}

// 构建函数类型
func (p *wat2wasmWorker) buildFuncType(in *ast.FuncType) error {
	t := &wasm.FunctionType{}
	for _, x := range in.Params {
		t.Params = append(t.Params, p.buildValueType(x.Type))
	}
	for _, x := range in.Results {
		t.Results = append(t.Results, p.buildValueType(x))
	}
	p.registerFuncType(t)
	return nil
}

// 内部的 block/loop/if 也有类型
func (p *wat2wasmWorker) buildFuncBodyTypes(insList []ast.Instruction) error {
	for _, insX := range insList {
		switch insX.Token() {
		case token.INS_BLOCK:
			ins := insX.(ast.Ins_Block)
			if len(ins.Results) > 1 {
				p.buildFuncType(&ast.FuncType{Results: ins.Results})
			}
			p.buildFuncBodyTypes(ins.List)

		case token.INS_LOOP:
			ins := insX.(ast.Ins_Loop)
			if len(ins.Results) > 1 {
				p.buildFuncType(&ast.FuncType{Results: ins.Results})
			}
			p.buildFuncBodyTypes(ins.List)

		case token.INS_IF:
			ins := insX.(ast.Ins_If)
			if len(ins.Results) > 1 {
				p.buildFuncType(&ast.FuncType{Results: ins.Results})
			}
			p.buildFuncBodyTypes(ins.Body)
			p.buildFuncBodyTypes(ins.Else)
		}
	}
	return nil
}
