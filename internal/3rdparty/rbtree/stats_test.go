// Copyright 2015, Hu Keping. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbtree

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestDeleteReturnValue(t *testing.T) {
	rbt := New()

	rbt.Insert(String("go"))
	rbt.Insert(String("lang"))

	if rbt.Len() != 2 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 2)
	}

	// go should be in the rbtree
	deletedGo := rbt.Delete(String("go"))
	if deletedGo != String("go") {
		t.Errorf("expect %v, got %v", "go", deletedGo)
	}

	// C should not be in the rbtree
	deletedC := rbt.Delete(String("C"))
	if deletedC != nil {
		t.Errorf("expect %v, got %v", nil, deletedC)
	}
}

func TestMin(t *testing.T) {
	rbt := New()

	m := 0
	n := 1000
	for m < n {
		rbt.Insert(Int(m))
		m++
	}
	if rbt.Len() != uint(n) {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), n)
	}

	for m >= 0 {
		rbt.Delete(rbt.Min())
		m--
	}

	if rbt.Len() != 0 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 0)
	}
}

// This test will first add 1000 numbers into a tree and then delete some
// from it.
//
// All the adding and deleting are in goroutine so that the delete action
// will not always succeed for example `delete(400)` but `add(400)` has not
// been executed yet.
//
// Finally, we'll add back all the deleted nodes to see if there is
// anything wrong.
func TestWithGoroutine(t *testing.T) {
	var Cache struct {
		rbt    *Rbtree
		locker *sync.Mutex
	}

	var DeletedArray struct {
		array  []Item
		locker *sync.Mutex
	}

	// Init the rbtree and the locker for later use
	Cache.rbt = New()
	Cache.locker = new(sync.Mutex)

	// Init the locker for later use
	DeletedArray.locker = new(sync.Mutex)

	i, m := 0, 0
	j, n := 1000, 1000

	var expected []Item

	// This loop will add intergers [m~n) to the rbtree.
	for m < n {
		expected = append(expected, Int(m))

		go func(x Int) {
			Cache.locker.Lock()
			Cache.rbt.Insert(x)
			Cache.locker.Unlock()
		}(Int(m))
		m++
	}

	// This loop will try to delete the even integers in [m~n),
	// Be noticed that the delete will not always succeeds since we are
	// in the goroutines.
	// We will record which ones have been removed.
	for i < j {
		if i%2 == 0 {
			go func(x Int) {
				Cache.locker.Lock()
				value := Cache.rbt.Delete(x)
				Cache.locker.Unlock()

				DeletedArray.locker.Lock()
				DeletedArray.array = append(DeletedArray.array, value)
				DeletedArray.locker.Unlock()
			}(Int(i))
		}
		i++
	}

	// Let's give a little time to those goroutines to finish their job.
	time.Sleep(time.Second * 1)

	// Add deleted Items back
	cnt := 0
	DeletedArray.locker.Lock()
	for _, v := range DeletedArray.array {
		if v != nil {
			Cache.locker.Lock()
			Cache.rbt.Insert(v)
			Cache.locker.Unlock()
			cnt++
		}
	}
	DeletedArray.locker.Unlock()
	fmt.Printf("In TestWithGoroutine(), we have deleted [%v] nodes.\n", cnt)

	var ret []Item
	Cache.rbt.Ascend(Cache.rbt.Min(), func(item Item) bool {
		ret = append(ret, item)
		return true
	})

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("expected %v but got %v", expected, ret)
	}
}
