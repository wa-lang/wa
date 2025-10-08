// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

// This example demonstrates how to inspect the AST of a Go program.
func ExampleInspect() {
	// src is the input for which we want to inspect the AST.
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})

	// Output:
	// src.go:2:9:	p
	// src.go:3:7:	c
	// src.go:3:11:	1.0
	// src.go:4:5:	X
	// src.go:4:9:	f
	// src.go:4:11:	3.14
	// src.go:4:17:	2
	// src.go:4:21:	c
}

// This example shows what an AST looks like when printed for debugging.
func ExamplePrint() {
	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
	println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)

	// Output:
	//      0  *ast.File {
	//      1  .  W2Mode: false
	//      2  .  Shebang: ""
	//      3  .  Package: 2:1
	//      4  .  Name: *ast.Ident {
	//      5  .  .  NamePos: 2:9
	//      6  .  .  Name: "main"
	//      7  .  }
	//      8  .  Decls: []ast.Decl (len = 1) {
	//      9  .  .  0: *ast.FuncDecl {
	//     10  .  .  .  Name: *ast.Ident {
	//     11  .  .  .  .  NamePos: 3:6
	//     12  .  .  .  .  Name: "main"
	//     13  .  .  .  .  Obj: *ast.Object {
	//     14  .  .  .  .  .  Kind: func
	//     15  .  .  .  .  .  Name: "main"
	//     16  .  .  .  .  .  Decl: *(obj @ 9)
	//     17  .  .  .  .  }
	//     18  .  .  .  }
	//     19  .  .  .  Type: *ast.FuncType {
	//     20  .  .  .  .  TokPos: 3:1
	//     21  .  .  .  .  Tok: func
	//     22  .  .  .  .  Params: *ast.FieldList {
	//     23  .  .  .  .  .  Opening: 3:10
	//     24  .  .  .  .  .  Closing: 3:11
	//     25  .  .  .  .  }
	//     26  .  .  .  .  ArrowPos: -
	//     27  .  .  .  }
	//     28  .  .  .  Body: *ast.BlockStmt {
	//     29  .  .  .  .  Lbrace: 3:13
	//     30  .  .  .  .  List: []ast.Stmt (len = 1) {
	//     31  .  .  .  .  .  0: *ast.ExprStmt {
	//     32  .  .  .  .  .  .  X: *ast.CallExpr {
	//     33  .  .  .  .  .  .  .  Fun: *ast.Ident {
	//     34  .  .  .  .  .  .  .  .  NamePos: 4:2
	//     35  .  .  .  .  .  .  .  .  Name: "println"
	//     36  .  .  .  .  .  .  .  }
	//     37  .  .  .  .  .  .  .  Lparen: 4:9
	//     38  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
	//     39  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
	//     40  .  .  .  .  .  .  .  .  .  ValuePos: 4:10
	//     41  .  .  .  .  .  .  .  .  .  Kind: STRING
	//     42  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
	//     43  .  .  .  .  .  .  .  .  }
	//     44  .  .  .  .  .  .  .  }
	//     45  .  .  .  .  .  .  .  Ellipsis: -
	//     46  .  .  .  .  .  .  .  Rparen: 4:25
	//     47  .  .  .  .  .  .  }
	//     48  .  .  .  .  .  }
	//     49  .  .  .  .  }
	//     50  .  .  .  .  Rbrace: 5:1
	//     51  .  .  .  }
	//     52  .  .  }
	//     53  .  }
	//     54  .  Scope: *ast.Scope {
	//     55  .  .  Objects: map[string]*ast.Object (len = 1) {
	//     56  .  .  .  "main": *(obj @ 13)
	//     57  .  .  }
	//     58  .  }
	//     59  .  Unresolved: []*ast.Ident (len = 1) {
	//     60  .  .  0: *(obj @ 33)
	//     61  .  }
	//     62  }
}

// This example illustrates how to remove a variable declaration
// in a Go program while maintaining correct comment association
// using an ast.CommentMap.
func ExampleCommentMap() {
	// src is the input for which we create the AST that we
	// are going to manipulate.
	src := `
// This is the package comment.
package main

// This comment is associated with the hello constant.
const hello = "Hello, World!" // line comment 1

// This comment is associated with the foo variable.
var foo = hello // line comment 2

// This comment is associated with the main function.
func main() {
	fmt.Println(hello) // line comment 3
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, "src.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Create an ast.CommentMap from the ast.File's comments.
	// This helps keeping the association between comments
	// and AST nodes.
	cmap := ast.NewCommentMap(fset, f, f.Comments)

	// Remove the first variable declaration from the list of declarations.
	for i, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
			copy(f.Decls[i:], f.Decls[i+1:])
			f.Decls = f.Decls[:len(f.Decls)-1]
			break
		}
	}

	// Use the comment map to filter comments that don't belong anymore
	// (the comments associated with the variable declaration), and create
	// the new comments list.
	f.Comments = cmap.Filter(f).Comments()

	// Print the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	// Output:
	// // This is the package comment.
	// package main
	//
	// // This comment is associated with the hello constant.
	// const hello = "Hello, World!" // line comment 1
	//
	// // This comment is associated with the main function.
	// func main() {
	// 	fmt.Println(hello) // line comment 3
	// }
}
