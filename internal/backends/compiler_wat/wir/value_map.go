// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Map:
**************************************/
type Map struct {
	tCommon
	Key        ValueType
	Elem       ValueType
	underlying *Ref

	updateFnName, lookatFnName, lookatCommaokFnName, delFnName string
}

func (module *Module) GenValueType_Map(key_type, elem_type ValueType, name string, ei_type ValueType) *Map {
	map_t := Map{Key: key_type, Elem: elem_type}
	if len(name) > 0 {
		map_t.name = name
	} else {
		map_t.name = "runtime.map." + key_type.Named() + "." + elem_type.Named()
	}

	t, ok := module.findValueType(map_t.name)
	if ok {
		return t.(*Map)
	}

	mit, _ := module.findValueType("runtime.mapImp")
	map_t.underlying = module.GenValueType_Ref(mit)

	_, k_is_iface := key_type.(*Interface)
	_, v_is_iface := elem_type.(*Interface)

	// update
	{
		var f Function
		m := NewLocal("m", map_t.underlying)
		k := NewLocal("k", key_type)
		v := NewLocal("v", elem_type)
		f.Params = append(f.Params, m)
		f.Params = append(f.Params, k)
		f.Params = append(f.Params, v)

		ki := NewLocal("ki", ei_type)
		vi := NewLocal("vi", ei_type)
		f.Locals = append(f.Locals, ki)
		f.Locals = append(f.Locals, vi)

		if k_is_iface {
			f.Insts = append(f.Insts, module.EmitGenChangeInterface(k, ei_type)...)
		} else {
			f.Insts = append(f.Insts, module.EmitGenMakeInterface(k, ei_type)...)
		}
		f.Insts = append(f.Insts, ki.EmitPop()...)

		if v_is_iface {
			f.Insts = append(f.Insts, module.EmitGenChangeInterface(v, ei_type)...)
		} else {
			f.Insts = append(f.Insts, module.EmitGenMakeInterface(v, ei_type)...)
		}
		f.Insts = append(f.Insts, vi.EmitPop()...)

		f.Insts = append(f.Insts, m.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, ki.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, vi.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, wat.NewInstCall("runtime.mapUpdate"))

		f.Insts = append(f.Insts, vi.EmitRelease()...)
		f.Insts = append(f.Insts, ki.EmitRelease()...)

		f.InternalName = map_t.name + ".update"
		map_t.updateFnName = f.InternalName

		module.AddFunc(&f)
	}

	// lookup
	{
		var f Function
		m := NewLocal("m", map_t.underlying)
		k := NewLocal("k", key_type)
		f.Params = append(f.Params, m)
		f.Params = append(f.Params, k)
		f.Results = []ValueType{elem_type}

		ki := NewLocal("ki", ei_type)
		f.Locals = append(f.Locals, ki)
		vi := NewLocal("vi", ei_type)
		f.Locals = append(f.Locals, vi)
		ok := NewLocal("ok", module.BOOL)
		f.Locals = append(f.Locals, ok)

		if k_is_iface {
			f.Insts = append(f.Insts, module.EmitGenChangeInterface(k, ei_type)...)
		} else {
			f.Insts = append(f.Insts, module.EmitGenMakeInterface(k, ei_type)...)
		}
		f.Insts = append(f.Insts, ki.EmitPop()...)

		f.Insts = append(f.Insts, m.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, ki.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, wat.NewInstCall("runtime.mapLookup"))
		f.Insts = append(f.Insts, ok.EmitPop()...)
		f.Insts = append(f.Insts, vi.EmitPop()...)

		{
			f.Insts = append(f.Insts, ok.EmitPushNoRetain()...)
			if_found := wat.NewInstIf(nil, nil, nil)
			if_found.Ret = elem_type.Raw()

			zv := NewLocal("v", elem_type)
			f.Locals = append(f.Locals, zv)
			if_found.False = append(if_found.False, zv.EmitPush()...)

			if v_is_iface {
				if_found.True = append(if_found.True, module.EmitGenChangeInterface(vi, elem_type)...)
			} else {
				if_found.True = append(if_found.True, vi.(*aInterface).emitGetData(elem_type, false)...)
			}

			f.Insts = append(f.Insts, if_found)
		}

		f.Insts = append(f.Insts, vi.EmitRelease()...)
		f.Insts = append(f.Insts, ki.EmitRelease()...)

		f.InternalName = map_t.name + ".lookup"
		map_t.lookatFnName = f.InternalName

		module.AddFunc(&f)
	}

	// lookupCommaOk
	{
		var f Function
		m := NewLocal("m", map_t.underlying)
		k := NewLocal("k", key_type)
		f.Params = append(f.Params, m)
		f.Params = append(f.Params, k)
		f.Results = []ValueType{elem_type, module.BOOL}

		ki := NewLocal("ki", ei_type)
		vi := NewLocal("vi", ei_type)
		ok := NewLocal("ok", module.BOOL)
		f.Locals = append(f.Locals, ki)
		f.Locals = append(f.Locals, vi)
		f.Locals = append(f.Locals, ok)

		if k_is_iface {
			f.Insts = append(f.Insts, module.EmitGenChangeInterface(k, ei_type)...)
		} else {
			f.Insts = append(f.Insts, module.EmitGenMakeInterface(k, ei_type)...)
		}
		f.Insts = append(f.Insts, ki.EmitPop()...)

		f.Insts = append(f.Insts, m.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, ki.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, wat.NewInstCall("runtime.mapLookup"))
		f.Insts = append(f.Insts, ok.EmitPop()...)
		f.Insts = append(f.Insts, vi.EmitPop()...)

		{
			f.Insts = append(f.Insts, ok.EmitPushNoRetain()...)
			if_found := wat.NewInstIf(nil, nil, nil)
			if_found.Ret = elem_type.Raw()

			zv := NewLocal("zv", elem_type)
			f.Locals = append(f.Locals, zv)
			if_found.False = append(if_found.False, zv.EmitPush()...)

			if v_is_iface {
				if_found.True = append(if_found.True, module.EmitGenChangeInterface(vi, elem_type)...)
			} else {
				if_found.True = append(if_found.True, vi.(*aInterface).emitGetData(elem_type, false)...)
			}

			f.Insts = append(f.Insts, if_found)
		}

		f.Insts = append(f.Insts, ok.EmitPush()...)

		f.Insts = append(f.Insts, vi.EmitRelease()...)
		f.Insts = append(f.Insts, ki.EmitRelease()...)

		f.InternalName = map_t.name + ".lookupcommaok"
		map_t.lookatCommaokFnName = f.InternalName

		module.AddFunc(&f)
	}

	// next
	{
		var f Function
		iter_type, _ := module.findValueType("runtime.mapIter")
		iter := NewLocal("iter", iter_type)
		f.Params = append(f.Params, iter)
		f.Results = []ValueType{module.BOOL, key_type, elem_type, module.INT}

		ok := NewLocal("ok", module.BOOL)
		ki := NewLocal("ki", ei_type)
		vi := NewLocal("vi", ei_type)
		next_pos := NewLocal("np", module.INT)
		f.Locals = []Value{ok, ki, vi, next_pos}

		f.Insts = append(f.Insts, iter.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, wat.NewInstCall("runtime.mapNext"))
		f.Insts = append(f.Insts, next_pos.EmitPop()...)
		f.Insts = append(f.Insts, vi.EmitPop()...)
		f.Insts = append(f.Insts, ki.EmitPop()...)
		f.Insts = append(f.Insts, ok.EmitPop()...)

		f.Insts = append(f.Insts, ok.EmitPush()...)
		{
			f.Insts = append(f.Insts, ok.EmitPushNoRetain()...)
			ifok := wat.NewInstIf(nil, nil, nil)
			ifok.Ret = append(ifok.Ret, key_type.Raw()...)
			ifok.Ret = append(ifok.Ret, elem_type.Raw()...)

			zk := NewLocal("zk", key_type)
			zv := NewLocal("zv", elem_type)
			f.Locals = append(f.Locals, zk)
			f.Locals = append(f.Locals, zv)
			ifok.False = append(ifok.False, zk.EmitPush()...)
			ifok.False = append(ifok.False, zv.EmitPush()...)

			if k_is_iface {
				ifok.True = append(ifok.True, module.EmitGenChangeInterface(ki, key_type)...)
			} else {
				ifok.True = append(ifok.True, ki.(*aInterface).emitGetData(key_type, false)...)
			}

			if v_is_iface {
				ifok.True = append(ifok.True, module.EmitGenChangeInterface(vi, elem_type)...)
			} else {
				ifok.True = append(ifok.True, vi.(*aInterface).emitGetData(elem_type, false)...)
			}

			f.Insts = append(f.Insts, ifok)
		}
		f.Insts = append(f.Insts, next_pos.EmitPush()...)

		f.Insts = append(f.Insts, vi.EmitRelease()...)
		f.Insts = append(f.Insts, ki.EmitRelease()...)

		f.InternalName = "runtime.map." + key_type.Named() + "." + elem_type.Named() + ".next"
		module.AddFunc(&f)
	}

	// delete
	{
		var f Function
		m := NewLocal("m", map_t.underlying)
		k := NewLocal("k", key_type)
		f.Params = append(f.Params, m)
		f.Params = append(f.Params, k)

		ki := NewLocal("ki", ei_type)
		f.Locals = append(f.Locals, ki)

		if k_is_iface {
			f.Insts = append(f.Insts, module.EmitGenChangeInterface(k, ei_type)...)
		} else {
			f.Insts = append(f.Insts, module.EmitGenMakeInterface(k, ei_type)...)
		}
		f.Insts = append(f.Insts, ki.EmitPop()...)

		f.Insts = append(f.Insts, m.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, ki.EmitPushNoRetain()...)
		f.Insts = append(f.Insts, wat.NewInstCall("runtime.mapDelete"))

		f.Insts = append(f.Insts, ki.EmitRelease()...)

		f.InternalName = map_t.name + ".delete"
		map_t.delFnName = f.InternalName

		module.AddFunc(&f)
	}

	module.addValueType(&map_t)
	return &map_t
}

