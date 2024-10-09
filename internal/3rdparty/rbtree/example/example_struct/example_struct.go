// Copyright 2015, Hu Keping. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"wa-lang.org/wa/internal/3rdparty/rbtree"
)

// Var is the node of a struct
type Var struct {
	Expiry time.Time `json:"expiry,omitempty"`
	ID     string    `json:"id,omitempty"`
}

// Less will order the node by `Time`
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
	rbt.Ascend(rbt.Get(tmp), print)
}

func print(item rbtree.Item) bool {
	i, ok := item.(Var)
	if !ok {
		return false
	}
	fmt.Printf("%+v\n", i)
	return true
}
