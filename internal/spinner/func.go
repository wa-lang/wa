// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"fmt"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/constant"
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
		fn.Body.EmitStore(loc, recv, int(r.Pos()))
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
				loc := fn.Body.AddLocal("$"+name, typ, int(p.Pos()), p)
				fn.Body.EmitStore(loc, param, int(p.Pos()))
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
		return

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

	case *ast.ExprStmt:
		v := b.expr(s.X, block)
		block.EmitStore(nil, v, int(s.Pos()))

	case *ast.IncDecStmt:
		op := wire.ADD
		if s.Tok == token.DEC {
			op = wire.SUB
		}
		loc := b.location(s.X, 0, block)
		ov := block.NewLoad(loc, int(s.Pos()))
		nv := block.NewBiop(ov, b.constval(constant.MakeInt64(1), b.info.TypeOf(s.X), int(s.Pos())), op, int(s.Pos()))
		block.EmitStore(loc, nv, int(s.Pos()))

	case *ast.AssignStmt:
		switch s.Tok {
		case token.ASSIGN, token.DEFINE:
			b.assignStmt(s, block)

		default: // += 等操作符
			loc := b.location(s.Lhs[0], 0, block)
			x := block.NewLoad(loc, int(s.Pos()))
			y := b.expr(s.Rhs[0], block)
			var op wire.OpCode
			switch s.Tok {
			case token.ADD_ASSIGN: // +=
				op = wire.ADD
			case token.SUB_ASSIGN: // -=
				op = wire.SUB
			case token.MUL_ASSIGN: // *=
				op = wire.MUL
			case token.QUO_ASSIGN: // /=
				op = wire.QUO
			case token.REM_ASSIGN: // %=
				op = wire.REM
			case token.AND_ASSIGN: // &=
				op = wire.AND
			case token.OR_ASSIGN: // |=
				op = wire.OR
			case token.XOR_ASSIGN: // ^=
				op = wire.XOR
			case token.SHL_ASSIGN: // <<=
				op = wire.SHL
			case token.SHR_ASSIGN: // >>=
				op = wire.SHR
			case token.AND_NOT_ASSIGN: // &^=
				op = wire.ANDNOT
			default:
				panic(fmt.Sprintf("Unknown OpCode: %v", s.Tok))
			}
			nv := block.NewBiop(x, y, op, int(s.Pos()))
			block.EmitStore(loc, nv, int(s.Pos()))
		}

	case *ast.ReturnStmt:
		b.returnStmt(s, f, block)

	case *ast.IfStmt:
		b.ifStmt(s, f, block)

	case *ast.ForStmt:
		b.forStmt(s, f, block)

	case *ast.DeferStmt:
		panic("Todo")

	case *ast.BranchStmt:
		panic("Todo")

	case *ast.SwitchStmt:
		panic("Todo")

	case *ast.TypeSwitchStmt:
		panic("Todo")

	case *ast.RangeStmt:
		panic("Todo")

	default:
		panic(fmt.Sprintf("unexpected statement kind: %T", s))
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
	var locs []wire.Location
	for _, v := range spec.Names {
		var loc wire.Location = nil
		if !isBlankIdent(v) {
			loc = b.addLocalForIdent(v, block)
		}
		locs = append(locs, loc)
	}

	if len(spec.Values) > 0 {
		b.assignN(locs, spec.Values, int(spec.Pos()), block)
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
				loc = b.location(lh, 0, block) // 非逃逸
			}
			locs[i] = loc
		}
	}

	b.assignN(locs, s.Rhs, int(s.Pos()), block)
}

func (b *Builder) returnStmt(s *ast.ReturnStmt, f *Func, block *wire.Block) {
	var results []wire.Expr

	if len(s.Results) == 1 && f.sig.Results().Len() > 1 {
		// 返回元组
		v := b.exprN(s.Results[0], block)
		results = append(results, v)
	} else {
		// 1:1返回
		for _, r := range s.Results {
			v := b.expr(r, block)
			results = append(results, v)
		}
	}

	if f.namedResults != nil {
		// 函数包含具名返回值，具名返回值赋值
		block.EmitStoreN(f.namedResults, results, int(s.TokPos))
	}

	// Todo：Defer

	if f.namedResults != nil {
		// 重新装载具名返回值
		results = results[:0]
		for _, r := range f.namedResults {
			results = append(results, block.NewLoad(r, int(s.TokPos)))
		}
	}

	block.EmitReturn(results, int(s.TokPos))
}