func (t *Map) Size() int            { return t.underlying.Size() }
func (t *Map) align() int           { return t.underlying.align() }
func (t *Map) Kind() TypeKind       { return kMap }
func (t *Map) OnFree() int          { return t.underlying.OnFree() }
func (t *Map) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Map) Equal(u ValueType) bool {
	if ut, ok := u.(*Map); ok {
		return t.Key.Equal(ut.Key) && t.Elem.Equal(ut.Elem)
	}
	return false
}

func (t *Map) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Map) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

func (t *Map) emitGenMake() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstCall("runtime.mapMake"))
	return
}

/**************************************
aMap:
**************************************/
type aMap struct {
	aRef
	typ *Map
}

func newValue_Map(name string, kind ValueKind, typ *Map) *aMap {
	var v aMap
	v.typ = typ
	v.aRef = *newValue_Ref(name, kind, typ.underlying)
	return &v
}

func (v *aMap) Type() ValueType { return v.typ }

func (v *aMap) emitEq(r Value) ([]wat.Inst, bool) {
	if r.Kind() != ValueKindConst {
		return nil, false
	}

	r_s, ok := r.(*aMap)
	if !ok {
		logger.Fatal("r is not a Map")
	}

	if r_s.Name() != "0" {
		logger.Fatal("r should be nil")
	}

	return v.aRef.emitEq(&r_s.aRef)
}

