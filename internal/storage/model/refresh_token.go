package model

// RefreshToken is the model of a refresh token
// It enforces the relationship between a public key and a refresh token
type RefreshToken struct {
	// ID is the unique identifier of the refresh token
	ID string
	// PubKeyID is the ID of the public key associated with the refresh token
	PubKeyID string
	// RefreshToken is the actual refresh token value
	RefreshToken string
}
