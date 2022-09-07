// 版权 @2022 凹语言 作者。保留所有权利。

package wat

// 模块对象
type Module struct {
	Name    string
	Imports []Import
	Funcs   []*Function

	BaseWat string
}

func (m *Module) String() string {
	s := "(module $__walang__\n"

	for _, i := range m.Imports {
		s += i.Format("  ") + "\n"
	}

	s += m.BaseWat + "\n"

	for _, f := range m.Funcs {
		s += f.Format("  ") + "\n"
	}

	s += ") ;;module"
	return s
}
