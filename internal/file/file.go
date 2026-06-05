// Package file provides utility functions for file operations.
package file

import (
	"os"
)

// RWOpenCreatePrivate opens or creates a file with read/write permissions and private access.
func RWOpenCreatePrivate(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0o600)
}
