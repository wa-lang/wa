// 版权 @2021 凹语言 作者。保留所有权利。

package compiler

import (
	"fmt"

	"github.com/wa-lang/wa/internal/ssa"
)

func (p *Compiler) compileInstruction(ins ssa.Instruction) error {
	switch ins := ins.(type) {
	default:
		return p.compileInstruction_value(ins)
	case *ssa.DebugRef:
		return p.compileInstruction_DebugRef(ins)
	case *ssa.Defer:
		return p.compileInstruction_Defer(ins)
	case *ssa.If:
		return p.compileInstruction_If(ins)
	case *ssa.Jump:
		return p.compileInstruction_Jump(ins)
	case *ssa.MapUpdate:
		return p.compileInstruction_MapUpdate(ins)
	case *ssa.Panic:
		return p.compileInstruction_Panic(ins)
	case *ssa.Return:
		return p.compileInstruction_Return(ins)
	case *ssa.RunDefers:
		return p.compileInstruction_RunDefers(ins)
	case *ssa.Store:
		return p.compileInstruction_Store(ins)
	}
}

func (p *Compiler) compileInstruction_value(ins ssa.Instruction) error {
	if ssaValue, ok := ins.(ssa.Value); ok {
		p.curLocals[ssaValue] = p.getValue(ssaValue)
		return nil
	}
	return fmt.Errorf("unknown Instruction(%T)", ins)
}

func (p *Compiler) compileInstruction_DebugRef(ins *ssa.DebugRef) error {
	panic("TODO: *ssa.DebugRef")
}

func (p *Compiler) compileInstruction_Defer(ins *ssa.Defer) error {
	panic("TODO: *ssa.Defer")
}

func (p *Compiler) compileInstruction_If(ins *ssa.If) error {
	llirBlock := p.curBlockEntries[ins.Block()]

	cond := p.getValue(ins.Cond)
	blockTrue := p.curBlockEntries[ins.Block().Succs[0]]
	blockElse := p.curBlockEntries[ins.Block().Succs[1]]
	llirBlock.NewCondBr(cond, blockTrue, blockElse)
	return nil
}

func (p *Compiler) compileInstruction_Jump(ins *ssa.Jump) error {
	llirBlock := p.curBlockEntries[ins.Block()]

	blockJump := p.curBlockEntries[ins.Block().Succs[0]]
	llirBlock.NewBr(blockJump)
	return nil
}

func (p *Compiler) compileInstruction_MapUpdate(ins *ssa.MapUpdate) error {
	panic("TODO: *ssa.MapUpdate")
}

func (p *Compiler) compileInstruction_Panic(ins *ssa.Panic) error {
	panic("TODO: *ssa.Panic")
}

func (p *Compiler) compileInstruction_Return(ins *ssa.Return) error {
	llirBlock := p.curBlockEntries[ins.Block()]

	switch {
	case len(ins.Results) == 0:
		llirBlock.NewRet(nil)
	case len(ins.Results) == 1:
		llirBlock.NewRet(p.getValue(ins.Results[0]))
	default:
		return fmt.Errorf("todo: Multiple return values")
	}
	return nil
}

func (p *Compiler) compileInstruction_RunDefers(ins *ssa.RunDefers) error {
	panic("TODO: *ssa.RunDefers")
}

func (p *Compiler) compileInstruction_Store(ins *ssa.Store) error {
	llirBlock := p.curBlockEntries[ins.Block()]

	dst := p.getValue(ins.Addr)
	src := p.getValue(ins.Val)

	llirBlock.NewStore(src, dst)
	return nil
}
