package wir

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
)

type fnSigWrap struct {
	name     string
	typeAddr int
}

/*
*************************************
Module:
*************************************
*/
type Module struct {
	VOID, BOOL, RUNE, U8, U16, I32, U32, UPTR, I64, U64, INT, UINT, F32, F64, COMPLEX64, COMPLEX128, STRING, BYTES ValueType

	stkSize int // 栈大小

	types_map         map[string]ValueType
	usedConcreteTypes []ValueType
	usedInterfaces    []ValueType

	imports []wat.Import

	fnSigs     []*FnSig
	fnSigsName map[string]fnSigWrap

	Funcs     []*Function
	funcs_map map[string]*Function

	table     []string
	table_map map[string]int

	Globals           []Global
	globalsMapByValue map[ssa.Value]int
	globalsMapByName  map[string]int

	DataSeg *wat.DataSeg

	BaseWat string
}

func NewModule(stkSize int) *Module {
	var m Module
	m.types_map = make(map[string]ValueType)
	m.fnSigsName = make(map[string]fnSigWrap)
	m.funcs_map = make(map[string]*Function)
	m.table_map = make(map[string]int)
	m.globalsMapByValue = make(map[ssa.Value]int)
	m.globalsMapByName = make(map[string]int)

	m.VOID = m.GenValueType_void("")
	m.BOOL = m.GenValueType_bool("")
	m.RUNE = m.GenValueType_rune("")
	m.U8 = m.GenValueType_u8("")
	m.U16 = m.GenValueType_u16("")
	m.I32 = m.GenValueType_i32("")
	m.U32 = m.GenValueType_u32("")
	m.I64 = m.GenValueType_i64("")
	m.U64 = m.GenValueType_u64("")
	m.F32 = m.GenValueType_f32("")
	m.F64 = m.GenValueType_f64("")
	m.COMPLEX64 = m.GenValueType_complex64("")
	m.COMPLEX128 = m.GenValueType_complex128("")
	m.INT = m.I32
	m.UINT = m.U32
	m.UPTR = m.U32

	m.STRING = m.GenValueType_string("")
	m.BYTES = m.GenValueType_Slice(m.U8, "")

	m.stkSize = stkSize

	//table中先行插入一条记录，防止产生0值（无效值）id
	m.table = append(m.table, "")

	//data_seg中先插入标志，防止产生0值
	m.DataSeg = wat.NewDataSeg(stkSize + 32)
	m.DataSeg.Append([]byte("$$wads$$"), 8)

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
		return s.name
	}
	return ""
}

func (m *Module) AddFnSig(sig *FnSig) string {
	if s, ok := m.fnSigsName[sig.String()]; ok {
		return s.name
	}

	m.fnSigs = append(m.fnSigs, sig)
	s := "$$fnSig" + strconv.Itoa(len(m.fnSigs))

	m.fnSigsName[sig.String()] = fnSigWrap{name: s}
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
		m.Funcs = append(m.Funcs, f)
		m.funcs_map[f.InternalName] = f
	}
}

func (m *Module) AddGlobal(name_internal string, name_export string, typ ValueType, is_pointer bool, ssa_value ssa.Value) Value {
	var v Global

	v.Name = name_internal
	v.Name_exp = name_export
	v.Type = typ
	if is_pointer {
		t_ref, ok := typ.(*Ref)
		if !ok {
			logger.Fatal("typ should be *Ref")
		}
		gptr := m.DataSeg.Alloc(t_ref.Base.Size(), t_ref.Base.align())
		v.val = t_ref.newConstRef(gptr)
	} else {
		v.val = NewGlobal(name_internal, typ)
	}

	if ssa_value != nil {
		m.globalsMapByValue[ssa_value] = len(m.Globals)
	}
	m.globalsMapByName[name_internal] = len(m.Globals)
	m.Globals = append(m.Globals, v)
	return v.val
}

func (m *Module) FindGlobalByName(name string) Value {
	id, ok := m.globalsMapByName[name]
	if !ok {
		return nil
	}

	return m.Globals[id].val
}

func (m *Module) FindGlobalByValue(v ssa.Value) Value {
	id, ok := m.globalsMapByValue[v]
	if !ok {
		return nil
	}

	return m.Globals[id].val
}

