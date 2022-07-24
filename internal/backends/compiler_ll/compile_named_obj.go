package compiler

import "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"

func (p *Compiler) llGetNamedType(mangledName string) lltypes.Type {
	return p.llNamedTypeMap[mangledName]
}

func (p *Compiler) llGetTypeName(typ lltypes.Type) string {
	return p.llTypeNameMap[typ]
}

func (p *Compiler) llRegTypeName(typ lltypes.Type, name string) {
	p.llTypeNameMap[typ] = name
	p.llNamedTypeMap[name] = typ
}
