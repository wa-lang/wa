package natsort

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestStrings(t *testing.T) {
	golden := []struct {
		in   []string
		want []string
	}{
		{
			in: []string{"abc5", "abc1", "abc01", "ab", "abc10", "abc2"},
			want: []string{
				"ab",
				"abc1",
				"abc01",
				"abc2",
				"abc5",
				"abc10",
			},
		},
		{
			in: []string{"foo20", "foo.bar", "foo2", "foo.10", "foo.1", "foo.20", "foo.11", "foo1", "foobar", "foo21", "foo10", "foo11", "foo.21", "foo.2"},
			want: []string{
				"foo.1",
				"foo.2",
				"foo.10",
				"foo.11",
				"foo.20",
				"foo.21",
				"foo.bar",
				"foo1",
				"foo2",
				"foo10",
				"foo11",
				"foo20",
				"foo21",
				"foobar",
			},
		},
	}
	for _, g := range golden {
		Strings(g.in)
		if !reflect.DeepEqual(g.want, g.in) {
			t.Errorf("Error: sort failed, expected: %#q, got: %#q", g.want, g.in)
		}
	}
}

func TestLess(t *testing.T) {
	testset := []struct {
		s1, s2 string
		less   bool
	}{
		{"0", "00", true},
		{"00", "0", false},
		{"aa", "ab", true},
		{"ab", "abc", true},
		{"abc", "ad", true},
		{"ab1", "ab2", true},
		{"ab1c", "ab1c", false},
		{"ab12", "abc", true},
		{"ab2a", "ab10", true},
		{"a0001", "a0000001", true},
		{"a10", "abcdefgh2", true},
		{"аб2аб", "аб10аб", true},
		{"2аб", "3аб", true},
		//
		{"a1b", "a01b", true},
		{"a01b", "a1b", false},
		{"ab01b", "ab010b", true},
		{"ab010b", "ab01b", false},
		{"a01b001", "a001b01", true},
		{"a001b01", "a01b001", false},
		{"a1", "a1x", true},
		{"1ax", "1b", true},
		{"1b", "1ax", false},
		//
		{"082", "83", true},
		//
		{"083a", "9a", false},
		{"9a", "083a", true},
		//
		{"foo.bar", "foo123", true},
		{"foo123", "foo.bar", false},
	}
	for _, v := range testset {
		if res := Less(v.s1, v.s2); res != v.less {
			t.Errorf("Compared %#q to %#q: expected %v, got %v",
				v.s1, v.s2, v.less, res)
		}
	}
}

func BenchmarkStdStrings(b *testing.B) {
	set := testSet(300)
	arr := make([]string, len(set[0]))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, list := range set {
			b.StopTimer()
			copy(arr, list)
			b.StartTimer()

			sort.Strings(arr)
		}
	}
}

func BenchmarkStrings(b *testing.B) {
	set := testSet(300)
	arr := make([]string, len(set[0]))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, list := range set {
			b.StopTimer()
			copy(arr, list)
			b.StartTimer()

			Strings(arr)
		}
	}
}

func BenchmarkStdLess(b *testing.B) {
	set := testSet(300)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range set[0] {
			k := (j + 1) % len(set[0])
			_ = set[0][j] < set[0][k]
		}
	}
}

func BenchmarkLess(b *testing.B) {
	set := testSet(300)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range set[0] {
			k := (j + 1) % len(set[0])
			_ = Less(set[0][j], set[0][k])
		}
	}
}

// Get 1000 arrays of 10000-string-arrays (less if -short is specified).
func testSet(seed int) [][]string {
	gen := &generator{
		src: rand.New(rand.NewSource(
			int64(seed),
		)),
	}
	n := 1000
	if testing.Short() {
		n = 1
	}
	set := make([][]string, n)
	for i := range set {
		strings := make([]string, 10000)
		for idx := range strings {
			// Generate a random string
			strings[idx] = gen.NextString()
		}
		set[i] = strings
	}
	return set
}

type generator struct {
	src *rand.Rand
}

func (g *generator) NextInt(max int) int {
	return g.src.Intn(max)
}

// Gets random random-length alphanumeric string.
func (g *generator) NextString() (str string) {
	// Random-length 3-8 chars part
	strlen := g.src.Intn(6) + 3
	// Random-length 1-3 num
	numlen := g.src.Intn(3) + 1
	// Random position for num in string
	numpos := g.src.Intn(strlen + 1)
	// Generate the number
	var num string
	for i := 0; i < numlen; i++ {
		num += strconv.Itoa(g.src.Intn(10))
	}
	// Put it all together
	for i := 0; i < strlen+1; i++ {
		if i == numpos {
			str += num
		} else {
			c := 'a' + g.src.Intn(16)
			str += fmt.Sprintf("%c", c)
		}
	}
	return str
}
