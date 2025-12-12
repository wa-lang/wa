// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
Instruction: 指令接口
**************************************/
type Instruction interface {
	// 获取该指令的伪代码
	String() string

	// 获取该指令在源码中的位置
	Pos() int

	// 获取该指令所属的域
	ParentScope() Scope

	// 格式化输出
	Format(tab string, sb *strings.Builder)

	// 设置该指令所属的域
	setScope(Scope)
}

/**************************************
InstAlloc: Alloc 指令，分配一个变量位置，该对象同时实现了 Location 接口
**************************************/
type InstAlloc struct {
	aImv
	Location LocationKind
	dataType Type
	refType  Type
	object   interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object
}

func (i *InstAlloc) String() string {
	op := "alloc.stack"
	switch i.Location {
	case LocationKindRegister:
		op = "alloc.register"

	case LocationKindHeap:
		op = "alloc.heap"
	}

	return op
}

func (i *InstAlloc) Type() Type {
	if i.Location == LocationKindRegister {
		return i.dataType
	} else {
		return i.refType
	}
}

// Location 接口相关
func (i *InstAlloc) LocationKind() LocationKind { return i.Location }
func (i *InstAlloc) DataType() Type             { return i.dataType }
func (i *InstAlloc) Object() interface{}        { return i.object }

/**************************************
InstLoad: Load 指令，装载 Loc 处的变量
**************************************/
type InstLoad struct {
	aImv
	Loc Location
	typ Type
}

func (i *InstLoad) String() string {
	return fmt.Sprintf("*%s", i.Loc.Name())
}

