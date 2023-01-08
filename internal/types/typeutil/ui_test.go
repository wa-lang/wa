// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeutil_test

import (
	"fmt"
	"strings"
	"testing"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/types/typeutil"
)

func TestIntuitiveMethodSet(t *testing.T) {
	const source = `
package P
type A int
func (A) f()
func (*A) g()
`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, "hello.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}

	var conf types.Config
	pkg, err := conf.Check("P", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatal(err)
	}
	qual := types.RelativeTo(pkg)

	for _, test := range []struct {
		expr string // type expression
		want string // intuitive method set
	}{
		{"A", "(A).f (*A).g"},
		{"*A", "(*A).f (*A).g"},
		{"error", "(error).Error"},
		{"*error", ""},
		{"struct{A}", "(struct{A}).f (*struct{A}).g"},
		{"*struct{A}", "(*struct{A}).f (*struct{A}).g"},
	} {
		tv, err := types.Eval(fset, pkg, 0, test.expr)
		if err != nil {
			t.Errorf("Eval(%s) failed: %v", test.expr, err)
		}
		var names []string
		for _, m := range typeutil.IntuitiveMethodSet(tv.Type, nil) {
			name := fmt.Sprintf("(%s).%s", types.TypeString(m.Recv(), qual), m.Obj().Name())
			names = append(names, name)
		}
		got := strings.Join(names, " ")
		if got != test.want {
			t.Errorf("IntuitiveMethodSet(%s) = %q, want %q", test.expr, got, test.want)
		}
	}
}
