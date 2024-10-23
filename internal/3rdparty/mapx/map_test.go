package mapx

import "testing"

func TestInsertAndDelete(t *testing.T) {
	rbt := MakeMap()

	m := 0
	n := 1000
	for m < n {
		rbt.Update(m, m)
		m++
	}
	if rbt.Len() != n {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), n)
	}

	for m > 0 {
		rbt.Delete(m)
		m--
	}
	if rbt.Len() != 1 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 1)
	}
}

type testStruct struct {
	id   int
	text string
}

func TestInsertOrGet(t *testing.T) {
	rbt := MakeMap()

	items := []*testStruct{
		{1, "this"},
		{2, "is"},
		{3, "a"},
		{4, "test"},
	}

	for _, x := range items {
		rbt.Update(x.id, x.text)
	}

	newItem := &testStruct{items[0].id, "not"}
	text, ok := rbt.Lookup(newItem.id)
	if !ok {
		t.Fatal("failed")
	}
	if text != items[0].text {
		t.Errorf("tree.InsertOrGet = {id: %d, text: %s}, expect {id %d, text %s}", newItem.id, newItem.text, items[0].id, items[0].text)
	}

	newItem = &testStruct{5, "new"}
	rbt.Update(newItem.id, newItem.text)
	text, ok = rbt.Lookup(5)
	if !ok {
		t.Fatal("failed")
	}
	if text != "new" {
		t.Errorf("tree.InsertOrGet = {id: %d, text: %s}, expect {id %d, text %s}", newItem.id, newItem.text, 5, "new")
	}
}

func TestInsertString(t *testing.T) {
	rbt := MakeMap()

	rbt.Update("wa", "wa")
	rbt.Update("lang", "lang")

	if rbt.Len() != 2 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 2)
	}
}

// Test for duplicate
func TestInsertDup(t *testing.T) {
	rbt := MakeMap()

	rbt.Update("go", "wa")
	rbt.Update("go", "wa")
	rbt.Update("go", "wa")

	if rbt.Len() != 1 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 1)
	}
}

func TestGet(t *testing.T) {
	rbt := MakeMap()

	rbt.Update(1, 11)
	rbt.Update(2, 22)
	rbt.Update(3, 33)

	no, _ := rbt.Lookup(100)
	ok, _ := rbt.Lookup(1)

	if no != nil {
		t.Errorf("100 is expect not exists")
	}

	if ok == nil {
		t.Errorf("1 is expect exists")
	}
}
