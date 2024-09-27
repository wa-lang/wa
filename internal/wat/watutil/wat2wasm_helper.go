// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/leb128"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2wasmWorker) lookupTokenOpcode(tok token.Token) wasm.Opcode {
	return tokOpcodeMap[tok]
}

func (p *wat2wasmWorker) findLabelIndex(label string) wasm.Index {
	if idx, err := strconv.Atoi(label); err == nil {
		return wasm.Index(idx)
	}
	for i := 0; i < len(p.labelScope); i++ {
		if s := p.labelScope[len(p.labelScope)-i-1]; s == label {
			return wasm.Index(i)
		}
	}
	panic(fmt.Sprintf("wat2wasm: unknown label %q", label))
}
func (p *wat2wasmWorker) enterLabelScope(label string) {
	p.labelScope = append(p.labelScope, label)
}
func (p *wat2wasmWorker) leaveLabelScope() {
	p.labelScope = p.labelScope[:len(p.labelScope)-1]
}

func (p *wat2wasmWorker) findFuncIndex(ident string) wasm.Index {
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
	panic(fmt.Sprintf("wat2wasm: unknown func %q", ident))
}
func (p *wat2wasmWorker) findFuncLocalIndex(fn *ast.Func, ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}

	for i, x := range fn.Type.Params {
		if x.Name == ident {
			return wasm.Index(i)
		}
	}
	for i, x := range fn.Body.Locals {
		if x.Name == ident {
			return wasm.Index(len(fn.Type.Params) + i)
		}
	}
	panic(fmt.Sprintf("wat2wasm: unknown func local %q", ident))
}

func (p *wat2wasmWorker) findTableIndex(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}
	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.TABLE {
			if x.TableName == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	if p.mWat.Table.Name == ident {
		return wasm.Index(importCount)
	}

	panic(fmt.Sprintf("wat2wasm: unknown table %q", ident))
}
func (p *wat2wasmWorker) findMemoryIndex(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}

	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.MEMORY {
			if x.Memory.Name == ident {
				return wasm.Index(importCount)
			}
			importCount++
		}
	}
	if p.mWat.Memory.Name == ident {
		return wasm.Index(importCount)
	}
	panic(fmt.Sprintf("wat2wasm: unknown memory %q", ident))
}

func (p *wat2wasmWorker) findGlobalIndex(ident string) wasm.Index {
	if idx, err := strconv.Atoi(ident); err == nil {
		return wasm.Index(idx)
	}

	var importCount int
	for _, x := range p.mWat.Imports {
		if x.ObjKind == token.GLOBAL {
			if x.GlobalName == ident {
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

	panic(fmt.Sprintf("wat2wasm: unknown global %q", ident))
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

func (p *wat2wasmWorker) encodeUint32(i uint32) []byte {
	return leb128.EncodeUint32(i)
}

func (p *wat2wasmWorker) encodeUint64(i uint64) []byte {
	return leb128.EncodeUint64(i)
}

func (p *wat2wasmWorker) encodeFloat32(i float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(i))
	return b
}

func (p *wat2wasmWorker) encodeFloat64(i float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(i))
	return b
}
func (p *wat2wasmWorker) encodeAlign(align uint) byte {
	switch align {
	case 1:
		return 0
	case 2:
		return 1
	case 4:
		return 2
	case 8:
		return 3
	case 16:
		return 4
	}
	panic("unreachable")
}
