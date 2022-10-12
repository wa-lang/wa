// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strconv"

// 模块对象
type Module struct {
	Name string
	//Imports []Import

	Tables    Table
	FuncTypes []FuncType
	Globals   []Global
	Funcs     []*Function

	BaseWat string
}

func (m *Module) String() string {
	s := "(module $__walang__\n"

	//for _, i := range m.Imports {
	//	s += i.Format("  ") + "\n"
	//}

	s += m.BaseWat

	s += m.Tables.String()

	for _, ft := range m.FuncTypes {
		s += ft.String()
	}

	for _, g := range m.Globals {
		s += "(global $"
		s += g.V.Name()
		if g.IsMut {
			s += " (mut " + g.V.Type().Name() + ")"
		} else {
			s += " (" + g.V.Type().Name() + ")"
		}
		if len(g.InitValue) > 0 {
			s += " (" + g.V.Type().Name() + ".const " + g.InitValue + ")"
		} else {
			s += " (" + g.V.Type().Name() + ".const 0)"
		}
		s += ")\n"
	}

	for _, f := range m.Funcs {
		s += "\n\n" + f.Format("")
	}

	s += "\n) ;;module"
	return s
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
