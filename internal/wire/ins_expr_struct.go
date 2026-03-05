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
	X  Expr
	Id int

	types       *Types
	member_type Type
	member_name string
}

func (i *MemberLocation) Name() string   { return i.String() }
func (i *MemberLocation) Type() Type     { return i.member_type }
func (i *MemberLocation) retained() bool { return false }
func (i *MemberLocation) String() string {
	return fmt.Sprintf("member_location(%s, %s)", i.X.Name(), i.member_name)
}

func (i *MemberLocation) renewUnderlying() Expr {
	if x, ok := i.X.(Var); ok && x.Kind() == Register && unname(x.Type()).Kind() == TypeKindStruct {
		return newMember(x, i.Id, i.pos)
	} else {
		return newMemberAddr(i.X, i.Id, i.pos, i.types)
	}
}

func (b *Block) NewMemberLocation(x Expr, id int, pos int) *MemberLocation {
	v := &MemberLocation{X: x, Id: id, types: b.types}
	v.Stringer = v
	v.pos = pos

	xt := x.Type()

	switch t := xt.(type) {
	case *Ref:
		if st, ok := unname(t.Base).(*Struct); ok {
			m := st.At(id)
			v.member_type = b.types.GenRef(m.Type)
			v.member_name = m.Name
		} else {
			panic(fmt.Sprintf("Invalid X.Type():%s", xt.Name()))
		}

	case *Ptr:
		if st, ok := unname(t.Base).(*Struct); ok {
			m := st.At(id)
			v.member_type = b.types.genPtr(m.Type)
			v.member_name = m.Name
		} else {
			panic(fmt.Sprintf("Invalid X.Type():%s", xt.Name()))
		}

	default:
		if st, ok := unname(t).(*Struct); ok {
			m := st.At(id)
			v.member_type = b.types.GenRef(m.Type)
			v.member_name = m.Name
		} else {
			panic(fmt.Sprintf("Invalid X.Type():%s", xt.Name()))
		}
	}

	return v
}

/**************************************
Member: 结构体成员访问指令，Member 实现了 Expr、Var 接口
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
MemberAddr: 获取结构体成员地址，Member 实现了 Expr 接口
  - X 的类型应为 Ref(T)、或 Ptr(T)，其中 T 为匿名或具名结构体
**************************************/

type MemberAddr struct {
	aStmt
	X  Expr
	Id int

	typ   Type
	mname string
}

func (i *MemberAddr) Name() string   { return i.String() }
func (i *MemberAddr) Type() Type     { return i.typ }
func (i *MemberAddr) retained() bool { return false }
func (i *MemberAddr) String() string { return fmt.Sprintf("%s->%s", i.X.Name(), i.mname) }

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

		if xkind == TypeKindRef {
			v.typ = types.GenRef(st.At(id).Type)
		} else {
			v.typ = types.genPtr(st.At(id).Type)
		}
		v.mname = st.At(id).Name

		return v
	}
}