// location 方法返回一个左值表达式的位置
// escaping 为逃逸等级，0为非逃逸，1为取地址但未逃逸，2为取地址且逃逸（分别对应Local、Stack、Heap）
// 如以下代码：
//   i := S{
//      m: 1
//   }
//   return &i.m
// 将导致 i 逃逸
func (b *Builder) location(e ast.Expr, escaping wire.LocationKind, block *wire.Block) wire.Location {
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

// assign 将表达式 e 赋值给位置 loc 的动作降解并追加至 block 中
// loc 为 nil 是合法的，发生于向匿名变量 _ 赋值时
func (b *Builder) assign(loc wire.Location, e ast.Expr, pos int, block *wire.Block) {
	val := b.expr(e, block)
	block.EmitStore(loc, val, pos)
}

// assign 的多重赋值版本
func (b *Builder) assignN(locs []wire.Location, exprs []ast.Expr, pos int, block *wire.Block) {
	var vals []wire.Expr
	for _, e := range exprs {
		val := b.expr(e, block)
		vals = append(vals, val)
	}

	block.EmitStoreN(locs, vals, pos)
}

// expr 方法将表达式 e 降解为 wire 指令并追加至 block 中，返回 e 对应的 wire.Value
func (b *Builder) expr(e ast.Expr, block *wire.Block) wire.Expr {
	e = astutil.Unparen(e)

	tv := b.info.Types[e]

	// 常量？
	if tv.Value != nil {
		return b.constval(tv.Value, tv.Type, int(e.Pos()))
	}

	var v wire.Expr
	if tv.Addressable() {
		loc := b.location(e, 0, block)
		v = block.NewLoad(loc, int(e.Pos()))
	} else {
		v = b.expr1(e, tv, block)
	}

	return v
}

func (b *Builder) expr1(e ast.Expr, tv types.TypeAndValue, block *wire.Block) wire.Expr {
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
			loc := block.Lookup(obj, wire.LocationKindLocal)
			return block.NewLoad(loc, int(obj.Pos()))
		}

		// 函数
		panic("Todo")
		// :*ast.Ident

	case *ast.UnaryExpr: // 一元算符
		switch e.Op {
		case token.AND: // &x, 逃逸
			loc := b.location(e.X, wire.LocationKindHeap, block)
			if _, ok := unparen(e.X).(*ast.StarExpr); ok {
				// Todo: p 为空时，&*p 应panic
				panic("Todo")
			}
			return loc
		case token.ADD: // +x, 等价于 x
			return b.expr(e.X, block)
		case token.NOT: // !x
			x := b.expr(e.X, block)
			return block.NewUnop(x, wire.NOT, int(e.OpPos))
		case token.SUB: // -x
			x := b.expr(e.X, block)
			return block.NewUnop(x, wire.NEG, int(e.OpPos))
		case token.XOR: // ^x
			x := b.expr(e.X, block)
			return block.NewUnop(x, wire.XOR, int(e.OpPos))

		default:
			panic(e.Op)
		} // :*ast.UnaryExpr */

	case *ast.BinaryExpr: // 二元算符
		x := b.expr(e.X, block)
		y := b.expr(e.Y, block)
		switch e.Op {
		case token.XOR:
			return block.NewBiop(x, y, wire.XOR, int(e.OpPos))
		case token.LAND:
			return block.NewBiop(x, y, wire.LAND, int(e.OpPos))
		case token.LOR:
			return block.NewBiop(x, y, wire.LOR, int(e.OpPos))
		case token.SHL:
			return block.NewBiop(x, y, wire.SHL, int(e.OpPos))
		case token.SHR:
			return block.NewBiop(x, y, wire.SHR, int(e.OpPos))
		case token.ADD:
			return block.NewBiop(x, y, wire.ADD, int(e.OpPos))
		case token.SUB:
			return block.NewBiop(x, y, wire.SUB, int(e.OpPos))
		case token.MUL:
			return block.NewBiop(x, y, wire.MUL, int(e.OpPos))
		case token.QUO:
			return block.NewBiop(x, y, wire.QUO, int(e.OpPos))
		case token.REM:
			return block.NewBiop(x, y, wire.REM, int(e.OpPos))
		case token.AND:
			return block.NewBiop(x, y, wire.AND, int(e.OpPos))
		case token.OR:
			return block.NewBiop(x, y, wire.OR, int(e.OpPos))
		case token.AND_NOT:
			return block.NewBiop(x, y, wire.ANDNOT, int(e.OpPos))

		case token.EQL:
			return block.NewBiop(x, y, wire.EQL, int(e.OpPos))
		case token.NEQ:
			return block.NewBiop(x, y, wire.NEQ, int(e.OpPos))
		case token.GTR:
			return block.NewBiop(x, y, wire.GTR, int(e.OpPos))
		case token.LSS:
			return block.NewBiop(x, y, wire.LSS, int(e.OpPos))
		case token.GEQ:
			return block.NewBiop(x, y, wire.GEQ, int(e.OpPos))
		case token.LEQ:
			return block.NewBiop(x, y, wire.LEQ, int(e.OpPos))
		case token.SPACESHIP:
			return block.NewBiop(x, y, wire.LEG, int(e.OpPos))
		default:
			panic("illegal op in BinaryExpr: " + e.Op.String())
		} // :*ast.BinaryExpr */

	case *ast.CallExpr:
		if b.info.Types[e.Fun].IsType() {
			// 显式类型转换
			panic("Todo")
		}

		var callCommon wire.CallCommon

		if id, ok := unparen(e.Fun).(*ast.Ident); ok {
			switch obj := b.info.Uses[id].(type) {
			case *types.Builtin:
				// 内置函数调用，如 make、panic 等，runtime.SetFinalizer 也在此处理
				println(obj)
				panic("Todo")

			case *types.Func:
				callCommon.FnName = obj.FullName()
				callCommon.Sig = b.buildSig(obj.Type().(*types.Signature))
			}
		}

		if selector, ok := unparen(e.Fun).(*ast.SelectorExpr); ok {
			sel, ok := b.info.Selections[selector]
			if ok && sel.Kind() == types.MethodVal {
				obj := sel.Obj().(*types.Func)
				recvType := obj.Type().(*types.Signature).Recv().Type()
				if types.IsInterface(recvType) {
					// 接口调用
					panic("Todo")
				} else {
					// 对象方法调用
					panic("Todo")
				}
			} else {
				panic("")
			}
		}

		for _, v := range e.Args {
			param := b.expr(v, block)
			callCommon.Args = append(callCommon.Args, param)
		}
		callCommon.Pos = int(e.Pos())

		sc := wire.StaticCall{CallCommon: callCommon}
		return block.NewCall(&sc) //*/

	case *ast.FuncLit:
		// Todo
		panic("")

	case *ast.TypeAssertExpr:
		// Todo
		panic("")

	case *ast.SliceExpr:
		// Todo
		panic("")

	case *ast.SelectorExpr:
		// Todo
		panic("")

	case *ast.IndexExpr:
		// Todo
		panic("")

	case *ast.CompositeLit:
		// Todo
		panic("")

	case *ast.StarExpr:
		// Todo
		panic("")
	}

	panic(fmt.Sprintf("unexpected expr: %T", e))
}

