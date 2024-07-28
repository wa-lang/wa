// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

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

// 内部的 block/loop 也有类型
func (p *wat2wasmWorker) appendFuncBodyTypes(dst []*wasm.FunctionType, insList []ast.Instruction) []*wasm.FunctionType {
	for _, ins := range insList {
		switch ins.Token() {
		case token.INS_BLOCK:
			ins := ins.(ast.Ins_Block)
			dst = append(dst, p.buildFuncType(&ast.FuncType{Results: ins.Results}))
			dst = append(dst, p.appendFuncBodyTypes(dst, ins.List)...)
		case token.INS_LOOP:
			ins := ins.(ast.Ins_Loop)
			dst = append(dst, p.buildFuncType(&ast.FuncType{Results: ins.Results}))
			dst = append(dst, p.appendFuncBodyTypes(dst, ins.List)...)
		case token.INS_IF:
			ins := ins.(ast.Ins_If)
			dst = append(dst, p.buildFuncType(&ast.FuncType{Results: ins.Results}))
			dst = append(dst, p.appendFuncBodyTypes(dst, ins.Body)...)
			dst = append(dst, p.appendFuncBodyTypes(dst, ins.Else)...)
		}
	}
	return dst
}
