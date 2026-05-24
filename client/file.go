package main

import (
	"io"
	"os"
)

// storager is an interface for a file-like storage object.
type storager interface {
	io.ReadWriteSeeker
	io.Closer
	Truncate(size int64) error
}

// UNIX permissions
const (
	// user:  Write Read _
	// group: _     Read _
	// other: _     Read _
	filePerm os.FileMode = 0o600

	// user:  Write Read Exec
	// group: _     _		 _
	// other: _     _		 _
	dirPerm os.FileMode = 0o700
)
