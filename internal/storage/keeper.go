package storage

//go:generate moq -out storage_mock.go . Storage
type Keeper interface {
}
