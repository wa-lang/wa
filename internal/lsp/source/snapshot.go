// 版权 @2024 凹语言 作者。保留所有权利。

package source

import (
	"context"
	"go/ast"
	"go/token"

	"wa-lang.org/wa/internal/lsp/span"
)

// Snapshot represents the current state for the given view.
type Snapshot interface {
	ID() uint64

	// View returns the View associated with this snapshot.
	View() View

	// BackgroundContext returns a context used for all background processing
	// on behalf of this snapshot.
	BackgroundContext() context.Context

	// Fileset returns the Fileset used to parse all the Go files in this snapshot.
	FileSet() *token.FileSet

	// ValidBuildConfiguration returns true if there is some error in the
	// user's workspace. In particular, if they are both outside of a module
	// and their GOPATH.
	ValidBuildConfiguration() bool

	// FindFile returns the FileHandle for the given URI, if it is already
	// in the given snapshot.
	FindFile(uri span.URI) VersionedFileHandle

	// GetVersionedFile returns the VersionedFileHandle for a given URI,
	// initializing it if it is not already part of the snapshot.
	GetVersionedFile(ctx context.Context, uri span.URI) (VersionedFileHandle, error)

	// GetFile returns the FileHandle for a given URI, initializing it if it is
	// not already part of the snapshot.
	GetFile(ctx context.Context, uri span.URI) (FileHandle, error)

	// AwaitInitialized waits until the snapshot's view is initialized.
	AwaitInitialized(ctx context.Context)

	// IsOpen returns whether the editor currently has a file open.
	IsOpen(uri span.URI) bool

	// ParseWa returns the parsed AST for the file.
	// If the file is not available, returns nil and an error.
	ParseWa(ctx context.Context, fh FileHandle) (*ParsedWaFile, error)

	// PosToField is a cache of *ast.Fields by token.Pos. This allows us
	// to quickly find corresponding *ast.Field node given a *types.Var.
	// We must refer to the AST to render type aliases properly when
	// formatting signatures and other types.
	PosToField(ctx context.Context, pkg Package, pos token.Pos) (*ast.Field, error)

	// PosToDecl maps certain objects' positions to their surrounding
	// ast.Decl. This mapping is used when building the documentation
	// string for the objects.
	PosToDecl(ctx context.Context, pkg Package, pos token.Pos) (ast.Decl, error)

	/*
		// DiagnosePackage returns basic diagnostics, including list, parse, and type errors
		// for pkg, grouped by file.
		DiagnosePackage(ctx context.Context, pkg Package) (map[span.URI][]*Diagnostic, error)

		// Analyze runs the analyses for the given package at this snapshot.
		Analyze(ctx context.Context, pkgID string, analyzers []*Analyzer) ([]*Diagnostic, error)

		// RunGoCommandPiped runs the given `go` command, writing its output
		// to stdout and stderr. Verb, Args, and WorkingDir must be specified.
		RunGoCommandPiped(ctx context.Context, mode InvocationFlags, inv *gocommand.Invocation, stdout, stderr io.Writer) error

		// RunGoCommandDirect runs the given `go` command. Verb, Args, and
		// WorkingDir must be specified.
		RunGoCommandDirect(ctx context.Context, mode InvocationFlags, inv *gocommand.Invocation) (*bytes.Buffer, error)

		// RunGoCommands runs a series of `go` commands that updates the go.mod
		// and go.sum file for wd, and returns their updated contents.
		RunGoCommands(ctx context.Context, allowNetwork bool, wd string, run func(invoke func(...string) (*bytes.Buffer, error)) error) (bool, []byte, []byte, error)

		// RunProcessEnvFunc runs fn with the process env for this snapshot's view.
		// Note: the process env contains cached module and filesystem state.
		RunProcessEnvFunc(ctx context.Context, fn func(*imports.Options) error) error

		// ModFiles are the go.mod files enclosed in the snapshot's view and known
		// to the snapshot.
		ModFiles() []span.URI

		// ParseMod is used to parse go.mod files.
		ParseMod(ctx context.Context, fh FileHandle) (*ParsedModule, error)

		// ModWhy returns the results of `go mod why` for the module specified by
		// the given go.mod file.
		ModWhy(ctx context.Context, fh FileHandle) (map[string]string, error)

		// ModTidy returns the results of `go mod tidy` for the module specified by
		// the given go.mod file.
		ModTidy(ctx context.Context, pm *ParsedModule) (*TidiedModule, error)

		// GoModForFile returns the URI of the go.mod file for the given URI.
		GoModForFile(uri span.URI) span.URI

		// BuiltinFile returns information about the special builtin package.
		BuiltinFile(ctx context.Context) (*ParsedGoFile, error)

		// IsBuiltin reports whether uri is part of the builtin package.
		IsBuiltin(ctx context.Context, uri span.URI) bool

		// PackagesForFile returns the packages that this file belongs to, checked
		// in mode.
		PackagesForFile(ctx context.Context, uri span.URI, mode TypecheckMode) ([]Package, error)

		// PackageForFile returns a single package that this file belongs to,
		// checked in mode and filtered by the package policy.
		PackageForFile(ctx context.Context, uri span.URI, mode TypecheckMode, selectPackage PackageFilter) (Package, error)

		// GetActiveReverseDeps returns the active files belonging to the reverse
		// dependencies of this file's package, checked in TypecheckWorkspace mode.
		GetReverseDependencies(ctx context.Context, id string) ([]Package, error)

		// CachedImportPaths returns all the imported packages loaded in this
		// snapshot, indexed by their import path and checked in TypecheckWorkspace
		// mode.
		CachedImportPaths(ctx context.Context) (map[string]Package, error)

		// KnownPackages returns all the packages loaded in this snapshot, checked
		// in TypecheckWorkspace mode.
		KnownPackages(ctx context.Context) ([]Package, error)

		// WorkspacePackages returns the snapshot's top-level packages.
		WorkspacePackages(ctx context.Context) ([]Package, error)

		// GetCriticalError returns any critical errors in the workspace.
		GetCriticalError(ctx context.Context) *CriticalError

		// BuildGoplsMod generates a go.mod file for all modules in the workspace.
		// It bypasses any existing gopls.mod.
		BuildGoplsMod(ctx context.Context) (*modfile.File, error)
	*/
}
