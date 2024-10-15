// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

//
// Red-Black tree properties:  http://en.wikipedia.org/wiki/mapImp
//
//  1) A node is either red or black
//  2) The root is black
//  3) All leaves (NULL) are black
//  4) Both children of every red node are black
//  5) Every simple path from root to leaves contains the same number
//     of black nodes.
//

const (
	mapRED   = 0
	mapBLACK = 1
)

// mapImp represents a Red-Black tree.
type mapImp struct {
	NIL   *mapNode
	root  *mapNode
	count uint
}

type mapNode struct {
	Left   *mapNode
	Right  *mapNode
	Parent *mapNode
	Color  uint

	// for use by client.
	mapItem
}

type mapItem struct {
	k, v interface{}
}

func mapLess(x, y mapItem) bool {
	return Compare(x.k, y.k) < 0
}

func MakeMap() *mapImp {
	node := &mapNode{nil, nil, nil, mapBLACK, mapItem{nil, nil}}
	return &mapImp{
		NIL:   node,
		root:  node,
		count: 0,
	}
}

func (t *mapImp) leftRotate(x *mapNode) {
	// Since we are doing the left rotation, the right child should *NOT* nil.
	if x.Right == t.NIL {
		return
	}

	//
	// The illation of left rotation
	//
	//          |                                  |
	//          X                                  Y
	//         / \         left rotate            / \
	//        α  Y       ------------->         X   γ
	//           / \                            / \
	//          β  γ                         α  β
	//
	// It should be note that during the rotating we do not change
	// the Nodes' color.
	//
	y := x.Right
	x.Right = y.Left
	if y.Left != t.NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (t *mapImp) rightRotate(x *mapNode) {
	// Since we are doing the right rotation, the left child should *NOT* nil.
	if x.Left == t.NIL {
		return
	}

	//
	// The illation of right rotation
	//
	//          |                                  |
	//          X                                  Y
	//         / \         right rotate           / \
	//        Y   γ      ------------->         α  X
	//       / \                                    / \
	//      α  β                                 β  γ
	//
	// It should be note that during the rotating we do not change
	// the Nodes' color.
	//
	y := x.Left
	x.Left = y.Right
	if y.Right != t.NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}

func (t *mapImp) insert(z *mapNode) *mapNode {
	x := t.root
	y := t.NIL

	for x != t.NIL {
		y = x
		if mapLess(z.mapItem, x.mapItem) {
			x = x.Left
		} else if mapLess(x.mapItem, z.mapItem) {
			x = x.Right
		} else {
			return x
		}
	}

	z.Parent = y
	if y == t.NIL {
		t.root = z
	} else if mapLess(z.mapItem, y.mapItem) {
		y.Left = z
	} else {
		y.Right = z
	}

	t.count++
	t.insertFixup(z)
	return z
}

func (t *mapImp) insertFixup(z *mapNode) {
	for z.Parent.Color == mapRED {
		//
		// Howerver, we do not need the assertion of non-nil grandparent
		// because
		//
		//  2) The root is black
		//
		// Since the color of the parent is mapRED, so the parent is not root
		// and the grandparent must be exist.
		//
		if z.Parent == z.Parent.Parent.Left {
			// Take y as the uncle, although it can be NIL, in that case
			// its color is mapBLACK
			y := z.Parent.Parent.Right
			if y.Color == mapRED {
				//
				// Case 1:
				// Parent and uncle are both mapRED, the grandparent must be mapBLACK
				// due to
				//
				//  4) Both children of every red node are black
				//
				// Since the current node and its parent are all mapRED, we still
				// in violation of 4), So repaint both the parent and the uncle
				// to mapBLACK and grandparent to mapRED(to maintain 5)
				//
				//  5) Every simple path from root to leaves contains the same
				//     number of black nodes.
				//
				z.Parent.Color = mapBLACK
				y.Color = mapBLACK
				z.Parent.Parent.Color = mapRED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					//
					// Case 2:
					// Parent is mapRED and uncle is mapBLACK and the current node
					// is right child
					//
					// A left rotation on the parent of the current node will
					// switch the roles of each other. This still leaves us in
					// violation of 4).
					// The continuation into Case 3 will fix that.
					//
					z = z.Parent
					t.leftRotate(z)
				}
				//
				// Case 3:
				// Parent is mapRED and uncle is mapBLACK and the current node is
				// left child
				//
				// At the very beginning of Case 3, current node and parent are
				// both mapRED, thus we violate 4).
				// Repaint parent to mapBLACK will fix it, but 5) does not allow
				// this because all paths that go through the parent will get
				// 1 more black node. Then repaint grandparent to mapRED (as we
				// discussed before, the grandparent is mapBLACK) and do a right
				// rotation will fix that.
				//
				z.Parent.Color = mapBLACK
				z.Parent.Parent.Color = mapRED
				t.rightRotate(z.Parent.Parent)
			}
		} else { // same as then clause with "right" and "left" exchanged
			y := z.Parent.Parent.Left
			if y.Color == mapRED {
				z.Parent.Color = mapBLACK
				y.Color = mapBLACK
				z.Parent.Parent.Color = mapRED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rightRotate(z)
				}
				z.Parent.Color = mapBLACK
				z.Parent.Parent.Color = mapRED
				t.leftRotate(z.Parent.Parent)
			}
		}
	}
	t.root.Color = mapBLACK
}

