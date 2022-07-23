// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains test cases for short valid and invalid programs.

package parser

import "testing"

var valids = []string{
	"package p\n",
	`package p;`,
	`package p; import "fmt"; fn f() { fmt.Println("Hello, World!") };`,
	`package p; fn f() { if f(T{}) {} };`,
	`package p; fn f(fn() fn() fn());`,
	`package p; fn f(...T);`,
	`package p; fn f(float, ...int);`,
	`package p; fn f(x int, a ...int) { f(0, a...); f(1, a...,) };`,
	`package p; fn f(int,) {};`,
	`package p; fn f(...int,) {};`,
	`package p; fn f(x ...int,) {};`,
	`package p; type T []int; let a []bool; fn f() { if a[T{42}[0]] {} };`,
	`package p; type T []int; fn g(int) bool { return true }; fn f() { if g(T{42}[0]) {} };`,
	`package p; type T []int; fn f() { for _ = range []int{T{42}[0]} {} };`,
	`package p; let a = T{{1, 2}, {3, 4}}`,
	`package p; fn f() { if ; true {} };`,
	`package p; fn f() { switch ; {} };`,
	`package p; fn f() { for _ = range "foo" + "bar" {} };`,
	`package p; fn f() { let s []int; g(s[:], s[i:], s[:j], s[i:j], s[i:j:k], s[:j:k]) };`,
	`package p; let ( _ = (struct {*T}).m; _ = (interface {T}).m )`,
	`package p; fn ((T),) m() {}`,
	`package p; fn ((*T),) m() {}`,
	`package p; fn (*(T),) m() {}`,
	`package p; fn _(x []int) { for range x {} }`,
	`package p; fn _() { if [T{}.n]int{} {} }`,
	`package p; fn _() { map[int]int{}[0]++; map[int]int{}[0] += 1 }`,
	`package p; fn _(x interface{f()}) { interface{f()}(x).f() }`,
	`package p; const (x = 0; y; z)`, // issue 9639
	`package p; let _ = map[P]int{P{}:0, {}:1}`,
	`package p; let _ = map[*P]int{&P{}:0, {}:1}`,
	`package p; type T = int`,
	`package p; type (T = p.T; _ = struct{}; x = *T)`,
}

func TestValid(t *testing.T) {
	for _, src := range valids {
		checkErrors(nil, t, src, src)
	}
}

var invalids = []string{
	// `foo /* ERROR "expected 'package'" */ !`,
	`package p; fn f() { if { /* ERROR "missing condition" */ } };`,
	`package p; fn f() { if ; /* ERROR "missing condition" */ {} };`,
	`package p; fn f() { if f(); /* ERROR "missing condition" */ {} };`,
	`package p; fn f() { if _ = range /* ERROR "expected operand" */ x; true {} };`,
	`package p; fn f() { switch _ /* ERROR "expected switch expression" */ = range x; true {} };`,
	`package p; fn f() { for _ = range x ; /* ERROR "expected '{'" */ ; {} };`,
	`package p; fn f() { for ; ; _ = range /* ERROR "expected operand" */ x {} };`,
	`package p; fn f() { for ; _ /* ERROR "expected boolean or range expression" */ = range x ; {} };`,
	`package p; fn f() { switch t = /* ERROR "expected ':=', found '='" */ t.(type) {} };`,
	`package p; fn f() { switch t /* ERROR "expected switch expression" */ , t = t.(type) {} };`,
	`package p; fn f() { switch t /* ERROR "expected switch expression" */ = t.(type), t {} };`,
	`package p; let a = [ /* ERROR "expected expression" */ 1]int;`,
	`package p; let a = [ /* ERROR "expected expression" */ ...]int;`,
	`package p; let a = struct /* ERROR "expected expression" */ {}`,
	`package p; let a = fn /* ERROR "expected expression" */ ();`,
	`package p; let a = interface /* ERROR "expected expression" */ {}`,
	`package p; let a = [ /* ERROR "expected expression" */ ]int`,
	`package p; let a = map /* ERROR "expected expression" */ [int]int`,
	`package p; let a = []int{[ /* ERROR "expected expression" */ ]int};`,
	`package p; let a = ( /* ERROR "expected expression" */ []int);`,
	`package p; let a = a[[ /* ERROR "expected expression" */ ]int:[]int];`,
	`package p; fn f() { let t []int; t /* ERROR "expected identifier on left side of :=" */ [0] := 0 };`,
	`package p; fn f() { if x := g(); x /* ERROR "expected boolean expression" */ = 0 {}};`,
	`package p; fn f() { _ = x = /* ERROR "expected '=='" */ 0 {}};`,
	`package p; fn f() { _ = 1 == fn()int { let x bool; x = x = /* ERROR "expected '=='" */ true; return x }() };`,
	`package p; fn f() { let s []int; _ = s[] /* ERROR "expected operand" */ };`,
	`package p; fn f() { let s []int; _ = s[i:j: /* ERROR "3rd index required" */ ] };`,
	`package p; fn f() { let s []int; _ = s[i: /* ERROR "2nd index required" */ :k] };`,
	`package p; fn f() { let s []int; _ = s[i: /* ERROR "2nd index required" */ :] };`,
	`package p; fn f() { let s []int; _ = s[: /* ERROR "2nd index required" */ :] };`,
	`package p; fn f() { let s []int; _ = s[: /* ERROR "2nd index required" */ ::] };`,
	`package p; fn f() { let s []int; _ = s[i:j:k: /* ERROR "expected ']'" */ l] };`,
	`package p; fn f() { for x /* ERROR "boolean or range expression" */ = []string {} }`,
	`package p; fn f() { for x /* ERROR "boolean or range expression" */ := []string {} }`,
	`package p; fn f() { for i /* ERROR "boolean or range expression" */ , x = []string {} }`,
	`package p; fn f() { for i /* ERROR "boolean or range expression" */ , x := []string {} }`,
	`package p; fn f() { defer fn() {} /* ERROR HERE "function must be invoked" */ }`,
	`package p; fn f(x fn(), u v fn /* ERROR "missing ','" */ ()){}`,

	// issue 8656
	`package p; fn f() (a b string /* ERROR "missing ','" */ , ok bool)`,

	// issue 9639
	`package p; let x /* ERROR "missing variable type or initialization" */ , y, z;`,
	`package p; const x /* ERROR "missing constant value" */ ;`,
	`package p; const x /* ERROR "missing constant value" */ int;`,
	`package p; const (x = 0; y; z /* ERROR "missing constant value" */ int);`,

	// issue 12437
	`package p; let _ = struct { x int, /* ERROR "expected ';', found ','" */ }{};`,
	`package p; let _ = struct { x int, /* ERROR "expected ';', found ','" */ y float }{};`,

	// issue 11611
	`package p; type _ struct { int, } /* ERROR "expected type, found '}'" */ ;`,
	`package p; type _ struct { int, float } /* ERROR "expected type, found '}'" */ ;`,
	`package p; type _ struct { ( /* ERROR "expected anonymous field" */ int) };`,
	`package p; fn _()(x, y, z ... /* ERROR "expected '\)', found '...'" */ int){}`,
	`package p; fn _()(... /* ERROR "expected type, found '...'" */ int){}`,

	// issue 13475
	`package p; fn f() { if true {} else ; /* ERROR "expected if statement or block" */ }`,
	`package p; fn f() { if true {} else defer /* ERROR "expected if statement or block" */ f() }`,
}

func TestInvalid(t *testing.T) {
	for _, src := range invalids {
		checkErrors(nil, t, src, src)
	}
}
