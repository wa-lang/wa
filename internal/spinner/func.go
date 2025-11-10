// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"fmt"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wire"
)

/**************************************
本文用于将函数转为 wire.Function
**************************************/

//-------------------------------------

/**************************************
Func: 创建 wire.Function 时的中间结构
**************************************/
type Func struct {
	sig          *types.Signature
	body         ast.Stmt
	namedResults []wire.Location
}

//-------------------------------------

// 将函数声明（ast.FuncDecl）转换为 wire.Function
func (b *Builder) BuildFuncDecl(decl *ast.FuncDecl) *wire.Function {
	fndecl := b.info.Defs[decl.Name]
	fn := Func{
		sig:  fndecl.Type().(*types.Signature),
		body: decl.Body,
	}

	return b.buildFunc(&fn)
}

// 将函数字面值（ast.FuncLit）转换为 wire.Function
func (b *Builder) BuildFuncLit(lit *ast.FuncLit) *wire.Function {
	fn := Func{
		sig:  b.info.TypeOf(lit).(*types.Signature),
		body: lit.Body,
	}

	return b.buildFunc(&fn)
}

// 从 Func 中创建 wire.Function
func (b *Builder) buildFunc(f *Func) (fn *wire.Function) {
	fn = b.module.NewFunction()

	if f.body == nil {
		if recv := f.sig.Recv(); recv != nil {
			rt := b.BuildType(recv.Type())
			r := b.module.NewParam(recv.Name(), rt, int(recv.Pos()))
			fn.Params = append(fn.Params, r)
		}

		if params := f.sig.Params(); params != nil {
			for i := 0; i < params.Len(); i++ {
				p := params.At(i)
				typ := b.BuildType(p.Type())

				name := p.Name()
				if len(name) == 0 {
					name = fmt.Sprintf("$arg%d", i)
				}
				param := b.module.NewParam(name, typ, int(p.Pos()))
				fn.Params = append(fn.Params, param)
			}
		}

		if results := f.sig.Results(); results != nil {
			for i := 0; i < f.sig.Results().Len(); i++ {
				r := f.sig.Results().At(i)
				typ := b.BuildType(r.Type())

				fn.Results = append(fn.Results, typ)
			}
		}

		return
	} // if f.body == nil

	fn.StartBody()

	if r := f.sig.Recv(); r != nil {
		typ := b.BuildType(r.Type())
		name := r.Name() // 接收器不会是匿名的
		recv := b.module.NewParam(name, typ, int(r.Pos()))
		fn.Params = append(fn.Params, recv)
		loc := fn.Body.AddLocal(name, typ, int(r.Pos()), r)
		fn.Body.EmitStore(loc, recv, int(recv.Pos()))
	}

	if params := f.sig.Params(); params != nil {
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)
			typ := b.BuildType(p.Type())

			name := p.Name()
			if len(name) == 0 {
				// 匿名参数
				name = fmt.Sprintf("$arg%d", i)
				param := b.module.NewParam(name, typ, int(p.Pos()))
				fn.Params = append(fn.Params, param)
			} else {
				// 具名参数
				param := b.module.NewParam(name, typ, int(p.Pos()))
				fn.Params = append(fn.Params, param)
				loc := fn.Body.AddLocal(name, typ, int(p.Pos()), p)
				fn.Body.EmitStore(loc, param, param.Pos())
			}
		}
	}

	if results := f.sig.Results(); results != nil {
		for i := 0; i < f.sig.Results().Len(); i++ {
			r := f.sig.Results().At(i)
			typ := b.BuildType(r.Type())
			fn.Results = append(fn.Results, typ)

			name := r.Name()
			if len(name) != 0 {
				// 分配具名返回值
				loc := fn.Body.AddLocal(name, typ, int(r.Pos()), r)
				f.namedResults = append(f.namedResults, loc)
			}
		}
	}

	b.stmt(f.body, f, fn.Body)

	fn.EndBody()
	return
}

