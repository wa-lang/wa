// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strconv"

/**************************************
InstCall:
**************************************/
type InstCall struct {
	anInstruction
	name string
}

func NewInstCall(name string) *InstCall         { return &InstCall{name: name} }
func (i *InstCall) Format(indent string) string { return indent + "call $" + i.name }

/**************************************
InstBlock:
**************************************/
type InstBlock struct {
	anInstruction
	name  string
	Insts []Inst
}

func NewInstBlock(name string) *InstBlock { return &InstBlock{name: name} }
func (i *InstBlock) Format(indent string) string {
	s := indent + "(block $"
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
InstLoop:
**************************************/
type InstLoop struct {
	anInstruction
	name  string
	Insts []Inst
}

func NewInstLoop(name string) *InstLoop { return &InstLoop{name: name} }
func (i *InstLoop) Format(indent string) string {
	s := indent + "(loop $"
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
InstBr:
**************************************/
type InstBr struct {
	anInstruction
	Name string
}

func NewInstBr(name string) *InstBr           { return &InstBr{Name: name} }
func (i *InstBr) Format(indent string) string { return indent + "br $" + i.Name }

/**************************************
InstBrTable:
**************************************/
type InstBrTable struct {
	anInstruction
	Table []int
}

func NewInstBrTable(t []int) *InstBrTable { return &InstBrTable{Table: t} }
func (i *InstBrTable) Format(indent string) string {
	s := indent + "br_table"
	for _, v := range i.Table {
		s += " " + strconv.Itoa(v)
	}
	return s
}

/**************************************
InstIf:
**************************************/
type InstIf struct {
	anInstruction
	True  []Inst
	False []Inst
	Ret   []ValueType
}

func NewInstIf(instsTrue, instsFalse []Inst, ret []ValueType) *InstIf {
	return &InstIf{True: instsTrue, False: instsFalse, Ret: ret}
}
func (i *InstIf) Format(indent string) string {
	s := indent + "if"
	if len(i.Ret) > 0 {
		s += " (result"
		for _, r := range i.Ret {
			s += " " + r.Name()
		}
		s += ")"
	}
	s += "\n"

	for _, v := range i.True {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + "else\n"
	for _, v := range i.False {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + "end"
	return s
}

/**************************************
InstReturn:
**************************************/
type InstReturn struct {
	anInstruction
}

func NewInstReturn() *InstReturn                  { return &InstReturn{} }
func (i *InstReturn) Format(indent string) string { return indent + "return" }
