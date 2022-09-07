package wat

func (f *Function) Format(indent string) string {
	s := indent + "(func $" + f.Name + " (export \"" + f.Name + "\")"

	for _, param := range f.Params {
		s += " (param $" + param.Name() + " " + param.Type().Name() + ")"
	}

	if len(f.Results) > 0 {
		s += " (result"
		for _, r := range f.Results {
			s += " " + r.Name()
		}
		s += ")"
	}
	s += "\n"

	for _, local := range f.Locals {
		s += indent
		s += "  (local $" + local.Name() + " " + local.Type().Name() + ")"
		s += "\n"
	}

	for _, inst := range f.Insts {
		s += inst.Format(indent+"  ") + "\n"
	}

	s += indent + ") ;;" + f.Name
	return s
}

func (sig *FuncSig) String() string {
	str := ""
	for _, param := range sig.Params {
		str += " (param " + param.Name() + ")"
	}

	for _, ret := range sig.Results {
		str += " (result " + ret.Name() + ")"
	}
	return str
}