func (m *Module) SetGlobalInitValue(name string, val Value) {
	id, ok := m.globalsMapByName[name]
	if !ok {
		logger.Fatalf("Global not found:%s", name)
	}

	m.Globals[id].init_val = val
}

func (m *Module) ToWatModule() *wat.Module {
	m.buildItab()
	m.buildTypesInfo()

	var wat_module wat.Module
	wat_module.Imports = m.imports
	wat_module.BaseWat = m.BaseWat

	{
		var onfree_type wat.FuncType
		onfree_type.Name = "$OnFree"
		onfree_type.Params = m.I32.Raw()
		wat_module.FuncTypes = append(wat_module.FuncTypes, onfree_type)
	}

	{
		var comp_type wat.FuncType
		comp_type.Name = "$wa.runtime.comp"
		ptr_t := m.GenValueType_Ptr(m.VOID)
		comp_type.Params = append(comp_type.Params, ptr_t.Raw()...)
		comp_type.Params = append(comp_type.Params, ptr_t.Raw()...)
		comp_type.Results = append(comp_type.Results, m.I32.Raw()...)
		wat_module.FuncTypes = append(wat_module.FuncTypes, comp_type)
	}

	for _, t := range m.fnSigs {
		var fn_type wat.FuncType
		fn_type.Name = m.fnSigsName[t.String()].name
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

	for _, f := range m.Funcs {
		wat_module.Funcs = append(wat_module.Funcs, f.ToWatFunc())
	}

	for _, g := range m.Globals {
		if g.val.Kind() == ValueKindConst {
			g_v := newValue(g.Name, ValueKindGlobal, g.val.Type())
			raw_v := g_v.raw()
			raw_c := g.val.raw()
			var raw_e []wat.Value
			if len(g.Name_exp) > 0 {
				g_v_exp := newValue(g.Name_exp, ValueKindGlobal, g.val.Type())
				raw_e = g_v_exp.raw()
			}
			for i, r := range raw_v {
				var wat_global wat.Global
				wat_global.V = r
				wat_global.IsMut = false
				wat_global.InitValue = raw_c[i].Name()
				if i < len(raw_e) {
					wat_global.NameExp = raw_e[i].Name()
				}
				wat_module.Globals = append(wat_module.Globals, wat_global)
			}
		} else if g.init_val != nil {
			raw_v := g.val.raw()
			raw_c := g.init_val.raw()
			var raw_e []wat.Value
			if len(g.Name_exp) > 0 {
				g_v_exp := newValue(g.Name_exp, ValueKindGlobal, g.val.Type())
				raw_e = g_v_exp.raw()
			}
			for i, r := range raw_v {
				var wat_global wat.Global
				wat_global.V = r
				wat_global.IsMut = true
				wat_global.InitValue = raw_c[i].Name()
				if i < len(raw_e) {
					wat_global.NameExp = raw_e[i].Name()
				}
				wat_module.Globals = append(wat_module.Globals, wat_global)
			}
		} else {
			raw_v := g.val.raw()
			var raw_e []wat.Value
			if len(g.Name_exp) > 0 {
				g_v_exp := newValue(g.Name_exp, ValueKindGlobal, g.val.Type())
				raw_e = g_v_exp.raw()
			}
			for i, r := range raw_v {
				var wat_global wat.Global
				wat_global.V = r
				wat_global.IsMut = true
				if i < len(raw_e) {
					wat_global.NameExp = raw_e[i].Name()
				}
				wat_module.Globals = append(wat_module.Globals, wat_global)
			}
		}
	}

	{
		// todo:
		var heap_base wat.Global
		heap_base.V = wat.NewVar("__heap_base", wat.I32{})
		heap_base.IsMut = false
		heap_base.InitValue = strconv.Itoa(makeAlign(m.DataSeg.Size()+16, 16))
		wat_module.Globals = append(wat_module.Globals, heap_base)
	}

	wat_module.DataSeg = m.DataSeg

	return &wat_module
}

func (m *Module) addValueType(t ValueType) {
	_, ok := m.types_map[t.Named()]
	if ok {
		logger.Fatalf("ValueType:%T already registered.", t)
	}
	m.types_map[t.Named()] = t
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

	for _, concrete := range m.usedConcreteTypes {
		for _, iface := range m.usedInterfaces {
			fits := true

			vtable := make([]int, iface.NumMethods())

			for mid := 0; mid < iface.NumMethods(); mid++ {
				method := iface.Method(mid)
				found := false
				for fid := 0; fid < concrete.NumMethods(); fid++ {
					d := concrete.Method(fid)
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
				header := NewConst("0", t_itab).(*aStruct)
				header.setFieldConstValue("dhash", NewConst(strconv.Itoa(concrete.Hash()), m.I32))
				header.setFieldConstValue("ihash", NewConst(strconv.Itoa(iface.Hash()), m.I32))
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
	m.SetGlobalInitValue("$wa.runtime._itabsPtr", NewConst(strconv.Itoa(itabs_ptr), m.FindGlobalByName("$wa.runtime._itabsPtr").Type()))
	m.SetGlobalInitValue("$wa.runtime._interfaceCount", NewConst(strconv.Itoa(len(m.usedInterfaces)), m.FindGlobalByName("$wa.runtime._interfaceCount").Type()))
	m.SetGlobalInitValue("$wa.runtime._concretTypeCount", NewConst(strconv.Itoa(len(m.usedConcreteTypes)), m.FindGlobalByName("$wa.runtime._concretTypeCount").Type()))
}

func (m *Module) buildTypesInfo() {
	return
	for name, t := range m.types_map {
		if strings.HasPrefix(name, "runtime.") {
			continue
		}
		m.buildTypeInfo(t)
	}
}

func (m *Module) buildTypeInfo(t ValueType) int {
	if t.typeInfoAddr() != 0 {
		return t.typeInfoAddr()
	}

	_type := NewConst("0", m.types_map["runtime._type"]).(*aStruct)
	_type.setFieldConstValue("size", NewConst(strconv.Itoa(t.Size()), m.U32))
	_type.setFieldConstValue("hash", NewConst(strconv.Itoa(t.Hash()), m.I32))
	_type.setFieldConstValue("kind", NewConst(strconv.Itoa(int(t.Kind())), m.U8))
	_type.setFieldConstValue("align", NewConst(strconv.Itoa(t.align()), m.U8))
	_type.setFieldConstValue("flag", NewConst("0", m.U16))
	_type.setFieldConstValue("name", NewConst(t.Named(), m.STRING))

	switch typ := t.(type) {
	case *Void:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *I8:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *U8:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *I16:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *U16:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *I32:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *U32:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *I64:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *U64:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *F32:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *F64:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Rune:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *String:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Ptr:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Block:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Array:
		_array := NewConst("0", m.types_map["runtime._arrayType"]).(*aStruct)
		_array.setFieldConstValue("$_type", _type)
		_array.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.buildTypeInfo(typ.Base)), m.UPTR))
		_array.setFieldConstValue("cap", NewConst(strconv.Itoa(typ.Capacity), m.UPTR))
		typ.addr = m.DataSeg.Append(_array.Bin(), 8)
		return typ.addr

	case *Slice:
		_slice := NewConst("0", m.types_map["runtime._arrayType"]).(*aStruct)
		_slice.setFieldConstValue("$_type", _type)
		_slice.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.buildTypeInfo(typ.Base)), m.UPTR))
		typ.addr = m.DataSeg.Append(_slice.Bin(), 8)
		return typ.addr

	case *Ref:
		_sptr := NewConst("0", m.types_map["runtime._sptrType"]).(*aStruct)
		typ.addr = m.DataSeg.Alloc(len(_sptr.Bin()), 8)

		_sptr.setFieldConstValue("$_type", _type)
		_sptr.setFieldConstValue("elemType", NewConst(strconv.Itoa(m.buildTypeInfo(typ.Base)), m.UPTR))
		if len(typ.methods) > 0 {
			_uncommon := NewConst("0", m.types_map["runtime._uncommonType"]).(*aStruct)
			_uncommon.setFieldConstValue("methodCount", NewConst(strconv.Itoa(len(typ.methods)), m.U32))
			_uncommon_bin := _uncommon.Bin()
			for _, method := range typ.methods {
				_method := NewConst("0", m.types_map["runtime._method"]).(*aStruct)
				_method.setFieldConstValue("name", NewConst(method.Name, m.STRING))
				_method.setFieldConstValue("fnType", NewConst(strconv.Itoa(m.buildFnTypeInfo(&method.Sig)), m.UPTR))
				_method.setFieldConstValue("fnID", NewConst(strconv.Itoa(m.AddTableElem(method.FullFnName)), m.U32))
				_uncommon_bin = append(_uncommon_bin, _method.Bin()...)
			}
			_sptr.setFieldConstValue("uncommon", NewConst(strconv.Itoa(m.DataSeg.Append(_uncommon_bin, 8)), m.UPTR))
		}

		m.DataSeg.Set(_sptr.Bin(), typ.addr)
		return typ.addr

	case *Closure:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Tuple:
		typ.addr = m.DataSeg.Append(_type.Bin(), 8)
		return typ.addr

	case *Struct:
		_struct := NewConst("0", m.types_map["runtime._structType"]).(*aStruct)
		_structField := NewConst("0", m.types_map["runtime._structField"]).(*aStruct)
		typ.addr = m.DataSeg.Alloc(len(_struct.Bin())+len(_structField.Bin())*len(typ.fields), 8)

		_struct.setFieldConstValue("$_type", _type)
		_struct.setFieldConstValue("fieldCount", NewConst(strconv.Itoa(len(typ.fields)), m.U32))
		_struct_bin := _struct.Bin()
		for _, f := range typ.fields {
			_structField.setFieldConstValue("name", NewConst(f.Name(), m.STRING))
			_structField.setFieldConstValue("typ", NewConst(strconv.Itoa(m.buildTypeInfo(f.Type())), m.UPTR))
			_struct_bin = append(_struct_bin, _structField.Bin()...)
		}

		m.DataSeg.Set(_struct_bin, typ.addr)
		return typ.addr

	case *Interface:
		_interface := NewConst("0", m.types_map["runtime._interfaceType"]).(*aStruct)
		_imethod := NewConst("0", m.types_map["runtime._imethod"]).(*aStruct)
		typ.addr = m.DataSeg.Alloc(len(_interface.Bin())+len(_imethod.Bin())*typ.NumMethods(), 8)

		_interface.setFieldConstValue("methodCount", NewConst(strconv.Itoa(typ.NumMethods()), m.U32))
		_interface_bin := _interface.Bin()
		for _, method := range typ.methods {
			_imethod.setFieldConstValue("name", NewConst(method.Name, m.STRING))
			_imethod.setFieldConstValue("fnType", NewConst(strconv.Itoa(m.buildFnTypeInfo(&method.Sig)), m.UPTR))
			_interface_bin = append(_interface_bin, _imethod.Bin()...)
		}

		m.DataSeg.Set(_interface_bin, typ.addr)
		return typ.addr

	default:
		logger.Fatalf("Todo: %t", t)
		return 0
	}
}

func (m *Module) buildFnTypeInfo(sig *FnSig) int {
	s, ok := m.fnSigsName[sig.String()]
	if ok && s.typeAddr != 0 {
		return s.typeAddr
	}

	_fnType := NewConst("0", m.types_map["runtime._fnType"]).(*aStruct)
	_fnType.setFieldConstValue("paramCount", NewConst(strconv.Itoa(len(sig.Params)), m.U32))
	_fnType.setFieldConstValue("resultCount", NewConst(strconv.Itoa(len(sig.Results)), m.U32))
	_fnType_bin := _fnType.Bin()

	for _, p := range sig.Params {
		typaddr := m.buildTypeInfo(p)
		typaddr_bin := NewConst(strconv.Itoa(typaddr), m.UPTR).Bin()
		_fnType_bin = append(_fnType_bin, typaddr_bin...)
	}

	for _, p := range sig.Results {
		typaddr := m.buildTypeInfo(p)
		typaddr_bin := NewConst(strconv.Itoa(typaddr), m.UPTR).Bin()
		_fnType_bin = append(_fnType_bin, typaddr_bin...)
	}

	s.typeAddr = m.DataSeg.Append(_fnType_bin, 8)
	m.fnSigsName[sig.String()] = s
	return s.typeAddr
}
