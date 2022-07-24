package types

import "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"

type PackageInstance struct {
	backingType
	name  string
	funcs map[string]*Function
}

func (p *PackageInstance) SetName(name string) {
	p.name = name
}

func (p *PackageInstance) SetFunc(name string, val *Function) {
	if p.funcs == nil {
		p.funcs = make(map[string]*Function)
	}
	p.funcs[name] = val
}

func (p *PackageInstance) GetFunc(name string) (*Function, bool) {
	v, ok := p.funcs[name]
	return v, ok
}

func (PackageInstance) LLVM() lltypes.Type {
	// TODO: Packages are not values, and should be represented some other way
	// Maybe via LLVM IR modules?
	panic("Package does not have LLVM defined")
}

func (p PackageInstance) Name() string {
	return p.name
}
