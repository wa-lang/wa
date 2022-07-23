// Package osutil implements OS utility functions.
package osutil

import (
	"log"
	"os"
)

// Exists reports whether the given file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Printf("unable to stat path %q; %v", path, err)
	return false
}
