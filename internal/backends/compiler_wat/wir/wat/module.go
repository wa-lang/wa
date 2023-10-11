// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"strconv"
	"strings"
)

// 模块对象
type Module struct {
	Name    string
	Imports []Import

	Tables    Table
	FuncTypes []FuncType
	Globals   []Global
	Funcs     []*Function
	DataSeg   *DataSeg

	BaseWat string
}

func (m *Module) String() string {
	var sb strings.Builder

	sb.WriteString("(module $__walang__\n")

	for _, i := range m.Imports {
		sb.WriteString(i.Format("  "))
		sb.WriteByte('\n')
	}

	sb.WriteString(m.BaseWat)

	if len(m.DataSeg.data) > 0 {
		sb.WriteString("(data (i32.const ")
		sb.WriteString(strconv.Itoa(m.DataSeg.start))
		sb.WriteString(") \"")
		for _, d := range m.DataSeg.data {
			sb.WriteByte('\\')
			i := strconv.FormatInt(int64(d), 16)
			if len(i) < 2 {
				i = "0" + i
			}
			sb.WriteString(i)
		}
		sb.WriteString("\")\n")
	}

	sb.WriteString(m.Tables.String())

	for _, ft := range m.FuncTypes {
		sb.WriteString(ft.String())
	}

	for _, g := range m.Globals {
		sb.WriteString("(global $")
		sb.WriteString(g.V.Name())
		if len(g.NameExp) > 0 {
			sb.WriteString(" (export \"")
			sb.WriteString(g.NameExp)
			sb.WriteString("\")")
		}
		if g.IsMut {
			sb.WriteString(" (mut ")
			sb.WriteString(g.V.Type().Name())
			sb.WriteByte(')')
		} else {
			sb.WriteByte(' ')
			sb.WriteString(g.V.Type().Name())
		} //
		if len(g.InitValue) > 0 {
			sb.WriteString(" (")
			sb.WriteString(g.V.Type().Name())
			sb.WriteString(".const ")
			sb.WriteString(g.InitValue)
			sb.WriteByte(')')
		} else {
			sb.WriteString(" (")
			sb.WriteString(g.V.Type().Name())
			sb.WriteString(".const 0)")
		}
		sb.WriteString(")\n")
	}

	for _, f := range m.Funcs {
		sb.WriteString("\n\n")
		f.Format(&sb)
	}

	sb.WriteString("\n) ;;module")
	return sb.String()
}

func (t Table) String() string {
	s := "(table " + strconv.Itoa(len(t.Elems)) + " funcref)\n"

	for i, e := range t.Elems {
		if len(e) == 0 {
			continue
		}

		s += "(elem (i32.const " + strconv.Itoa(i) + ") $" + e + ")\n"
	}

	return s
}

func (t FuncType) String() string {
	s := "(type $" + t.Name + " (func"

	for _, p := range t.Params {
		s += " (param " + p.Name() + ")"
	}

	for _, r := range t.Results {
		s += " (result " + r.Name() + ")"
	}

	s += "))\n"
	return s
}
