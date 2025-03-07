package wat2c

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) Tracef(foramt string, a ...interface{}) {
	if p.trace {
		fmt.Printf(foramt, a...)
	}
}

func (p *wat2cWorker) getHostFuncCRetType(fnType *ast.FuncType, fnName string) string {
	switch len(fnType.Results) {
	case 0:
		return "void"
	case 1:
		switch fnType.Results[0] {
		case token.I32:
			return "i32_t"
		case token.I64:
			return "i64_t"
		case token.F32:
			return "f32_t"
		case token.F64:
			return "f64_t"
		default:
			panic("unreachable")
		}
	default:
		return fmt.Sprintf("host_fn_%s_ret_t", toCName(fnName))
	}
}

func (p *wat2cWorker) getFuncCRetType(fnType *ast.FuncType, fnName string) string {
	switch len(fnType.Results) {
	case 0:
		return "void"
	case 1:
		switch fnType.Results[0] {
		case token.I32:
			return "i32_t"
		case token.I64:
			return "i64_t"
		case token.F32:
			return "f32_t"
		case token.F64:
			return "f64_t"
		default:
			panic("unreachable")
		}
	default:
		return fmt.Sprintf("fn_%s_ret_t", toCName(fnName))
	}
}

func (p *wat2cWorker) ifUseMathX(ins token.Token) bool {
	switch ins {
	case token.INS_I32_CLZ:
		return true
	case token.INS_I32_CTZ:
		return true
	case token.INS_I32_POPCNT:
		return true
	case token.INS_I32_ROTL:
		return true
	case token.INS_I32_ROTR:
		return true

	case token.INS_I64_CLZ:
		return true
	case token.INS_I64_CTZ:
		return true
	case token.INS_I64_POPCNT:
		return true
	case token.INS_I64_ROTL:
		return true
	case token.INS_I64_ROTR:
		return true
	}
	return false
}

func (p *wat2cWorker) findGlobalType(ident string) token.Token {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	// 不支持导入的全局变量
	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Globals) {
			panic(fmt.Sprintf("wat2c: unknown global %q", ident))
		}
		return p.m.Globals[idx].Type
	}
	for _, g := range p.m.Globals {
		if g.Name == ident {
			return g.Type
		}
	}
	panic("unreachable")
}

func (p *wat2cWorker) findLocalType(fn *ast.Func, ident string) token.Token {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2c: unknown local %q", ident))
		}
		return p.localTypes[idx]
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return p.localTypes[idx]
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return p.localTypes[len(fn.Type.Params)+idx]
		}
	}
	panic("unreachable")
}

func (p *wat2cWorker) findLocalName(fn *ast.Func, ident string) string {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2c: unknown local %q", ident))
		}
		return p.localNames[idx]
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return p.localNames[idx]
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return p.localNames[len(fn.Type.Params)+idx]
		}
	}
	panic("unreachable")
}

func (p *wat2cWorker) findType(ident string) *ast.FuncType {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Types) {
			panic(fmt.Sprintf("wat2c: unknown type %q", ident))
		}
		return p.m.Types[idx].Type
	}
	for _, x := range p.m.Types {
		if x.Name == ident {
			return x.Type
		}
	}
	panic(fmt.Sprintf("wat2c: unknown type %q", ident))
}

func (p *wat2cWorker) findFuncType(ident string) *ast.FuncType {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	idx := p.findFuncIndex(ident)
	if idx < len(p.m.Imports) {
		return p.m.Imports[idx].FuncType
	}

	return p.m.Funcs[idx-len(p.m.Imports)].Type
}

func (p *wat2cWorker) findFuncIndex(ident string) int {
	if ident == "" {
		panic("wat2c: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		return idx
	}

	var importCount int
	for _, x := range p.m.Imports {
		if x.ObjKind == token.FUNC {
			if x.FuncName == ident {
				return importCount
			}
			importCount++
		}
	}
	for i, fn := range p.m.Funcs {
		if fn.Name == ident {
			return importCount + i
		}
	}
	panic(fmt.Sprintf("wat2c: unknown func %q", ident))
}

func (p *wat2cWorker) findLabelName(label string) string {
	if label == "" {
		panic("wat2c: empty label")
	}

	idx := p.findLabelIndex(label)
	if idx < len(p.scopeLabels) {
		return p.scopeLabels[len(p.scopeLabels)-idx-1]
	}
	panic(fmt.Sprintf("wat2c: unknown label %q", label))
}

func (p *wat2cWorker) findLabelIndex(label string) int {
	if label == "" {
		panic("wat2c: empty label")
	}

	if idx, err := strconv.Atoi(label); err == nil {
		return idx
	}
	for i := 0; i < len(p.scopeLabels); i++ {
		if s := p.scopeLabels[len(p.scopeLabels)-i-1]; s == label {
			return i
		}
	}
	panic(fmt.Sprintf("wat2c: unknown label %q", label))
}

func (p *wat2cWorker) enterLabelScope(stkBase int, label string) {
	p.scopeLabels = append(p.scopeLabels, label)
	p.scopeStackBases = append(p.scopeStackBases, stkBase)
}
func (p *wat2cWorker) leaveLabelScope() {
	p.scopeLabels = p.scopeLabels[:len(p.scopeLabels)-1]
}
