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
	VOID, RUNE, I8, U8, I16, U16, I32, U32, UPTR, I64, U64, F32, F64, STRING ValueType

	types_map         map[string]ValueType
	usedConcreteTypes []ValueType
	usedInterfaces    []ValueType

	imports []wat.Import

	fnSigs     []*FnSig
	fnSigsName map[string]string

	funcs     []*Function
	funcs_map map[string]*Function

	table     []string
	table_map map[string]int

	globals []struct {
		v        Value
		init_val string
	}
	globalsMapByValue map[ssa.Value]int
	globalsMapByName  map[string]int

	constGlobals []wat.Global

	DataSeg *DataSeg

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
	m.UPTR = m.U32
	m.I64 = &tI64{}
	m.U64 = &tU64{}
	m.F32 = &tF32{}
	m.F64 = &tF64{}

	m.types_map = make(map[string]ValueType)
	m.addValueType(m.VOID)
	m.addValueType(m.RUNE)
	m.addValueType(m.I8)
	m.addValueType(m.U8)
	m.addValueType(m.I16)
	m.addValueType(m.U16)
	m.addValueType(m.I32)
	m.addValueType(m.U32)
	m.addValueType(m.I64)
	m.addValueType(m.U64)
	m.addValueType(m.F32)
	m.addValueType(m.F64)

	m.STRING = m.GenValueType_String()

	m.fnSigsName = make(map[string]string)

	m.funcs_map = make(map[string]*Function)

	//table中先行插入一条记录，防止产生0值（无效值）id
	m.table = append(m.table, "")
	m.table_map = make(map[string]int)

	//data_seg中先插入标志，防止产生0值
	m.DataSeg = newDataSeg(0)
	m.DataSeg.Append([]byte("$$wads$$"), 8)

	m.globalsMapByValue = make(map[ssa.Value]int)
	m.globalsMapByName = make(map[string]int)
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

func (m *Module) FindFnSig(sig *FnSig) string {
	if s, ok := m.fnSigsName[sig.String()]; ok {
		return s
	}
	return ""
}

func (m *Module) AddFnSig(sig *FnSig) string {
	if s, ok := m.fnSigsName[sig.String()]; ok {
		return s
	}

	m.fnSigs = append(m.fnSigs, sig)
	s := "$$fnSig" + strconv.Itoa(len(m.fnSigs))

	m.fnSigsName[sig.String()] = s
	return s
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
	v := struct {
		v        Value
		init_val string
	}{v: NewGlobal(name, typ, is_pointer)}

	if ssa_value != nil {
		m.globalsMapByValue[ssa_value] = len(m.globals)
	}
	m.globalsMapByName[name] = len(m.globals)
	m.globals = append(m.globals, v)
	return v.v
}

func (m *Module) AddConstGlobal(name string, init_val string, typ ValueType) {
	var v wat.Global
	v.V = wat.NewVar(name, toWatType(typ))
	v.IsMut = false
	v.InitValue = init_val

	m.constGlobals = append(m.constGlobals, v)
}

func (m *Module) FindGlobalByName(name string) Value {
	id, ok := m.globalsMapByName[name]
	if !ok {
		return nil
	}

	return m.globals[id].v
}

func (m *Module) FindGlobalByValue(v ssa.Value) Value {
	id, ok := m.globalsMapByValue[v]
	if !ok {
		return nil
	}

	return m.globals[id].v
}

func (m *Module) SetGlobalInitValue(name string, val string) {
	id, ok := m.globalsMapByName[name]
	if !ok {
		logger.Fatalf("Global not found:%s", name)
	}

	m.globals[id].init_val = val
}

func (m *Module) genGlobalAlloc() *Function {
	var f Function
	f.InternalName = "$waGlobalAlloc"

	for _, g := range m.globals {
		if g.v.Kind() != ValueKindGlobal_Pointer {
			continue
		}

		ref := g.v.(*aRef)
		t := ref.Type().(*Ref).Base
		f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Size())))
		f.Insts = append(f.Insts, wat.NewInstCall("$waHeapAlloc"))
		f.Insts = append(f.Insts, ref.Extract("data").EmitPop()...)
	}

	return &f
}

func (m *Module) ToWatModule() *wat.Module {
	m.buildItab()

	var wat_module wat.Module
	wat_module.Imports = m.imports
	wat_module.BaseWat = m.BaseWat

	{
		var onfree_type wat.FuncType
		onfree_type.Name = "$onFree"
		onfree_type.Params = m.I32.Raw()
		wat_module.FuncTypes = append(wat_module.FuncTypes, onfree_type)
	}
	for _, t := range m.fnSigs {
		var fn_type wat.FuncType
		fn_type.Name = m.fnSigsName[t.String()]
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
		raw := g.v.raw()
		for _, r := range raw {
			var wat_global wat.Global
			wat_global.V = r
			wat_global.IsMut = true
			wat_global.InitValue = g.init_val
			wat_module.Globals = append(wat_module.Globals, wat_global)
		}
	}

	wat_module.Globals = append(wat_module.Globals, m.constGlobals...)

	wat_module.DataSeg = m.DataSeg.data

	return &wat_module
}

