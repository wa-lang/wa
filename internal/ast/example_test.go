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
	//      1  .  Shebang: ""
	//      2  .  Package: 2:1
	//      3  .  Name: *ast.Ident {
	//      4  .  .  NamePos: 2:9
	//      5  .  .  Name: "main"
	//      6  .  }
	//      7  .  Decls: []ast.Decl (len = 1) {
	//      8  .  .  0: *ast.FuncDecl {
	//      9  .  .  .  Name: *ast.Ident {
	//     10  .  .  .  .  NamePos: 3:6
	//     11  .  .  .  .  Name: "main"
	//     12  .  .  .  .  Obj: *ast.Object {
	//     13  .  .  .  .  .  Kind: func
	//     14  .  .  .  .  .  Name: "main"
	//     15  .  .  .  .  .  Decl: *(obj @ 8)
	//     16  .  .  .  .  }
	//     17  .  .  .  }
	//     18  .  .  .  Type: *ast.FuncType {
	//     19  .  .  .  .  Func: 3:1
	//     20  .  .  .  .  Params: *ast.FieldList {
	//     21  .  .  .  .  .  Opening: 3:10
	//     22  .  .  .  .  .  Closing: 3:11
	//     23  .  .  .  .  }
	//     24  .  .  .  .  ArrowPos: -
	//     25  .  .  .  }
	//     26  .  .  .  Body: *ast.BlockStmt {
	//     27  .  .  .  .  Lbrace: 3:13
	//     28  .  .  .  .  List: []ast.Stmt (len = 1) {
	//     29  .  .  .  .  .  0: *ast.ExprStmt {
	//     30  .  .  .  .  .  .  X: *ast.CallExpr {
	//     31  .  .  .  .  .  .  .  Fun: *ast.Ident {
	//     32  .  .  .  .  .  .  .  .  NamePos: 4:2
	//     33  .  .  .  .  .  .  .  .  Name: "println"
	//     34  .  .  .  .  .  .  .  }
	//     35  .  .  .  .  .  .  .  Lparen: 4:9
	//     36  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
	//     37  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
	//     38  .  .  .  .  .  .  .  .  .  ValuePos: 4:10
	//     39  .  .  .  .  .  .  .  .  .  Kind: STRING
	//     40  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
	//     41  .  .  .  .  .  .  .  .  }
	//     42  .  .  .  .  .  .  .  }
	//     43  .  .  .  .  .  .  .  Ellipsis: -
	//     44  .  .  .  .  .  .  .  Rparen: 4:25
	//     45  .  .  .  .  .  .  }
	//     46  .  .  .  .  .  }
	//     47  .  .  .  .  }
	//     48  .  .  .  .  Rbrace: 5:1
	//     49  .  .  .  }
	//     50  .  .  }
	//     51  .  }
	//     52  .  Scope: *ast.Scope {
	//     53  .  .  Objects: map[string]*ast.Object (len = 1) {
	//     54  .  .  .  "main": *(obj @ 12)
	//     55  .  .  }
	//     56  .  }
	//     57  .  Unresolved: []*ast.Ident (len = 1) {
	//     58  .  .  0: *(obj @ 31)
	//     59  .  }
	//     60  }
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
