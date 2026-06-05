package service

import "context"

type ctxKey int

const CtxKeyPublicKeyID ctxKey = 1

func PubKeyIDFromCtx(ctx context.Context) (string, bool) {
	publicKeyID, ok := ctx.Value(CtxKeyPublicKeyID).(string)
	return publicKeyID, ok
}
