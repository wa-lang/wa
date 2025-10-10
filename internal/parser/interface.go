// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the exported entry points for invoking the parser.

package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser/w2parser"
	"wa-lang.org/wa/internal/token"
)

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func readSource(vfs fs.FS, filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return ioutil.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}
	if vfs != nil {
		return fs.ReadFile(vfs, filename)
	}
	return os.ReadFile(filename)
}

// A Mode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
type Mode = ast.ParserMode

const (
	PackageClauseOnly = ast.PackageClauseOnly // stop parsing after package clause
	ImportsOnly       = ast.ImportsOnly       // stop parsing after import declarations
	ParseComments     = ast.ParseComments     // parse comments and add them to AST
	Trace             = ast.Trace             // print a trace of parsed productions
	DeclarationErrors = ast.DeclarationErrors // report declaration errors
	SpuriousErrors    = ast.SpuriousErrors    // same as AllErrors, for backward-compatibility
	AllErrors         = SpuriousErrors        // report all errors (not just the first 10 on different lines)
)

// ParseFile parses the source code of a single Go source file and returns
// the corresponding ast.File node. The source code may be provided via
// the filename of the source file, or via the src parameter.
//
// If src != nil, ParseFile parses the source from src and the filename is
// only used when recording position information. The type of the argument
// for the src parameter must be string, []byte, or io.Reader.
// If src == nil, ParseFile parses the file specified by filename.
//
// The mode parameter controls the amount of source text parsed and other
// optional parser functionality. Position information is recorded in the
// file set fset, which must not be nil.
//
// If the source couldn't be read, the returned AST is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors
// are returned via a scanner.ErrorList which is sorted by file position.
func ParseFile(vfs fs.FS, fset *token.FileSet, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	if strings.HasSuffix(filename, ".wz") {
		return w2parser.ParseFile(vfs, fset, filename, src, mode)
	}

	// get source
	text, err := readSource(vfs, filename, src)
	if err != nil {
		return nil, err
	}

	var p parser
	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}

		// set result values
		if f == nil {
			// source is not a valid Go source file - satisfy
			// ParseFile API and return a valid (but) empty
			// *ast.File
			f = &ast.File{
				Name:  new(ast.Ident),
				Scope: ast.NewScope(nil),
			}
		}

		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse source
	p.init(fset, filename, text, mode)
	f = p.parseFile()

	return
}

// ParseDir calls ParseFile for all files with names ending in ".wa" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with os.FileInfo entries passing through
// the filter (and ending in ".wa") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
func ParseDir(vfs fs.FS, fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var list []os.FileInfo
	{
		allList, err := fd.Readdir(-1)
		if err != nil {
			return nil, err
		}

		// 一个包只能有一种中英文模式, 禁止混用
		var waList []os.FileInfo
		var w2List []os.FileInfo
		for _, d := range allList {
			if strHasSuffix(d.Name(), ".wa") {
				waList = append(waList, d)
			}
			if strHasSuffix(d.Name(), ".wz") {
				w2List = append(w2List, d)
			}
		}
		if len(waList) > 0 && len(w2List) > 0 {
			err = fmt.Errorf("%s donot support wa and wz mode in same package", path)
			return nil, err
		}
		if len(waList) > 0 {
			list = waList
		}
		if len(w2List) > 0 {
			list = w2List
		}
	}

	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if strHasSuffix(d.Name(), ".wa", ".wz") {
			if filter == nil || filter(d) {
				filename := filepath.Join(path, d.Name())
				if src, err := ParseFile(vfs, fset, filename, nil, mode); err == nil {
					name := src.Name.Name
					pkg, found := pkgs[name]
					if !found {
						pkg = &ast.Package{
							W2Mode: strHasSuffix(d.Name(), ".wz"),
							Name:   name,
							Files:  make(map[string]*ast.File),
						}
						pkgs[name] = pkg
					}
					assert(pkg.W2Mode == src.W2Mode, "package only support same mode, wa or w2")
					pkg.Files[filename] = src
				} else if first == nil {
					first = err
				}
			}
		}
	}

	return
}

func strHasSuffix(s string, ext ...string) bool {
	for _, suffix := range ext {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// ParseExprFrom is a convenience function for parsing an expression.
// The arguments have the same meaning as for ParseFile, but the source must
// be a valid Go (type or value) expression. Specifically, fset must not
// be nil.
func ParseExprFrom(fset *token.FileSet, filename string, src interface{}, mode Mode) (ast.Expr, error) {
	if fset == nil {
		panic("parser.ParseExprFrom: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(nil, filename, src)
	if err != nil {
		return nil, err
	}

	var p parser
	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}
		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse expr
	p.init(fset, filename, text, mode)
	// Set up pkg-level scopes to avoid nil-pointer errors.
	// This is not needed for a correct expression x as the
	// parser will be ok with a nil topScope, but be cautious
	// in case of an erroneous x.
	p.openScope()
	p.pkgScope = p.topScope
	e := p.parseRhsOrType()
	p.closeScope()
	assert(p.topScope == nil, "unbalanced scopes")

	// If a semicolon was inserted, consume it;
	// report an error if there's more tokens.
	if p.tok == token.SEMICOLON && p.lit == "\n" {
		p.next()
	}
	p.expect(token.EOF)

	if p.errors.Len() > 0 {
		p.errors.Sort()
		return nil, p.errors.Err()
	}

	return e, nil
}

// ParseExpr is a convenience function for obtaining the AST of an expression x.
// The position information recorded in the AST is undefined. The filename used
// in error messages is the empty string.
func ParseExpr(x string) (ast.Expr, error) {
	return ParseExprFrom(token.NewFileSet(), "", []byte(x), 0)
}
