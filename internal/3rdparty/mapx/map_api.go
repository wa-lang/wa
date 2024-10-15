// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

func (t *mapImp) Len() uint { return t.count }

func (t *mapImp) Insert(k, v interface{}) {
	t.insert(&mapNode{t.NIL, t.NIL, t.NIL, mapRED, mapItem{k, v}})
}

func (t *mapImp) Delete(k interface{}) {
	t.delete(&mapNode{t.NIL, t.NIL, t.NIL, mapRED, mapItem{k: k}})
}

func (t *mapImp) Lookup(k interface{}) (interface{}, bool) {
	ret := t.search(&mapNode{t.NIL, t.NIL, t.NIL, mapRED, mapItem{k: k}})
	if ret == nil {
		return nil, false
	}

	return ret.mapItem.v, true
}
