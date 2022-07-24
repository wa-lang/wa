// 版权 @2021 凹语言 作者。保留所有权利。

package compiler

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/ssa"
)

type Compiler struct {
	prog *loader.Program

	module *llir.Module

	curFunc         *llir.Func
	curLocals       map[ssa.Value]llvalue.Value
	curBlockEntries map[*ssa.BasicBlock]*llir.Block // a *ssa.BasicBlock may be split up
	phis            map[*ssa.Phi]*llir.InstPhi

	llTypeNameMap  map[lltypes.Type]string
	llNamedTypeMap map[string]lltypes.Type

	llFuncNameMap  map[*llir.Func]string
	llNamedFuncMap map[string]*llir.Func

	compilePkgDoneMap map[*ssa.Package]bool
}

func New() *Compiler {
	return &Compiler{}
}

func (p *Compiler) Compile(prog *loader.Program) (output string, err error) {
	p.prog = prog
	p.module = llir.NewModule()

	p.llTypeNameMap = make(map[lltypes.Type]string)
	p.llNamedTypeMap = make(map[string]lltypes.Type)

	p.llFuncNameMap = make(map[*llir.Func]string)
	p.llNamedFuncMap = make(map[string]*llir.Func)

	p.compilePkgDoneMap = make(map[*ssa.Package]bool)

	if !p.prog.Cfg.Debug {
		defer func() {
			//if r := recover(); r != nil {
			//	err = fmt.Errorf("%v", r)
			//}
		}()
	}

	// 定义 builtin 类型和函数
	if err := p.defineBuiltin(); err != nil {
		return "", err
	}

	// 编译全部包
	for _, pkg := range prog.SSAProgram.AllPackages() {
		if err := p.compilePackage(pkg); err != nil {
			return "", err
		}
	}
	// main 包入口
	{
		main_init := p.getFunc("main.init")
		main_main := p.getFunc("main.main")

		mainFunc := p.module.NewFunc("main", lltypes.I32)
		mainFuncBlock := mainFunc.NewBlock("entry")
		mainFuncBlock.NewCall(main_init)
		mainFuncBlock.NewCall(main_main)
		mainFuncBlock.NewRet(llconstant.NewInt(lltypes.I32, 0))
	}

	// 输出 LLIR 文件
	return p.module.String(), nil
}

func (p *Compiler) compilePackage(pkg *ssa.Package) (err error) {
	// 已经编译过
	if p.compilePkgDoneMap[pkg] {
		return nil
	}

	// 成功编译后记录状态
	defer func() {
		if err == nil {
			p.compilePkgDoneMap[pkg] = true
		}
	}()

	// 先编译导入的包
	for _, importPkg := range pkg.Pkg.Imports() {
		ssaPkg := p.prog.Pkgs[importPkg.Path()].SSAPkg
		if err := p.compilePackage(ssaPkg); err != nil {
			return err
		}
	}

	// 定义类型
	for _, member := range pkg.Members {
		if typ, ok := member.(*ssa.Type); ok {
			if err := p.compileType(pkg, typ); err != nil {
				return err
			}
		}
	}

	// 定义常量
	for _, member := range pkg.Members {
		if namedConst, ok := member.(*ssa.NamedConst); ok {
			if err := p.compileNamedConst(pkg, namedConst); err != nil {
				return err
			}
		}
	}

	// 定义全局变量
	for _, member := range pkg.Members {
		if g, ok := member.(*ssa.Global); ok {
			if err := p.compileGlobal(pkg, g); err != nil {
				return err
			}
		}
	}

	// 声明函数
	for _, member := range pkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			if err := p.compileFuncDeclare(pkg, fn); err != nil {
				return err
			}
		}
	}

	// 定义函数
	for _, member := range pkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			if err := p.compileFunc(fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Compiler) compileNamedConst(pkg *ssa.Package, namedConst *ssa.NamedConst) error {
	return nil
}
func (p *Compiler) compileGlobal(pkg *ssa.Package, g *ssa.Global) error {
	p.defGlobal(g)
	return nil
}

func (p *Compiler) compileFuncDeclare(pkg *ssa.Package, fn *ssa.Function) error {
	logger.Tracef(&config.EnableTrace_compiler, "pkgpath:%s, name=%s", fn.Pkg.Pkg.Path(), fn.Pkg.Pkg.Name())

	var mangledName string
	if pkgName := fn.Pkg.Pkg.Name(); pkgName == "main" || pkgName == "" {
		mangledName = fmt.Sprintf("%s.%s", "main", fn.Name())
	} else {
		mangledName = fmt.Sprintf("%s.%s", pkg.Pkg.Path(), fn.Name())
	}

	p.curFunc = p.getFunc(mangledName)
	if p.curFunc == nil {
		fnType := p.toLLFuncType(fn.Signature)
		fnParams := []*llir.Param{}
		for i, typ := range fnType.Params {
			fnParams = append(fnParams, llir.NewParam(fmt.Sprint(i), typ))
		}
		p.curFunc = p.module.NewFunc(mangledName, fnType.RetType, fnParams...)
	}
	return nil
}

func (p *Compiler) compileFunc(fn *ssa.Function) error {
	logger.Tracef(&config.EnableTrace_compiler, "pkgpath=%s, name=%s", fn.Pkg.Pkg.Path(), fn.Pkg.Pkg.Name())

	var mangledName string
	if pkgName := fn.Pkg.Pkg.Name(); pkgName == "main" || pkgName == "" {
		mangledName = fmt.Sprintf("%s.%s", "main", fn.Name())
	} else {
		mangledName = fmt.Sprintf("%s.%s", fn.Pkg.Pkg.Path(), fn.Name())
	}

	p.curFunc = p.getFunc(mangledName)
	if p.curFunc == nil {
		fnType := p.toLLFuncType(fn.Signature)
		fnParams := []*llir.Param{}
		for i, typ := range fnType.Params {
			fnParams = append(fnParams, llir.NewParam(fmt.Sprint(i), typ))
		}
		p.curFunc = p.module.NewFunc(mangledName, fnType.RetType, fnParams...)
	}

	p.curLocals = make(map[ssa.Value]llvalue.Value)
	p.curBlockEntries = make(map[*ssa.BasicBlock]*llir.Block)
	p.phis = make(map[*ssa.Phi]*llir.InstPhi)

	for _, block := range fn.Blocks {
		p.curBlockEntries[block] = p.curFunc.NewBlock(fmt.Sprintf("block_%04d", block.Index))
	}

	for _, block := range fn.Blocks {
		if err := p.compileBlock(block); err != nil {
			return err
		}
	}

	// fix phis
	for expr, phi := range p.phis {
		var incomings []*llir.Incoming
		for i, edge := range expr.Edges {
			x := p.getValue(edge)
			pred := p.curBlockEntries[expr.Block().Preds[i]]
			incomings = append(incomings, llir.NewIncoming(x, pred))
		}
		phi.Incs = incomings
	}

	return nil
}

func (p *Compiler) compileBlock(block *ssa.BasicBlock) error {
	llirBlock := p.curBlockEntries[block]
	if llirBlock == nil {
		panic("unreachable")
	}

	for _, ins := range block.Instrs {
		if err := p.compileInstruction(ins); err != nil {
			return err
		}
	}
	if llirBlock.Term == nil {
		llirBlock.Term = llir.NewUnreachable()
	}

	return nil
}
