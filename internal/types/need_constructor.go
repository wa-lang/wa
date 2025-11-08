package types

var _ = (*Checker).needConstructor

func (check *Checker) needConstructor(typ Type) error {
	return nil
}
