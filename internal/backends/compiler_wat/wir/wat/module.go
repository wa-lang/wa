// 版权 @2022 凹语言 作者。保留所有权利。

package wat

// 模块对象
type Module struct {
	Name string
	//Imports []Import
	Globals []Global
	Funcs   []*Function

	BaseWat string
}

type Global struct {
	V         Value
	IsMut     bool
	InitValue string
}

func (m *Module) String() string {
	s := "(module $__walang__\n"

	//for _, i := range m.Imports {
	//	s += i.Format("  ") + "\n"
	//}

	for _, g := range m.Globals {
		s += "(global "
		s += g.V.Name()
		if g.IsMut {
			s += "(mut " + g.V.Type().Name() + ") "
		} else {
			s += "(" + g.V.Type().Name() + ") "
		}
		if len(g.InitValue) > 0 {
			s += "(" + g.V.Type().Name() + ".const " + g.InitValue + ")"
		}
		s += ")\n"
	}

	s += m.BaseWat

	for _, f := range m.Funcs {
		s += "\n\n" + f.Format("")
	}

	s += "\n) ;;module"
	return s
}
