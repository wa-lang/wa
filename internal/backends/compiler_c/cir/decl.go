package cir

import (
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
)

/**************************************
VarDecl: 变量声明
**************************************/
type VarDecl struct {
	Var
	InitVal Expr
}

func NewVarDecl(name string, t ctypes.Type) *VarDecl {
	return &VarDecl{Var: *NewVar(name, t)}
}

func (v *VarDecl) CIRString() string {
	buf := &strings.Builder{}
	buf.WriteString(v.Type().CIRString())
	buf.WriteString(" ")
	buf.WriteString(v.Var.CIRString())
	if v.InitVal != nil {
		buf.WriteString(" = ")
		buf.WriteString(v.InitVal.CIRString())
	}
	return buf.String()
}

/**************************************
FuncDecl: 函数声明
**************************************/
type FuncDecl struct {
	Ident
	Result ctypes.Type
	Params []VarDecl
	Body   *Scope
}

func NewFuncDecl(name string, result ctypes.Type, params []VarDecl) *FuncDecl {
	return &FuncDecl{Ident: *NewIdent(name), Result: result, Params: params}
}

func (f *FuncDecl) CIRString() string {
	buf := &strings.Builder{}

	buf.WriteString(f.Result.CIRString())
	buf.WriteString(" ")
	buf.WriteString(f.Ident.CIRString())
	buf.WriteString("(")
	for i, v := range f.Params {
		buf.WriteString(v.CIRString())
		if i < len(f.Params)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")")

	if f.Body == nil {
		buf.WriteString(";")
	} else {
		buf.WriteString(" ")
		buf.WriteString(f.Body.CIRString())
	}

	return buf.String()
}

func (f *FuncDecl) isStmt() {}

/**************************************
StructDecl: 结构体声明
**************************************/
type StructDecl struct {
	ctypes.Struct
}

func NewStructDecl(name string, m []Var) *StructDecl {
	var f []ctypes.Field
	for _, i := range m {
		f = append(f, *ctypes.NewField(i.Name, i.Type()))
	}
	return &StructDecl{Struct: *ctypes.NewStruct(name, f)}
}

func (s *StructDecl) CIRString() string {
	buf := &strings.Builder{}
	buf.WriteString("struct ")
	buf.WriteString(s.Struct.CIRString())
	buf.WriteString(" {\n")

	for _, v := range s.Members {
		buf.WriteString("\t")
		buf.WriteString(v.Type().CIRString())
		buf.WriteString(" ")
		buf.WriteString(v.CIRString())
		buf.WriteString(";\n")
	}

	buf.WriteString("}")
	return buf.String()
}

func (s *StructDecl) isStmt() {}
