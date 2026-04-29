package model

import "time"

// RefreshToken is the model of a refresh token
// It enforces the relationship between a public key and a refresh token
type RefreshToken struct {
	// PubKeyID is the ID of the public key associated with the refresh token
	PubKeyID string
	// ID is the actual refresh token value
	ID        RefreshTokenValue
	TTL       time.Duration
	RevokedAt time.Time
}

type RefreshTokenValue interface {
	Time() time.Time
}

func (r *RefreshToken) IsExpired() bool {
	expiresAt := r.ID.Time().Add(r.TTL)
	if time.Now().After(expiresAt) {
		return true
	}

	return false
}

func (r *RefreshToken) IsRevoked() bool {
	issuedAt := r.ID.Time()
	if issuedAt.After(r.RevokedAt) {
		return true
	}
	return false
}
