package transform

func ValueToPtr[T any](V T) *T {
	return &V
}
