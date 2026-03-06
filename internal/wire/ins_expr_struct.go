// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
MemberLocation: 结构体成员位置指令，MemberLocation 实现了 Expr
**************************************/

type MemberLocation struct {
	aStmt
	X  Location
	Id int

	types  *Types
	member StructMember
}

func (i *MemberLocation) Name() string   { return i.String() }
func (i *MemberLocation) DataType() Type { return i.member.Type }
func (i *MemberLocation) String() string {
	return fmt.Sprintf("member_location(%s, %s)", i.X.Name(), i.member.Name)
}

func (b *Block) NewMemberLocation(x Location, id int, pos int) *MemberLocation {
	v := &MemberLocation{X: x, Id: id, types: b.types}
	v.Stringer = v
	v.pos = pos

	if st, ok := unname(x.DataType()).(*Struct); ok {
		v.member = st.At(id)
	} else {
		panic(fmt.Sprintf("Invalid X.DataType():%s", x.DataType().Name()))
	}

	return v
}

/**************************************
Member: 寄存器型结构体成员访问指令，Member 实现了 Expr、Var 接口
  - X 应为 匿名或具名 Struct 类型的 Register 变量
**************************************/

type Member struct {
	aStmt
	X      Var
	Id     int
	member StructMember
}

func (i *Member) Name() string   { return i.String() }
func (i *Member) Type() Type     { return i.member.Type }
func (i *Member) retained() bool { return false }
func (i *Member) String() string {
	if tank := i.Tank(); tank != nil {
		return fmt.Sprintf("%s.%s:%s", i.X.Name(), i.member.Name, tank.String())
	} else {
		return fmt.Sprintf("%s.%s", i.X.Name(), i.member.Name)
	}
}
func (i *Member) Kind() LocationKind { return Register }
func (i *Member) DataType() Type     { return i.Type() }
func (i *Member) Tank() *tank {
	tank := i.X.Tank()
	if tank != nil {
		return i.X.Tank().member[i.Id]
	}
	return nil
}

func newMember(x Var, id int, pos int) *Member {
	if x.Kind() != Register {
		panic("Member.X must be Register")
	}

	if st, ok := unname(x.Type()).(*Struct); !ok {
		panic(fmt.Sprintf("Invalid X.Type():%s", x.Type().Name()))
	} else {
		if id < 0 || id >= st.Len() {
			panic(fmt.Sprintf("Member.Id out of range: %d", id))
		}

		v := &Member{X: x, Id: id}
		v.Stringer = v
		v.pos = pos
		v.member = st.At(id)
		return v
	}
}

/**************************************
MemberAddr: 获取结构体成员地址，Member 实现了 Expr、Location接口
  - X 的类型应为 Ref(T)、或 Ptr(T)，其中 T 为匿名或具名结构体
**************************************/

type MemberAddr struct {
	aStmt
	X  Expr
	Id int

	member StructMember
	typ    Type
}

func (i *MemberAddr) Name() string   { return i.String() }
func (i *MemberAddr) Type() Type     { return i.typ }
func (i *MemberAddr) retained() bool { return false }
func (i *MemberAddr) String() string { return fmt.Sprintf("%s->%s", i.X.Name(), i.member.Name) }
func (i *MemberAddr) DataType() Type { return i.member.Type }

func newMemberAddr(x Expr, id int, pos int, types *Types) *MemberAddr {
	xkind := x.Type().Kind()
	if xkind != TypeKindRef && xkind != TypeKindPtr {
		panic("MemberAddr.X must be a Ref or Ptr")
	}

	if st, ok := unname(deref(x.Type())).(*Struct); !ok {
		panic(fmt.Sprintf("Invalid X.Type():%s", x.Type().Name()))
	} else {
		if id < 0 || id >= st.Len() {
			panic(fmt.Sprintf("MemberAddr.Id out of range: %d", id))
		}

		v := &MemberAddr{X: x, Id: id}
		v.Stringer = v
		v.pos = pos

		v.member = st.At(id)
		if xkind == TypeKindRef {
			v.typ = types.GenRef(v.member.Type)
		} else {
			v.typ = types.genPtr(v.member.Type)
		}

		return v
	}
}