func (v *aMap) emitCompare(r Value) (insts []wat.Inst) {
	if !r.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	t := r.(*aMap)

	return v.aRef.emitCompare(&t.aRef)
}

func (v *aMap) emitUpdate(key, elem Value) (insts []wat.Inst) {
	insts = append(insts, v.aRef.EmitPushNoRetain()...)
	insts = append(insts, key.EmitPushNoRetain()...)
	insts = append(insts, elem.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall(v.typ.updateFnName))

	return
}

func (v *aMap) emitLookup(key Value, CommaOk bool) (insts []wat.Inst) {
	insts = append(insts, v.aRef.EmitPushNoRetain()...)
	insts = append(insts, key.EmitPushNoRetain()...)

	if CommaOk {
		insts = append(insts, wat.NewInstCall(v.typ.lookatCommaokFnName))
	} else {
		insts = append(insts, wat.NewInstCall(v.typ.lookatFnName))
	}

	return
}

func (v *aMap) emitLen() (insts []wat.Inst) {
	insts = append(insts, v.aRef.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("runtime.mapLen"))

	return
}

func (v *aMap) emitDelete(key Value) (insts []wat.Inst) {
	insts = append(insts, v.aRef.EmitPushNoRetain()...)
	insts = append(insts, key.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall(v.typ.delFnName))

	return
}