// stmt 将 AST 语句降解为 wire 指令，追加至 wire.Block
func (b *Builder) stmt(s ast.Stmt, f *Func, block *wire.Block) {
	switch s := s.(type) {
	case *ast.EmptyStmt:
		//

	case *ast.BlockStmt:
		newblock := block.EmitBlock("", int(s.Pos()))
		b.blockStmt(s.List, f, newblock)

	case *ast.DeclStmt: // Con, Var or Typ
		d := s.Decl.(*ast.GenDecl)
		if d.Tok == token.VAR || d.Tok == token.GLOBAL {
			for _, spec := range d.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					b.localValueSpec(vs, block)
				}
			}
		}

	case *ast.AssignStmt:
		switch s.Tok {
		case token.ASSIGN, token.DEFINE:
			b.assignStmt(s, block)

		default: // += 等操作符
			panic("Todo")
		}

	case *ast.ReturnStmt:
		b.returnStmt(s, f, block)

	default:
		panic("Todo")
	}

}

// blockStmt 将一组 AST 语句降解为 wire 指令，追加至 wire.Block
func (b *Builder) blockStmt(list []ast.Stmt, f *Func, block *wire.Block) {
	for _, s := range list {
		b.stmt(s, f, block)
	}
}

// localValueSpec 将局部变量声明降解，追加至 wire.Block
func (b *Builder) localValueSpec(spec *ast.ValueSpec, block *wire.Block) {
	switch {
	case len(spec.Values) == len(spec.Names):
		// 1：1赋值，如：x, y := 1, 2
		for i, id := range spec.Names {
			var loc wire.Location = nil
			if !isBlankIdent(id) {
				loc = b.addLocalForIdent(id, block)
			}
			b.assign(loc, spec.Values[i], int(id.Pos()), block, nil)
		}

	case len(spec.Values) == 0:
		// 未指定初始值，默认 0 值初始化
		for _, id := range spec.Names {
			if !isBlankIdent(id) {
				b.addLocalForIdent(id, block)
			}
		}

	default:
		// 元组赋值给多个变量，如： x, y := swap(x, y)
		tuple := b.exprN(spec.Values[0], block)
		for i, id := range spec.Names {
			if !isBlankIdent(id) {
				loc := b.addLocalForIdent(id, block)
				val := block.EmitExtract(tuple, i, int(id.Pos()))
				block.EmitStore(loc, val, int(id.Pos()))
			}
		}
	}
}

// 在 block 中分配一个由 id 定义的局部变量
func (b *Builder) addLocalForIdent(id *ast.Ident, block *wire.Block) wire.Location {
	obj := b.info.Defs[id]
	typ := b.BuildType(obj.Type())
	return block.AddLocal(obj.Name(), typ, int(obj.Pos()), obj)
}

// 向 block 中添加赋值（ = , := ）操作
func (b *Builder) assignStmt(s *ast.AssignStmt, block *wire.Block) {
	isDef := false
	if s.Tok == token.DEFINE {
		isDef = true
	}

	locs := make([]wire.Location, len(s.Lhs))
	for i, lh := range s.Lhs {
		if !isBlankIdent(lh) {
			var loc wire.Location
			if isDef {
				loc = b.addLocalForIdent(lh.(*ast.Ident), block)
			} else {
				loc = b.location(lh, false, block) // 非逃逸
			}
			locs[i] = loc
		}
	}

	if len(s.Lhs) == len(s.Rhs) {
		// 左右部数量相等，注意右部表达式可能引用左部，因此赋值操作需先暂存
		var sb storebuf
		for i := range s.Lhs {
			b.assign(locs[i], s.Rhs[i], int(s.Pos()), block, &sb)
		}
		sb.emit(block)
	} else {
		// 将元组赋值给多个变量， 如 x, y = swap(x, y)
		tuple := b.exprN(s.Rhs[0], block)
		for i, loc := range locs {
			rh := block.EmitExtract(tuple, i, int(s.Pos()))
			block.EmitStore(loc, rh, int(s.Pos()))
		}
	}
}

func (b *Builder) returnStmt(s *ast.ReturnStmt, f *Func, block *wire.Block) {
	var results []wire.Value

	if len(s.Results) == 1 && f.sig.Results().Len() > 1 {
		// 返回元组
		panic("Todo")
	} else {
		// 1:1返回
		for _, r := range s.Results {
			v := b.expr(r, block)
			results = append(results, v)
		}
	}

	if f.namedResults != nil {
		// 函数包含具名返回值，具名返回值赋值
		for i, r := range results {
			block.EmitStore(f.namedResults[i], r, int(s.TokPos))
		}
	}

	// Todo：Defer

	if f.namedResults != nil {
		// 重新装载具名返回值
		results = results[:0]
		for _, r := range f.namedResults {
			results = append(results, block.EmitLoad(r, r.DataType(), int(s.TokPos)))
		}
	}

	block.EmitReturn(results, int(s.TokPos))
}

