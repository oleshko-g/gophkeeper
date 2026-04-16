package storage

//go:generate moq -out storage_mock.go . Depositor
type Depositor interface {
}