func (m *Module) addValueType(t ValueType) {
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

func (m *Module) markConcreteTypeUsed(t ValueType) {
	if t.Hash() != 0 {
		return
	}

	m.usedConcreteTypes = append(m.usedConcreteTypes, t)
	t.SetHash(len(m.usedConcreteTypes))
}

func (m *Module) markInterfaceUsed(t ValueType) {
	if t.Hash() != 0 {
		return
	}

	m.usedInterfaces = append(m.usedInterfaces, t)
	t.SetHash(-len(m.usedInterfaces))
}

func (m *Module) buildItab() {
	var itabs []byte
	t_itab := m.types_map["runtime._itab"]

	for _, conrete := range m.usedConcreteTypes {
		for _, iface := range m.usedInterfaces {
			fits := true

			vtable := make([]int, iface.NumMethods())

			for mid := 0; mid < iface.NumMethods(); mid++ {
				method := iface.Method(mid)
				found := false
				for fid := 0; fid < conrete.NumMethods(); fid++ {
					d := conrete.Method(fid)
					if d.Name == method.Name && d.Sig.Equal(&method.Sig) {
						found = true
						vtable[mid] = m.AddTableElem(d.FullFnName)
						break
					}
				}

				if !found {
					fits = false
					break
				}
			}

			var addr int
			if fits {
				var itab_bin []byte
				header := NewConst("0", t_itab)
				itab_bin = append(itab_bin, header.Bin()...)
				for _, v := range vtable {
					fnid := NewConst(strconv.Itoa(v), m.U32)
					itab_bin = append(itab_bin, fnid.Bin()...)
				}

				addr = m.DataSeg.Append(itab_bin, 8)
			}

			itabs = append(itabs, NewConst(strconv.Itoa(addr), m.U32).Bin()...)
		}
	}

	itabs_ptr := m.DataSeg.Append(itabs, 8)
	m.SetGlobalInitValue("$wa.RT._itabsPtr", strconv.Itoa(itabs_ptr))
	m.SetGlobalInitValue("$wa.RT._interfaceCount", strconv.Itoa(len(m.usedInterfaces)))
	m.SetGlobalInitValue("$wa.RT._concretTypeCount", strconv.Itoa(len(m.usedConcreteTypes)))
}

func (m *Module) genTypeInfo(t ValueType) int {

	_type := NewConst("0", m.types_map["runtime._type"]).(*aStruct)
	_type.setFieldConstValue("size", NewConst(strconv.Itoa(t.Size()), m.U32))
	_type.setFieldConstValue("hash", NewConst(strconv.Itoa(t.Hash()), m.I32))
	_type.setFieldConstValue("kind", NewConst(strconv.Itoa(int(t.Kind())), m.U8))
	_type.setFieldConstValue("align", NewConst(strconv.Itoa(t.align()), m.U8))
	_type.setFieldConstValue("flag", NewConst("0", m.U16))
	_type.setFieldConstValue("name", NewConst(t.Name(), m.STRING))

	switch typ := t.(type) {
	case *tI8:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tU8:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tI16:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tU16:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tI32:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tU32:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tI64:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tU64:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tF32:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tF64:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *tRune:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *String:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Ptr:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Block:
		if typ.addr != 0 {
			return typ.addr
		}
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Array:
		if typ.addr != 0 {
			return typ.addr
		}
		_array := NewConst("0", m.types_map["runtime._arrayType"]).(*aStruct)
		_array.setFieldConstValue("$_type", _type)
		_array.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.genTypeInfo(typ.Base)), m.UPTR))
		_array.setFieldConstValue("cap", NewConst(strconv.Itoa(typ.Capacity), m.UPTR))
		typ.addr = m.DataSeg.Append(_array.Bin(), 8)
		return typ.addr

	case *Slice:
		if typ.addr != 0 {
			return typ.addr
		}
		_slice := NewConst("0", m.types_map["runtime._arrayType"]).(*aStruct)
		_slice.setFieldConstValue("$_type", _type)
		_slice.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.genTypeInfo(typ.Base)), m.UPTR))
		typ.addr = m.DataSeg.Append(_slice.Bin(), 8)
		return typ.addr

		/*
			case *Ref:
				if typ.addr != 0 {
					return typ.addr
				}
				_ref := NewConst("0", m.types_map["runtime._refType"]).(*aStruct)
				_ref.setFieldConstValue("$_type", _type)
				_ref.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.genTypeInfo(typ.Base)), m.UPTR))
				if len(typ.methods) > 0 {
					_uncommon := NewConst("0", m.types_map["runtime._uncommonType"]).(*aStruct)
					_uncommon.setFieldConstValue("methodCount", NewConst(strconv.Itoa(len(typ.methods)), m.U32))
					_uncommon_bin := _uncommon.Bin()
					for _, method := range typ.methods {
						_method := NewConst("0", m.types_map["runtime._method"]).(*aStruct)
						_method.setFieldConstValue("name", NewConst(method.Name, m.STRING))
					}

				}  //*/

	default:
		logger.Fatalf("Todo: %t", t)
		return 0
	}
}
