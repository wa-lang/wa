package wir

func (m *Module) String() string {
	s := "(module $__walang__\n"

	s += m.BaseWat + "\n"

	for _, i := range m.Imports {
		s += i.Format("  ") + "\n"
	}

	for _, f := range m.Funcs {
		s += f.Format("  ") + "\n"
	}

	s += ") ;;module"
	return s
}
