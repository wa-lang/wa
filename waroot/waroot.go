// 版权 @2023 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
)

//go:embed bin/*
//go:embed cplog/*
//go:embed docs/*
//go:embed src/*
//go:embed changelog.md
//go:embed CONTRIBUTORS
//go:embed hello.wa
//go:embed LICENSE
//go:embed logo.png
//go:embed README.md
//go:embed VERSION
var _warootFS embed.FS

func GetRootFS() fs.FS {
	return _warootFS
}
