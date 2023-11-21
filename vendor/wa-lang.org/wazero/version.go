package wazero

import "wa-lang.org/wazero/internal/version"

// wazeroVersion holds the current version of wazero.
var wazeroVersion string

func init() {
	wazeroVersion = version.GetWazeroVersion()
}
