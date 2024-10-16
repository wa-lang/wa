// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

type mapIter struct {
	m   *mapImp    // map
	v   int64      // map version
	stk []*mapNode // stack
}

func MakeMapIter(m *mapImp) *mapIter {
	return &mapIter{
		m:   m,
		v:   m.version,
		stk: []*mapNode{m.root},
	}
}

func (this *mapIter) HasNext() (ok bool) {
	if this.v != this.m.version {
		return
	}
	if len(this.stk) == 0 {
		return
	}

	h := this.stk[len(this.stk)-1]
	if h == nil || h == this.m.NIL {
		return
	}

	return true
}

func (this *mapIter) KeyValue() (k, v interface{}) {
	if this.v != this.m.version {
		return
	}
	if len(this.stk) == 0 {
		return
	}

	h := this.stk[len(this.stk)-1]
	if h == nil || h == this.m.NIL {
		return
	}

	k = h.mapItem.k
	v = h.mapItem.v
	return
}

func (this *mapIter) Next() (ok bool, k, v interface{}) {
	if this.v != this.m.version {
		return
	}
	if len(this.stk) == 0 {
		return
	}

	h := this.stk[len(this.stk)-1]
	this.stk = this.stk[:len(this.stk)-1]

	if h == nil || h == this.m.NIL {
		return
	}

	ok = true
	k = h.mapItem.k
	v = h.mapItem.v

	if h.Right != nil && h.Right != this.m.NIL {
		this.stk = append(this.stk, h.Right)
	}

	if h.Left != nil && h.Left != this.m.NIL {
		this.stk = append(this.stk, h.Left)
	}

	return
}
