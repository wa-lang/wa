package wazero

import "wa-lang.org/wa/internal/3rdparty/wazero/internal/version"

// wazeroVersion holds the current version of wazero.
var wazeroVersion string

func init() {
	wazeroVersion = version.GetWazeroVersion()
}
