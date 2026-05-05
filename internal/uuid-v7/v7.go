// Package uuidv7 provides a UUID generator based on the v7 specification.
package uuidv7

import (
	"time"

	"github.com/google/uuid"
)

type UUID[T timeID] struct {
	Value  T
	String string
}

// timeID is the interface that an ID value must implement.
// It's intended to be:
//   - instantiated with [uuid.UUID].
//   - populated with a UUID v7 value. So the ID iself could answer at what timestamp the refresh token has been issued.
type timeID interface {
	Time() uuid.Time
	Version() uuid.Version
	String() string
}

// New generates a new UUID based on the v7 specification as a uuid.UUID.
func New() UUID[uuid.UUID] {
	id, _ := uuid.NewV7()
	return UUID[uuid.UUID]{Value: id, String: id.String()}
}

// Time returns the timestamp at which the Value has been created.
// It uses the UUID v7 timestamp to determine the issuance time.
func (u *UUID[T]) Time() time.Time {
	sec, nsec := u.Value.Time().UnixTime()
	return time.Date(1970, 01, 01, 0, 0, int(sec), int(nsec), time.UTC)
}

func FromString(s string) UUID[uuid.UUID] {
	id, _ := uuid.Parse(s)
	return UUID[uuid.UUID]{Value: id}
}
