package secret

import (
	"github.com/google/uuid"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	"github.com/oleshko-g/gophkeeper/internal/db/pgx/queries"
	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

// DepositedSecret is a stored secret.
type DepositedSecret struct {
	ID       uuidv7.UUID[uuid.UUID]
	PubKeyID uuidv7.UUID[uuid.UUID]
	Data
}

// Data is the secret.
type Data []byte

// FromProto converts a protobuf Secret to a [secret.Value] model.
func FromProto(in *pb.DepositSecretRequest) Data {
	return Data(in.GetPayload())
}

func FromQueryResult(in queries.DepositedSecret) *DepositedSecret {
	return &DepositedSecret{
		ID: uuidv7.UUID[uuid.UUID]{
			Value: in.ID,
		},
		PubKeyID: uuidv7.UUID[uuid.UUID]{
			Value: in.DepositorPubKeyID,
		},
		Data: in.EncryptedData,
	}
}