// exprN 方法将多返回值表达式 e 降解并追加至 block，返回的 wire.Value 是元组
// 除自定义的的多返回值函数外，"v, ok" 形式的类型断言、Map查找等也属于多返回值表达式
func (b *Builder) exprN(e ast.Expr, block *wire.Block) wire.Expr {
	switch e := e.(type) {
	case *ast.ParenExpr:
		return b.exprN(e.X, block)
	}
	panic(fmt.Sprintf("exprN(%T)", e))
}

func (b *Builder) buildSig(s *types.Signature) (d wire.FnSig) {
	for i := 0; i < s.Params().Len(); i++ {
		t := b.BuildType(s.Params().At(i).Type())
		d.Params = append(d.Params, t)
	}
	d.Results = b.BuildType(s.Results())
	return
}

func (b *Builder) ifStmt(s *ast.IfStmt, f *Func, block *wire.Block) {
	if s.Init != nil {
		block = block.EmitBlock("if.begin", int(s.Pos()))
		b.stmt(s.Init, f, block)
	}

	cond := b.expr(s.Cond, block)
	i := block.EmitIf(cond, int(s.Cond.Pos()))

	if s.Body != nil {
		b.blockStmt(s.Body.List, f, i.True)
	}

	if s.Else != nil {
		if bs, ok := s.Else.(*ast.BlockStmt); ok {
			b.blockStmt(bs.List, f, i.False)
		} else {
			b.stmt(s.Else, f, i.False)
		}
	}

}

func (b *Builder) forStmt(s *ast.ForStmt, f *Func, block *wire.Block) {
	block = block.EmitBlock("for.begin", int(s.Pos()))
	if s.Init != nil {
		b.stmt(s.Init, f, block)
	}

	cond := b.expr(s.Cond, block)
	i := block.EmitLoop(cond, "", int(s.Cond.Pos()))

	if s.Body != nil {
		b.blockStmt(s.Body.List, f, i.Body)
	}

	if s.Post != nil {
		if bs, ok := s.Post.(*ast.BlockStmt); ok {
			b.blockStmt(bs.List, f, i.Post)
		} else {
			b.stmt(s.Post, f, i.Post)
		}
	}
}
