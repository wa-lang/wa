// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

type mapIter struct {
	m   *mapImp
	pos int
}

func MakeMapIter(m *mapImp) *mapIter {
	return &mapIter{m: m}
}

func (this *mapIter) Next() (ok bool, k, v interface{}) {
	if this.pos >= this.m.Len() {
		return
	}

	this.pos++
	node := this.m.nodes[this.pos]

	ok = true
	k = node.Key
	v = node.Val

	return
}
