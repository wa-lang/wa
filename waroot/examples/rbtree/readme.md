# BUG

main.wa:

```
import "rbtree/mapx"

func main {
	m := mapx.MakeMap()
	m.Update("three", 3)
	m.Update("one", 1)

	m.Dump()
}
```

map.wa

```
func mapImp.insert(z: *mapNode) => *mapNode {
	println("insert", z, z.Key.(string))

	x := this.root
	y := this.NIL

	println("insert: y.v0=", y)

	for x != this.NIL {
		y = x
		println("insert: y.v1=", y)

		if mapLess(z, x) {
			x = x.Left
		} else if mapLess(x, z) {
			x = x.Right
		} else {
			return x
		}
	}

	println("insert: y.v2=", y) // BUG: y的值没有变化

	z.Parent = y
	if y == this.NIL {
		this.root = z
	} else if mapLess(z, y) {
		y.Left = z
	} else {
		y.Right = z
	}

	this.count++
	this.insertFixup(z)
	return z
}
```

BUG: 其中y在for循环中被修改，跳出循环后y值没有变化。

```
$ make
go run ../../../main.go run .
insert (0x3fffdd0,8,0x3fffde0) three
insert: y.v0= (0x3ffffb8,18,0x3ffffc8)
insert: y.v2= (0x3ffffb8,21,0x3ffffc8)
l > r: three one
l < r: one three
insert (0x3fffd38,8,0x3fffd48) one
insert: y.v0= (0x3ffffb8,19,0x3ffffc8)
insert: y.v1= (0x3fffdd0,3,0x3fffde0)
l < r: one three
insert: y.v2= (0x3ffffb8,23,0x3ffffc8)
(0x3ffffb8,7,0x3ffffc8) nil
(0x3fffd38,8,0x3fffd48) (0x3ffffb8,8,0x3ffffc8) (0x3ffffb8,8,0x3ffffc8) one 1 1
(0x3ffffb8,10,0x3ffffc8) nil
```