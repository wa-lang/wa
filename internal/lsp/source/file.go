// 版权 @2024 凹语言 作者。保留所有权利。

package source

import (
	"context"
	"go/ast"
	"go/scanner"
	"go/token"

	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/lsp/span"
)

// FileKind describes the kind of the file in question.
// It can be one of Go, mod, or sum.
type FileKind int

const (
	// UnknownKind is a file type we don't know about.
	UnknownKind = FileKind(iota)

	// Wa is a normal wa source file.
	Wa
	// Mod is a wa.mod file.
	Mod
)

// FileIdentity uniquely identifies a file at a version from a FileSystem.
type FileIdentity struct {
	URI span.URI

	// Identifier represents a unique identifier for the file's content.
	Hash string

	// Kind is the file's kind.
	Kind FileKind
}

// FileHandle represents a handle to a specific version of a single file.
type FileHandle interface {
	URI() span.URI
	Kind() FileKind

	// FileIdentity returns a FileIdentity for the file, even if there was an
	// error reading it.
	FileIdentity() FileIdentity
	// Read reads the contents of a file.
	// If the file is not available, returns a nil slice and an error.
	Read() ([]byte, error)
	// Saved reports whether the file has the same content on disk.
	Saved() bool
}

type VersionedFileHandle interface {
	FileHandle
	Version() int32
	Session() string

	// LSPIdentity returns the version identity of a file.
	VersionedFileIdentity() VersionedFileIdentity
}

type VersionedFileIdentity struct {
	URI span.URI

	// SessionID is the ID of the LSP session.
	SessionID string

	// Version is the version of the file, as specified by the client. It should
	// only be set in combination with SessionID.
	Version int32
}

type FileAction int

const (
	UnknownFileAction = FileAction(iota)
	Open
	Change
	Close
	Save
	Create
	Delete
	InvalidateMetadata
)

// Overlay is the type for a file held in memory on a session.
type Overlay interface {
	VersionedFileHandle
}

// FileModification represents a modification to a file.
type FileModification struct {
	URI    span.URI
	Action FileAction

	// OnDisk is true if a watched file is changed on disk.
	// If true, Version will be -1 and Text will be nil.
	OnDisk bool

	// Version will be -1 and Text will be nil when they are not supplied,
	// specifically on textDocument/didClose and for on-disk changes.
	Version int32
	Text    []byte

	// LanguageID is only sent from the language client on textDocument/didOpen.
	LanguageID string
}

// Session represents a single connection from a client.
// This is the level at which things like open files are maintained on behalf
// of the client.
// A session may have many active views at any given time.
type Session interface {
	// ID returns the unique identifier for this session on this server.
	ID() string
	// NewView creates a new View, returning it and its first snapshot. If a
	// non-empty tempWorkspace directory is provided, the View will record a copy
	// of its gopls workspace module in that directory, so that client tooling
	// can execute in the same main module.
	NewView(ctx context.Context, name string, folder, tempWorkspace span.URI) (View, Snapshot, func(), error)

	// Cache returns the cache that created this session, for debugging only.
	Cache() interface{}

	// View returns a view with a matching name, if the session has one.
	View(name string) View

	// ViewOf returns a view corresponding to the given URI.
	ViewOf(uri span.URI) (View, error)

	// Views returns the set of active views built by this session.
	Views() []View

	// Shutdown the session and all views it has created.
	Shutdown(ctx context.Context)

	// GetFile returns a handle for the specified file.
	GetFile(ctx context.Context, uri span.URI) (FileHandle, error)

	// DidModifyFile reports a file modification to the session. It returns
	// the new snapshots after the modifications have been applied, paired with
	// the affected file URIs for those snapshots.
	DidModifyFiles(ctx context.Context, changes []FileModification) (map[Snapshot][]span.URI, []func(), error)

	// ExpandModificationsToDirectories returns the set of changes with the
	// directory changes removed and expanded to include all of the files in
	// the directory.
	ExpandModificationsToDirectories(ctx context.Context, changes []FileModification) []FileModification

	// Overlays returns a slice of file overlays for the session.
	Overlays() []Overlay
}

// A FileSource maps uris to FileHandles. This abstraction exists both for
// testability, and so that algorithms can be run equally on session and
// snapshot files.
type FileSource interface {
	// GetFile returns the FileHandle for a given URI.
	GetFile(ctx context.Context, uri span.URI) (FileHandle, error)
}

// A ParsedWaFile contains the results of parsing a Wa file.
type ParsedWaFile struct {
	URI  span.URI
	File *ast.File
	Tok  *token.File
	// Source code used to build the AST. It may be different from the
	// actual content of the file if we have fixed the AST.
	Src      []byte
	Mapper   *protocol.ColumnMapper
	ParseErr scanner.ErrorList
}