// Just traverse the node from root to left recursively until left is NIL.
// The node whose left is NIL is the node with minimum value.
func (t *mapImp) min(x *mapNode) *mapNode {
	if x == t.NIL {
		return t.NIL
	}

	for x.Left != t.NIL {
		x = x.Left
	}

	return x
}

// Just traverse the node from root to right recursively until right is NIL.
// The node whose right is NIL is the node with maximum value.
func (t *mapImp) max(x *mapNode) *mapNode {
	if x == t.NIL {
		return t.NIL
	}

	for x.Right != t.NIL {
		x = x.Right
	}

	return x
}

func (t *mapImp) search(x *mapNode) *mapNode {
	p := t.root

	for p != t.NIL {
		if mapLess(p.mapItem, x.mapItem) {
			p = p.Right
		} else if mapLess(x.mapItem, p.mapItem) {
			p = p.Left
		} else {
			break
		}
	}

	return p
}

// TODO: Need Document
func (t *mapImp) successor(x *mapNode) *mapNode {
	if x == t.NIL {
		return t.NIL
	}

	// Get the minimum from the right sub-tree if it existed.
	if x.Right != t.NIL {
		return t.min(x.Right)
	}

	y := x.Parent
	for y != t.NIL && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

// TODO: Need Document
func (t *mapImp) delete(key *mapNode) *mapNode {
	z := t.search(key)

	if z == t.NIL {
		return t.NIL
	}
	ret := &mapNode{t.NIL, t.NIL, t.NIL, z.Color, z.mapItem}

	var y *mapNode
	var x *mapNode

	if z.Left == t.NIL || z.Right == t.NIL {
		y = z
	} else {
		y = t.successor(z)
	}

	if y.Left != t.NIL {
		x = y.Left
	} else {
		x = y.Right
	}

	// Even if x is NIL, we do the assign. In that case all the NIL nodes will
	// change from {nil, nil, nil, mapBLACK, nil} to {nil, nil, ADDR, mapBLACK, nil},
	// but do not worry about that because it will not affect the compare
	// between mapNode-X with mapNode-NIL
	x.Parent = y.Parent

	if y.Parent == t.NIL {
		t.root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}

	if y != z {
		z.mapItem = y.mapItem
	}

	if y.Color == mapBLACK {
		t.deleteFixup(x)
	}

	t.count--

	return ret
}

func (t *mapImp) deleteFixup(x *mapNode) {
	for x != t.root && x.Color == mapBLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == mapRED {
				w.Color = mapBLACK
				x.Parent.Color = mapRED
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == mapBLACK && w.Right.Color == mapBLACK {
				w.Color = mapRED
				x = x.Parent
			} else {
				if w.Right.Color == mapBLACK {
					w.Left.Color = mapBLACK
					w.Color = mapRED
					t.rightRotate(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = mapBLACK
				w.Right.Color = mapBLACK
				t.leftRotate(x.Parent)
				// this is to exit while loop
				x = t.root
			}
		} else { // the code below is has left and right switched from above
			w := x.Parent.Left
			if w.Color == mapRED {
				w.Color = mapBLACK
				x.Parent.Color = mapRED
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Left.Color == mapBLACK && w.Right.Color == mapBLACK {
				w.Color = mapRED
				x = x.Parent
			} else {
				if w.Left.Color == mapBLACK {
					w.Right.Color = mapBLACK
					w.Color = mapRED
					t.leftRotate(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = mapBLACK
				w.Left.Color = mapBLACK
				t.rightRotate(x.Parent)
				x = t.root
			}
		}
	}
	x.Color = mapBLACK
}
