package storage

import "errors"

var ErrEmptyInput = errors.New("empty input")
var ErrNotFound = errors.New("not found")
var ErrOwnerMismatch = errors.New("owner mismatch")