// location 方法返回一个左值表达式的位置
// 当逃逸标志 escaping 为 true 时，本方法会将左值表达式的基变量标记为逃逸，如以下代码：
//   i := S{
//      m: 1
//   }
//   return &i.m
// 将导致 i 逃逸
func (b *Builder) location(e ast.Expr, escaping bool, block *wire.Block) wire.Location {
	switch e := e.(type) {
	case *ast.Ident:
		if isBlankIdent(e) {
			return nil
		}

		obj := b.info.ObjectOf(e)
		v, ok := b.module.Globals[obj]
		if !ok {
			v = block.Lookup(obj, escaping)
		}

		return v
	}

	panic(fmt.Sprintf("unexpected address expression: %T", e))
}

type store struct {
	loc wire.Location
	val wire.Value
	pos int
}

type storebuf struct{ stores []store }

func (sb *storebuf) store(loc wire.Location, val wire.Value, pos int) {
	sb.stores = append(sb.stores, store{loc, val, pos})
}

func (sb *storebuf) emit(block *wire.Block) {
	for _, s := range sb.stores {
		block.EmitStore(s.loc, s.val, s.pos)
	}
}

// assign 方法将表达式 e 赋值给位置 loc 的动作降解并追加至 block 中
// loc 为 nil 是合法的，发生于向匿名变量 _ 赋值时
// 存储缓冲 sb 不为空，则 Store 动作将保存在 sb 中，而非追加至 block 中，
// 该情况出现在类似 x, y = inc(x), x+y 的多重赋值时，避免右值对左值交叉应用
func (b *Builder) assign(loc wire.Location, e ast.Expr, pos int, block *wire.Block, sb *storebuf) {
	val := b.expr(e, block)

	if loc == nil {
		return
	}

	if sb != nil {
		sb.store(loc, val, pos)
	} else {
		block.EmitStore(loc, val, pos)
	}
}

// expr 方法将表达式 e 降解为 wire 指令并追加至 block 中，返回 e 对应的 wire.Value
func (b *Builder) expr(e ast.Expr, block *wire.Block) wire.Value {
	e = astutil.Unparen(e)

	tv := b.info.Types[e]

	// 常量？
	if tv.Value != nil {
		return b.constval(tv.Value, tv.Type, int(e.Pos()))
	}

	var v wire.Value
	if tv.Addressable() {
		loc := b.location(e, false, block)
		typ := b.BuildType(deref(tv.Type))
		v = block.EmitLoad(loc, typ, int(e.Pos()))
	} else {
		v = b.expr0(e, tv, block)
	}

	return v
}

func (b *Builder) expr0(e ast.Expr, tv types.TypeAndValue, block *wire.Block) wire.Value {
	switch e := e.(type) {
	case *ast.Ident:
		obj := b.info.Uses[e]
		switch obj := obj.(type) {
		case *types.Builtin:
			// 内建函数
			panic("Todo" + obj.Name())

		case *types.Nil:
			// nil
			return b.nilConst(tv.Type, int(e.Pos()))
		}

		if _, ok := obj.(*types.Var); ok {
			loc := block.Lookup(obj, false)
			typ := b.BuildType(obj.Type())
			return block.EmitLoad(loc, typ, int(obj.Pos()))
		}

		// 函数
		panic("Todo")
	}

	panic(fmt.Sprintf("unexpected expr: %T", e))
}

// exprN 方法将多返回值表达式 e 降解并追加至 block，返回的 wire.Value 是元组
// 除自定义的的多返回值函数外，"v, ok" 形式的类型断言、Map查找等也属于多返回值表达式
func (b *Builder) exprN(e ast.Expr, block *wire.Block) wire.Value {
	switch e := e.(type) {
	case *ast.ParenExpr:
		return b.exprN(e.X, block)
	}
	panic(fmt.Sprintf("exprN(%T)", e))
}
