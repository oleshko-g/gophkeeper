// Package model contains the data models used by the gophkeeper application.
package depositor

import (
	"time"

	"github.com/google/uuid"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

// PubKey is the model of a registered pub key.
// It enforces the relationship between a public key and a refresh token.
type PubKey struct {
	ID    uuidv7.UUID[uuid.UUID] `json:"pub_key_id"`
	Value string                 `json:"-"`
}

// IDValue is the interface that a refresh token value must implement.
// It's intended to be:
//   - instantiated with [uuid.UUID].
//   - populated with a UUID v7 value. So the ID iself could answer at what timestamp the refresh token has been issued.
type IDValue interface {
	Time() uuid.Time
	Version() uuid.Version
	String() string
}

// RegisteredAt returns the timestamp at which the pub key has been registered.
// It uses the UUID v7 timestamp to determine the issuance time.
func (r *PubKey) RegisteredAt() time.Time {
	return r.ID.Time()
}
