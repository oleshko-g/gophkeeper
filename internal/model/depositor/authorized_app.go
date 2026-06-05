package depositor

import (
	"time"

	"github.com/google/uuid"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

type AuthorizedApp struct {
	ID uuidv7.UUID[uuid.UUID] `json:"authorized_app_id"`
	PubKey
}

func (a *AuthorizedApp) AuthorizedAt() time.Time {
	return a.ID.Time()
}
