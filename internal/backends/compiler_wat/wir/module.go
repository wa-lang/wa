package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
)

/**************************************
Module:
**************************************/
type Module struct {
	VOID, RUNE, I8, U8, I16, U16, I32, U32, I64, U64, F32, F64, STRING ValueType

	types_map map[string]ValueType

	imports []wat.Import

	fn_types     []*FnType
	fn_types_map map[string]int

	funcs     []*Function
	funcs_map map[string]*Function

	table     []string
	table_map map[string]int

	globals     []Value
	Globals_map map[ssa.Value]Value

	data_seg []byte

	BaseWat string
}

func NewModule() *Module {
	var m Module

	m.VOID = &tVoid{}
	m.RUNE = &tRune{}
	m.I8 = &tI8{}
	m.U8 = &tU8{}
	m.I16 = &tI16{}
	m.U16 = &tU16{}
	m.I32 = &tI32{}
	m.U32 = &tU32{}
	m.I64 = &tI64{}
	m.U64 = &tU64{}
	m.F32 = &tF32{}
	m.F64 = &tF64{}

	m.types_map = make(map[string]ValueType)
	m.regValueType(m.VOID)
	m.regValueType(m.RUNE)
	m.regValueType(m.I8)
	m.regValueType(m.U8)
	m.regValueType(m.I16)
	m.regValueType(m.U16)
	m.regValueType(m.I32)
	m.regValueType(m.U32)
	m.regValueType(m.I64)
	m.regValueType(m.U64)
	m.regValueType(m.F32)
	m.regValueType(m.F64)

	m.STRING = m.GenValueType_String()

	m.fn_types_map = make(map[string]int)
	{
		var free_type FnType
		free_type.Name = "$onFree"
		free_type.Params = []ValueType{m.I32}
		m.addFnType(&free_type)
	}

	m.funcs_map = make(map[string]*Function)

	//table中先行插入一条记录，防止产生0值（无效值）id
	m.table = append(m.table, "")
	m.table_map = make(map[string]int)

	//data_seg中先插入标志，防止产生0值
	m.data_seg = append(m.data_seg, []byte("$$wads$$")...)

	m.Globals_map = make(map[ssa.Value]Value)
	return &m
}

func (m *Module) AddImportFunc(moduleName string, objName string, funcName string, sig FnSig) {
	var wat_sig wat.FuncSig
	for _, i := range sig.Params {
		wat_sig.Params = append(wat_sig.Params, i.Raw()...)
	}
	for _, r := range sig.Results {
		wat_sig.Results = append(wat_sig.Results, r.Raw()...)
	}

	m.imports = append(m.imports, wat.NewImpFunc(moduleName, objName, funcName, wat_sig))
}

func (m *Module) findFnType(name string) int {
	if i, ok := m.fn_types_map[name]; ok {
		return i
	}
	return 0
}

func (m *Module) addFnType(typ *FnType) int {
	if i := m.findFnType(typ.Name); i != 0 {
		return i
	}
	i := len(m.fn_types)
	m.fn_types = append(m.fn_types, typ)
	m.fn_types_map[typ.Name] = i
	return i
}

func (m *Module) findTableElem(elem string) int {
	if i, ok := m.table_map[elem]; ok {
		return i
	}
	return 0
}

func (m *Module) AddTableElem(elem string) int {
	if i := m.findTableElem(elem); i != 0 {
		return i
	}
	i := len(m.table)
	m.table = append(m.table, elem)
	m.table_map[elem] = i
	return i
}

func (m *Module) FindFunc(fn_name string) *Function {
	if f, ok := m.funcs_map[fn_name]; ok {
		return f
	}
	return nil
}

func (m *Module) AddFunc(f *Function) {
	if m.FindFunc(f.InternalName) == nil {
		m.funcs = append(m.funcs, f)
		m.funcs_map[f.InternalName] = f
	}
}

func (m *Module) AddGlobal(name string, typ ValueType, is_pointer bool, ssa_value ssa.Value) Value {
	v := NewGlobal(name, typ, is_pointer)
	m.globals = append(m.globals, v)
	if ssa_value != nil {
		m.Globals_map[ssa_value] = v
	}
	return v
}

func (m *Module) FindGlobal(name string) Value {
	for _, g := range m.globals {
		if g.Name() == name {
			return g
		}
	}
	return nil
}

func (m *Module) AddDataSeg(data []byte) (ptr int) {
	ptr = len(m.data_seg)
	m.data_seg = append(m.data_seg, data...)
	return
}

func (m *Module) genGlobalAlloc() *Function {
	var f Function
	f.InternalName = "$waGlobalAlloc"

	for _, g := range m.globals {
		if g.Kind() != ValueKindGlobal_Pointer {
			continue
		}

		ref := g.(*aRef)
		t := ref.Type().(*Ref).Base
		f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.size())))
		f.Insts = append(f.Insts, wat.NewInstCall("$waHeapAlloc"))
		f.Insts = append(f.Insts, ref.Extract("data").EmitPop()...)
	}

	return &f
}

func (m *Module) ToWatModule() *wat.Module {
	var wat_module wat.Module
	wat_module.Imports = m.imports
	wat_module.BaseWat = m.BaseWat

	for _, t := range m.fn_types {
		var fn_type wat.FuncType
		fn_type.Name = t.Name
		for _, i := range t.Params {
			fn_type.Params = append(fn_type.Params, i.Raw()...)
		}
		for _, r := range t.Results {
			fn_type.Results = append(fn_type.Results, r.Raw()...)
		}
		wat_module.FuncTypes = append(wat_module.FuncTypes, fn_type)
	}

	{
		wat_module.Tables.Elems = m.table
	}

	wat_module.Funcs = append(wat_module.Funcs, m.genGlobalAlloc().ToWatFunc())

	for _, f := range m.funcs {
		wat_module.Funcs = append(wat_module.Funcs, f.ToWatFunc())
	}

	for _, g := range m.globals {
		raw := g.raw()
		for _, r := range raw {
			var wat_global wat.Global
			wat_global.V = r
			wat_global.IsMut = true
			wat_module.Globals = append(wat_module.Globals, wat_global)
		}
	}

	wat_module.DataSeg = m.data_seg

	return &wat_module
}

func (m *Module) regValueType(t ValueType) {
	_, ok := m.types_map[t.Name()]
	if ok {
		logger.Fatalf("ValueType:%T already registered.", t)
	}
	m.types_map[t.Name()] = t
}

func (m *Module) findValueType(name string) (ValueType, bool) {
	t, ok := m.types_map[name]
	return t, ok
}
