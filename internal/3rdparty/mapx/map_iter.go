// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

import "unsafe"

type mapIter struct {
	m   *mapImp
	pos int // node ptr
}

func mapNext(iter mapIter) (ok bool, k, v interface{}, pos int) {
	ptr := (*mapNode)(unsafe.Pointer(uintptr(iter.pos)))
	if ptr == nil {
		// todo: 第一次调用，初始化
	}
	if ptr == iter.m.NIL {
		// todo: 已经遍历完了
		return false, nil, nil, 0
	}

	panic("TODO")
}

func (t *mapImp) walk(n *mapNode, f func(k, v interface{})) {
	if n == t.NIL {
		return
	}
	t.walk(n.Left, f)
	f(n.k, n.v)
	t.walk(n.Right, f)
}
