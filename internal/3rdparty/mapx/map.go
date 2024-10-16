// 版权 @2024 凹语言 作者。保留所有权利。

package mapx

//
// Red-Black tree properties:  http://en.wikipedia.org/wiki/rbtree
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

type mapImp struct {
	NIL     *mapNode
	root    *mapNode
	count   uint
	version int64
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

func (this *mapImp) Len() uint { return this.count }

func (this *mapImp) Update(k, v interface{}) {
	this.version++
	this.insert(&mapNode{this.NIL, this.NIL, this.NIL, mapRED, mapItem{k, v}})
}

func (this *mapImp) Lookup(k interface{}) (interface{}, bool) {
	ret := this.search(&mapNode{this.NIL, this.NIL, this.NIL, mapRED, mapItem{k: k}})
	if ret == nil {
		return nil, false
	}

	return ret.mapItem.v, true
}

func (this *mapImp) Delete(k interface{}) {
	this.version++
	this.delete(&mapNode{this.NIL, this.NIL, this.NIL, mapRED, mapItem{k: k}})
}

func (this *mapImp) leftRotate(x *mapNode) {
	// Since we are doing the left rotation, the right child should *NOT* nil.
	if x.Right == this.NIL {
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
	if y.Left != this.NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == this.NIL {
		this.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (this *mapImp) rightRotate(x *mapNode) {
	// Since we are doing the right rotation, the left child should *NOT* nil.
	if x.Left == this.NIL {
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
	if y.Right != this.NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == this.NIL {
		this.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}

func (this *mapImp) insert(z *mapNode) *mapNode {
	x := this.root
	y := this.NIL

	for x != this.NIL {
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
	if y == this.NIL {
		this.root = z
	} else if mapLess(z.mapItem, y.mapItem) {
		y.Left = z
	} else {
		y.Right = z
	}

	this.count++
	this.insertFixup(z)
	return z
}

func (this *mapImp) insertFixup(z *mapNode) {
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
					this.leftRotate(z)
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
				this.rightRotate(z.Parent.Parent)
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
					this.rightRotate(z)
				}
				z.Parent.Color = mapBLACK
				z.Parent.Parent.Color = mapRED
				this.leftRotate(z.Parent.Parent)
			}
		}
	}
	this.root.Color = mapBLACK
}

// Just traverse the node from root to left recursively until left is NIL.
// The node whose left is NIL is the node with minimum value.
func (this *mapImp) min(x *mapNode) *mapNode {
	if x == this.NIL {
		return this.NIL
	}

	for x.Left != this.NIL {
		x = x.Left
	}

	return x
}

// Just traverse the node from root to right recursively until right is NIL.
// The node whose right is NIL is the node with maximum value.
func (this *mapImp) max(x *mapNode) *mapNode {
	if x == this.NIL {
		return this.NIL
	}

	for x.Right != this.NIL {
		x = x.Right
	}

	return x
}

func (this *mapImp) search(x *mapNode) *mapNode {
	p := this.root

	for p != this.NIL {
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
func (this *mapImp) successor(x *mapNode) *mapNode {
	if x == this.NIL {
		return this.NIL
	}

	// Get the minimum from the right sub-tree if it existed.
	if x.Right != this.NIL {
		return this.min(x.Right)
	}

	y := x.Parent
	for y != this.NIL && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

// TODO: Need Document
func (this *mapImp) delete(key *mapNode) *mapNode {
	z := this.search(key)

	if z == this.NIL {
		return this.NIL
	}
	ret := &mapNode{this.NIL, this.NIL, this.NIL, z.Color, z.mapItem}

	var y *mapNode
	var x *mapNode

	if z.Left == this.NIL || z.Right == this.NIL {
		y = z
	} else {
		y = this.successor(z)
	}

	if y.Left != this.NIL {
		x = y.Left
	} else {
		x = y.Right
	}

	// Even if x is NIL, we do the assign. In that case all the NIL nodes will
	// change from {nil, nil, nil, mapBLACK, nil} to {nil, nil, ADDR, mapBLACK, nil},
	// but do not worry about that because it will not affect the compare
	// between mapNode-X with mapNode-NIL
	x.Parent = y.Parent

	if y.Parent == this.NIL {
		this.root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}

	if y != z {
		z.mapItem = y.mapItem
	}

	if y.Color == mapBLACK {
		this.deleteFixup(x)
	}

	this.count--

	return ret
}

func (this *mapImp) deleteFixup(x *mapNode) {
	for x != this.root && x.Color == mapBLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == mapRED {
				w.Color = mapBLACK
				x.Parent.Color = mapRED
				this.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == mapBLACK && w.Right.Color == mapBLACK {
				w.Color = mapRED
				x = x.Parent
			} else {
				if w.Right.Color == mapBLACK {
					w.Left.Color = mapBLACK
					w.Color = mapRED
					this.rightRotate(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = mapBLACK
				w.Right.Color = mapBLACK
				this.leftRotate(x.Parent)
				// this is to exit while loop
				x = this.root
			}
		} else { // the code below is has left and right switched from above
			w := x.Parent.Left
			if w.Color == mapRED {
				w.Color = mapBLACK
				x.Parent.Color = mapRED
				this.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Left.Color == mapBLACK && w.Right.Color == mapBLACK {
				w.Color = mapRED
				x = x.Parent
			} else {
				if w.Left.Color == mapBLACK {
					w.Right.Color = mapBLACK
					w.Color = mapRED
					this.leftRotate(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = mapBLACK
				w.Left.Color = mapBLACK
				this.rightRotate(x.Parent)
				x = this.root
			}
		}
	}
	x.Color = mapBLACK
}
