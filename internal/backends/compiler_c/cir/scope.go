package cir

import (
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
)

/**************************************
Scope: CIR的域结构
Scope可近似理解为由大括号对{}括起来的语句块
**************************************/
type Scope struct {
	//父域指针，Scope.NewScope()设置该值
	Parent *Scope

	//在该域中定义的结构体
	Structs []*StructDecl

	locals []*VarDecl

	//在该域中定义的函数
	Funcs []*FuncDecl

	//域中包含的语句
	Stmts []Stmt

	//在该域中定义的变量（含局部变量和临时变量，局部变量位于域头部结构体之后，临时变量可出现于域内任意位置）
	Vars []*VarDecl
}

//实现Stmt接口
func (s *Scope) isStmt() {}

func (s *Scope) Type() ctypes.Type {
	return ctypes.Void
}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent}
}

func (s *Scope) AddStructDecl(name string, m []Var) *StructDecl {
	decl := NewStructDecl(name, m)
	s.Structs = append(s.Structs, decl)
	return decl
}

func (s *Scope) AddTupleDecl(t ctypes.Tuple) *StructDecl {
	var decl StructDecl
	decl.Struct = t.GenStruct()
	s.Structs = append(s.Structs, &decl)
	return &decl
}

func (s *Scope) AddTempVarDecl(name string, t ctypes.Type) *VarDeclStmt {
	decl := NewVarDeclStmt(name, t)
	s.Stmts = append(s.Stmts, decl)
	s.Vars = append(s.Vars, &decl.VarDecl)
	return decl
}

func (s *Scope) AddLocalVarDecl(name string, t ctypes.Type) *VarDecl {
	decl := NewVarDecl(name, t)
	s.locals = append(s.locals, decl)
	s.Vars = append(s.Vars, decl)
	return decl
}

func (s *Scope) AddFuncDecl(name string, result ctypes.Type, params []VarDecl) *FuncDecl {
	decl := NewFuncDecl(name, result, params)
	s.Funcs = append(s.Funcs, decl)
	return decl
}

func (s *Scope) AddAssignStmt(l, r Expr) *AssignStmt {
	stmt := NewAssignStmt(l, r)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddStoreStmt(addr, value Expr) *StoreStmt {
	stmt := NewStoreStmt(addr, value)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddLabel(name string) *Label {
	l := NewLabel(name)
	s.Stmts = append(s.Stmts, l)
	return l
}

func (s *Scope) AddBlackLine() *BlackLine {
	l := NewBlackLine()
	s.Stmts = append(s.Stmts, l)
	return l
}

func (s *Scope) AddReturnStmt(v []Expr) *ReturnStmt {
	stmt := NewReturnStmt(v)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddIfStmt(cond Expr, succs [2]int) *IfStmt {
	stmt := NewIfStmt(cond, succs)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddJumpStmt(name string) *JumpStmt {
	stmt := NewJumpStmt(name)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddPhiStmt(incoming Expr, dest Expr, edges []PhiEdge) *PhiStmt {
	stmt := NewPhiStmt(incoming, dest, edges)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddExprStmt(x Expr) Stmt {
	if st, ok := x.(Stmt); ok {
		s.Stmts = append(s.Stmts, st)
		return st
	}
	stmt := NewExprStmt(x)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) AddScope() *Scope {
	stmt := NewScope(s)
	s.Stmts = append(s.Stmts, stmt)
	return stmt
}

func (s *Scope) Lookup(name string) (*Scope, Expr) {
	for ; s != nil; s = s.Parent {
		for i := range s.Vars {
			if s.Vars[i].Name == name {
				return s, &s.Vars[i].Var
			}
		}
	}

	return nil, nil
}

func (s *Scope) CIRString() string {
	indent := 0
	for sc := s.Parent; sc != nil; sc = sc.Parent {
		indent++
	}

	buf := &strings.Builder{}
	if s.Parent != nil {
		buf.WriteString("{\n")
	}

	for _, v := range s.Structs {
		buf.WriteString(v.CIRString())
		buf.WriteString(";\n")
	}

	for _, v := range s.locals {
		buf.WriteString(genIndent(indent))
		buf.WriteString(v.CIRString())
		buf.WriteString(";\n")
	}

	for _, v := range s.Funcs {
		b := v.Body
		v.Body = nil
		buf.WriteString(genIndent(indent))
		buf.WriteString(v.CIRString())
		buf.WriteString("\n")
		v.Body = b
	}
	for _, v := range s.Funcs {
		if v.Body == nil {
			continue
		}
		buf.WriteString(genIndent(indent))
		buf.WriteString(v.CIRString())
		buf.WriteString("\n")
	}

	for _, v := range s.Stmts {
		buf.WriteString(genIndent(indent))
		buf.WriteString(v.CIRString())
		buf.WriteString("\n")
	}

	if s.Parent != nil {
		buf.WriteString(genIndent(indent - 1))
		buf.WriteString("}")
	}
	return buf.String()
}

func genIndent(c int) string {
	s := ""
	for i := 0; i < c; i++ {
		s += "\t"
	}
	return s
}
