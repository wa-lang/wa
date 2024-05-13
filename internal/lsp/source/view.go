// 版权 @2024 凹语言 作者。保留所有权利。

package source

import (
	"context"

	"wa-lang.org/wa/internal/lsp/span"
)

// View represents a single workspace.
// This is the level at which we maintain configuration like working directory
// and build tags.
type View interface {
	// Name returns the name this view was constructed with.
	Name() string

	// Folder returns the folder with which this view was created.
	Folder() span.URI

	// TempWorkspace returns the folder this view uses for its temporary
	// workspace module.
	TempWorkspace() span.URI

	// Shutdown closes this view, and detaches it from its session.
	Shutdown(ctx context.Context)

	// Snapshot returns the current snapshot for the view.
	Snapshot(ctx context.Context) (Snapshot, func())

	// Rebuild rebuilds the current view, replacing the original view in its session.
	Rebuild(ctx context.Context) (Snapshot, func(), error)
}
