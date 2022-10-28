package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/ssa"
)

/**************************************
Module:
**************************************/
type Module struct {
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

	m.funcs_map = make(map[string]*Function)

	//table中先行插入一条记录，防止产生0值（无效值）id
	m.table = append(m.table, "")
	m.table_map = make(map[string]int)

	//data_seg中先插入标志，防止产生0值
	m.data_seg = append(m.data_seg, []byte("$$wads$$")...)

	m.Globals_map = make(map[ssa.Value]Value)
	return &m
}

func (m *Module) findTableElem(elem string) int {
	if i, ok := m.table_map[elem]; ok {
		return i
	}
	return 0
}

func (m *Module) addTableFunc(f *Function) int {
	if i := m.findTableElem(f.Name); i != 0 {
		return i
	}

	m.AddFunc(f)
	i := len(m.table)
	m.table = append(m.table, f.Name)
	m.table_map[f.Name] = i
	return i
}

func (m *Module) findFunc(fn_name string) *Function {
	if f, ok := m.funcs_map[fn_name]; ok {
		return f
	}
	return nil
}

func (m *Module) AddFunc(f *Function) {
	if m.findFunc(f.Name) == nil {
		m.funcs = append(m.funcs, f)
		m.funcs_map[f.Name] = f
	}
}

func (m *Module) AddGlobal(name string, typ ValueType, is_pointer bool, ssa_value ssa.Value) Value {
	v := NewGlobal(name, typ, is_pointer)
	m.globals = append(m.globals, v)
	m.Globals_map[ssa_value] = v
	return v
}

func (m *Module) AddDataSeg(data []byte) (ptr int) {
	ptr = len(m.data_seg)
	m.data_seg = append(m.data_seg, data...)
	return
}

func (m *Module) genGlobalAlloc() *Function {
	var f Function
	f.Name = "$waGlobalAlloc"
	f.Result = VOID{}

	for _, g := range m.globals {
		if g.Kind() != ValueKindGlobal_Pointer {
			continue
		}

		t := g.Type().(Pointer).Base
		f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.size())))
		f.Insts = append(f.Insts, wat.NewInstCall("$waHeapAlloc"))
		f.Insts = append(f.Insts, g.EmitPop()...)
	}

	return &f
}

func (m *Module) ToWatModule() *wat.Module {
	var wat_module wat.Module
	wat_module.BaseWat = m.BaseWat

	{
		var onFreeFnType wat.FuncType
		onFreeFnType.Name = "$onFree"
		onFreeFnType.Params = append(onFreeFnType.Params, wat.I32{})
		wat_module.FuncTypes = append(wat_module.FuncTypes, onFreeFnType)
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
