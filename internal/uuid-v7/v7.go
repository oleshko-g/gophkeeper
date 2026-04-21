// Package uuid provides a UUID generator based on the v7 specification.
package uuidv7

import "github.com/google/uuid"

// NewString generates a new UUID based on the v7 specification as a string.
func NewString() string {
	id, _ := uuid.NewV7()
	return id.String()
}
