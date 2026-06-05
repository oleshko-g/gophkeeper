package keeper

import (
	"errors"
)

var errEmptyPubKeyID = errors.New("empty public key ID")
var errEmptySecretID = errors.New("empty secret ID")