func (i *InstLoad) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 Load 指令
func (b *Block) EmitLoad(loc Location, typ Type, pos int) *InstLoad {
	v := &InstLoad{}
	v.Stringer = v
	v.Loc = loc
	v.typ = typ
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstStore: Store 指令，将 Val 存储到 Loc 指定的位置
**************************************/
type InstStore struct {
	anInstruction
	Loc Location
	Val Value
}

func (i *InstStore) String() string {
	return fmt.Sprintf("*%s = %s", i.Loc.Name(), i.Val.Name())
}

// 在 Block 中添加一条 Store 指令
func (b *Block) EmitStore(loc Location, val Value, pos int) *InstStore {
	v := &InstStore{}
	v.Stringer = v
	v.Loc = loc
	v.Val = val
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstExtract: Extract 指令，提取元组 Tuple 的第 Index 个元素
**************************************/
type InstExtract struct {
	aImv
	X     Value
	Index int
}

func (i *InstExtract) String() string {
	return fmt.Sprintf("extract %s #%d", i.X.Name(), i.Index)
}

func (i *InstExtract) Type() Type {
	return i.X.Type().(*Tuple).fields[i.Index]
}

// 在 Block 中添加一条 Extract 指令
func (b *Block) EmitExtract(tuple Value, index int, pos int) *InstExtract {
	v := &InstExtract{}
	v.Stringer = v
	v.X = tuple
	v.Index = index
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstReturn: Return 指令
**************************************/
type InstReturn struct {
	anInstruction
	Results []Value
}

func (i *InstReturn) String() string {
	s := "return"
	for _, r := range i.Results {
		s += " "
		s += r.Name()
	}
	return s
}

// 在 Block 中添加一条 Return 指令
func (b *Block) EmitReturn(results []Value, pos int) *InstReturn {
	v := &InstReturn{}
	v.Stringer = v
	v.Results = results
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstUnopNot:  一元非指令，!x
**************************************/
type InstUnopNot struct {
	aImv
	X Value
}

func (i *InstUnopNot) String() string {
	return fmt.Sprintf("!%s", i.X.Name())
}

func (i *InstUnopNot) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 UnopNot 指令
func (b *Block) EmitUnopNot(x Value, pos int) *InstUnopNot {
	v := &InstUnopNot{X: x}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstUnopSub:  取负指令，-x
**************************************/
type InstUnopSub struct {
	aImv
	X Value
}

func (i *InstUnopSub) String() string {
	return fmt.Sprintf("-%s", i.X.Name())
}

func (i *InstUnopSub) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 UnopSub 指令
func (b *Block) EmitUnopSub(x Value, pos int) *InstUnopSub {
	v := &InstUnopSub{X: x}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstUnopXor:  一元异或指令，^x
**************************************/
type InstUnopXor struct {
	aImv
	X Value
}

func (i *InstUnopXor) String() string {
	return fmt.Sprintf("^%s", i.X.Name())
}

func (i *InstUnopXor) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 UnopXor 指令
func (b *Block) EmitUnopXor(x Value, pos int) *InstUnopXor {
	v := &InstUnopXor{X: x}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstLAND:  逻辑与指令，x && y
**************************************/
type InstLAND struct {
	aImv
	X, Y Value
}

func (i *InstLAND) Type() Type {
	return i.X.Type()
}

func (i *InstLAND) String() string {
	return fmt.Sprintf("%s && %s", i.X.Name(), i.Y.Name())
}

// 在 Block 中添加一条 BiopLAND 指令
func (b *Block) EmitInstLAND(x, y Value, pos int) *InstLAND {
	v := &InstLAND{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstLOR:  逻辑或指令，x || y
**************************************/
type InstLOR struct {
	aImv
	X, Y Value
}

func (i *InstLOR) String() string {
	return fmt.Sprintf("%s || %s", i.X.Name(), i.Y.Name())
}

func (i *InstLOR) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 BiopLOR 指令
func (b *Block) EmitInstLOR(x, y Value, pos int) *InstLOR {
	v := &InstLOR{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstSHL:  左移指令，x << y
**************************************/
type InstSHL struct {
	aImv
	X, Y Value
}

func (i *InstSHL) Type() Type {
	return i.X.Type()
}

func (i *InstSHL) String() string {
	return fmt.Sprintf("%s << %s", i.X.Name(), i.Y.Name())
}

// 在 Block 中添加一条 InstSHL 指令
func (b *Block) EmitInstSHL(x, y Value, pos int) *InstSHL {
	v := &InstSHL{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstSHR:  右移指令，x >> y
**************************************/
type InstSHR struct {
	aImv
	X, Y Value
}

func (i *InstSHR) String() string {
	return fmt.Sprintf("%s >> %s", i.X.Name(), i.Y.Name())
}

func (i *InstSHR) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstSHR 指令
func (b *Block) EmitInstSHR(x, y Value, pos int) *InstSHR {
	v := &InstSHR{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstADD:  加指令，x + y
**************************************/
type InstADD struct {
	aImv
	X, Y Value
}

func (i *InstADD) String() string {
	return fmt.Sprintf("%s + %s", i.X.Name(), i.Y.Name())
}

func (i *InstADD) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstADD 指令
func (b *Block) EmitInstADD(x, y Value, pos int) *InstADD {
	v := &InstADD{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstSUB:  减指令，x - y
**************************************/
type InstSUB struct {
	aImv
	X, Y Value
}

func (i *InstSUB) String() string {
	return fmt.Sprintf("%s - %s", i.X.Name(), i.Y.Name())
}

func (i *InstSUB) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstSUB 指令
func (b *Block) EmitInstSUB(x, y Value, pos int) *InstSUB {
	v := &InstSUB{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstMUL:  乘指令，x * y
**************************************/
type InstMUL struct {
	aImv
	X, Y Value
}

func (i *InstMUL) Type() Type {
	return i.X.Type()
}

func (i *InstMUL) String() string {
	return fmt.Sprintf("%s * %s", i.X.Name(), i.Y.Name())
}

// 在 Block 中添加一条 InstMUL 指令
func (b *Block) EmitInstMUL(x, y Value, pos int) *InstMUL {
	v := &InstMUL{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstQUO:  除指令，x / y
**************************************/
type InstQUO struct {
	aImv
	X, Y Value
}

func (i *InstQUO) Type() Type {
	return i.X.Type()
}

func (i *InstQUO) String() string {
	return fmt.Sprintf("%s / %s", i.X.Name(), i.Y.Name())
}

// 在 Block 中添加一条 InstQUO 指令
func (b *Block) EmitInstQUO(x, y Value, pos int) *InstQUO {
	v := &InstQUO{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstREM:  取余指令，x % y
**************************************/
type InstREM struct {
	aImv
	X, Y Value
}

func (i *InstREM) String() string {
	return fmt.Sprintf("%s %% %s", i.X.Name(), i.Y.Name())
}

func (i *InstREM) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstREM 指令
func (b *Block) EmitInstREM(x, y Value, pos int) *InstREM {
	v := &InstREM{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstAND:  与指令，x & y
**************************************/
type InstAND struct {
	aImv
	X, Y Value
}

func (i *InstAND) String() string {
	return fmt.Sprintf("%s & %s", i.X.Name(), i.Y.Name())
}

func (i *InstAND) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstAND 指令
func (b *Block) EmitInstAND(x, y Value, pos int) *InstAND {
	v := &InstAND{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstOR:  或指令，x & y
**************************************/
type InstOR struct {
	aImv
	X, Y Value
}

func (i *InstOR) String() string {
	return fmt.Sprintf("%s | %s", i.X.Name(), i.Y.Name())
}

func (i *InstOR) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstOR 指令
func (b *Block) EmitInstOR(x, y Value, pos int) *InstOR {
	v := &InstOR{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstXOR:  异或指令，x ^ y
**************************************/
type InstXOR struct {
	aImv
	X, Y Value
}

func (i *InstXOR) String() string {
	return fmt.Sprintf("%s ^ %s", i.X.Name(), i.Y.Name())
}

func (i *InstXOR) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstXOR 指令
func (b *Block) EmitInstXOR(x, y Value, pos int) *InstXOR {
	v := &InstXOR{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstANDNOT:  与或指令，x &^ y
**************************************/
type InstANDNOT struct {
	aImv
	X, Y Value
}

func (i *InstANDNOT) String() string {
	return fmt.Sprintf("%s &^ %s", i.X.Name(), i.Y.Name())
}

func (i *InstANDNOT) Type() Type {
	return i.X.Type()
}

// 在 Block 中添加一条 InstANDNOT 指令
func (b *Block) EmitInstANDNOT(x, y Value, pos int) *InstANDNOT {
	v := &InstANDNOT{X: x, Y: y}
	v.Stringer = v
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
InstEQL:  等于指令，x == y
**************************************/
type InstEQL struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstEQL) String() string {
	return fmt.Sprintf("%s == %s", i.X.Name(), i.Y.Name())
}

func (i *InstEQL) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstEQL 指令
func (b *Block) EmitInstEQL(x, y Value, pos int) *InstEQL {
	v := &InstEQL{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstNEQ:  不等于指令，x != y
**************************************/
type InstNEQ struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstNEQ) String() string {
	return fmt.Sprintf("%s != %s", i.X.Name(), i.Y.Name())
}

func (i *InstNEQ) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstNEQ 指令
func (b *Block) EmitInstNEQ(x, y Value, pos int) *InstNEQ {
	v := &InstNEQ{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstGTR:  大于指令，x > y
**************************************/
type InstGTR struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstGTR) String() string {
	return fmt.Sprintf("%s > %s", i.X.Name(), i.Y.Name())
}

func (i *InstGTR) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstGTR 指令
func (b *Block) EmitInstGTR(x, y Value, pos int) *InstGTR {
	v := &InstGTR{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstLSS:  小于指令，x < y
**************************************/
type InstLSS struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstLSS) String() string {
	return fmt.Sprintf("%s < %s", i.X.Name(), i.Y.Name())
}

func (i *InstLSS) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstLSS 指令
func (b *Block) EmitInstLSS(x, y Value, pos int) *InstLSS {
	v := &InstLSS{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstGEQ:  大等于指令，x >= y
**************************************/
type InstGEQ struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstGEQ) String() string {
	return fmt.Sprintf("%s >= %s", i.X.Name(), i.Y.Name())
}

func (i *InstGEQ) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstGEQ 指令
func (b *Block) EmitInstGEQ(x, y Value, pos int) *InstGEQ {
	v := &InstGEQ{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstLEQ:  小等于指令，x <= y
**************************************/
type InstLEQ struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstLEQ) String() string {
	return fmt.Sprintf("%s <= %s", i.X.Name(), i.Y.Name())
}

func (i *InstLEQ) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstLEQ 指令
func (b *Block) EmitInstLEQ(x, y Value, pos int) *InstLEQ {
	v := &InstLEQ{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Bool

	b.emit(v)
	return v
}

/**************************************
InstCOMP:  比较指令，x <=> y
**************************************/
type InstCOMP struct {
	aImv
	X, Y Value
	typ  Type
}

func (i *InstCOMP) String() string {
	return fmt.Sprintf("%s <=> %s", i.X.Name(), i.Y.Name())
}

func (i *InstCOMP) Type() Type {
	return i.typ
}

// 在 Block 中添加一条 InstCOMP 指令
func (b *Block) EmitInstCOMP(x, y Value, pos int) *InstCOMP {
	v := &InstCOMP{X: x, Y: y}
	v.Stringer = v
	v.pos = pos
	v.typ = b.types.Int

	b.emit(v)
	return v
}

/**************************************
InstCall:  函数调用指令，Call 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一
**************************************/
type InstCall struct {
	aImv
	Call Value
}

func (i *InstCall) Type() Type     { return i.Call.Type() }
func (i *InstCall) Pos() int       { return i.Call.Pos() }
func (i *InstCall) String() string { return i.Call.(fmt.Stringer).String() }

// 在 Block 中添加一条 InstCall 指令，Call 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一
func (b *Block) EmitInstCall(call Value) *InstCall {
	v := &InstCall{Call: call}
	v.Stringer = v
	v.pos = call.Pos()

	b.emit(v)
	return v
}

///**************************************
//InstCallBuiltin:  调用内置函数指令
//**************************************/
//type InstCallBuiltin struct {
//	aImv
//	BuiltinCall
//}
//
//// 在 Block 中添加一条 InstCallBuiltin 指令
//func (b *Block) EmitInstCallBuiltin(call BuiltinCall) *InstCallBuiltin {
//	v := &InstCallBuiltin{BuiltinCall: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
//
///**************************************
//InstCallStatic:  调用静态函数指令
//**************************************/
//type InstCallStatic struct {
//	aImv
//	StaticCall
//}
//
//// 在 Block 中添加一条 InstCallStatic 指令
//func (b *Block) EmitInstCallStatic(call StaticCall) *InstCallStatic {
//	v := &InstCallStatic{StaticCall: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
//
///**************************************
//InstCallMethod:  调用对象方法指令
//**************************************/
//type InstCallMethod struct {
//	aImv
//	MethodCall
//}
//
//func (i *InstCallMethod) String() string {
//	return i.Recv.Name() + i.Call.String()
//}
//
//// 在 Block 中添加一条 InstCallMethod 指令
//func (b *Block) EmitInstCallMethod(call MethodCall) *InstCallMethod {
//	v := &InstCallMethod{MethodCall: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
//
///**************************************
//InstCallClosure:  调用闭包指令
//**************************************/
//type InstCallClosure struct {
//	aImv
//	ClosureCall
//}
//
//func (i *InstCallClosure) String() string {
//	return i.Closure.Name() + i.Call.String()
//}
//
//// 在 Block 中添加一条 InstCallMethod 指令
//func (b *Block) EmitInstCallClosure(call ClosureCall) *InstCallClosure {
//	v := &InstCallClosure{ClosureCall: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
//
///**************************************
//InstInvoke:  调用接口方法指令
//**************************************/
//type InstInvoke struct {
//	aImv
//	InterfaceCall
//}
//
//func (i *InstInvoke) String() string {
//	return i.Interface.Name() + i.Call.String()
//}
//
//// 在 Block 中添加一条 InstInvoke 指令
//func (b *Block) EmitInstInvoke(call InterfaceCall) *InstInvoke {
//	v := &InstInvoke{InterfaceCall: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
