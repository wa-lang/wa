// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wasm

import (
	"bytes"

	"github.com/wa-lang/wa/internal/3rdparty/wasm/encoding"
	"github.com/wa-lang/wa/internal/3rdparty/wasm/instruction"
	"github.com/wa-lang/wa/internal/3rdparty/wasm/module"
	"github.com/wa-lang/wa/internal/3rdparty/wasm/types"
	"github.com/wa-lang/wa/internal/ast"
	"github.com/wa-lang/wa/internal/loader"
)

type Compiler struct {
	stages []func() error

	prog *loader.Program

	globalNames []string
	globalNodes map[string]*ast.Ident

	module   *module.Module
	mainBody *module.CodeEntry

	tinyMainIndex  uint32
	tinyReadIndex  uint32
	tinyWriteIndex uint32
}

func New() *Compiler {
	p := new(Compiler)
	p.stages = []func() error{
		p.reset,
		p.initModule,
		p.emitGlobals,
		p.emitMainBody,
	}
	return p
}

func (p *Compiler) Compile(prog *loader.Program) (output string, err error) {
	p.prog = prog

	for _, stage := range p.stages {
		if err := stage(); err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	if err := encoding.WriteModule(&buf, p.module); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *Compiler) reset() error {
	p.globalNames = []string{}
	p.globalNodes = make(map[string]*ast.Ident)
	return nil
}

func (p *Compiler) initModule() error {
	p.module = &module.Module{
		Names: module.NameSection{
			Module: "walang",
		},
		Memory: module.MemorySection{
			Memorys: []module.Memory{
				{InitPages: 1, MaxPages: 1},
			},
		},
		Export: module.ExportSection{
			Exports: []module.Export{
				{
					Name: "memory",
					Descriptor: module.ExportDescriptor{
						Type:  module.MemoryExportType,
						Index: 0,
					},
				},
			},
		},
	}

	// init types
	p.tinyReadIndex = 0
	p.tinyWriteIndex = 1
	p.tinyMainIndex = 2

	p.module.Type.Functions = []module.FunctionType{
		// func __tiny_read() int32
		{Results: []types.ValueType{types.I32}},
		// func __tiny_write(x int32)
		{Params: []types.ValueType{types.I32}},
		// func _start
		{},
	}

	// import
	p.module.Import.Imports = []module.Import{
		{
			Module: "env",
			Name:   "__tiny_read",
			Descriptor: module.FunctionImport{
				Func: p.tinyReadIndex,
			},
		},
		{
			Module: "env",
			Name:   "__tiny_write",
			Descriptor: module.FunctionImport{
				Func: p.tinyWriteIndex,
			},
		},
	}

	// _start func
	p.module.Function.TypeIndices = []uint32{
		p.tinyMainIndex,
	}
	p.module.Names.Functions = append(p.module.Names.Functions, module.NameMap{
		Index: p.tinyMainIndex,
		Name:  "_start",
	})

	// _start func body
	{
		var entry = &module.CodeEntry{
			Func: module.Function{
				Locals: []module.LocalDeclaration{},
				Expr: module.Expr{
					Instrs: []instruction.Instruction{
						instruction.I32Const{Value: 42},
						instruction.Call{Index: p.tinyWriteIndex},
						instruction.Return{},
					},
				},
			},
		}

		var buf bytes.Buffer
		if err := encoding.WriteCodeEntry(&buf, entry); err != nil {
			return err
		}

		p.module.Code.Segments = append(p.module.Code.Segments, module.RawCodeSegment{
			Code: buf.Bytes(),
		})
	}

	// export _start
	p.module.Export.Exports = append(p.module.Export.Exports, module.Export{
		Name: "_start",
		Descriptor: module.ExportDescriptor{
			Type:  module.FunctionExportType,
			Index: p.tinyMainIndex,
		},
	})

	return nil
}

func (p *Compiler) emitGlobals() error {
	if len(p.module.Global.Globals) > 0 {
		return nil
	}
	panic("TODO")
}

func (p *Compiler) emitMainBody() error {
	p.mainBody = &module.CodeEntry{
		Func: module.Function{
			Locals: []module.LocalDeclaration{},
			Expr: module.Expr{
				Instrs: []instruction.Instruction{
					// instruction.I32Const{Value: 42 + 1},
					// instruction.Call{Index: p.tinyWriteIndex},
				},
			},
		},
	}

	if err := p.compileNone(&p.mainBody.Func.Expr, nil); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := encoding.WriteCodeEntry(&buf, p.mainBody); err != nil {
		return err
	}

	// replace main body
	p.module.Code.Segments = []module.RawCodeSegment{
		{Code: buf.Bytes()},
	}

	return nil
}

func (p *Compiler) compileNone(ctx *module.Expr, node ast.Node) error {
	if node == nil {
		return nil
	}

	panic("TODO")
}

func (p *Compiler) globalIndexByName(name string) uint32 {
	for i, s := range p.globalNames {
		if s == name {
			return uint32(i) + 1
		}
	}
	return 0
}

func (p *Compiler) emit(ctx *module.Expr, x ...instruction.Instruction) {
	ctx.Instrs = append(ctx.Instrs, x...)
}
