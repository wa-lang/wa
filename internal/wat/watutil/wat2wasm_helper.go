// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/leb128"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2wasmWorker) findFuncIdx(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}

	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.FUNC {
			if x.FuncName == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	for i, fn := range p.mWat.Funcs {
		if fn.Name == ident {
			return wasm.Index(importCount + i)
		}
	}
	return 0
}
func (p *wat2wasmWorker) findTableIdx(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}
	if !strings.HasPrefix(ident, "$") {
		panic("invalid ident:" + ident)
	}
	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.TABLE {
			if x.FuncName == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	if p.mWat.Table.Name == ident {
		return wasm.Index(importCount)
	}

	return 0
}
func (p *wat2wasmWorker) findMemoryIdx(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}

	if !strings.HasPrefix(ident, "$") {
		panic("invalid ident:" + ident)
	}
	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.MEMORY {
			if x.FuncName == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	if p.mWat.Memory.Name == ident {
		return wasm.Index(importCount)
	}
	return 0
}

func (p *wat2wasmWorker) findGlobalIdx(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}
	if !strings.HasPrefix(ident, "$") {
		panic("invalid ident:" + ident)
	}

	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.GLOBAL {
			if x.FuncName == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	for i, g := range p.mWat.Globals {
		if g.Name == ident {
			return wasm.Index(importCount + i)
		}
	}
	return 0
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
		x.Data = p.encodeInt32(g.I32Value)
	case token.I64:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = p.encodeInt64(g.I64Value)
	case token.F32:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = p.encodeFloat32(g.F32Value)
	case token.F64:
		x.Opcode = wasm.OpcodeI32Const
		x.Data = p.encodeFloat64(g.F64Value)
	default:
		panic("unreachable")
	}
	return x
}

func (p *wat2wasmWorker) encodeInt32(i int32) []byte {
	return leb128.EncodeInt32(i)
}

func (p *wat2wasmWorker) encodeInt64(i int64) []byte {
	return leb128.EncodeInt64(i)
}

func (p *wat2wasmWorker) encodeFloat32(i float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(i))
	return b
}

func (p *wat2wasmWorker) encodeFloat64(i float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(nil, math.Float64bits(i))
	return b
}
