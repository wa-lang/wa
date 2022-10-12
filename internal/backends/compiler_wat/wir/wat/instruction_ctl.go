// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strconv"

/**************************************
instCall:
**************************************/
type instCall struct {
	anInstruction
	name string
}

func NewInstCall(name string) *instCall         { return &instCall{name: name} }
func (i *instCall) Format(indent string) string { return indent + "call $" + i.name }

/**************************************
instCallIndirect:
**************************************/
type instCallIndirect struct {
	anInstruction
	func_type string
}

func NewInstCallIndirect(func_type string) *instCallIndirect {
	return &instCallIndirect{func_type: func_type}
}
func (i *instCallIndirect) Format(indent string) string {
	return indent + "call_indirect (type $" + i.func_type + ")"
}

/**************************************
instBlock:
**************************************/
type instBlock struct {
	anInstruction
	name  string
	Insts []Inst
}

func NewInstBlock(name string) *instBlock { return &instBlock{name: name} }
func (i *instBlock) Format(indent string) string {
	s := indent + "(block $"
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
instLoop:
**************************************/
type instLoop struct {
	anInstruction
	name  string
	Insts []Inst
}

func NewInstLoop(name string) *instLoop { return &instLoop{name: name} }
func (i *instLoop) Format(indent string) string {
	s := indent + "(loop $"
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
instBr:
**************************************/
type instBr struct {
	anInstruction
	Name string
}

func NewInstBr(name string) *instBr           { return &instBr{Name: name} }
func (i *instBr) Format(indent string) string { return indent + "br $" + i.Name }

/**************************************
instBrTable:
**************************************/
type instBrTable struct {
	anInstruction
	Table []int
}

func NewInstBrTable(t []int) *instBrTable { return &instBrTable{Table: t} }
func (i *instBrTable) Format(indent string) string {
	s := indent + "br_table"
	for _, v := range i.Table {
		s += " " + strconv.Itoa(v)
	}
	return s
}

/**************************************
instIf:
**************************************/
type instIf struct {
	anInstruction
	True  []Inst
	False []Inst
	Ret   []ValueType
}

func NewInstIf(instsTrue, instsFalse []Inst, ret []ValueType) *instIf {
	return &instIf{True: instsTrue, False: instsFalse, Ret: ret}
}
func (i *instIf) Format(indent string) string {
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
instReturn:
**************************************/
type instReturn struct {
	anInstruction
}

func NewInstReturn() *instReturn                  { return &instReturn{} }
func (i *instReturn) Format(indent string) string { return indent + "return" }
