package cir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Stmt:
**************************************/
type Stmt interface {
	CIRStringer
	isStmt()
}

/**************************************
ExprStmt:
**************************************/
type ExprStmt struct {
	X Expr
}

func NewExprStmt(x Expr) *ExprStmt {
	return &ExprStmt{X: x}
}

func (s *ExprStmt) CIRString() string {
	return s.X.CIRString() + ";"
}

func (s *ExprStmt) isStmt() {}

/**************************************
VarDeclStmt:
**************************************/
type VarDeclStmt struct {
	VarDecl
}

func NewVarDeclStmt(name string, t ctypes.Type) *VarDeclStmt {
	return &VarDeclStmt{VarDecl: *NewVarDecl(name, t)}
}

func (s *VarDeclStmt) CIRString() string {
	return s.VarDecl.CIRString() + ";"
}

func (s *VarDeclStmt) isStmt() {}

/**************************************
AssignStmt:
**************************************/
type AssignStmt struct {
	Lh Expr
	Rh Expr
}

func NewAssignStmt(l, r Expr) *AssignStmt {
	return &AssignStmt{Lh: l, Rh: r}
}

func (s *AssignStmt) CIRString() string {
	return s.Lh.CIRString() + " = " + s.Rh.CIRString() + ";"
}

func (s *AssignStmt) isStmt() {}

/**************************************
StoreStmt:
**************************************/
type StoreStmt struct {
	addr  Expr
	value Expr
}

func NewStoreStmt(addr, value Expr) *StoreStmt {
	return &StoreStmt{addr: addr, value: value}
}

func (s *StoreStmt) CIRString() string {
	switch s.addr.Type().(type) {
	case *ctypes.PointerType:
		return "*(" + s.addr.CIRString() + ")" + " = " + s.value.CIRString() + ";"

	case *ctypes.RefType:
		return s.addr.CIRString() + ".SetValue(" + s.value.CIRString() + ");"
	}

	logger.Fatal("addr is not a valid pointer/reference:", s.addr)
	return ""
}

func (s *StoreStmt) isStmt() {}

/**************************************
Label:
**************************************/
type Label struct {
	Name string
}

func NewLabel(name string) *Label {
	return &Label{Name: name}
}

func (s *Label) CIRString() string {
	return s.Name + ":"
}

func (s *Label) isStmt() {}

/**************************************
BlackLine:
**************************************/
type BlackLine struct {
}

func NewBlackLine() *BlackLine {
	return &BlackLine{}
}

func (s *BlackLine) CIRString() string {
	return ""
}

func (s *BlackLine) isStmt() {}

/**************************************
ReturnStmt:
**************************************/
type ReturnStmt struct {
	Vars []Expr
}

func NewReturnStmt(v []Expr) *ReturnStmt {
	return &ReturnStmt{Vars: v}
}

func (s *ReturnStmt) CIRString() string {
	switch len(s.Vars) {
	case 0:
		return "return;"

	case 1:
		return "return " + s.Vars[0].CIRString() + ";"

	default:
		str := "return {"
		for i, v := range s.Vars {
			str += v.CIRString()
			if i < len(s.Vars)-1 {
				str += ", "
			}
		}
		str += "};"
		return str
	}
}

func (s *ReturnStmt) isStmt() {}

/**************************************
IfStmt:
**************************************/
type IfStmt struct {
	Cond  Expr
	Succs [2]int
}

func NewIfStmt(cond Expr, succs [2]int) *IfStmt {
	return &IfStmt{Cond: cond, Succs: succs}
}

func (s *IfStmt) CIRString() string {
	str := "if ("
	str += s.Cond.CIRString()
	str += ") goto $Block_"
	str += strconv.Itoa(s.Succs[0])
	str += "; else goto $Block_"
	str += strconv.Itoa(s.Succs[1])
	str += ";"
	return str
}

func (s *IfStmt) isStmt() {}

/**************************************
JumpStmt:
**************************************/
type JumpStmt struct {
	Name string
}

func NewJumpStmt(name string) *JumpStmt {
	return &JumpStmt{Name: name}
}

func (s *JumpStmt) CIRString() string {
	return "goto " + s.Name + ";"
}

func (s *JumpStmt) isStmt() {}

/**************************************
PhiStmt:
**************************************/
type PhiEdge struct {
	Incoming int
	Value    Expr
}
type PhiStmt struct {
	IncomingBlockId Expr
	Dest            Expr
	Edges           []PhiEdge
}

func NewPhiStmt(incoming Expr, dest Expr, edges []PhiEdge) *PhiStmt {
	return &PhiStmt{IncomingBlockId: incoming, Dest: dest, Edges: edges}
}

func (s *PhiStmt) CIRString() string {
	str := "if "

	for i, v := range s.Edges {
		str += "("
		str += s.IncomingBlockId.CIRString()
		str += " == "
		str += strconv.Itoa(v.Incoming)
		str += ") "
		str += s.Dest.CIRString()
		str += " = "
		str += v.Value.CIRString()
		str += ";"

		if i < len(s.Edges)-1 {
			str += " else if "
		}
	}

	return str
}

func (s *PhiStmt) isStmt() {}
