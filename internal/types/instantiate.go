// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements instantiation of generic types
// through substitution of type parameters by type arguments.

package types

import (
	"fmt"

	"wa-lang.org/wa/internal/token"
)

// Instantiate instantiates the type orig with the given type arguments targs.
// orig must be a *Named or a *Signature type. If there is no error, the
// resulting Type is an instantiated type of the same kind (either a *Named or
// a *Signature). Methods attached to a *Named type are also instantiated, and
// associated with a new *Func that has the same position as the original
// method, but nil function scope.
//
// If ctxt is non-nil, it may be used to de-duplicate the instance against
// previous instances with the same identity. As a special case, generic
// *Signature origin types are only considered identical if they are pointer
// equivalent, so that instantiating distinct (but possibly identical)
// signatures will yield different instances.
//
// If validate is set, Instantiate verifies that the number of type arguments
// and parameters match, and that the type arguments satisfy their
// corresponding type constraints. If verification fails, the resulting error
// may wrap an *ArgumentError indicating which type argument did not satisfy
// its corresponding type parameter constraint, and why.
//
// If validate is not set, Instantiate does not verify the type argument count
// or whether the type arguments satisfy their constraints. Instantiate is
// guaranteed to not return an error, but may panic. Specifically, for
// *Signature types, Instantiate will panic immediately if the type argument
// count is incorrect; for *Named types, a panic may occur later inside the
// *Named API.
func Instantiate(ctxt *Context, orig Type, targs []Type, validate bool) (Type, error) {
	if validate {
		var tparams []*TypeParam
		switch t := orig.(type) {
		case *Named:
			tparams = t.TypeParams().list()
		case *Signature:
			tparams = t.TypeParams().list()
		}
		if len(targs) != len(tparams) {
			return nil, fmt.Errorf("got %d type arguments but %s has %d type parameters", len(targs), orig, len(tparams))
		}
		//if i, err := (*Checker)(nil).verify(token.NoPos, tparams, targs); err != nil {
		//	return nil, &ArgumentError{i, err}
		//}
	}

	inst := (*Checker)(nil).instance(token.NoPos, orig, targs, ctxt)
	return inst, nil
}

// instance creates a type or function instance using the given original type
// typ and arguments targs. For Named types the resulting instance will be
// unexpanded.
func (check *Checker) instance(pos token.Pos, orig Type, targs []Type, ctxt *Context) (res Type) {
	var h string
	if ctxt != nil {
		h = ctxt.instanceHash(orig, targs)
		// typ may already have been instantiated with identical type arguments. In
		// that case, re-use the existing instance.
		if inst := ctxt.lookup(h, orig, targs); inst != nil {
			return inst
		}
	}

	switch orig := orig.(type) {
	case *Named:
		panic("todo")
		/*
			tname := NewTypeName(pos, orig.obj.pkg, orig.obj.name, nil)
			named := check.newNamed(tname, orig, nil, nil, nil) // underlying, tparams, and methods are set when named is resolved
			named.targs = newTypeList(targs)
			named.resolver = func(ctxt *Context, n *Named) (*TypeParamList, Type, *methodList) {
				return expandNamed(ctxt, n, pos)
			}
			res = named
		*/

	case *Signature:
		panic("todo")
		/*
			tparams := orig.TypeParams()
			if !check.validateTArgLen(pos, tparams.Len(), len(targs)) {
				return Typ[Invalid]
			}
			if tparams.Len() == 0 {
				return orig // nothing to do (minor optimization)
			}
			sig := check.subst(pos, orig, makeSubstMap(tparams.list(), targs), ctxt).(*Signature)
			// If the signature doesn't use its type parameters, subst
			// will not make a copy. In that case, make a copy now (so
			// we can set tparams to nil w/o causing side-effects).
			if sig == orig {
				copy := *sig
				sig = &copy
			}
			// After instantiating a generic signature, it is not generic
			// anymore; we need to set tparams to nil.
			sig.tparams = nil
			res = sig
		*/
	default:
		// only types and functions can be generic
		panic(fmt.Sprintf("%v: cannot instantiate %v", pos, orig))
	}

}
