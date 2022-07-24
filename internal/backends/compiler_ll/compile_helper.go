// 版权 @2021 凹语言 作者。保留所有权利。

package compiler

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) getLLType(typName string) lltypes.Type {
	for _, typ := range p.module.TypeDefs {
		if typ.Name() == typName {
			return typ
		}
	}
	return nil
}

func (p *Compiler) getLLFunc(fnName string) *llir.Func {
	for _, x := range p.module.Funcs {
		if x.Name() == fnName {
			return x
		}
	}
	return nil
}

// 定义全局变量
func (p *Compiler) defGlobal(g *ssa.Global) {
	mangledName := fmt.Sprintf("%s.%s", g.Pkg.Pkg.Path(), g.Name())
	// 全局变量类似一个地址, 需要取其指向的元素类型
	typ := p.toLLType(g.Type().(*types.Pointer).Elem())
	glob := p.module.NewGlobal(mangledName, typ)
	glob.Init = llconstant.NewZeroInitializer(typ)
}

// 根据名字查询全局变量
func (p *Compiler) getGlobal(g *ssa.Global) *llir.Global {
	mangledName := fmt.Sprintf("%s.%s", g.Pkg.Pkg.Path(), g.Name())
	for _, x := range p.module.Globals {
		if x.Name() == mangledName {
			return x
		}
	}
	return nil
}

func (p *Compiler) getLocal(expr ssa.Value) llvalue.Value {
	x, ok := p.curLocals[expr]
	if !ok {
		logger.Debugln(expr, "not found")
	}
	return x
}

func (p *Compiler) defLocal(block *llir.Block, expr ssa.Value) llvalue.Value {
	if x, ok := p.curLocals[expr]; ok {
		return x
	}
	x := block.NewAlloca(p.toLLType(expr.Type()))
	p.curLocals[expr] = x
	return x
}

func (p *Compiler) getFunc(name string) *llir.Func {
	for _, fn := range p.module.Funcs {
		if fn.Name() == name {
			return fn
		}
	}
	return nil
}
func (p *Compiler) mustGetFunc(name string) *llir.Func {
	if fn := p.getFunc(name); fn != nil {
		return fn
	}
	panic(fmt.Sprintf("mustGetFunc: %s not found", name))
}

func (p *Compiler) getBlock(name string) *llir.Block {
	if p.curFunc == nil {
		return nil
	}
	for _, block := range p.curFunc.Blocks {
		if block.Name() == name {
			return block
		}
	}
	return nil
}
