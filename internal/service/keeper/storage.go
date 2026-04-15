package keeper

//go:generate moq -out storage_mock.go . Storage
type Storage interface {
}
