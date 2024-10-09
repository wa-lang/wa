# Rbtree  [![GoDoc](https://godoc.org/github.com/HuKeping/rbtree?status.svg)](https://godoc.org/github.com/HuKeping/rbtree)

This is an implementation of Red-Black tree written by Golang which does **not** support `duplicate keys`.

## Installation

With a healthy Go language installed, simply run `go get github.com/HuKeping/rbtree`

## Example
All you have to do is to implement a comparison function `Less() bool` for your Item
which will be store in the Red-Black tree, here are some examples.
#### A simple case for `int` items.
	package main
	
	import (
		"fmt"
		"github.com/HuKeping/rbtree"
	)
	
	func main() {
		rbt := rbtree.New()
	
		m := 0
		n := 10
	
		for m < n {
			rbt.Insert(rbtree.Int(m))
			m++
		}
	
		m = 0
		for m < n {
			if m%2 == 0 {
				rbt.Delete(rbtree.Int(m))
			}
			m++
		}

		// 1, 3, 5, 7, 9 were expected.
		rbt.Ascend(rbt.Min(), Print)
	}
	
	func Print(item rbtree.Item) bool {
		i, ok := item.(rbtree.Int)
		if !ok {
			return false
		}
		fmt.Println(i)
		return true
	}

#### A simple case for `string` items.
	package main
	
	import (
		"fmt"
		"github.com/HuKeping/rbtree"
	)
	
	func main() {
		rbt := rbtree.New()
	
		rbt.Insert(rbtree.String("Hello"))
		rbt.Insert(rbtree.String("World"))

		rbt.Ascend(rbt.Min(), Print)
	}
	
	func Print(item rbtree.Item) bool {
		i, ok := item.(rbtree.String)
		if !ok {
			return false
		}
		fmt.Println(i)
		return true
	}

#### A quite interesting case for `struct` items.
	package main
	
	import (
		"fmt"
		"github.com/HuKeping/rbtree"
		"time"
	)
	
	type Var struct {
		Expiry time.Time `json:"expiry,omitempty"`
		ID     string    `json:"id",omitempty`
	}
	
	// We will order the node by `Time`
	func (x Var) Less(than rbtree.Item) bool {
		return x.Expiry.Before(than.(Var).Expiry)
	}
	
	func main() {
		rbt := rbtree.New()
	
		var1 := Var{
			Expiry: time.Now().Add(time.Second * 10),
			ID:     "var1",
		}
		var2 := Var{
			Expiry: time.Now().Add(time.Second * 20),
			ID:     "var2",
		}
		var3 := Var{
			Expiry: var2.Expiry,
			ID:     "var2-dup",
		}
		var4 := Var{
			Expiry: time.Now().Add(time.Second * 40),
			ID:     "var4",
		}
		var5 := Var{
			Expiry: time.Now().Add(time.Second * 50),
			ID:     "var5",
		}
	
		rbt.Insert(var1)
		rbt.Insert(var2)
		rbt.Insert(var3)
		rbt.Insert(var4)
		rbt.Insert(var5)
	
		tmp := Var{
			Expiry: var4.Expiry,
			ID:     "This field is not the key factor",
		}
	
		// var4 and var5 were expected
		rbt.Ascend(rbt.Get(tmp), Print)
	}
	
	func Print(item rbtree.Item) bool {
		i, ok := item.(Var)
		if !ok {
			return false
		}
		fmt.Printf("%+v\n", i)
		return true
	}
