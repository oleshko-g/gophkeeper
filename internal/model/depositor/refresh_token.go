// Package model contains the data models used by the gophkeeper application.
package depositor

import (
	"time"

	"github.com/google/uuid"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

// RefreshToken is the model of a refresh token.
// It enforces the relationship between a public key and a refresh token.
type RefreshToken struct {
	// PubKeyID is the ID of the public key associated with the refresh token
	PubKeyID uuidv7.UUID[uuid.UUID]
	// ID is the actual refresh token value
	ID uuidv7.UUID[uuid.UUID]
	// TTL is the time-to-live of the refresh token
	TTL time.Duration
	// RevokedAt is the timestamp at which the refresh token has been revoked
	RevokedAt time.Time
}

// RefreshTokenValue is the interface that a refresh token value must implement.
// It's intended to be:
//   - instantiated with [uuid.UUID].
//   - populated with a UUID v7 value. So the ID iself could answer at what timestamp the refresh token has been issued.
type RefreshTokenValue interface {
	Time() uuid.Time
	Version() uuid.Version
	String() string
}

// IssuedAt returns the timestamp at which the refresh token has been issued.
// It uses the UUID v7 timestamp to determine the issuance time.
func (r *RefreshToken) IssuedAt() time.Time {
	return r.ID.Time()
}

// ExpiresAt returns the timestamp at which the refresh token expires.
func (r *RefreshToken) ExpiresAt() time.Time {
	return r.IssuedAt().Add(r.TTL)
}

// IsExpired returns true if the refresh token has expired.
func (r *RefreshToken) IsExpired() bool {
	if time.Now().After(r.ExpiresAt()) {
		return true
	}

	return false
}

// IsRevoked returns true if the refresh token has been revoked.
func (r *RefreshToken) IsRevoked() bool {
	if r.IssuedAt().After(r.RevokedAt) {
		return true
	}
	return false
}
